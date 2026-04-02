package main

import (
	"fmt"
	"sync"
)

func generator(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func worker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func fanOut(in <-chan int, n int) []<-chan int {
	workers := make([]<-chan int, 0, n)
	for i := 0; i < n; i++ {
		workers = append(workers, worker(in))
	}
	return workers
}

func fanIn(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	multiplex := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			out <- n
		}
	}

	wg.Add(len(chans))
	for _, ch := range chans {
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in := generator([]int{1, 2, 3, 4, 5, 6, 7, 8})
	for result := range fanIn(fanOut(in, 3)...) {
		fmt.Println(result)
	}
}
