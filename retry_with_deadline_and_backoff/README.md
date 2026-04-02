# Retry With Deadline And Backoff

Implement a retry helper:

```go
func Retry(ctx context.Context, attempts int, backoff time.Duration, fn func(context.Context) error) error
```

Requirements:

- Stop on the first success.
- Respect `ctx.Done()`.
- Wait between attempts using `time.Timer` or `time.After`.
- If every attempt fails, return the last error from `fn`.

Interview angle:

- This is a good follow-up to “what would you do if the call is flaky?”

Source materials:

- [`go-resiliency` retrier and deadline patterns](https://github.com/eapache/go-resiliency)
- [Antonz, `Go Concurrency`](https://antonz.org/go-concurrency/)
