package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func process(ctx context.Context, values []int, limit int) (<-chan int, <-chan error) {
	ctx, cancel := context.WithCancel(ctx)
	results := make(chan int)
	errs := make(chan error, 1)
	sem := make(chan struct{}, limit)
	var (
		wg       sync.WaitGroup
		launchWg sync.WaitGroup
	)

	sendErr := func(err error) {
		select {
		case errs <- err:
			cancel()
		default:
		}
	}

	launchWg.Add(1)
	go func() {
		defer launchWg.Done()
		for _, value := range values {
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(v int) {
				defer wg.Done()
				defer func() { <-sem }()

				if v < 0 {
					sendErr(errors.New("negative value encountered"))
					return
				}

				time.Sleep(40 * time.Millisecond)

				// Let already-started work publish its result even if a peer failed.
				results <- v * 2
			}(value)
		}
	}()

	go func() {
		launchWg.Wait()
		wg.Wait()
		cancel()
		close(results)
		close(errs)
	}()

	return results, errs
}

func collectProcess(ctx context.Context, values []int, limit int) ([]int, error) {
	results, errs := process(ctx, values, limit)
	var out []int
	var firstErr error

	for results != nil || errs != nil {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			out = append(out, result)
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return out, firstErr
}

func main() {
	ctx := context.Background()
	results, err := collectProcess(ctx, []int{1, 2, 3, -1, 5, 6}, 3)
	for _, result := range results {
		fmt.Println("result:", result)
	}
	fmt.Println("error:", err)
}
