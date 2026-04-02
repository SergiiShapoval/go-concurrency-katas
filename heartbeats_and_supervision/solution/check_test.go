package main

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestMonitorCollectsPulsesAndResults(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		time.AfterFunc(250*time.Millisecond, func() { close(done) })

		pulses, results, timedOut := monitor(done, 20*time.Millisecond, 200*time.Millisecond)
		if timedOut {
			t.Fatal("did not expect timeout")
		}
		if pulses == 0 {
			t.Fatal("expected at least one pulse")
		}
		if len(results) == 0 {
			t.Fatal("expected at least one result")
		}
	})
}

func TestMonitorReportsTimeoutWhenWorkerStalls(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		pulses, results, timedOut := monitor(done, 200*time.Millisecond, 40*time.Millisecond)
		if !timedOut {
			t.Fatal("expected timeout")
		}
		if pulses != 0 {
			t.Fatalf("expected no pulses before timeout, got %d", pulses)
		}
		if len(results) != 0 {
			t.Fatalf("expected no results before timeout, got %v", results)
		}
	})
}
