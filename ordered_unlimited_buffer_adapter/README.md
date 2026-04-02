# Ordered Unlimited-Buffer Adapter

Implement:

```go
func BufferOrdered[T any](in <-chan T) <-chan T
```

The returned channel should behave like a long-lived ordered buffer:

- preserve input order
- accept bursts even when the consumer is slower
- drain remaining buffered values after the input channel closes

Requirements:

- Use one goroutine.
- Use channel semantics plus an internal queue.
- Disable the send case with a `nil` channel when the queue is empty.

Source materials:

- [blogtitle, `Go advanced concurrency patterns: part 4 (unlimited buffer channels)`](https://blogtitle.github.io/go-advanced-concurrency-patterns-part-4-unlimited-buffer-channels/)
- [luk4z7, `Queuing`](https://github.com/luk4z7/go-concurrency-guide)
- [Bryan C. Mills, `Rethinking Classical Concurrency Patterns`](https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view)

Expected behavior:

- A fast producer and slow consumer should still complete without deadlocking.

Bonus:

- Generalize the queue behind a small interface.
