package main

import (
	"fmt"
	"time"
)

func BufferOrdered[T any](in <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		var (
			queue []T
			next  T
			nout  chan<- T
		)

		for in != nil || len(queue) > 0 {
			if len(queue) > 0 {
				next = queue[0]
				nout = out
			} else {
				nout = nil
			}

			select {
			case v, ok := <-in:
				if !ok {
					in = nil
					continue
				}
				queue = append(queue, v)
			case nout <- next:
				queue = queue[1:]
			}
		}
	}()

	return out
}

func main() {
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 8; i++ {
			in <- i
			fmt.Println("produced", i)
		}
	}()

	for v := range BufferOrdered(in) {
		fmt.Println("consumed", v)
		time.Sleep(80 * time.Millisecond)
	}
}
