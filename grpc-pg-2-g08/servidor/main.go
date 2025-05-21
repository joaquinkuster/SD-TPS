package main

import (
	"fmt"
	"grpc-pg-2/proto"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedMonitorServer
	mu             sync.Mutex
	ultimaVista    map[string]time.Time
	estadoAnterior map[string]string
}

func (s *servidor) EnviarHeartbeat(stream proto.Monitor_EnviarHeartbeatServer) error {
	var nodoId string

	for {
		hb, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Ack{Mensaje: "Stream cerrado"})
		}
		if err != nil {
			log.Printf("Error en stream: %v", err)
			return err
		}

		nodoId = hb.NodoId

		s.mu.Lock()
		s.ultimaVista[nodoId] = time.Unix(hb.MarcaTiempo, 0)
		s.mu.Unlock()

		log.Printf("[HEARTBEAT] %v %v", nodoId, hb.MarcaTiempo)
	}
}

func (s *servidor) detectorFallas(intervalo time.Duration) {
	for {
		time.Sleep(intervalo)
		s.mu.Lock()
		ahora := time.Now()

		for nodo, ultimo := range s.ultimaVista {
			segundos := ahora.Sub(ultimo).Seconds()
			estado := ""

			if segundos < 3*intervalo.Seconds() {
				estado = "activo"
			} else {
				estado = "inactivo"
			}

			estadoPrevio, existe := s.estadoAnterior[nodo]

			if !existe {
				log.Printf("🔵 Nodo %v: CONECTADO por primera vez hace %.0fs", nodo, segundos)
			} else if estadoPrevio == "inactivo" && estado == "activo" {
				log.Printf("🟣 Nodo %v: REACTIVADO hace %.0fs", nodo, segundos)
			} else if estadoPrevio == "activo" && estado == "inactivo" {
				log.Printf("🔴 Nodo %v: CAÍDO hace %.0fs", nodo, segundos)
			} else if estado == "inactivo" {
				log.Printf("🟠 Nodo %v: INACTIVO desde hace %.0fs", nodo, segundos)
			} else {
				log.Printf("🟢 Nodo %v: ACTIVO desde hace %.0fs", nodo, segundos)
			}

			// Actualizamos el estado actual para el próximo ciclo
			s.estadoAnterior[nodo] = estado
		}

		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	servidor := &servidor{
		ultimaVista:    make(map[string]time.Time),
		estadoAnterior: make(map[string]string),
	}

	proto.RegisterMonitorServer(s, servidor)

	go servidor.detectorFallas(5 * time.Second)

	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
