package main

import "testing"

func TestFibonacciOne(t *testing.T) {
	actual := fibonacciOne(6)

	if actual != 8 {
		t.Error("expected value of 8 at posiion 6.")
	}
}

func TestFibonacciTwo(t *testing.T) {
	actual := fibonacciTwo(6, nil)

	if actual != 8 {
		t.Error("expected value of 8 at posiion 6.")
	}
}

func TestFibonacciThree(t *testing.T) {
	actual := fibonacciThree(6)

	if actual != 8 {
		t.Error("expected value of 8 at posiion 6.")
	}
}
