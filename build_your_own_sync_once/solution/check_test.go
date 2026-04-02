package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOnceRunsFunctionOnlyOnce(t *testing.T) {
	var once Once
	var calls atomic.Int32
	start := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			once.Do(func() {
				calls.Add(1)
			})
		}()
	}

	close(start)
	wg.Wait()

	if got := calls.Load(); got != 1 {
		t.Fatalf("call count mismatch: got %d want 1", got)
	}
}

func TestOnceWaitsForInitializationToFinish(t *testing.T) {
	var once Once
	started := make(chan struct{})
	release := make(chan struct{})
	secondReturned := make(chan struct{})

	go func() {
		once.Do(func() {
			close(started)
			<-release
		})
	}()

	<-started

	go func() {
		once.Do(func() {})
		close(secondReturned)
	}()

	select {
	case <-secondReturned:
		t.Fatal("second caller returned before initialization finished")
	case <-time.After(20 * time.Millisecond):
	}

	close(release)

	select {
	case <-secondReturned:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second caller did not return after initialization finished")
	}
}

func TestOnceDoesNotRerunAfterPanic(t *testing.T) {
	var once Once
	var calls atomic.Int32

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected first Do call to panic")
			}
		}()

		once.Do(func() {
			calls.Add(1)
			panic("boom")
		})
	}()

	once.Do(func() {
		calls.Add(1)
	})

	if got := calls.Load(); got != 1 {
		t.Fatalf("call count mismatch after panic: got %d want 1", got)
	}
}

func TestOnceWaitersDoNotRerunAfterPanic(t *testing.T) {
	var once Once
	var calls atomic.Int32
	started := make(chan struct{})
	waiterReturned := make(chan struct{})
	panicReturned := make(chan struct{})

	go func() {
		defer close(panicReturned)
		defer func() {
			if recover() == nil {
				t.Error("expected first Do call to panic")
			}
		}()

		once.Do(func() {
			calls.Add(1)
			close(started)
			time.Sleep(20 * time.Millisecond)
			panic("boom")
		})
	}()

	<-started

	go func() {
		once.Do(func() {
			calls.Add(1)
		})
		close(waiterReturned)
	}()

	select {
	case <-waiterReturned:
		t.Fatal("waiter returned before panicking call finished")
	case <-time.After(5 * time.Millisecond):
	}

	<-panicReturned

	select {
	case <-waiterReturned:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("waiter did not return after panic completed")
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("call count mismatch after concurrent panic: got %d want 1", got)
	}
}
