# Merge With `WaitGroup` And Close Discipline

Implement:

```go
func merge[T any](chs ...<-chan T) <-chan T
```

Requirements:

- One goroutine per input channel.
- Close the output channel exactly once, after all input readers finish.
- No sends to closed channels.

Interview angle:

- This is a classic fan-in exercise.
- The real point is ownership and close timing, not syntax.

Source materials:

- [luk4z7, `Fan-in and Fan-out`](https://github.com/luk4z7/go-concurrency-guide)
- [kat-co, `fig-fan-out-naive-prime-finder.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/concurrency-patterns-in-go/fan-out-fan-in/fig-fan-out-naive-prime-finder.go)
- [VictoriaMetrics, `sync.WaitGroup`](https://victoriametrics.com/blog/go-sync-waitgroup/)
