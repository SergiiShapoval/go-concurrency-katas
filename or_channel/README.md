# `or` Channel

Implement an `or` helper that combines multiple done channels into one:

```go
func or(channels ...<-chan struct{}) <-chan struct{}
```

Contract:

- If no channels are provided, return `nil`.
- If one channel is provided, return it directly.
- Otherwise, return a channel that closes as soon as any input channel closes.
- The returned channel must close exactly once.

Why this matters in real code:

- combining multiple cancellation sources
- racing several stop conditions while exposing one unified done channel
- building cancellation-safe helpers such as `tee`, `bridge`, and `orDone`

Source materials:

- [luk4z7, `OR Channel`](https://github.com/luk4z7/go-concurrency-guide/blob/main/README.md#or-channel)
- [Go memory model: channel close synchronization](https://go.dev/ref/mem)
