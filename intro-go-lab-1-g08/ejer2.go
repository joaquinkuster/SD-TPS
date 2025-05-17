package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ContarPalabras(txt string) int {
	return len(strings.Fields(txt))
}

func ejer2() {
	var txt string
	fmt.Println("Digite un texto:")
	// Usar bufio.Scanner para leer toda la línea
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	txt = scanner.Text()

	fmt.Println("Cant de palabras:", ContarPalabras(txt))
}
