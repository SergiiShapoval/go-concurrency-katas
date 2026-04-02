package main

import (
	"slices"
	"testing"
)

func TestMergeCollectsAllValues(t *testing.T) {
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

	var got []int
	for v := range merge(makeStream(1, 2), makeStream(3, 4), makeStream(5, 6)) {
		got = append(got, v)
	}

	slices.Sort(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if !slices.Equal(got, want) {
		t.Fatalf("merge mismatch: got %v want %v", got, want)
	}
}
