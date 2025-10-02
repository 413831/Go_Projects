package main

import "fmt"

// shape define una interfaz que requiere implementar el método getArea
// Cualquier tipo que implemente este método satisface la interfaz shape
type shape interface {
	getArea() float64
}

// triangle representa un triángulo con base y altura
type triangle struct {
	base   float64 // Longitud de la base del triángulo
	height float64 // Altura del triángulo
}

// getArea calcula el área del triángulo usando la fórmula: (base * altura) / 2
// Implementa el método de la interfaz shape
func (t triangle) getArea() float64 {
	return 0.5 * t.base * t.height
}

// square representa un cuadrado con un lado
type square struct {
	lenght float64 // Longitud de un lado del cuadrado
}

// getArea calcula el área del cuadrado usando la fórmula: lado * lado
// Implementa el método de la interfaz shape
func (s square) getArea() float64 {
	return s.lenght * s.lenght
}

// printArea es una función genérica que acepta cualquier tipo que implemente la interfaz shape
// Recibe un parámetro de tipo shape e imprime su área
func printArea(s shape) {
	fmt.Println(s.getArea())
}

// main demuestra el uso de interfaces con formas geométricas
// Crea instancias de triángulo y cuadrado, y calcula sus áreas
func main() {
	// Crea un triángulo con base y altura de 10
	t := triangle{
		base:   10,
		height: 10,
	}

	// Crea un cuadrado con lado de 10
	s := square{lenght: 10}

	// Calcula e imprime las áreas usando la función genérica
	printArea(t)
	printArea(s)
}
