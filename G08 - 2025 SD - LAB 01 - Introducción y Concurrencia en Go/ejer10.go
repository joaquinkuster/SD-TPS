package main

import (
	"fmt"
	"sync"
)

var (
	x       int
	cerrojo sync.Mutex
)

func incrementar(wg *sync.WaitGroup) {
	cerrojo.Lock()
	defer cerrojo.Unlock()
	x += 5
	wg.Done()
}

func ejer10() {
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go incrementar(&wg)
	}

	wg.Wait()
	fmt.Println("El valor de x es:", x)
}
