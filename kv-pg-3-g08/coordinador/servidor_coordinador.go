package main

import (
	"context"
	"flag"
	"grpc-pg-3/proto"
	"log"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServidorCoordinador implementa proto.CoordinadorServer.
type ServidorCoordinador struct {
	proto.UnimplementedCoordinadorServer
	listaReplicas []string                       // ej: []string{":50051", ":50052", ":50053"}
	clientes      map[string]proto.ReplicaClient // stubs gRPC a las réplicas
	indiceRR      uint64                         // contador atómico para round-robin
}

// NewServidorCoordinador crea un Coordinador con direcciones de réplica.
func NewServidorCoordinador(replicas []string) *ServidorCoordinador {
	clientes := make(map[string]proto.ReplicaClient, len(replicas))
	for _, addr := range replicas {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("No se pudo conectar a réplica %s: %v", addr, err)
		}
		// NOTA: no cerramos conn aquí; lo dejamos vivo mientras dure el servidor.
		clientes[addr] = proto.NewReplicaClient(conn)
	}
	return &ServidorCoordinador{
		listaReplicas: replicas,
		clientes:      clientes,
		indiceRR:      0,
	}
}

// elegirReplicaParaEscritura: round-robin simple (ignora la clave).
func (c *ServidorCoordinador) elegirReplicaParaEscritura(_ string) string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// elegirReplicaParaLectura: también round-robin.
func (c *ServidorCoordinador) elegirReplicaParaLectura() string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// Obtener: redirige petición de lectura a una réplica.
func (c *ServidorCoordinador) Obtener(ctx context.Context, req *proto.SolicitudObtener) (*proto.RespuestaObtener, error) {
	// TODO: Implementar lógica de redirección a réplica
	replica := c.elegirReplicaParaLectura()
	cliente := c.clientes[replica]

	// Redirigir la solicitud a la réplica
	log.Printf("Redirigiendo petición de lectura a la réplica: %s", replica)
	respuesta, err := cliente.ObtenerLocal(ctx, req)
	if err != nil {
		log.Printf("Error al obtener de la réplica %s: %v", replica, err)
		return nil, err
	}

	return respuesta, nil
}

// Guardar: redirige petición de escritura a una réplica elegida.
func (c *ServidorCoordinador) Guardar(ctx context.Context, req *proto.SolicitudGuardar) (*proto.RespuestaGuardar, error) {
	// TODO: Implementar lógica de redirección a réplica para guardar
	replica := c.elegirReplicaParaEscritura(req.Clave)
	cliente := c.clientes[replica]

	// Redirigir la solicitud a la réplica
	log.Printf("Redirigiendo petición de escritura a la réplica: %s", replica)
	respuesta, err := cliente.GuardarLocal(ctx, req)
	if err != nil {
		log.Printf("Error al guardar en la réplica %s: %v", replica, err)
		return nil, err
	}

	return respuesta, nil
}

// Eliminar: redirige petición de eliminación a una réplica elegida.
func (c *ServidorCoordinador) Eliminar(ctx context.Context, req *proto.SolicitudEliminar) (*proto.RespuestaEliminar, error) {
	// TODO: Implementar lógica de redirección a réplica para eliminar
	replica := c.elegirReplicaParaEscritura(req.Clave)
	cliente := c.clientes[replica]

	// Redirigir la solicitud a la réplica
	log.Printf("Redirigiendo petición de eliminación a la réplica: %s", replica)
	respuesta, err := cliente.EliminarLocal(ctx, req)
	if err != nil {
		log.Printf("Error al eliminar en la réplica %s: %v", replica, err)
		return nil, err
	}

	return respuesta, nil
}

func main() {
	// Definir bandera para la dirección de escucha del Coordinador.
	listen := flag.String("listen", ":6000", "dirección para que escuche el Coordinador (p.ej., :6000)")
	flag.Parse()
	replicas := flag.Args()
	if len(replicas) < 3 {
		log.Fatalf("Debe proveer al menos 3 direcciones de réplicas, p.ej.: go run servidor_coordinador.go -listen :6000 :50051 :50052 :50053")
	}

	// Crear instancia del coordinador
	coordinador := NewServidorCoordinador(replicas)

	// Inicializar servidor gRPC
	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("No se pudo escuchar en %s: %v", *listen, err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterCoordinadorServer(grpcServer, coordinador)
	log.Printf("Coordinador escuchando en %s", *listen)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
