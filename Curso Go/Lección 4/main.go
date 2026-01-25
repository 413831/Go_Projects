package main

import "fmt"

// * -> puntero | operador de indirección
// & -> dirección de memoria | operador de dirección

var a = 1

func incrementar(numero *int) {
	*numero++
}

func main() {
	valor := 10
	fmt.Println("Valor antes de incrementar: ", valor)

	incrementar(&valor)

	fmt.Println("Valor después de incrementar: ", valor)

	// new()
	puntero := new(int) // puntero int inicializado en 0
	fmt.Println("Puntero: ", puntero)
	fmt.Println("Valor inicial con new: ", *puntero)
}
