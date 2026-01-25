package main

import "sync"

type Buffer struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func NewBuffer() *Buffer {
	b := &Buffer{items: []int{}}
	b.cond = sync.NewCond(&b.mu)
	return b
}
