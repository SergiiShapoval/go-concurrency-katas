package main

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

func TestBufferOrderedPreservesOrder(t *testing.T) {
	testCases := []struct {
		name          string
		input         []int
		consumerDelay time.Duration
	}{
		{
			name:          "fast producer slow consumer",
			input:         []int{1, 2, 3, 4, 5, 6},
			consumerDelay: 5 * time.Millisecond,
		},
		{
			name:          "single item",
			input:         []int{42},
			consumerDelay: 5 * time.Millisecond,
		},
		{
			name:          "already ordered burst",
			input:         []int{10, 20, 30, 40},
			consumerDelay: 20 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				in := make(chan int)
				go func() {
					defer close(in)
					for _, v := range tc.input {
						in <- v
					}
				}()

				var got []int
				for v := range BufferOrdered(in) {
					got = append(got, v)
					time.Sleep(tc.consumerDelay)
				}

				if !slices.Equal(got, tc.input) {
					t.Fatalf("order mismatch: got %v want %v", got, tc.input)
				}
			})
		})
	}
}

func TestBufferOrderedDrainsBufferedValuesBeforeInputCloses(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := make(chan int, 3)
		in <- 1
		in <- 2
		in <- 3

		out := BufferOrdered(in)

		firstConsumed := make(chan struct{})
		gotCh := make(chan []int, 1)

		go func() {
			var got []int

			got = append(got, <-out)
			close(firstConsumed)

			// Give the adapter time to queue the remaining values.
			time.Sleep(10 * time.Millisecond)

			got = append(got, <-out)
			got = append(got, <-out)
			gotCh <- got
		}()

		<-firstConsumed

		select {
		case got := <-gotCh:
			close(in)
			want := []int{1, 2, 3}
			if !slices.Equal(got, want) {
				t.Fatalf("drain mismatch: got %v want %v", got, want)
			}
		case <-time.After(50 * time.Millisecond):
			close(in)
			t.Fatal("buffered values did not drain while input channel remained open")
		}
	})
}
