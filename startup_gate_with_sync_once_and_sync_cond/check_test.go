package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestServiceStartUnblocksWaiters(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		service := NewService()

		var readyCount atomic.Int32
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				service.WaitUntilReady()
				readyCount.Add(1)
			}()
		}

		time.Sleep(20 * time.Millisecond)
		if got := readyCount.Load(); got != 0 {
			t.Fatalf("waiters should still be blocked, got %d", got)
		}

		service.Start()
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("waiters did not unblock")
		}

		if got := readyCount.Load(); got != 4 {
			t.Fatalf("unblocked waiter count mismatch: got %d want 4", got)
		}
	})
}
