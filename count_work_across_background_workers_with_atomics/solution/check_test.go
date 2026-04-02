package main

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestRunWorkersWaitsForStopAndProcessesWork(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		stopCh := make(chan struct{})
		done := make(chan int64, 1)

		go func() {
			done <- RunWorkers(3, stopCh)
		}()

		select {
		case <-done:
			close(stopCh)
			t.Fatal("RunWorkers returned before stop was requested")
		case <-time.After(30 * time.Millisecond):
		}

		close(stopCh)

		select {
		case processed := <-done:
			if processed <= 0 {
				t.Fatalf("expected processed work > 0, got %d", processed)
			}
		case <-time.After(50 * time.Millisecond):
			t.Fatal("RunWorkers did not stop promptly after stop was requested")
		}
	})
}

func TestRunWorkersWithZeroWorkersReturnsZeroAfterStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		stopCh := make(chan struct{})
		done := make(chan int64, 1)

		go func() {
			done <- RunWorkers(0, stopCh)
		}()

		close(stopCh)

		select {
		case processed := <-done:
			if processed != 0 {
				t.Fatalf("expected zero processed work, got %d", processed)
			}
		case <-time.After(20 * time.Millisecond):
			t.Fatal("RunWorkers with zero workers did not return after stop")
		}
	})
}
