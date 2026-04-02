package main

import "fmt"

func orDone[T any](done <-chan struct{}, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		// TODO: forward until done or close
	}()
	return out
}

func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		// TODO: flatten the input stream-of-streams
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

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

	chanStream := make(chan (<-chan int))
	go func() {
		defer close(chanStream)
		chanStream <- makeStream(1, 2)
		chanStream <- makeStream(3)
		chanStream <- makeStream(4, 5, 6)
	}()

	for v := range bridge(done, chanStream) {
		fmt.Println(v)
	}
}
