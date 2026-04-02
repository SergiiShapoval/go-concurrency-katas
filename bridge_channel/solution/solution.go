package main

import "fmt"

func orDone[T any](done <-chan struct{}, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for stream := range orDone(done, chanStream) {
			for v := range orDone(done, stream) {
				for {
					select {
					case <-done:
						return
					default:
					}

					select {
					case out <- v:
						goto nextValue
					default:
					}
				}
			nextValue:
			}
		}
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
