# Count Work Across Background Workers With Atomics

Model a simple service shutdown flow:

- start several background workers
- each worker keeps doing small units of work
- when a stop signal arrives, all workers stop
- return how many work units were processed in total

Implement:

```go
func RunWorkers(workerCount int, stopCh <-chan struct{}) int64
```

Contract:

- `workerCount` workers start immediately
- each worker loops until `stopCh` is closed
- each loop iteration increments one shared atomic counter
- when `stopCh` is closed, the function should stop workers, wait for them, and return the final processed count

Why atomics fit here:

- the only shared mutable state is the processed counter
- each worker can stop by observing `stopCh` directly

Source materials:

- [Go memory model: `sync/atomic`](https://go.dev/ref/mem)
- [Anton Zhiyanov, `Go Concurrency`](https://antonz.org/go-concurrency/)
