package main

import (
	"fmt"
)

func SumarPares(slice []int) (r int) {
	for _, valor := range slice {
		if valor%2 == 0 {
			r += valor
		}
	}
	return
}

func ejer1() {
	var n int
	fmt.Println("Digite cant de valores:")
	fmt.Scan(&n)
	ValidarNumero(n)

	slice := make([]int, n)

	for i := 0; i < n; i++ {
		fmt.Println("Digite un valor:")
		fmt.Scan(&slice[i])
	}

	fmt.Println("Suma de los pares:", SumarPares(slice))
}
