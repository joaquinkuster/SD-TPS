package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"net"

	"sync"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedServicioServer
	personas []string
	cerrojo  sync.RWMutex
}

func (s *servidor) Hola(ctx context.Context, req *proto.Requerimiento) (*proto.Respuesta, error) {
	log.Printf("Recibido: %s", req.Nombre)

	s.cerrojo.Lock()
	defer s.cerrojo.Unlock()
	s.personas = append(s.personas, req.Nombre)

	return &proto.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
}

func (s *servidor) ListadoPersonas(ctx context.Context, req *proto.Vacio) (*proto.Lista, error) {

	s.cerrojo.RLock()
	defer s.cerrojo.RUnlock()

	return &proto.Lista{Personas: s.personas}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterServicioServer(s, &servidor{})
	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
