# Graceful Pipeline Shutdown

Build a three-stage pipeline:

1. `producer` emits words.
2. `toLower` transforms them.
3. `sink` prints them.

Requirements:

- Thread a `context.Context` through all stages.
- The goroutine that creates a channel must close it.
- Every blocking send or receive must be cancellable via `ctx.Done()`.
- Stop cleanly after a timeout without leaking goroutines.

Source materials:

- [AMBOSS, `A Simple Pipeline` and `Graceful Shutdown With Context`](https://medium.com/amboss/applying-modern-go-concurrency-patterns-to-data-pipelines-b3b5327908d4#50e3)
- [Rob Pike, `Go Concurrency Patterns` (2012)](https://go.dev/talks/2012/concurrency.slide?utm_source=golangweekly&utm_medium=email#1)
- [Anton Zhiyanov, `Go Concurrency` overview and chapter links](https://antonz.org/go-concurrency/)

Expected behavior:

- The program prints a few lowercase words.
- It exits without deadlocking.

Bonus:

- Add a second transform stage.
- Swap the timeout for manual cancellation.
