package main

import (
	"fmt"
	"time"
)

// Y, en este otro caso, no se tomó main como una gorutina (nodo) más

func ejer6_2() {
	done := make(chan struct{})
	canales := make([]chan string, 5)

	go func() {
		time.Sleep(1 * time.Minute)
		close(done)
	}()

	for i := 0; i < 5; i++ {
		canales[i] = make(chan string, 1)
	}

	WgAnillo.Add(5)
	go NodoAnilo(0, canales[4], canales[0], done)
	for i := 1; i < 5; i++ {
		go NodoAnilo(i, canales[i-1], canales[i], done)
	}

	canales[0] <- "Ping desde nodo 1"

	WgAnillo.Wait()
	fmt.Println("Terminó el tiempo")
}
