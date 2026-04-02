package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestRetryEventuallySucceeds(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var attempts atomic.Int32
		err := Retry(context.Background(), 5, time.Millisecond, func(context.Context) error {
			if attempts.Add(1); attempts.Load() < 3 {
				return errors.New("try again")
			}
			return nil
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := attempts.Load(); got != 3 {
			t.Fatalf("attempt count mismatch: got %d want 3", got)
		}
	})
}

func TestRetryReturnsLastFailureWhenAllAttemptsFail(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		wantErr := errors.New("still failing")
		var attempts atomic.Int32

		err := Retry(context.Background(), 3, time.Millisecond, func(context.Context) error {
			attempts.Add(1)
			return wantErr
		})

		if !errors.Is(err, wantErr) {
			t.Fatalf("unexpected error: got %v want %v", err, wantErr)
		}
		if got := attempts.Load(); got != 3 {
			t.Fatalf("attempt count mismatch: got %d want 3", got)
		}
	})
}

func TestRetryStopsWhenContextIsCancelledBeforeSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var attempts atomic.Int32

		errCh := make(chan error, 1)
		go func() {
			errCh <- Retry(ctx, 5, 10*time.Millisecond, func(context.Context) error {
				if attempts.Add(1) == 1 {
					cancel()
				}
				return errors.New("transient failure")
			})
		}()

		err := <-errCh
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("unexpected error: got %v want %v", err, context.Canceled)
		}
		if got := attempts.Load(); got != 1 {
			t.Fatalf("expected retry loop to stop after cancellation, got %d attempts", got)
		}
	})
}
