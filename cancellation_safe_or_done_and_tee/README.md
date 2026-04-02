# Cancellation-Safe `orDone` And `tee`

Implement two helpers:

```go
func orDone[T any](done <-chan struct{}, in <-chan T) <-chan T
func tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T)
```

Requirements:

- `orDone` should stop forwarding when `done` is closed.
- `tee` should duplicate each value to both output channels.
- Both helpers must avoid goroutine leaks.
- Reuse `orDone` inside `tee` instead of re-implementing the same cancellation-safe receive loop twice.

Source materials:

- [luk4z7, `Or done channel` and `Tee channel`](https://github.com/luk4z7/go-concurrency-guide)
