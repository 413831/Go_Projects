package algorithms

// Busca un elemento dividiendo la lista ordenada en dos
func binarySearch(arr []int, target int) int {
	var (
		low  = 0
		high = len(arr) - 1
	)

	for low <= high {
		mid := (low + high) / 2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}
