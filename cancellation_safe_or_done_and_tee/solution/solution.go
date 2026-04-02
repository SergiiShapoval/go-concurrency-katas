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

func tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)
		for v := range orDone(done, in) {
			ch1, ch2 := out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case ch1 <- v:
					ch1 = nil
				case ch2 <- v:
					ch2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	done := make(chan struct{})
	in := make(chan int)

	go func() {
		defer close(in)
		for i := 1; i <= 4; i++ {
			in <- i
		}
	}()

	a, b := tee(done, in)
	for v := range a {
		fmt.Println("a:", v, "b:", <-b)
		if v == 3 {
			close(done)
			return
		}
	}
}
