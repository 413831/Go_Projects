package algorithms

import "fmt"

func bubbleSort(arr []int) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				fmt.Printf("J : %d", arr[j])
				fmt.Printf("J + 1 : %d", arr[j])
				arr[j], arr[j+1] = arr[j+1], arr[j]

				fmt.Printf("swapped J : %d", arr[j])
				fmt.Printf("swapped J + 1 : %d", arr[j])
			}
		}
	}
}
