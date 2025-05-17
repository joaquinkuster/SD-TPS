package main

import (
	"fmt"
	"os"
)

func ejer5() {
	if len(os.Args) < 2 {
		fmt.Println("Cant de arg inválida")
		os.Exit(1)
	}

	archivo, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error al abrir el archivo %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}

	fmt.Println("Contenido:", string(archivo))
}
