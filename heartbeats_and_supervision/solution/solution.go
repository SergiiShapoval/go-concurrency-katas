package main

import (
	"fmt"
	"time"
)

func doWork(done <-chan struct{}, pulseInterval time.Duration) (<-chan struct{}, <-chan int) {
	heartbeat := make(chan struct{})
	results := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.NewTicker(pulseInterval)
		work := time.NewTicker(2 * pulseInterval)
		defer pulse.Stop()
		defer work.Stop()

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		value := 0
		for {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			case <-work.C:
				value++
				for {
					select {
					case <-done:
						return
					case <-pulse.C:
						sendPulse()
					case results <- value:
						goto next
					}
				}
			}
		next:
		}
	}()

	return heartbeat, results
}

func monitor(done <-chan struct{}, pulseInterval, timeout time.Duration) (int, []int, bool) {
	heartbeat, results := doWork(done, pulseInterval)
	pulses := 0
	var out []int

	for {
		select {
		case <-done:
			return pulses, out, false
		case _, ok := <-heartbeat:
			if !ok {
				return pulses, out, false
			}
			pulses++
		case result, ok := <-results:
			if !ok {
				return pulses, out, false
			}
			out = append(out, result)
		case <-time.After(timeout):
			return pulses, out, true
		}
	}
}

func main() {
	done := make(chan struct{})
	time.AfterFunc(1500*time.Millisecond, func() { close(done) })

	pulses, results, timedOut := monitor(done, 100*time.Millisecond, 350*time.Millisecond)
	fmt.Println("pulses:", pulses)
	for _, result := range results {
		fmt.Println("result:", result)
	}
	fmt.Println("timedOut:", timedOut)
}
