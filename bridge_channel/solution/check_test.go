package main

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

func TestBridgeFlattensChannelStream(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		makeStream := func(vals ...int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for _, v := range vals {
					out <- v
				}
			}()
			return out
		}

		chanStream := make(chan (<-chan int))
		go func() {
			defer close(chanStream)
			chanStream <- makeStream(1, 2)
			chanStream <- makeStream(3)
			chanStream <- makeStream(4, 5)
		}()

		var got []int
		for v := range bridge(done, chanStream) {
			got = append(got, v)
		}

		want := []int{1, 2, 3, 4, 5}
		if !slices.Equal(got, want) {
			t.Fatalf("bridge mismatch: got %v want %v", got, want)
		}
	})
}

func TestBridgeDrainsInnerChannelsSequentially(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		first := make(chan int, 2)
		second := make(chan int, 1)
		chanStream := make(chan (<-chan int), 2)

		first <- 1
		second <- 3
		chanStream <- first
		chanStream <- second
		close(chanStream)

		out := bridge(done, chanStream)

		if got := <-out; got != 1 {
			t.Fatalf("first value mismatch: got %d want 1", got)
		}

		select {
		case got := <-out:
			t.Fatalf("bridge emitted %d from a later channel before the first channel finished", got)
		case <-time.After(20 * time.Millisecond):
		}

		first <- 2
		close(first)
		close(second)

		var got []int
		for v := range out {
			got = append(got, v)
		}

		want := []int{2, 3}
		if !slices.Equal(got, want) {
			t.Fatalf("tail mismatch: got %v want %v", got, want)
		}
	})
}

func TestBridgeStopsPromptlyWhenDoneClosesDuringBlockedSend(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		inner := make(chan int, 2)
		chanStream := make(chan (<-chan int), 1)

		inner <- 1
		inner <- 2
		close(inner)
		chanStream <- inner
		close(chanStream)

		out := bridge(done, chanStream)

		if got := <-out; got != 1 {
			t.Fatalf("first value mismatch: got %d want 1", got)
		}

		close(done)

		select {
		case _, ok := <-out:
			if ok {
				t.Fatal("expected bridge output to close after cancellation")
			}
		case <-time.After(20 * time.Millisecond):
			t.Fatal("bridge did not stop after cancellation while blocked on output send")
		}
	})
}
