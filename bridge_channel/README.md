# Bridge Channel

Implement:

```go
func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T
```

`bridge` should flatten a stream of channels into one output stream.

Requirements:

- Read channels from `chanStream`.
- Drain each inner channel safely.
- Stop promptly when `done` is closed.
- Preserve the order implied by `chanStream`: fully drain the first inner channel before moving to the next one.

Interview angle:

- This tests whether you can compose small channel helpers into a bigger abstraction.

Source materials:

- [luk4z7, `Bridge channel`](https://github.com/luk4z7/go-concurrency-guide)
- [Bryan C. Mills, `Rethinking Classical Concurrency Patterns`](https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view)
