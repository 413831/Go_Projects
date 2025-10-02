package main

// main es la función principal que demuestra el uso del sistema de cartas
// Crea un nuevo mazo, lo mezcla y lo imprime
func main() {
	// -- CÓDIGO DE EJEMPLO COMENTADO --
	// Inicialización básica
	//var card string = "Ace of Spades"
	// card := newCard()
	//fmt.Println("Example" + card)

	// Uso de slices
	// cards := deck{"Ace of Diamonds", newCard()}
	// cards = append(cards, "Six of Spades")
	// cards.print()
	// ---------------------------------------------

	// Crea un nuevo mazo de cartas
	cards := newDeck()

	// Ejemplos de uso comentados:
	//hand, remainingDeck := deal(cards, 5)
	//hand.print()
	//remainingDeck.print()
	//
	// fmt.Println(cards.toString())
	// cards.saveToFile("my_cards")
	// cardsNew := newDeckFromFile("my_card")
	// cardsNew.print()
	//

	// Mezcla las cartas y las imprime
	cards.shuffle()
	cards.print()
}

/**
// newCard es una función de ejemplo que retorna una carta específica
func newCard() string {
	return "Five of Diamonds"
}
**/
