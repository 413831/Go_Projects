package main

import "fmt"

func fibonacciOne(position int) int {
	seq := []int{0, 1}

	for len(seq) <= position {
		seq = append(seq, seq[len(seq)-1]+seq[len(seq)-2])
	}

	return seq[len(seq)-1]
}

func fibonacciTwo(position int, seq []int) int {
	if seq == nil {
		seq = []int{0, 1}
	}

	for len(seq) <= position {
		seq = append(seq, seq[len(seq)-1]+seq[len(seq)-2])

		return fibonacciTwo(position, seq)
	}

	return seq[len(seq)-1]
}

func fibonacciThree(n int) int {
	fmt.Println(n)
	if n < 2 {
		return n
	}

	a := fibonacciThree(n - 1)
	b := fibonacciThree(n - 2)

	return a + b
}
