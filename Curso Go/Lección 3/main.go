package main

import (
	"errors"
	"fmt"
)

// Función clásica
func suma(a int, b int) int {
	return a + b
}

// Función que devuelve más de un valor
func dividir(a float64, b float64) (float64, error) {
	if b == 0 {
		fmt.Errorf("No se puede dividir por 0")
		return 0, errors.New("no se puede dividir por 0")
	}

	cociente := a / b

	return cociente, nil
}

// Función con número variable de argumentos
func imprimirNombres(nombres ...string) {
	for _, nombre := range nombres {
		fmt.Println(nombre)
	}
}

// Función closure
func contador() func() int {
	c := 0

	return func() int {
		c++
		return c
	}
}

type Rectangulo struct {
	Ancho float64
	Alto  float64
}

func (r Rectangulo) Area() float64 {
	return r.Alto * r.Ancho
}

func main() {
	cociente, err := dividir(10, 3)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Cociente de la división: %f\n", cociente)

	cont := contador()
	fmt.Println("contador: ", cont())
	fmt.Println("contador: ", cont())
	fmt.Println("contador: ", cont())
	fmt.Println("contador: ", cont())
	fmt.Println("contador: ", cont())
	fmt.Println("contador: ", cont())

	rect := Rectangulo{Ancho: 10, Alto: 5}
	fmt.Println("Area del rectangulo: ", rect.Area())
}
