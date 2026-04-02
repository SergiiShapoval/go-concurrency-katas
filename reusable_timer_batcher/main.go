package main

import (
	"fmt"
	"time"
)

func batch(in <-chan int, maxSize int, flushAfter time.Duration) <-chan []int {
	out := make(chan []int)
	go func() {
		defer close(out)
		// TODO: reuse one timer correctly
		var buffer []int
		ticker := time.NewTicker(flushAfter)
		defer ticker.Stop()
		for {
			select {
			case val, ok := <-in:
				if !ok {
					out <- buffer
					buffer = nil
					return
				}
				buffer = append(buffer, val)
				if len(buffer) >= maxSize {
					out <- buffer
					buffer = nil
				}
			case <-ticker.C:
				out <- buffer
				buffer = nil
				ticker.Reset(flushAfter)
			}
		}
	}()
	return out
}

func collectBatches(in <-chan int, maxSize int, flushAfter time.Duration) [][]int {
	// TODO: consume batch() and return all emitted batches

	results := batch(in, maxSize, flushAfter)
	var res [][]int
	for result := range results {
		res = append(res, result)
		fmt.Printf("received batch: %v\n", result)
	}
	return res
}

func main() {
	// TODO: feed a few bursts of integers into batch and print the slices

	in := make(chan int)
	go func() {

		for i := 0; i < 20; i++ {
			time.Sleep(time.Second)
			in <- i
		}
		close(in)
	}()

	res := collectBatches(in, 5, 2*time.Second)
	fmt.Printf("res: %v", res)

}
