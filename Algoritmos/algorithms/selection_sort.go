package algorithms

func selectionSort(arr []int) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		minIDx := i

		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIDx] {
				minIDx = j
			}
		}
		arr[i], arr[minIDx] = arr[minIDx], arr[i]
	}
}
