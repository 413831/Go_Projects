package main

import (
	"fmt"
)

type Address struct {
	Street, City, Country string
}

func (a *Address) DeepCopy() *Address {
	return &Address{
		Street:  a.Street,
		City:    a.City,
		Country: a.Country,
	}
}

type Person struct {
	Name    string
	Address *Address
	Friends []string
}

func (p *Person) DeepCopy() *Person {
	q := *p
	q.Address = p.Address.DeepCopy()
	copy(q.Friends, p.Friends)

	return &q
}

func main() {
	john := Person{
		"John",
		&Address{"123 London", "London", "UK"},
		[]string{"Chris", "Matt"},
	}

	// deep copying

	jane := john.DeepCopy()
	jane.Name = "Jane"
	jane.Address.Street = "321 Baker St"
	jane.Friends = append(jane.Friends, "Angela")

	fmt.Println(john, john.Address)
	fmt.Println(jane, jane.Address)
}
