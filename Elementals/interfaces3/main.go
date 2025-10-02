package main

import (
	"fmt"
	"io"
	"os"
)

// main demuestra el uso de interfaces para copiar archivos
// Abre un archivo especificado como argumento y lo copia a la salida estándar
func main() {
	// Abre el archivo especificado en el primer argumento de línea de comandos
	f, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	// Copia el contenido del archivo a la salida estándar
	// io.Copy usa interfaces para trabajar con cualquier tipo que implemente
	// io.Reader (el archivo) y io.Writer (os.Stdout)
	io.Copy(os.Stdout, f)
}
