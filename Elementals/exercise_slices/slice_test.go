package main

import "testing"

func TestEvenOddDetection(t *testing.T) {
	testCases := []struct {
		number   int
		expected string
	}{
		{1, "odd"},
		{2, "even"},
		{3, "odd"},
		{4, "even"},
		{5, "odd"},
		{6, "even"},
		{7, "odd"},
		{8, "even"},
		{10, "even"},
	}

	for _, tc := range testCases {
		result := isEven(tc.number)
		expectedBool := tc.expected == "even"

		if result != expectedBool {
			t.Errorf("isEven(%d) = %v, expected %v", tc.number, result, expectedBool)
		}
	}
}

// isEven es una función auxiliar para los tests
// Determina si un número es par
func isEven(n int) bool {
	return n%2 == 0
}
