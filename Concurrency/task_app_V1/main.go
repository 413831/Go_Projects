package main

import (
	"fmt"
	"sync"
)

var (
	completadas int
	mu          sync.Mutex
)

func tarea(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()
	completadas++
	mu.Unlock()

	fmt.Printf("Tarea %d completada\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go tarea(i, &wg)
	}

	wg.Wait()
	mu.Lock()
	fmt.Printf("%d tareas completadas\n", completadas)
	mu.Unlock()
}
