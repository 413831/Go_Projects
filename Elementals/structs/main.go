package main

import "fmt"

// contactInfo representa la información de contacto de una persona
type contactInfo struct {
	email   string // Dirección de correo electrónico
	zipCode int    // Código postal
}

// person representa una persona con su información personal y de contacto
type person struct {
	firstName string      // Nombre de la persona
	lastName  string      // Apellido de la persona
	contact   contactInfo // Información de contacto anidada
}

// main demuestra el uso de estructuras y métodos con punteros en Go
// Crea una persona, actualiza su nombre y la imprime
func main() {
	// Crea una instancia de person con información completa
	jim := person{
		firstName: "Jim",
		lastName:  "Party",
		contact: contactInfo{
			email:   "jim@gmail.com",
			zipCode: 94000,
		},
	}

	// Actualiza el nombre usando un método con puntero
	jim.updateName("Jimmy")

	// Imprime la información de la persona
	jim.print()
}

// updateName actualiza el nombre de la persona
// Recibe un puntero a person para poder modificar la estructura original
// Cambia el firstName por el nuevo nombre proporcionado
func (p *person) updateName(newFirstName string) {
	(*p).firstName = newFirstName
}

// print imprime la información completa de la persona
// Usa el formato %+v para mostrar los nombres de los campos y sus valores
func (p *person) print() {
	fmt.Printf("%+v", p)
}
