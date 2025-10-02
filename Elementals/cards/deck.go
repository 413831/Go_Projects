package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// deck es un tipo personalizado que representa un mazo de cartas
// Es un slice de strings donde cada string representa una carta
type deck []string

// newDeck crea un nuevo mazo de cartas completo
// Genera todas las combinaciones de palos y valores para crear 16 cartas
// Retorna un deck con todas las cartas generadas
func newDeck() deck {
	cards := deck{}
	cardSuits := []string{"Spades", "Diamond", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			cards = append(cards, value+" of "+suit)
		}
	}
	return cards
}

// print imprime todas las cartas del mazo con su índice
// Recibe un receiver de tipo deck y muestra cada carta numerada
func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

// deal divide un mazo en dos partes: mano y mazo restante
// Recibe el mazo original y el tamaño de la mano deseada
// Retorna la mano (primeras cartas) y el mazo restante
func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

// toString convierte el mazo en una cadena de texto separada por comas
// Útil para serializar el mazo antes de guardarlo en archivo
func (d deck) toString() string {
	return strings.Join(d, ",")
}

// saveToFile guarda el mazo en un archivo de texto
// Recibe el nombre del archivo donde guardar
// Retorna un error si la operación falla
func (d deck) saveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0666)
}

// newDeckFromFile carga un mazo desde un archivo de texto
// Recibe el nombre del archivo a leer
// Retorna un deck reconstruido desde el archivo
// Si hay error, termina el programa con os.Exit(1)
func newDeckFromFile(filename string) deck {
	bs, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	s := strings.Split(string(bs), ",")
	return deck(s)
}

// shuffle mezcla aleatoriamente las cartas del mazo
// Usa una semilla basada en el tiempo actual para generar aleatoriedad
// Intercambia cada carta con otra en una posición aleatoria
func (d deck) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := range d {
		newPosition := r.Intn(len(d) - 1)
		d[i], d[newPosition] = d[newPosition], d[i]
	}
}
