package main

import (
	"fmt"
	"time"
)

func BufferOrdered[T any](in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		// TODO: keep a FIFO queue and use a nil send channel when empty
	}()
	return out
}

func main() {
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 6; i++ {
			in <- i
			fmt.Println("produced", i)
		}
	}()

	for v := range BufferOrdered(in) {
		fmt.Println("consumed", v)
		time.Sleep(80 * time.Millisecond)
	}
}
