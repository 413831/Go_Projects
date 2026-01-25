package main

import "fmt"

func main() {
	defer fmt.Println("FIN")

	edad := 20

	// Condicional
	if edad >= 18 {
		fmt.Println("Eres mayor de edad")
	} else {
		fmt.Println("Eres menor de edad")
	}

	// Assertive / Negative programming
	if edad < 18 {
		fmt.Println("Eres menor de edad")
		//return
	}

	fmt.Println("Eres mayor de edad")

	// Bucle clásico
	for i := 0; i < 5; i++ {
		fmt.Printf("Iteración: %d \n", i)
	}

	// Bucle While
	n := 0
	for n < 3 {
		fmt.Printf("Iteración: %d\n", n)
		n++
	}

	// Bucle infinito
	n = 0

	for {
		n++

		if n == 5 {
			continue
		}

		fmt.Printf("n en bucle infinito: %d\n", n)

		if n >= 7 {
			break
		}
	}

	// Range
	slice := []string{"uno", "dos", "tres", "cuatro"}

	for index, value := range slice {
		fmt.Printf("Indice: %d, Valor: %s\n", index, value)
	}

	// Switch
	valor := 3

	switch valor {
	case 1:
		fmt.Println("Es 1")
	case 2:
		fmt.Println("Es 2")
	default:
		fmt.Println("No es 1 ni 2")
	}
}
