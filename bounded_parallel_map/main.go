package main

import (
	"fmt"
	"time"
)

func ParallelMap[T any, R any](in []T, limit int, fn func(T) R) []R {
	// TODO: preserve order while limiting concurrency
	return nil
}

func main() {
	result := ParallelMap([]int{1, 2, 3, 4, 5}, 2, func(v int) int {
		time.Sleep(30 * time.Millisecond)
		return v * v
	})
	fmt.Println(result)
}
