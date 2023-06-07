package main

import "fmt"

// OCP
// open for extension, closed for modification
// Specification
const (
	red Color = iota
	green
	blue
)

const (
	small Size = iota
	medium
	large
)

type Color int
type Size int

type Product struct {
	name  string
	color Color
	size  Size
}

type Filter struct {
}

func (f *Filter) FilterByColor(products []Product, color Color) []*Product {
	result := make([]*Product, 0)

	for i, v := range products {
		if v.color == color {
			result = append(result, &products[i])
		}
	}

	return result
}

func main() {
	apple := Product{
		"Apple",
		green,
		small,
	}
	tree := Product{
		name:  "Tree",
		color: green,
		size:  large,
	}
	house := Product{
		name:  "House",
		color: blue,
		size:  large,
	}

	products := []Product{apple, tree, house}

	fmt.Printf("Green products (old):\n")

	f := Filter{}

	for _, v := range f.FilterByColor(products, green) {
		fmt.Printf(" - %s is green\n", v.name)
	}
}
