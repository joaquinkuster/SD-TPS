package main

import (
	"fmt"
)

func ConvertirCelsiusFahrenheit(g float64) float64 {
	return (g * 9 / 5) + 32
}

func ConvertirFahrenheitCelsius(g float64) float64 {
	return (g - 32) * 5 / 9
}

func ejer4() {
	var op int
	fmt.Println("1. Celsius a Fahrenheit")
	fmt.Println("2. Fahrenheit a Celsius")
	fmt.Scan(&op)
	ValidarNumero(op)

	var g float64
	fmt.Println("Digite grados (temperatura):")
	fmt.Scan(&g)
	ValidarNumero(g)

	switch op {
	case 1:
		fmt.Printf("%.2f °C son %.2f °F", g, ConvertirCelsiusFahrenheit(g))
	case 2:
		fmt.Printf("%.2f °F son %.2f °C", g, ConvertirFahrenheitCelsius(g))
	default:
		fmt.Println("Opción inválida")
	}
}
