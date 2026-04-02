package main

import (
	"fmt"
	"time"
)

func batch(in <-chan int, maxSize int, flushAfter time.Duration) <-chan []int {
	out := make(chan []int)

	go func() {
		defer close(out)

		timer := time.NewTimer(flushAfter)
		if !timer.Stop() {
			<-timer.C
		}

		var buf []int
		timerActive := false

		flush := func() {
			if len(buf) == 0 {
				return
			}
			copyBuf := append([]int(nil), buf...)
			out <- copyBuf
			buf = buf[:0]
		}

		stopAndDrain := func() {
			if !timerActive {
				return
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timerActive = false
		}

		for {
			var timeout <-chan time.Time
			if timerActive {
				timeout = timer.C
			}

			select {
			case v, ok := <-in:
				if !ok {
					stopAndDrain()
					flush()
					return
				}

				buf = append(buf, v)
				if len(buf) == 1 {
					timer.Reset(flushAfter)
					timerActive = true
				}

				if len(buf) >= maxSize {
					stopAndDrain()
					flush()
				}
			case <-timeout:
				timerActive = false
				flush()
			}
		}
	}()

	return out
}

func collectBatches(in <-chan int, maxSize int, flushAfter time.Duration) [][]int {
	var out [][]int
	for b := range batch(in, maxSize, flushAfter) {
		out = append(out, b)
	}
	return out
}

func main() {
	in := make(chan int)
	go func() {
		defer close(in)
		for _, burst := range [][]int{{1, 2}, {3, 4, 5}, {6}, {7, 8}} {
			for _, v := range burst {
				in <- v
				time.Sleep(30 * time.Millisecond)
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()

	for _, b := range collectBatches(in, 3, 100*time.Millisecond) {
		fmt.Println(b)
	}
}
