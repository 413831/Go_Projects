package main

import "fmt"

type Person struct {
	Name     string
	Age      int
	EyeCount int
}

func NewPerson(name string, age int) *Person {
	if age < 16 {
		return nil
	}
	return &Person{name, age, 2}
}

func main() {
	p := NewPerson("Peter", 27)

	fmt.Println(p)
}
