package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Address struct {
	Suite                 int
	Street, City, Country string
}

func (a *Address) DeepCopy() *Address {
	return &Address{
		Street:  a.Street,
		City:    a.City,
		Country: a.Country,
	}
}

type Employee struct {
	Name   string
	Office Address
}

func (p *Employee) DeepCopy() *Employee {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	_ = e.Encode(p)

	fmt.Println(string(b.Bytes()))

	d := gob.NewDecoder(&b)
	result := Employee{}
	_ = d.Decode(&result)

	return &result
}

var mainOffice = Employee{
	Name: "",
	Office: Address{
		Street:  "123 East Dr",
		City:    "London",
		Country: "UK",
	},
}

var auxOffice = Employee{
	Name: "",
	Office: Address{
		Street:  "66 West Dr",
		City:    "London",
		Country: "UK",
	},
}

func newEmployee(prototype *Employee, name string, suite int) *Employee {
	result := prototype.DeepCopy()
	result.Name = name
	result.Office.Suite = suite

	return result
}

func NewMainOfficeEmployee(name string, suite int) *Employee {
	return newEmployee(&mainOffice, name, suite)
}

func NewAuxOfficeEmployee(name string, suite int) *Employee {
	return newEmployee(&auxOffice, name, suite)
}

func main() {
	john := NewMainOfficeEmployee("John", 100)
	jane := NewAuxOfficeEmployee("Jane", 80)

	fmt.Println(john)
	fmt.Println(jane)
}
