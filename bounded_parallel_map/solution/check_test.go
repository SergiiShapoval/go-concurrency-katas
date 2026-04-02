package main

import (
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestParallelMapPreservesOrder(t *testing.T) {
	testCases := []struct {
		name  string
		in    []int
		limit int
		want  []int
		delay time.Duration
	}{
		{
			name:  "bounded parallel order",
			in:    []int{1, 2, 3, 4, 5},
			limit: 2,
			want:  []int{1, 4, 9, 16, 25},
			delay: 5 * time.Millisecond,
		},
		{
			name:  "serial limit one",
			in:    []int{3, 1, 2},
			limit: 1,
			want:  []int{9, 1, 4},
			delay: 8 * time.Millisecond,
		},
		{
			name:  "limit exceeds input",
			in:    []int{2, 4},
			limit: 8,
			want:  []int{4, 16},
			delay: 3 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				got := ParallelMap(tc.in, tc.limit, func(v int) int {
					time.Sleep(tc.delay)
					return v * v
				})
				if !slices.Equal(got, tc.want) {
					t.Fatalf("parallel map mismatch: got %v want %v", got, tc.want)
				}
			})
		})
	}
}

func TestParallelMapLimitsActiveCalls(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var (
			active    atomic.Int32
			maxActive atomic.Int32
		)

		got := ParallelMap([]int{1, 2, 3, 4, 5, 6}, 2, func(v int) int {
			n := active.Add(1)
			for {
				cur := maxActive.Load()
				if n <= cur || maxActive.CompareAndSwap(cur, n) {
					break
				}
			}

			time.Sleep(10 * time.Millisecond)
			active.Add(-1)
			return v * 10
		})

		want := []int{10, 20, 30, 40, 50, 60}
		if !slices.Equal(got, want) {
			t.Fatalf("parallel map mismatch: got %v want %v", got, want)
		}
		if maxActive.Load() > 2 {
			t.Fatalf("active calls exceeded limit: got %d want <= 2", maxActive.Load())
		}
	})
}
