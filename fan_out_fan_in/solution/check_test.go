package main

import (
	"slices"
	"testing"
)

func TestFanOutFanInSquaresAllInputs(t *testing.T) {
	in := generator([]int{1, 2, 3, 4, 5})
	got := make([]int, 0, 5)
	for v := range fanIn(fanOut(in, 3)...) {
		got = append(got, v)
	}

	slices.Sort(got)
	want := []int{1, 4, 9, 16, 25}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected merged output: got %v want %v", got, want)
	}
}
