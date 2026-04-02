package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type provider func(context.Context) (string, error)

func FirstResult(timeout time.Duration, providers ...provider) (string, error) {
	// TODO: return the first successful result before the timeout
	return "", nil
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
