package main

// Ejecutar desde consola main.go seguido del ejercicio a probar
// go run main.go ejerX.go

import (
	"fmt"
	"os"
)

func ValidarNumero[T int | float64](n T) {
	if n == 0 {
		fmt.Println("Número inválido")
		os.Exit(1)
	}
}

func main() {
	//ejer1()
	//ejer2()
	//ejer3()
	//ejer4()
	//ejer5()
	//ejer6()
	//ejer6_2()
	//ejer7()
	//ejer8()
	//ejer9()
	//ejer9_2()
	//ejer10()
	//ejer11()
}
