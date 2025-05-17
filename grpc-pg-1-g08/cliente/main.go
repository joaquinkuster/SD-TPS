package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	c := proto.NewServicioClient(conn)

	var wg sync.WaitGroup

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Second)
			defer cancel()

			nombre := fmt.Sprintf("Cliente %d", id)

			r, err := c.Hola(ctx, &proto.Requerimiento{Nombre: nombre})
			if err != nil {
				log.Fatalf("Error al llamar al servidor: %v", err)
			}
			log.Printf("Respuesta: %s", r.Mensaje)
		}(i)
	}

	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	r, err := c.ListadoPersonas(ctx, &proto.Vacio{})
	if err != nil {
		log.Fatalf("Error al llamar al servidro: %v", err)
	}

	fmt.Printf("Listado de personas saludadas (%d):\n", len(r.Personas))
	for _, p := range r.Personas {
		fmt.Println(p)
	}
}
