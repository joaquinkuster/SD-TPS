package main

import (
	"context"
	"encoding/binary"
	"grpc-pg-3/proto"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServidorReplica implementa proto.ReplicaServer
type ServidorReplica struct {
	proto.UnimplementedReplicaServer
	mu           sync.Mutex
	almacen      map[string]ValorConVersion // map[clave]ValorConVersion
	relojVector  VectorReloj
	idReplica    int                   // 0, 1 o 2
	clientesPeer []proto.ReplicaClient // stubs gRPC a las otras réplicas
}

// NewServidorReplica crea una instancia de ServidorReplica
// idReplica: 0, 1 o 2
// peerAddrs: direcciones gRPC de los otros dos peers (ej.: []string{":50052", ":50053"})
func NewServidorReplica(idReplica int, peerAddrs []string) *ServidorReplica {
	peers := make([]proto.ReplicaClient, 0, len(peerAddrs))
	for _, addr := range peerAddrs {
		// establecemos la conexión una sola vez
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("No se pudo conectar a peer %s: %v", addr, err)
		}
		peers = append(peers, proto.NewReplicaClient(conn))
	}

	return &ServidorReplica{
		almacen:      make(map[string]ValorConVersion),
		relojVector:  VectorReloj{},
		idReplica:    idReplica,
		clientesPeer: peers,
	}
}

// VectorReloj representa un reloj vectorial de longitud 3 (tres réplicas).
type VectorReloj [3]uint64

// ValorConVersion guarda el valor y su reloj vectorial asociado.
type ValorConVersion struct {
	Valor       []byte
	RelojVector VectorReloj
}

// Incrementar aumenta en 1 el componente correspondiente a la réplica que llama.
func (vr *VectorReloj) Incrementar(idReplica int) {
	(*vr)[idReplica] += 1
}

// Fusionar toma el máximo elemento a elemento entre dos vectores.
func (vr *VectorReloj) Fusionar(otro VectorReloj) {
	for i := range *vr {
		if otro[i] > (*vr)[i] {
			(*vr)[i] = otro[i]
		}
	}
}

// AntesDe devuelve true si vr < otro en el sentido estricto (strictly less).
func (vr VectorReloj) AntesDe(otro VectorReloj) bool {
	strictlyLess := false
	for i := range vr {
		if vr[i] > otro[i] {
			return false
		}
		if vr[i] < otro[i] {
			strictlyLess = true
		}
	}
	return strictlyLess
}

// encodeVector serializa el VectorReloj a []byte para enviarlo por gRPC.
func encodeVector(vr VectorReloj) []byte {
	buf := make([]byte, 8*3)
	for i := 0; i < 3; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], vr[i])
	}
	return buf
}

// decodeVector convierte []byte a VectorReloj.
func decodeVector(b []byte) VectorReloj {
	var vr VectorReloj
	for i := 0; i < 3; i++ {
		vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return vr
}

// GuardarLocal recibe la petición del Coordinador para almacenar clave/valor.
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *proto.SolicitudGuardar) (*proto.RespuestaGuardar, error) {
	r.mu.Lock()

	// 1. Incrementar nuestro componente del reloj vectorial
	r.relojVector.Incrementar(r.idReplica)

	// 2. Guardar en el mapa local
	r.almacen[req.Clave] = ValorConVersion{
		Valor:       req.Valor,
		RelojVector: r.relojVector,
	}

	// 3. Construir mutación para replicar a peers
	mutacion := &proto.Mutacion{
		Tipo:        proto.Mutacion_GUARDAR,
		Clave:       req.Clave,
		Valor:       req.Valor,
		RelojVector: encodeVector(r.relojVector),
	}

	r.mu.Unlock()

	// 4. Replicar asíncronamente a cada peer
	wg := sync.WaitGroup{}
	for _, cliente := range r.clientesPeer {
		wg.Add(1)
		go func(m *proto.Mutacion, client proto.ReplicaClient) {
			defer wg.Done()
			ctx2, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			ack, err := client.ReplicarMutacion(ctx2, m)
			defer cancel()
			if err == nil && ack.Ok {
				relojVectorialPeer := decodeVector(ack.RelojVectorAck)
				r.mu.Lock()
				r.relojVector.Fusionar(relojVectorialPeer)
				r.mu.Unlock()
			}
		}(mutacion, cliente)
	}

	// 5. Responder al Coordinador con el nuevo reloj vectorial
	wg.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Log de la operación
	log.Printf("Replica %d: Guardar clave=%s, valor=%s, reloj=%v", r.idReplica, req.Clave, r.almacen[req.Clave].Valor, r.relojVector)

	respuestaGuardar := &proto.RespuestaGuardar{
		Exito:            true,
		NuevoRelojVector: encodeVector(r.relojVector),
	}

	return respuestaGuardar, nil
}

// EliminarLocal recibe la petición del Coordinador para eliminar una clave.
// EliminarLocal elimina la clave del mapa local y replica la mutación a los peers.
func (r *ServidorReplica) EliminarLocal(ctx context.Context, req *proto.SolicitudEliminar) (*proto.RespuestaEliminar, error) {
	r.mu.Lock()

	// 1. Incrementar nuestro componente del reloj vectorial
	r.relojVector.Incrementar(r.idReplica)

	// 2. Borrar del mapa local (si existe)
	delete(r.almacen, req.Clave)

	// 3. Construir mutación de eliminación
	mutacion := &proto.Mutacion{
		Tipo:        proto.Mutacion_ELIMINAR,
		Clave:       req.Clave,
		Valor:       nil,
		RelojVector: encodeVector(r.relojVector),
	}

	r.mu.Unlock()

	// 4. Replicar a peers
	wg := sync.WaitGroup{}
	for _, cliente := range r.clientesPeer {
		wg.Add(1)
		go func(m *proto.Mutacion, client proto.ReplicaClient) {
			defer wg.Done()
			ctx2, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			ack, err := client.ReplicarMutacion(ctx2, m)
			defer cancel()
			if err == nil && ack.Ok {
				relojVectorialPeer := decodeVector(ack.RelojVectorAck)
				r.mu.Lock()
				r.relojVector.Fusionar(relojVectorialPeer)
				r.mu.Unlock()
			}
		}(mutacion, cliente)
	}

	// 5. Responder al Coordinador
	wg.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Log de la operación
	log.Printf("Replica %d: Eliminar clave=%s, reloj=%v", r.idReplica, req.Clave, r.relojVector)

	respuestaEliminar := &proto.RespuestaEliminar{
		Exito:            true,
		NuevoRelojVector: encodeVector(r.relojVector),
	}

	return respuestaEliminar, nil
}

// ObtenerLocal retorna el valor y reloj vectorial de una clave en esta réplica.
func (r *ServidorReplica) ObtenerLocal(ctx context.Context, req *proto.SolicitudObtener) (*proto.RespuestaObtener, error) {
	// Implementación básica: buscar la clave y devolver el valor y el reloj vectorial, o error si no existe
	r.mu.Lock()
	defer r.mu.Unlock()

	valor, ok := r.almacen[req.Clave]
	if !ok {
		return &proto.RespuestaObtener{
			Valor:       nil,
			RelojVector: encodeVector(VectorReloj{}),
			Existe:      false,
		}, os.ErrNotExist
	}

	// Log de la operación
	log.Printf("Replica %d: Obtener clave=%s, valor=%s, reloj=%v", r.idReplica, req.Clave, string(valor.Valor), valor.RelojVector)

	return &proto.RespuestaObtener{
		Valor:       valor.Valor,
		RelojVector: encodeVector(valor.RelojVector),
		Existe:      true,
	}, nil
}

func (r *ServidorReplica) ReplicarMutacion(ctx context.Context, m *proto.Mutacion) (*proto.Reconocimiento, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Decodificar el reloj vectorial de la mutación
	relojVectorialRemoto := decodeVector(m.RelojVector)

	// Considere que si no existía, o la mutación es “más nueva”, sobrescribir
	// Si existe y nuestra versión local está “por delante”, ignoramos (conflicto resuelto a favor local).
	if valor, ok := r.almacen[m.Clave]; ok {
		if !valor.RelojVector.AntesDe(relojVectorialRemoto) {
			// Si el reloj vectorial remoto no está antes que el nuestro, hay un conflicto
			if !relojVectorialRemoto.AntesDe(valor.RelojVector) {
				log.Printf("Replica %d: Conflicto concurrente en clave=%s", r.idReplica, m.Clave)
			} else {
				// Si nuestra versión es más reciente, ignoramos la mutación
				log.Printf("Replica %d: Ignorar mutación de clave=%s, reloj local=%v, reloj remoto=%v", r.idReplica, m.Clave, valor.RelojVector, relojVectorialRemoto)
			}
			// Responder con nuestro reloj vectorial actualizado
			return &proto.Reconocimiento{
				Ok:             false,
				RelojVectorAck: encodeVector(r.relojVector),
			}, nil
		}
	}

	if m.Tipo == proto.Mutacion_GUARDAR {
		r.almacen[m.Clave] = ValorConVersion{Valor: m.Valor, RelojVector: relojVectorialRemoto}
	} else {
		delete(r.almacen, m.Clave)
	}

	// 2. Fusionar nuestro reloj vectorial con el remoto
	r.relojVector.Fusionar(relojVectorialRemoto)

	// 3. Responder con nuestro reloj actualizado

	// Log de la replicación
	log.Printf("Replica %d: ReplicarMutacion clave=%s, reloj local=%v, reloj remoto=%v", r.idReplica, m.Clave, r.relojVector, relojVectorialRemoto)

	return &proto.Reconocimiento{
		Ok:             true,
		RelojVectorAck: encodeVector(r.relojVector),
	}, nil
}

func main() {
	// Uso: go run servidor_replica.go <idReplica> <direccionEscucha> <peer1> <peer2>
	// Ejemplo: go run servidor_replica.go 0 :50051 :50052 :50053
	if len(os.Args) != 5 {
		log.Fatalf("Uso: %s <idReplica> <direccionEscucha> <peer1> <peer2>", os.Args[0])
		return
	}

	idReplica, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error convirtiendo idReplica: %v", err)
	}
	addr := os.Args[2]
	peer1 := os.Args[3]
	peer2 := os.Args[4]
	peerAddrs := []string{peer1, peer2}

	// 1. Inicializar servidor gRPC
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	// 2. Crear instancia de ServidorReplica
	grpcServer := grpc.NewServer()
	replica := NewServidorReplica(idReplica, peerAddrs)

	proto.RegisterReplicaServer(grpcServer, replica)

	// 3. Iniciar servidor
	log.Printf("Réplica %d escuchando en %s", idReplica, addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
