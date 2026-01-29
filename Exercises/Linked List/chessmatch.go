package main

// Move the element of list
type Move struct {
	data string
	next *Move
}

func NewMove(data string, next *Move) *Move {
	return &Move{
		data: data,
		next: next,
	}
}

// ChessMatch the linked list storing chess moves
type ChessMatch struct {
	head *Move
}

func NewChessMatch() *ChessMatch {
	return &ChessMatch{}
}

func (c *ChessMatch) getAt(index int) *Move {
	curr := c.head
	pos := 0

	for pos < index && curr != nil {
		curr = curr.next
		pos++
	}

	return curr
}

func (c *ChessMatch) getLast() *Move {
	return nil
}

func (c *ChessMatch) insertAt(index int, data string) {
	if c.head == nil {
		c.head = NewMove(data, nil)

		return
	}

	prev := c.getAt(index - 1)
	prev.next = NewMove(data, prev.next)
}

func (c *ChessMatch) removeAt(index int) {

}

func (c *ChessMatch) forEach(predicate func(*Move)) {

}
