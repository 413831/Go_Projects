package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Download struct {
	URL      string
	Filename string
}

func downloadFile(dl Download, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	// Download file data from URL
	resp, err := http.Get(dl.URL)
	if err != nil {
		results <- fmt.Sprintf("Error descargando %s: %v", dl.Filename, err)
		return
	}
	defer resp.Body.Close()

	// Create empty file to store downloaded data
	out, err := os.Create(dl.Filename)
	if err != nil {
		results <- fmt.Sprintf("Error creando %s: %v", dl.Filename, err)
	}
	defer out.Close()

	// Copy downloaded data to empty file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Error guardando %s: %v", dl.Filename, err)
	}

	results <- fmt.Sprintf("Descargado %s", dl.Filename)
}

func main() {
	fmt.Println("Hello, World!")

	downloads := []Download{
		{"https://jsonplaceholder.typicode.com/todos/1", "todo1.json"},
		{"https://jsonplaceholder.typicode.com/todos/2", "todo2.json"},
		{"https://jsonplaceholder.typicode.com/todos/3", "todo3.json"},
		{"https://jsonplaceholder.typicode.com/posts/1", "post1.json"},
		{"https://jsonplaceholder.typicode.com/posts/2", "post2.json"},
	}

	var wg sync.WaitGroup
	results := make(chan string, len(downloads))

	for _, dl := range downloads {
		wg.Add(1)
		go downloadFile(dl, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("Descargas completadas")
}
