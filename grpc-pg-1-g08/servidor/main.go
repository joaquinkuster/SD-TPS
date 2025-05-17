package main

import (
	"context"
	"grpc-pg-1/proto"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedServicioServer
	personas []string
	cerrojo  sync.RWMutex
}

func (s *servidor) Hola(ctx context.Context, req *proto.Requerimiento) (*proto.Respuesta, error) {
	s.cerrojo.Lock()
	s.personas = append(s.personas, req.Nombre)
	s.cerrojo.Unlock()

	log.Printf("Recibido: %s", req.Nombre)
	return &proto.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
}

func (s *servidor) ListadoPersonas(ctx context.Context, req *proto.Vacio) (*proto.Lista, error) {
	s.cerrojo.RLock()
	defer s.cerrojo.RUnlock()
	return &proto.Lista{Personas: s.personas}, nil
}

func main() {
	// 1) Abrir archivo para la salida del servidor
	f, err := os.Create("listado_servidor.txt")
	if err != nil {
		log.Fatalf("No se pudo crear listado_servidor.txt: %v", err)
	}
	defer f.Close()

	// 2) Redirigir todos los log.Printf al archivo
	log.SetOutput(f)

	// 3) Levantar gRPC
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterServicioServer(srv, &servidor{})

	log.Println("Servidor escuchando en :8000")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
