package main

import "testing"

func TestTriangleArea(t *testing.T) {
	triangle := triangle{
		base:   10,
		height: 5,
	}

	expected := 25.0
	actual := triangle.getArea()

	if actual != expected {
		t.Errorf("Expected triangle area %f, got %f", expected, actual)
	}
}

func TestSquareArea(t *testing.T) {
	square := square{
		lenght: 5,
	}

	expected := 25.0
	actual := square.getArea()

	if actual != expected {
		t.Errorf("Expected square area %f, got %f", expected, actual)
	}
}

func TestPrintArea(t *testing.T) {
	triangle := triangle{
		base:   4,
		height: 3,
	}

	square := square{
		lenght: 4,
	}

	// Test triangle area calculation
	triangleArea := triangle.getArea()
	if triangleArea != 6.0 {
		t.Errorf("Expected triangle area 6.0, got %f", triangleArea)
	}

	// Test square area calculation
	squareArea := square.getArea()
	if squareArea != 16.0 {
		t.Errorf("Expected square area 16.0, got %f", squareArea)
	}
}
