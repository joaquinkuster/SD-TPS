package main

import (
	"fmt"
	"sync"
	"time"
)

// Con las gorutinas definidas en el main

func ejer9() {
	var wg sync.WaitGroup

	done := make(chan struct{})

	canales := make([]chan string, 3)
	for i := 0; i < 3; i++ {
		canales[i] = make(chan string, 1)
	}

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			for {
				select {
				case <-done:
					fmt.Printf("Finalizó el suscriptor %d\n", id+1)
					wg.Done()
					return
				case msj := <-canales[id]:
					fmt.Printf("Suscriptor %d recibió: %s\n", id+1, msj)
				}
			}
		}(i)
	}

	wg.Add(1)
	go func() {
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
				for _, canal := range canales {
					canal <- fmt.Sprintf("Evento%d", i)
				}
			}
		}
	}()

	fmt.Println("Presione enter para terminar...")
	fmt.Scanln()
	close(done)

	wg.Wait()
	fmt.Println("Terminó el tiempo")
}
