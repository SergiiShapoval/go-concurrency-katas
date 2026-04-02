package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func producer(ctx context.Context, words []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case out <- word:
			}
		}
	}()
	return out
}

func toLower(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case word, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- strings.ToLower(word):
				}
			}
		}
	}()
	return out
}

func sink(ctx context.Context, in <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case word, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(word)
		}
	}
}

func runPipeline(ctx context.Context, words []string) []string {
	in := producer(ctx, words)
	out := toLower(ctx, in)

	var result []string
	for {
		select {
		case <-ctx.Done():
			return result
		case word, ok := <-out:
			if !ok {
				return result
			}
			result = append(result, word)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	words := []string{"Go", "CONCURRENCY", "Pipelines", "Context"}
	for _, word := range runPipeline(ctx, words) {
		fmt.Println(word)
	}
}
