package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
	done atomic.Bool
	mu   sync.Mutex
}

func (o *Once) Do(f func()) {
	if !o.done.Load() {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.done.Load() {
		defer o.done.Store(true)
		f()
	}
}

func main() {
	var once Once
	var wg sync.WaitGroup
	var loaded atomic.Int32

	load := func() {
		time.Sleep(20 * time.Millisecond)
		loaded.Add(1)
		fmt.Println("loaded shared resource")
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(load)
		}()
	}

	wg.Wait()
	fmt.Println("times loaded:", loaded.Load())
}
