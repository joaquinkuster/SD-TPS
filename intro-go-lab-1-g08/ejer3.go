package main

import (
	"fmt"
)

type Alumno struct {
	Nombre string
	Notas  []float64
}

func (a Alumno) Promedio() (r float64) {
	for _, nota := range a.Notas {
		r += nota
	}
	r /= float64(len(a.Notas))
	return
}

func ejer3() {
	var n int
	fmt.Println("Cant de alumnos:")
	fmt.Scan(&n)
	ValidarNumero(n)

	alumnos := make([]Alumno, n)

	n = 0
	fmt.Println("Cant de notas:")
	fmt.Scan(&n)
	ValidarNumero(n)

	for i := 0; i < len(alumnos); i++ {
		var nombre string
		notas := make([]float64, n)

		fmt.Printf("Nombre del alumno %d: ", i+1)
		fmt.Scan(&nombre)

		fmt.Printf("Digite las %d notas de %s: ", len(notas), nombre)
		for j := 0; j < len(notas); j++ {
			fmt.Scan(&notas[j])
		}

		a := Alumno{nombre, notas}
		alumnos[i] = a
	}

	for _, a := range alumnos {
		fmt.Printf("Promedio de %s: %.2f\n", a.Nombre, a.Promedio())
	}
}
