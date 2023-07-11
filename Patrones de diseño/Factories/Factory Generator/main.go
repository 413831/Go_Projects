package main

import "fmt"

type Employee struct {
	Name, Position string
	Income         int
}

// functional
func NewEmployeeFactory(position string, annualIncome int) func(name string) *Employee {
	return func(name string) *Employee {
		return &Employee{
			Name:     name,
			Position: position,
			Income:   annualIncome,
		}
	}
}

type EmployeeFactory struct {
	Position     string
	AnnualIncome int
}

func (f *EmployeeFactory) Create(name string) *Employee {
	return &Employee{
		Name:     name,
		Position: f.Position,
		Income:   f.AnnualIncome,
	}
}

func NewEmployeeFactory2(position string, annualIncome int) *EmployeeFactory {
	return &EmployeeFactory{
		Position:     position,
		AnnualIncome: annualIncome,
	}
}

func main() {
	developerFactory := NewEmployeeFactory("developer", 60000)
	managerFactory := NewEmployeeFactory("manager", 100000)

	developer := developerFactory("Adam")
	manager := managerFactory("Jane")

	fmt.Println(developer)
	fmt.Println(manager)

	bossFactory := NewEmployeeFactory2("CEO", 300000)
	boss := bossFactory.Create("Sam")
	fmt.Println(boss)
}
