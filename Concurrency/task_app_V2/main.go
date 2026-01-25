package main

import (
	"fmt"
	"time"
)

func main() {
	buf := NewBuffer()

	// consumer
	go func() {
		buf.mu.Lock()

		for len(buf.items) == 0 {
			buf.cond.Wait()
		}

		item := buf.items[0]
		buf.items = buf.items[1:] // pop
		buf.mu.Unlock()

		fmt.Printf("Consumiendo %d\n", item)
	}()

	time.Sleep(500 * time.Millisecond)
	buf.mu.Lock()
	buf.items = append(buf.items, 42) // push
	buf.mu.Unlock()
	buf.cond.Signal() // notify to consumer
}
