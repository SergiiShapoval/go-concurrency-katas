package main

import (
	"strings"
	"sync/atomic"
	"testing"
)

func TestComputeAllPreservesOrderAndCaches(t *testing.T) {
	var cache Cache
	var calls atomic.Int32
	keys := []string{"go", "go", "channels", "go"}

	got := computeAll(&cache, keys, func(s string) string {
		calls.Add(1)
		return strings.ToUpper(s)
	})

	want := []string{"GO", "GO", "CHANNELS", "GO"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("result %d mismatch: got %q want %q", i, got[i], want[i])
		}
	}
	if calls.Load() < 2 {
		t.Fatalf("expected at least 2 computations for unique keys, got %d", calls.Load())
	}
}
