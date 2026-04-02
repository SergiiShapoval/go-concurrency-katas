package main

import (
	"slices"
	"testing"
	"time"
)

func TestCollectBatchesFlushesBySizeAndClose(t *testing.T) {
	in := make(chan int)
	go func() {
		defer close(in)
		for _, v := range []int{1, 2, 3, 4, 5} {
			in <- v
		}
	}()

	got := collectBatches(in, 3, time.Hour)
	want := [][]int{{1, 2, 3}, {4, 5}}
	if len(got) != len(want) {
		t.Fatalf("batch count mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Fatalf("batch %d mismatch: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestCollectBatchesFlushesByTime(t *testing.T) {
	in := make(chan int)
	go func() {
		defer close(in)
		in <- 1
		time.Sleep(20 * time.Millisecond)
		in <- 2
	}()

	got := collectBatches(in, 10, 5*time.Millisecond)
	want := [][]int{{1}, {2}}
	if len(got) != len(want) {
		t.Fatalf("batch count mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Fatalf("batch %d mismatch: got %v want %v", i, got[i], want[i])
		}
	}
}

func TestCollectBatchesReturnsNoBatchesForClosedEmptyInput(t *testing.T) {
	in := make(chan int)
	close(in)

	got := collectBatches(in, 3, time.Millisecond)
	if len(got) != 0 {
		t.Fatalf("expected no batches, got %v", got)
	}
}
