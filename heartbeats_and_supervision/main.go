package main

import "time"

func doWork(done <-chan struct{}, pulseInterval time.Duration) (<-chan struct{}, <-chan int) {
	heartbeat := make(chan struct{})
	results := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(results)
		// TODO: emit pulses often and results less often
	}()
	return heartbeat, results
}

func monitor(done <-chan struct{}, pulseInterval, timeout time.Duration) (int, []int, bool) {
	// TODO: supervise doWork; return pulse count, collected results, and whether a timeout happened
	return 0, nil, false
}

func main() {
	// TODO: supervise the worker with a timeout
}
