package main

import "fmt"

// main es la función principal que demuestra el uso de slices en Go
// Itera sobre un slice de enteros y determina si cada número es par o impar
func main() {
	// Declara e inicializa un slice de enteros
	var intSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 10}

	// Itera sobre cada valor del slice
	for _, value := range intSlice {
		// Verifica si el número es par usando el operador módulo
		if value%2 == 0 {
			fmt.Printf("%d is even\n", value)
		} else {
			fmt.Printf("%d is odd\n", value)
		}
	}
}
