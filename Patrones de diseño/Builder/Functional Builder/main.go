package main

import "fmt"

type (
	Employee struct {
		name, position string
	}

	employeeMod func(*Employee)

	EmployeeBuilder struct {
		actions []employeeMod
	}
)

func (b *EmployeeBuilder) Called(name string) *EmployeeBuilder {
	b.actions = append(b.actions, func(e *Employee) {
		e.name = name
	})

	return b
}

func (b *EmployeeBuilder) Build() *Employee {
	e := Employee{}

	for _, a := range b.actions {
		a(&e)
	}

	return &e
}

func (b *EmployeeBuilder) WorksAsA(position string) *EmployeeBuilder {
	b.actions = append(b.actions, func(e *Employee) {
		e.position = position
	})

	return b
}

func main() {
	b := EmployeeBuilder{}
	e := b.Called("Dmitri").WorksAsA("Software Engineer").Build()

	fmt.Println(&e)
}
