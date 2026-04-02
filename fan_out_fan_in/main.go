package main

import "sync"

func generator(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// TODO: send all nums
	}()
	return out
}

func worker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// TODO: read numbers and send their squares
	}()
	return out
}

func fanOut(in <-chan int, n int) []<-chan int {
	workers := make([]<-chan int, 0, n)
	// TODO: append n worker channels
	return workers
}

func fanIn(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	// TODO: merge all channels into out and close out after wg.Wait()
	_ = wg
	return out
}

func main() {
	// TODO: wire generator, fanOut, fanIn, and print results
}
