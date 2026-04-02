package main

import (
	"testing"
	"testing/synctest"
)

func TestOrReturnsNilWhenNoChannelsProvided(t *testing.T) {
	if got := or(); got != nil {
		t.Fatalf("expected nil channel, got %v", got)
	}
}

func TestOrReturnsSingleChannelDirectly(t *testing.T) {
	ch := make(chan struct{})
	if got := or(ch); got != ch {
		t.Fatal("expected or to return the original channel for one input")
	}
}

func TestOrClosesWhenAnyInputCloses(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan struct{})
		ch2 := make(chan struct{})
		ch3 := make(chan struct{})
		out := or(ch1, ch2, ch3)

		select {
		case <-out:
			t.Fatal("or channel should not be closed before any input closes")
		default:
		}

		close(ch2)
		<-out
		close(ch1)
		close(ch3)
	})
}

func TestOrWithManyChannelsClosesPromptly(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan struct{})
		ch2 := make(chan struct{})
		ch3 := make(chan struct{})
		ch4 := make(chan struct{})
		ch5 := make(chan struct{})
		ch6 := make(chan struct{})
		out := or(ch1, ch2, ch3, ch4, ch5, ch6)

		select {
		case <-out:
			t.Fatal("or channel should not be closed before any input closes")
		default:
		}

		close(ch5)
		<-out
		close(ch1)
		close(ch2)
		close(ch3)
		close(ch4)
		close(ch6)
	})
}
