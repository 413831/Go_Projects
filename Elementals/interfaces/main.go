package main

import "fmt"

// bot define una interfaz que requiere implementar el método getGreeting
// Cualquier tipo que implemente este método satisface la interfaz bot
type bot interface {
	getGreeting() string
}

// englishBot es un tipo vacío que implementa la interfaz bot
type englishBot struct{}

// spanishBot es un tipo vacío que implementa la interfaz bot
type spanishBot struct{}

// main demuestra el uso de interfaces en Go
// Crea instancias de diferentes bots y los pasa a una función genérica
func main() {
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

// printGreeting es una función genérica que acepta cualquier tipo que implemente la interfaz bot
// Recibe un parámetro de tipo bot y llama a su método getGreeting
func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}

// getGreeting implementa el método de la interfaz bot para englishBot
// Retorna un saludo en inglés
func (englishBot) getGreeting() string {
	return "Hi There!"
}

// getGreeting implementa el método de la interfaz bot para spanishBot
// Retorna un saludo en español
func (spanishBot) getGreeting() string {
	return "¡Hola!"
}
