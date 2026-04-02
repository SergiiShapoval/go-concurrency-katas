package main

import (
	"context"
	"sync"
)

func process(ctx context.Context, values []int, limit int) (<-chan int, <-chan error) {
	results := make(chan int)
	errs := make(chan error, 1)

	// TODO: use a buffered channel as a semaphore
	// TODO: start goroutines only after acquiring a token
	// TODO: send value * 2 for successful jobs
	// TODO: cancel on first error and close results/errs when done

	var _ sync.WaitGroup
	return results, errs
}

func collectProcess(ctx context.Context, values []int, limit int) ([]int, error) {
	// TODO: consume process() until channels close, collect doubled results, and return the first error if any
	return nil, nil
}

func main() {
	// TODO: call process and print either results or the first error
}
