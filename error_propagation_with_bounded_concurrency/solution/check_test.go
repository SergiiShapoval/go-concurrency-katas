package main

import (
	"context"
	"slices"
	"testing"
)

func TestCollectProcessReturnsResultsUntilFirstError(t *testing.T) {
	got, err := collectProcess(context.Background(), []int{1, 2, 3, -1, 5}, 2)
	if err == nil {
		t.Fatal("expected an error for negative input")
	}

	slices.Sort(got)
	want := []int{2, 4, 6}
	if !slices.Equal(got, want) {
		t.Fatalf("results mismatch: got %v want %v", got, want)
	}
}

func TestCollectProcessDoublesSuccessfulValues(t *testing.T) {
	got, err := collectProcess(context.Background(), []int{2, -1}, 1)
	if err == nil {
		t.Fatal("expected an error for negative input")
	}

	want := []int{4}
	if !slices.Equal(got, want) {
		t.Fatalf("results mismatch: got %v want %v", got, want)
	}
}

func TestCollectProcessCompletesNormalFlow(t *testing.T) {
	got, err := collectProcess(context.Background(), []int{1, 2, 3, 4}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slices.Sort(got)
	want := []int{2, 4, 6, 8}
	if !slices.Equal(got, want) {
		t.Fatalf("results mismatch: got %v want %v", got, want)
	}
}
