package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupDoSharesInFlightCall(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T, results []string, shared []bool, calls int32)
	}{
		{
			name: "all callers receive the same value",
			check: func(t *testing.T, results []string, _ []bool, _ int32) {
				for i, v := range results {
					if v != "value" {
						t.Fatalf("result %d mismatch: %q", i, v)
					}
				}
			},
		},
		{
			name: "expensive function runs once",
			check: func(t *testing.T, _ []string, _ []bool, calls int32) {
				if calls != 1 {
					t.Fatalf("expensive function called %d times, want 1", calls)
				}
			},
		},
		{
			name: "only waiters are marked shared",
			check: func(t *testing.T, _ []string, shared []bool, _ int32) {
				var sharedTrue, sharedFalse int
				for _, v := range shared {
					if v {
						sharedTrue++
						continue
					}
					sharedFalse++
				}
				if sharedFalse != 1 {
					t.Fatalf("expected exactly one non-shared caller, got %d", sharedFalse)
				}
				if sharedTrue != len(shared)-1 {
					t.Fatalf("expected %d shared callers, got %d", len(shared)-1, sharedTrue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g Group
			var calls atomic.Int32

			fn := func() (string, error) {
				calls.Add(1)
				time.Sleep(25 * time.Millisecond)
				return "value", nil
			}

			var wg sync.WaitGroup
			results := make([]string, 8)
			shared := make([]bool, 8)
			for i := range results {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					v, err, isShared := g.Do("same-key", fn)
					if err != nil {
						t.Errorf("unexpected error: %v", err)
						return
					}
					results[i] = v
					shared[i] = isShared
				}(i)
			}
			wg.Wait()

			tt.check(t, results, shared, calls.Load())
		})
	}
}
