# Resource Lifecycle Evolution To `sync.Pool`

This exercise now follows the progression from Jesse Allen's
[`Jump in the Pool`](https://medium.com/@jessecarl/jump-in-the-pool-6c0013385b51):

1. allocate a fresh resource every time,
2. reuse resources manually with a free list,
3. switch to `sync.Pool` when the resource is short-lived memory that benefits from quick reuse.

The goal is not only to make `sync.Pool` work, but to make the learner see why it is the right fit for this specific kind of workload.

Requirements:

- Keep `AllocFormatter` as the baseline that allocates a fresh `*bytes.Buffer` for each record.
- Implement `FreeListFormatter` with a mutex-protected free list of reusable buffers.
- Implement `PoolFormatter` with `sync.Pool`.
- Keep the same resource lifecycle across all implementations:
  prepare a new resource, get a resource for use, use it, prepare it for reuse.
- `FormatAll` should process records concurrently, preserve output order, and bound concurrency with a semaphore.
- After formatting, copy the final bytes into a standalone string before the buffer can be reused.
- Reset a buffer before returning it to a reusable structure.

Source materials:

- [kat-co, `fig-sync-pool.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/the-sync-package/pool/fig-sync-pool.go)
- [Jesse Allen, `Jump in the Pool`](https://medium.com/@jessecarl/jump-in-the-pool-6c0013385b51)
  This is the primary source for the task design. The article walks through:
  worker pool, fixed free list, dynamic free list, and finally `sync.Pool`.
- [VictoriaMetrics, `sync.Pool`](https://victoriametrics.com/blog/go-sync-pool/)

Expected behavior:

- All formatter implementations produce the same formatted output.
- Reusing a buffer must not corrupt strings returned by earlier calls.
- `AllocFormatter` should keep allocating new buffers across batches.
- `FreeListFormatter` should reuse buffers across batches.
- `PoolFormatter` should return buffers to `sync.Pool`, but tests must allow the pool to drop idle buffers between batches.
- `go test -bench .` should let the learner compare baseline allocation, manual reuse, and `sync.Pool`.
- The benchmark includes both a retainable workload and an oversized workload.
- The oversized workload is intentionally large enough to produce formatted outputs that grow beyond `maxRetainedBufferCap`, so the learner can observe what happens when reusable buffers are dropped instead of retained.

Design notes:

- The sequential benchmark highlights the pure resource-management difference.
- The concurrent benchmark is the final step: it shows that the same lifecycle still applies when formatting runs in parallel, and that `sync.Pool` fits concurrent short-lived memory reuse well.
- A manual free list may outperform `sync.Pool` in small, stable workloads because it is simpler and more predictable.
- `sync.Pool` is still valuable here because it gives dynamic concurrent reuse without forcing the caller to manage the pool size directly, and the garbage collector may discard idle pooled buffers after bursts.
- If every record exceeds `maxRetainedBufferCap`, `sync.Pool` cannot demonstrate reuse at all, and may even look slightly worse than plain allocation because it adds pooling overhead without getting a reuse win.
- Once formatted output grows beyond `maxRetainedBufferCap`, both reusable strategies intentionally drop oversized buffers. That tradeoff reduces retained memory, but it also reduces the reuse benefit in the benchmarks.
- This task intentionally stays in the "temporary memory resource" slice of the article. It is not trying to teach worker pools or long-lived resource pools such as open files or network connections.

Bonus:

- Drop unusually large buffers instead of returning them to the free list or pool.
- Extend the benchmarks with a larger workload and compare `FreeListFormatter` versus `PoolFormatter`.
