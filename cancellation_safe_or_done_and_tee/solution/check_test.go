package main

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

func TestOrDoneStopsBlockedForwardOnCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		in := make(chan int)
		out := orDone(done, in)

		sent := make(chan struct{})
		go func() {
			in <- 1
			close(sent)
		}()

		<-sent // orDone has received the value and is now blocked trying to forward it.
		close(done)

		select {
		case v, ok := <-out:
			if ok {
				t.Fatalf("orDone forwarded %v after cancellation", v)
			}
		case <-time.After(50 * time.Millisecond):
			t.Fatal("orDone did not close output after cancellation")
		}
	})
}

func TestTeeDuplicatesValues(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		in := make(chan int)
		go func() {
			defer close(in)
			for _, v := range []int{1, 2, 3} {
				in <- v
			}
		}()

		a, b := tee(done, in)
		var gotA, gotB []int
		for v := range a {
			gotA = append(gotA, v)
			gotB = append(gotB, <-b)
		}

		want := []int{1, 2, 3}
		if !slices.Equal(gotA, want) || !slices.Equal(gotB, want) {
			t.Fatalf("tee mismatch: gotA=%v gotB=%v want=%v", gotA, gotB, want)
		}
	})
}

func TestTeeStopsIfOneOutputIsBlockedAndDoneCloses(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		in := make(chan int, 1)
		in <- 1
		close(in)

		a, b := tee(done, in)

		if v := <-a; v != 1 {
			t.Fatalf("first output value mismatch: got %d want 1", v)
		}

		close(done)

		select {
		case _, ok := <-a:
			if ok {
				t.Fatal("expected first output to be closed after cancellation")
			}
		case <-time.After(50 * time.Millisecond):
			// Unblock the broken implementation so synctest can finish after reporting the failure.
			<-b
			_, _ = <-a
			t.Fatal("tee did not stop after cancellation while second output was blocked")
		}
	})
}
