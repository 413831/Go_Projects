package main

import "fmt"

type shape interface {
	getArea() float64
}

type triangle struct {
	base   float64
	height float64
}

func (t triangle) getArea() float64 {
	return 0.5 * t.base * t.height
}

type square struct {
	lenght float64
}

func (s square) getArea() float64 {
	return s.lenght * s.lenght
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}

func main() {
	t := triangle{
		base:   10,
		height: 10,
	}

	s := square{lenght: 10}

	printArea(t)
	printArea(s)
}


