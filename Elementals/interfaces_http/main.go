package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// logWriter es un tipo personalizado que implementa la interfaz io.Writer
// Se usa para escribir datos de forma personalizada
type logWriter struct{}

// main demuestra el uso de interfaces con HTTP
// Realiza una petición HTTP y usa un writer personalizado para mostrar la respuesta
func main() {
	// Realiza una petición GET a Google
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Crea una instancia del writer personalizado
	lw := logWriter{}

	// Copia el cuerpo de la respuesta al writer personalizado
	io.Copy(lw, resp.Body)
}

// Write implementa la interfaz io.Writer para logWriter
// Recibe un slice de bytes y los imprime como string
// Retorna el número de bytes escritos y cualquier error
func (logWriter) Write(bs []byte) (int, error) {
	return fmt.Println(string(bs))
}
