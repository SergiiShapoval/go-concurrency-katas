package main

import (
	"fmt"
	"sync"
)

func merge[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	// TODO: forward all input values and close out after wg.Wait()
	_ = wg
	return out
}

func main() {
	makeStream := func(vals ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, v := range vals {
				out <- v
			}
		}()
		return out
	}

	for v := range merge(makeStream(1, 2), makeStream(3, 4), makeStream(5, 6)) {
		fmt.Println(v)
	}
}
