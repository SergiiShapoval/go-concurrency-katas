package main

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestSnapshotPublisherPublishesToOneWaiter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := make(chan struct{})
		snapshot, ready := startSnapshotPublisher(start)

		select {
		case <-ready:
			t.Fatal("snapshot publisher should not signal readiness before start")
		case <-time.After(10 * time.Millisecond):
		}

		close(start)
		<-ready

		if snapshot.Value != "snapshot published" {
			t.Fatalf("snapshot mismatch: got %q", snapshot.Value)
		}
	})
}

func TestCatalogPublisherBroadcastsToAllWaiters(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := make(chan struct{})
		catalog, ready := startCatalogPublisher(start)

		done1 := make(chan struct{})
		done2 := make(chan struct{})
		go func() {
			<-ready
			close(done1)
		}()
		go func() {
			<-ready
			close(done2)
		}()

		select {
		case <-done1:
			t.Fatal("catalog publisher should not wake waiters before start")
		case <-done2:
			t.Fatal("catalog publisher should not wake waiters before start")
		case <-time.After(10 * time.Millisecond):
		}

		close(start)
		<-done1
		<-done2

		if catalog.Value != "catalog published" {
			t.Fatalf("catalog mismatch: got %q", catalog.Value)
		}
	})
}
