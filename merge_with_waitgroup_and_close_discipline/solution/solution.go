package main

import (
	"fmt"
	"sync"
)

func merge[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	forward := func(ch <-chan T) {
		defer wg.Done()
		for v := range ch {
			out <- v
		}
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

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
