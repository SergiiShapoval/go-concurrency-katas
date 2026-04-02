package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoodGetReturnsSameInstance(t *testing.T) {
	goodOnce = sync.Once{}
	goodCfg = nil

	var wg sync.WaitGroup
	results := make([]*Config, 8)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = GoodGet()
		}(i)
	}
	wg.Wait()

	if results[0] == nil {
		t.Fatal("GoodGet returned nil")
	}
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Fatalf("GoodGet returned different pointers: %p vs %p", results[0], results[i])
		}
	}
}

func TestGoodGetInitializesOnlyOnce(t *testing.T) {
	goodOnce = sync.Once{}
	goodCfg = nil

	var initCalls atomic.Int32
	original := goodInit
	goodInit = func() *Config {
		initCalls.Add(1)
		time.Sleep(20 * time.Millisecond)
		return &Config{Value: "slow"}
	}
	defer func() {
		goodInit = original
		goodOnce = sync.Once{}
		goodCfg = nil
	}()

	const callers = 8
	results := make(chan *Config, callers)
	for i := 0; i < callers; i++ {
		go func() {
			results <- GoodGet()
		}()
	}

	first := <-results
	for i := 1; i < callers; i++ {
		if got := <-results; got != first {
			t.Fatalf("GoodGet returned different pointers: %p vs %p", first, got)
		}
	}

	if got := initCalls.Load(); got != 1 {
		t.Fatalf("initialization ran %d times, want 1", got)
	}
}
