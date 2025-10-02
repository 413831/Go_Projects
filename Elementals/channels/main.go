package main

import (
	"fmt"
	"net/http"
	"time"
)

// main es la función principal que demuestra el uso de goroutines y canales
// Verifica el estado de múltiples sitios web de forma concurrente
func main() {
	// Lista de URLs a verificar
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	// Canal para comunicación entre goroutines
	c := make(chan string)

	// Inicia una goroutine para cada link
	for _, link := range links {
		go checkLink(link, c)
	}

	// Escucha continuamente los resultados del canal
	// Cuando recibe un resultado, inicia una nueva verificación después de 5 segundos
	for l := range c {
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l)
	}
}

// checkLink verifica si un sitio web está disponible
// Recibe la URL a verificar y un canal para enviar el resultado
// Realiza una petición HTTP GET y reporta el estado
func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")
		c <- link
		return
	}

	fmt.Println(link, "is up!")
	c <- link
}
