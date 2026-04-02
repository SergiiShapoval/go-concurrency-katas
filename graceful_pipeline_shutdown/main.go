package main

import (
	"context"
)

func producer(ctx context.Context, words []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		// TODO: send words unless ctx is cancelled
	}()
	return out
}

func toLower(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		// TODO: receive from in, convert to lowercase, send to out
	}()
	return out
}

func sink(ctx context.Context, in <-chan string) {
	// TODO: print values until input is closed or ctx is cancelled
}

func runPipeline(ctx context.Context, words []string) []string {
	// TODO: wire producer and toLower, collect all output values, and return them
	return nil
}

func main() {
	// TODO: create a timeout context and wire the pipeline together
}
