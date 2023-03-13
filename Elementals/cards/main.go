package main

func main() {
	// -- EXAMPLE CODE --
	// Initialization
	//var card string = "Ace of Spades"
	// card := newCard()
	//fmt.Println("Example" + card)

	// Slices
	// cards := deck{"Ace of Diamonds", newCard()}
	// cards = append(cards, "Six of Spades")

	// cards.print()
	// ---------------------------------------------

	cards := newDeck()

	//hand, remainingDeck := deal(cards, 5)
	//hand.print()
	//remainingDeck.print()
	//
	// fmt.Println(cards.toString())
	// cards.saveToFile("my_cards")
	// cardsNew := newDeckFromFile("my_card")
	// cardsNew.print()
	//
	cards.shuffle()
	cards.print()
}

/**
func newCard() string {
	return "Five of Diamonds"
}**/
