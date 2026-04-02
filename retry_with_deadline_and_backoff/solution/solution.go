package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Retry(ctx context.Context, attempts int, backoff time.Duration, fn func(context.Context) error) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i == attempts-1 {
			break
		}

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return lastErr
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	attempt := 0
	err := Retry(ctx, 5, 40*time.Millisecond, func(context.Context) error {
		attempt++
		if attempt < 3 {
			return errors.New("temporary failure")
		}
		fmt.Println("success on attempt", attempt)
		return nil
	})
	fmt.Println("final error:", err)
}
