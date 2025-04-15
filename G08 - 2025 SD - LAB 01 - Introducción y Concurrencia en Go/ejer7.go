package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Log struct {
	eventos []string
	cerrojo sync.Mutex
}

func (log *Log) RegistrarEvento(e string) {
	log.cerrojo.Lock()
	defer log.cerrojo.Unlock()
	log.eventos = append(log.eventos, e)
}

func ejer7() {
	var e []string
	log := Log{eventos: e}

	var wg sync.WaitGroup
	done := make(chan struct{})

	eventos := []string{
		"Temperatura alta",
		"Pérdida de conexión",
		"Baja energía",
		"Reinicio inesperado",
		"Lectura inválida",
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for {
				select {
				case <-done:
					fmt.Printf("Se finalizó la gorutina %d\n", i+1)
					wg.Done()
					return
				default:
					log.RegistrarEvento(fmt.Sprintf("Nodo -%d: %s", id+1, eventos[rand.Intn(len(eventos))]))
					//fmt.Println(log.eventos[len(log.eventos)-1])
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(i)
	}

	fmt.Println("Presione enter para terminar...")
	fmt.Scanln()
	close(done)

	wg.Wait()

	fmt.Println("Eventos registrados:")
	for _, evento := range log.eventos {
		fmt.Println(evento)
	}
}
