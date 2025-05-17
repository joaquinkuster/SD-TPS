package main

import (
	"fmt"
	"sync"
	"time"
)

func Suscriptor(id int, in <-chan string, done <-chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			fmt.Printf("Finalizó el suscriptor %d\n", id+1)
			wg.Done()
			return
		case msj := <-in:
			fmt.Printf("Suscriptor %d recibió: %s\n", id+1, msj)
		}
	}
}

func Publicador(canales []chan string, done <-chan struct{}, wg *sync.WaitGroup) {
	var i int
	for {
		select {
		case <-done:
			fmt.Println("Finalizó el publicador")
			wg.Done()
			return
		default:
			time.Sleep(1 * time.Second)
			i++
			for _, out := range canales {
				out <- fmt.Sprintf("Evento%d", i)
			}
		}
	}
}

// Con las gorutinas definidas en funciones

func ejer9_2() {
	var wg sync.WaitGroup

	done := make(chan struct{})

	canales := make([]chan string, 3)
	for i := 0; i < 3; i++ {
		canales[i] = make(chan string, 1)
	}

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go Suscriptor(i, canales[i], done, &wg)
	}

	wg.Add(1)
	go Publicador(canales, done, &wg)

	fmt.Println("Presione enter para terminar...")
	fmt.Scanln()
	close(done)

	wg.Wait()
	fmt.Println("Terminó el tiempo")
}
