package main

import "fmt"

// School representa una estructura que contiene una lista de cursos
type School struct {
	courses []string
}

// main es la función principal que demuestra el uso de estructuras en Go
// Crea una instancia vacía de School y muestra la longitud de su slice de cursos
func main() {
	s := School{}
	fmt.Println(len(s.courses))
}
