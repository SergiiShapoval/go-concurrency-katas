package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
}

func (o *Once) Do(f func()) {
	// TODO: implement contract
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
