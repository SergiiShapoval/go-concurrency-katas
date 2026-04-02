# Bounded Parallel Map

Implement:

```go
func ParallelMap[T any, R any](in []T, limit int, fn func(T) R) []R
```

Requirements:

- Preserve output order.
- Limit concurrent workers to `limit`.
- Use a semaphore pattern.

Interview angle:

- This is a common “build a worker pool or bounded map” question.
- Be ready to discuss ordering, backpressure, and goroutine counts.

Source materials:

- [AMBOSS, `Maximum Efficiency With Semaphores`](https://medium.com/amboss/applying-modern-go-concurrency-patterns-to-data-pipelines-b3b5327908d4#50e3)
- [`go-resiliency` semaphore pattern](https://github.com/eapache/go-resiliency)
