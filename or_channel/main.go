package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan struct{}) <-chan struct{} {
	// TODO: return a channel that closes when any input channel closes
	return nil
}

func sig(after time.Duration) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(after)
	}()
	return done
}

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
	)
	fmt.Println("done after", time.Since(start).Round(time.Millisecond))
}
