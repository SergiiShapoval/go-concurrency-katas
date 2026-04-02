# Reusable Timer Batcher

Build a batcher that groups incoming integers and flushes them when:

- the batch reaches a max size
- a flush timer fires
- the input channel closes

Requirements:

- Reuse a single `time.Timer`.
- Stop and drain the timer correctly before resetting it.
- Preserve input order inside each batch.

Source materials:

- [blogtitle, `Go advanced concurrency patterns: part 2 (timers)`](https://blogtitle.github.io/go-advanced-concurrency-patterns-part-2-timers/)
- [blogtitle, `Go advanced concurrency patterns: part 1`](https://blogtitle.github.io/go-advanced-concurrency-patterns-part-1/)
- [Antonz, `Go Concurrency`](https://antonz.org/go-concurrency/)

Expected behavior:

- Bursty inputs should produce grouped batches.
- No stale timer ticks should trigger phantom flushes.

Bonus:

- Convert the batcher into a debouncer.
