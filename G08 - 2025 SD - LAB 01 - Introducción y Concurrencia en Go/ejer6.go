package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var WgAnillo sync.WaitGroup

func NodoAnilo(id int, in <-chan string, out chan<- string, done <-chan struct{}) {
	for {
		select {
		case msj := <-in:
			fmt.Println(msj)
			time.Sleep(1 * time.Second)
			out <- fmt.Sprintf("Ping desde nodo %d", id+1)
		case <-done:
			fmt.Printf("Terminó el tiempo en el nodo %d\n", id+1)
			//close(out)
			WgAnillo.Done()
			return
		}
	}
}

// En este caso, se tomó main como una gorutina más (nodo).
// Ya que, técnicamente se ejecuta sobre un hilo principal

func ejer6() {
	done := make(chan struct{})
	canales := make([]chan string, 5)

	go func() {
		time.Sleep(1 * time.Minute)
		close(done)
	}()

	WgAnillo.Add(4)
	for i := 0; i < 5; i++ {
		canales[i] = make(chan string, 1)
	}

	for i := 1; i < 5; i++ {
		go NodoAnilo(i, canales[i-1], canales[i], done)
	}

	canales[0] <- "Ping desde nodo 1"

	for {
		select {
		case msj := <-canales[4]:
			fmt.Println(msj)
			time.Sleep(1 * time.Second)
			canales[0] <- "Ping desde nodo 1"
		case <-done:
			fmt.Println("Terminó el tiempo en el nodo 1")
			WgAnillo.Wait()
			os.Exit(0)
		}
	}
}
