package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func RunWorkers(workerCount int, stopCh <-chan struct{}) int64 {
	var (
		processed atomic.Int64
		wg        sync.WaitGroup
	)

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
			}

			// TODO: count one processed work unit atomically
			time.Sleep(10 * time.Millisecond)
		}
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}

	wg.Wait()
	return processed.Load()
}

func main() {
	stopCh := make(chan struct{})
	go func() {
		time.Sleep(80 * time.Millisecond)
		close(stopCh)
	}()

	fmt.Println("processed:", RunWorkers(3, stopCh))
}
