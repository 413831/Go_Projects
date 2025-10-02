package main

import "testing"

func TestPrintMap(t *testing.T) {
	// Test map with known values
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4bf745",
		"white": "#ffffff",
		"blue":  "#0000ff",
	}

	// Test that all expected colors are present
	expectedColors := []string{"red", "green", "white", "blue"}
	for _, color := range expectedColors {
		if _, exists := colors[color]; !exists {
			t.Errorf("Expected color %s to be present in map", color)
		}
	}

	// Test specific color values
	if colors["red"] != "#ff0000" {
		t.Errorf("Expected red to be #ff0000, got %s", colors["red"])
	}

	if colors["green"] != "#4bf745" {
		t.Errorf("Expected green to be #4bf745, got %s", colors["green"])
	}

	if colors["white"] != "#ffffff" {
		t.Errorf("Expected white to be #ffffff, got %s", colors["white"])
	}
}

func TestMapOperations(t *testing.T) {
	colors := make(map[string]string)

	// Test adding elements
	colors["red"] = "#ff0000"
	colors["blue"] = "#0000ff"

	if len(colors) != 2 {
		t.Errorf("Expected map length 2, got %d", len(colors))
	}

	// Test accessing elements
	if colors["red"] != "#ff0000" {
		t.Errorf("Expected red to be #ff0000, got %s", colors["red"])
	}

	// Test deleting elements
	delete(colors, "red")
	if len(colors) != 1 {
		t.Errorf("Expected map length 1 after deletion, got %d", len(colors))
	}
}
