# Concurrent Cache With `sync.Map`

Build a cache that memoizes expensive string computations.

Requirements:

- Use `sync.Map`.
- Implement `GetOrCompute(key string, fn func() string) string`.
- Ensure only one value is stored per key.
- Keep the API simple enough to explain tradeoffs in an interview.

Interview angle:

- Be ready to explain when `sync.Map` is a good fit and when a `map` plus `sync.RWMutex` is better.

Source materials:

- [Rodaine, `Avoiding Concurrency Boilerplate With golang.org/x/sync`](https://rodaine.com/2018/08/x-files-sync-golang/)
- [VictoriaMetrics, `sync.Map`](https://victoriametrics.com/blog/go-sync-map/)
