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
	cards.print()

}

/**
func newCard() string {
	return "Five of Diamonds"
}**/
