# Publication Signals

Implement two helpers:

- `startSnapshotPublisher`
- `startCatalogPublisher`

API:

```go
func startSnapshotPublisher(start <-chan struct{}) (*Snapshot, <-chan struct{})
func startCatalogPublisher(start <-chan struct{}) (*Catalog, <-chan struct{})
```

How to read this API:

- `start` lets the test control when publication begins
- `*Snapshot` or `*Catalog` is the shared data being published
- the returned `ready` channel tells readers when the published value is safe to observe

Constraints:

- Both helpers should begin their publication work only after `start`.
- The published data must be safe to read after the readiness signal is observed.
- The two helpers intentionally use different readiness semantics.
- Use the tests and referenced materials to infer the appropriate signaling strategy.

Why this matters in real code:

- publishing one expensive snapshot to a consumer
- broadcasting readiness or shutdown to many goroutines

Source materials:

- [Go memory model](https://go.dev/ref/mem)
  Channel send, receive, and close create synchronization edges that make published writes visible.
- [Go blog: Pipelines](https://go.dev/blog/pipelines)
  The `done` channel pattern uses channel close as a broadcast signal to many goroutines.
