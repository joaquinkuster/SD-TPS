package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-pg-3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	coordinadorAddr := "localhost:6000"

	conn, err := grpc.Dial(coordinadorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al coordinador: %v", err)
	}
	defer conn.Close()

	cliente := proto.NewCoordinadorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Guardar la clave "usuario123" con valor "datosImportantes"
	valor := []byte("datosImportantes")
	clave := "usuario123"
	respuestaGuardar, err := cliente.Guardar(ctx, &proto.SolicitudGuardar{
		Clave:       clave,
		Valor:       valor,
		RelojVector: nil, // o []byte{}
	})
	if err != nil {
		log.Fatalf("Error al guardar: %v", err)
	} else if respuestaGuardar.Exito {
		log.Printf("Guardado OK.")
	} else {
		log.Fatalf("Error al guardar. La clave '%s' ya existe o error de replicación.", clave)
	}

	// 2. Obtener la clave "usuario123"
	respuestaObtener, err := cliente.Obtener(ctx, &proto.SolicitudObtener{
		Clave: clave,
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	} else if respuestaObtener.Existe {
		log.Printf("Obtenido: valor=%s", string(respuestaObtener.Valor))
	} else {
		log.Printf("La clave '%s' no existe.", clave)
	}

	// 3. Eliminar la clave, enviando el reloj vectorial recibido
	respuestaEliminar, err := cliente.Eliminar(ctx, &proto.SolicitudEliminar{
		Clave:       clave,
		RelojVector: respuestaObtener.RelojVector,
	})
	if err != nil {
		log.Fatalf("Error al eliminar: %v", err)
	} else if respuestaEliminar.Exito {
		fmt.Println("Eliminado OK.")
	} else {
		log.Fatalf("Error al eliminar. La clave '%s' ya existe o error de replicación.", clave)
	}

	// 4. Obtener nuevamente para verificar que ya no existe
	respuestaObtener, err = cliente.Obtener(ctx, &proto.SolicitudObtener{
		Clave: clave,
	})
	if err != nil {
		log.Fatalf("Error al obtener (tras eliminar): %v", err)
	} else if respuestaObtener.Existe {
		log.Printf("Obtenido tras eliminar: valor=%s", string(respuestaObtener.Valor))
	} else {
		fmt.Printf("La clave '%s' no existe tras eliminar.", clave)
	}
}
