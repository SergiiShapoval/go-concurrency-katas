# Build Your Own `sync.Once`

Implement a custom `Once` type:

```go
type Once struct { ... }
func (o *Once) Do(f func())
```

Contract:

- `Do(f)` must run `f` at most once, even with many concurrent callers.
- If one caller is already running `f`, later callers must not return until that execution has finished.
- After the first completed call, later calls should take a cheap fast path.
- Match `sync.Once` behavior when `f` panics.

Why this matters in real code:

- lazy-loading configuration
- creating shared clients or caches exactly once
- understanding why an atomic fast path alone is not enough for correct publication

Source materials:

- [VictoriaMetrics, `Go sync.Once is Simple... Does It Really?`](https://victoriametrics.com/blog/go-sync-once/)
- [Go memory model: `sync.Once`](https://go.dev/ref/mem)
- [Go memory model: incorrect synchronization](https://go.dev/ref/mem)
