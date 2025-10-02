package main

import "fmt"

// main demuestra el uso de mapas en Go
// Crea un mapa de colores con sus códigos hexadecimales y los imprime
func main() {
	// Declara e inicializa un mapa de string a string
	// Las claves son nombres de colores y los valores son códigos hexadecimales
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4bf745",
		"white": "#ffffff",
	}

	// Llama a la función para imprimir el mapa
	printMap(colors)
}

// printMap itera sobre un mapa de colores e imprime cada par clave-valor
// Recibe un mapa donde las claves son nombres de colores y los valores son códigos hex
func printMap(c map[string]string) {
	for color, hex := range c {
		fmt.Println("Hex code for", color, "is", hex)
	}
}
