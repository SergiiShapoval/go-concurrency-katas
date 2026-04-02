package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type provider func(context.Context) (string, error)

type result struct {
	val string
	err error
}

func FirstResult(timeout time.Duration, providers ...provider) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results := make(chan result, len(providers))
	for _, p := range providers {
		go func(fn provider) {
			val, err := fn(ctx)
			results <- result{val: val, err: err}
		}(p)
	}

	var lastErr error = errors.New("no provider succeeded")
	for range providers {
		select {
		case res := <-results:
			if res.err == nil {
				cancel()
				return res.val, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return "", errors.New("timed out waiting for providers")
		}
	}
	return "", lastErr
}

func main() {
	slow := func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(120 * time.Millisecond):
		}
		return "slow", nil
	}
	fast := func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(40 * time.Millisecond):
		}
		return "fast", nil
	}
	broken := func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		return "", errors.New("backend failed")
	}

	val, err := FirstResult(200*time.Millisecond, slow, fast, broken)
	fmt.Println("value:", val, "error:", err)
}
