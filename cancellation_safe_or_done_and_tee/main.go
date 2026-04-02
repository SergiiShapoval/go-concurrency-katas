package main

func orDone[T any](done <-chan struct{}, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		// TODO: forward values until done or input close
	}()
	return out
}

func tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		// TODO: duplicate each input value to both outputs
	}()
	return out1, out2
}

func main() {
	// TODO: demonstrate tee with early cancellation
}
