package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestFirstResultReturnsFastSuccessfulProvider(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		slow := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(80 * time.Millisecond):
			}
			return "slow", nil
		}
		fast := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(10 * time.Millisecond):
			}
			return "fast", nil
		}
		broken := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
			}
			return "", errors.New("nope")
		}

		got, err := FirstResult(200*time.Millisecond, slow, fast, broken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "fast" {
			t.Fatalf("unexpected result: got %q want %q", got, "fast")
		}
	})
}

func TestFirstResultReturnsTimeoutWhenNoProviderSucceedsInTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		slow := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(100 * time.Millisecond):
			}
			return "late", nil
		}

		got, err := FirstResult(20*time.Millisecond, slow)
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if got != "" {
			t.Fatalf("unexpected result: got %q want empty string", got)
		}
	})
}

func TestFirstResultReturnsLastFailureWhenAllProvidersFail(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		firstErr := errors.New("first failure")
		lastErr := errors.New("last failure")

		first := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(10 * time.Millisecond):
			}
			return "", firstErr
		}
		second := func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(20 * time.Millisecond):
			}
			return "", lastErr
		}

		got, err := FirstResult(200*time.Millisecond, first, second)
		if !errors.Is(err, lastErr) {
			t.Fatalf("unexpected error: got %v want %v", err, lastErr)
		}
		if got != "" {
			t.Fatalf("unexpected result: got %q want empty string", got)
		}
	})
}
