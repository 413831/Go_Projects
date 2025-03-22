package main

import "testing"

func TestGetAtRandomPosition(t *testing.T) {
	c := NewChessMatch()

	mData := "e4"
	m2Data := "Nf6"

	c.insertAt(0, mData)
	c.insertAt(1, m2Data)

	actual := c.getAt(1)

	if actual == nil || actual.data != m2Data {
		t.Errorf("expected element %s at position 1.", m2Data)
	}
}

func TestInsertAtRandomPosition(t *testing.T) {
	c := NewChessMatch()

	mData := "e4"
	m2Data := "Nf6"
	m3Data := "e6"
	m4Data := "Nd6"

	c.insertAt(0, mData)
	c.insertAt(1, m2Data)
	c.insertAt(2, m3Data)
	c.insertAt(3, m4Data)

	actual := c.getAt(2)

	if actual == nil || actual.data != m3Data {
		t.Errorf("expected element %s at position 2.", m3Data)
	}
}
