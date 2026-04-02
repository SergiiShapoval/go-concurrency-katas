# Startup Gate With `sync.Once` And `sync.Cond`

Simulate a service that must initialize exactly once, and only then allow waiting workers to proceed.

Requirements:

- Use `sync.Once` to perform expensive initialization exactly once.
- Use `sync.Cond` to broadcast readiness to waiting goroutines.
- Workers must block until initialization completes.

Source materials:

- [kat-co, `fig-sync-once.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/the-sync-package/once/fig-sync-once.go)
- [kat-co, `fig-cond-broadcast.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/the-sync-package/cond/fig-cond-broadcast.go)
- [VictoriaMetrics, `sync.Once`](https://victoriametrics.com/blog/go-sync-once/) and [VictoriaMetrics, `sync.Cond`](https://victoriametrics.com/blog/go-sync-cond/)

Expected behavior:

- Initialization runs once even if many goroutines call `Start`.
- All waiters are released after initialization finishes.

Bonus:

- Return initialization errors to waiters.
