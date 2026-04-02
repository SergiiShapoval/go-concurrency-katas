# Singleflight-Style Duplicate Suppression

Implement a tiny `Group` type with this API:

```go
func (g *Group) Do(key string, fn func() (string, error)) (string, error, bool)
```

The third return value should indicate whether the result was shared with another caller.

Requirements:

- Only one goroutine may execute `fn` for the same key at a time.
- Concurrent callers for the same key must wait and reuse the same result.
- Different keys may execute in parallel.
- Return `shared == false` for the goroutine that actually executed `fn`.
- Return `shared == true` only for callers that reused an in-flight result computed by another goroutine.

Source materials:

- [VictoriaMetrics, `Go Singleflight Melts in Your Code, Not in Your DB`](https://victoriametrics.com/blog/go-singleflight/index.html)
- [Rodaine, `Avoiding Concurrency Boilerplate With golang.org/x/sync`](https://rodaine.com/2018/08/x-files-sync-golang/)
- [VictoriaMetrics, `sync.Map`](https://victoriametrics.com/blog/go-sync-map/)

Expected behavior:

- If 10 goroutines ask for the same key, the expensive function should run once.

Bonus:

- Add `Forget(key string)`.
- Swap the internal map for a typed cache guarded by `sync.RWMutex`.
