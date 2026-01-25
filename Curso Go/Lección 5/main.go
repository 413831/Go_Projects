package main

import (
	"fmt"
	"math"
)

type Forma interface {
	Area() float64
}

type Circulo struct {
	Radio float64
}

func (c Circulo) Area() float64 {
	return math.Pi * c.Radio * c.Radio
}

func imprimirArea(f Forma) {
	fmt.Printf("El área es: %.2f\n", f.Area())
}

func main() {
	c := Circulo{Radio: 5}

	imprimirArea(c)
}
