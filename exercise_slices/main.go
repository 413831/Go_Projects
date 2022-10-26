package main

func main() {
	var intSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 10}

	for _, value := range intSlice {
		if value%2 == 0 {
			println("%d is even", value)
		} else {
			println("%d is odd", value)
		}
	}
}
