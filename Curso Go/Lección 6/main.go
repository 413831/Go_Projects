package main

import (
	"fmt"
	"sync"
	"time"
)

// go -> GoRoutine
// GoRoutine -> hilo de ejecución ligero virtual
// Channel -> Elemento para comunicar entre GoRoutine

func decirHola(canal chan<- string) {
	fmt.Println("Escribiendo mensaje...")
	time.Sleep(1 * time.Second)
	canal <- "Hola desde la GoRoutine"
}

func imprimirMensaje(canal <-chan string) {
	fmt.Println("Esperando mensaje...")
	msg := <-canal
	fmt.Println(msg)
}

func main() {
	// Ejemplo 1
	canal := make(chan string)
	go decirHola(canal)
	imprimirMensaje(canal)

	// Ejemplo 2
	canal2 := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			canal2 <- i
		}

		close(canal2)
	}()

	for num := range canal2 {
		fmt.Println("Número recibido: ", num)
	}

	// Ejemplo Mutex
	var contador int
	var mu sync.RWMutex

	// Writer
	go func() {
		for i := 0; i < 5; i++ {
			mu.Lock()
			contador++
			mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Reader
	for i := 0; i < 3; i++ {
		go func(id int) {
			for j := 0; j < 5; j++ {
				mu.RLock()
				fmt.Printf("Leyendo desde la GoRoutine %d: %d\n", id, contador)
				mu.RUnlock()
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	time.Sleep(2 * time.Second)
}
