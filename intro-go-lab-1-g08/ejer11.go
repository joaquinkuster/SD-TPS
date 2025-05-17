package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex

func a() {
	mutex.Lock()
	b()
	mutex.Unlock()
}

func b() {
	mutex.Lock() // esto genera el deadlock
	fmt.Println("Hola mundo")
	mutex.Unlock()
}

// Lo que sucede al ejecutar el programa es que cuando la función a() invoca
// a la función b(), el mutex ya se encuentra bloqueado. Entonces, b() intenta
// bloquear el mismo mutex, pero no puede acceder a él hasta que a() lo libere.
// Esto genera un interbloqueo (deadlock), ya que la función b() queda
// esperando indefinidamente por un recurso que nunca será liberado,
// dado que a() también está esperando a que b() termine.
//
// Como resultado, todas las goroutines quedan dormidas y el programa se
// detiene con el mensaje:
// fatal error: all goroutines are asleep - deadlock!

func ejer11() {
	a()
}
