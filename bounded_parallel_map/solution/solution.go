package main

import (
	"fmt"
	"sync"
	"time"
)

func ParallelMap[T any, R any](in []T, limit int, fn func(T) R) []R {
	out := make([]R, len(in))
	sem := make(chan struct{}, limit)

	var wg sync.WaitGroup
	for i, v := range in {
		sem <- struct{}{}
		wg.Add(1)
		go func(idx int, value T) {
			defer wg.Done()
			defer func() { <-sem }()
			out[idx] = fn(value)
		}(i, v)
	}
	wg.Wait()
	return out
}

func main() {
	result := ParallelMap([]int{1, 2, 3, 4, 5}, 2, func(v int) int {
		time.Sleep(30 * time.Millisecond)
		return v * v
	})
	fmt.Println(result)
}
