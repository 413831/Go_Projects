package main

import "fmt"

type Person struct {
	// address
	StreetAddress, Postcode, City string

	// job
	CompanyName, Position string
	AnnualIncome          int
}

type PersonBuilder struct {
	person *Person
}

func (b *PersonBuilder) Lives() *PersonAddressBuilder {
	return &PersonAddressBuilder{*b}
}

func (b *PersonBuilder) Works() *PersonJobBuilder {
	return &PersonJobBuilder{*b}
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{&Person{}}
}

type PersonAddressBuilder struct {
	PersonBuilder
}
type PersonJobBuilder struct {
	PersonBuilder
}

func (ab *PersonAddressBuilder) At(streetAddress string) *PersonAddressBuilder {
	ab.person.StreetAddress = streetAddress

	return ab
}

func (ab *PersonAddressBuilder) In(city string) *PersonAddressBuilder {
	ab.person.City = city

	return ab
}

func (ab *PersonAddressBuilder) WithPostcode(postcode string) *PersonAddressBuilder {
	ab.person.Postcode = postcode

	return ab
}

func (pjb *PersonJobBuilder) At(companyName string) *PersonJobBuilder {
	pjb.person.CompanyName = companyName

	return pjb
}

func (pjb *PersonJobBuilder) AsA(position string) *PersonJobBuilder {
	pjb.person.Position = position

	return pjb
}

func (pjb *PersonJobBuilder) Earning(annualIncome int) *PersonJobBuilder {
	pjb.person.AnnualIncome = annualIncome

	return pjb
}

func (b *PersonBuilder) Build() *Person {
	return b.person
}

func main() {
	pb := NewPersonBuilder()
	pb.Lives().At("Cucha Cucha 208").In("Buenos Aires").WithPostcode("3520").
		Works().At("Globant").AsA("Developer").Earning(80000)

	person := pb.Build()

	fmt.Println(person)
}
