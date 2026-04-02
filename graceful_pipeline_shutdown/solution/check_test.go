package main

import (
	"context"
	"slices"
	"testing"
)

func TestRunPipelineLowercasesAllWords(t *testing.T) {
	got := runPipeline(context.Background(), []string{"Go", "PIPE", "Line"})
	want := []string{"go", "pipe", "line"}
	if !slices.Equal(got, want) {
		t.Fatalf("pipeline mismatch: got %v want %v", got, want)
	}
}

func TestRunPipelineStopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := runPipeline(ctx, []string{"Go", "PIPE", "Line"})
	if len(got) != 0 {
		t.Fatalf("expected no values after cancellation, got %v", got)
	}
}
