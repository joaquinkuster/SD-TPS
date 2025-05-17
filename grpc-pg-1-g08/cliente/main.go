package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1) Abrir archivo para toda la salida del cliente
	f, err := os.Create("salida_cliente.txt")
	if err != nil {
		log.Fatalf("No se pudo crear salida_cliente.txt: %v", err)
	}
	defer f.Close()

	// 2) Redirigir stdout y log al archivo
	log.SetOutput(f)
	os.Stdout = f

	// 3) Conectar al servidor
	conn, err := grpc.Dial("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	c := proto.NewServicioClient(conn)
	var wg sync.WaitGroup

	// 4) Lanzar 1000 goroutines que envían saludos
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			nombre := fmt.Sprintf("Cliente %d", id)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.Hola(ctx, &proto.Requerimiento{Nombre: nombre})
			if err != nil {
				log.Printf("Error al saludar a %s: %v", nombre, err)
				return
			}
			log.Printf("Respuesta: %s", r.Mensaje)
		}(i)
	}

	wg.Wait()

	// 5) Al terminar, pedir el listado y mostrarlo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lista, err := c.ListadoPersonas(ctx, &proto.Vacio{})
	if err != nil {
		log.Fatalf("Error al obtener listado: %v", err)
	}

	fmt.Printf("\nListado de personas saludadas (%d):\n", len(lista.Personas))
	for _, p := range lista.Personas {
		fmt.Println(p)
	}
}
