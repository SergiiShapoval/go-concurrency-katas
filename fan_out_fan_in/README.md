# Fan-Out / Fan-In

Create a parallel pipeline stage that computes squares.

Requirements:

- `generator` emits integers.
- `worker` consumes integers and emits their squares.
- `fanOut` starts `n` workers sharing one input channel.
- `fanIn` merges all worker output channels into one result channel using `sync.WaitGroup`.

Source materials:

- [AMBOSS, `Adding Parallelism with Fan-Out and Fan-In`](https://medium.com/amboss/applying-modern-go-concurrency-patterns-to-data-pipelines-b3b5327908d4#50e3)
- [luk4z7, `Go Concurrency Guide`](https://github.com/luk4z7/go-concurrency-guide)
- [kat-co, `fig-fan-out-naive-prime-finder.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/concurrency-patterns-in-go/fan-out-fan-in/fig-fan-out-naive-prime-finder.go)

Expected behavior:

- All inputs are processed.
- Result order does not need to match input order.
- Output channel closes after all workers finish.

Bonus:

- Make worker count depend on `runtime.NumCPU()`.
- Add cancellation with `context.Context`.
