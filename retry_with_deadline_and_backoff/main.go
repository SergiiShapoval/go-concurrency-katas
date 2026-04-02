package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Retry(ctx context.Context, attempts int, backoff time.Duration, fn func(context.Context) error) error {
	// TODO: retry fn with delay while respecting ctx
	return nil
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
