# Error Propagation With Bounded Concurrency

Process jobs concurrently, but keep the number of active workers bounded.

Requirements:

- Use a token channel as a semaphore.
- Launch work only after acquiring a token.
- Cancel the whole pipeline on the first error.
- Wait for all started goroutines before closing result and error channels.

Suggested job rule:

- Return an error for negative inputs.
- For non-negative inputs, publish `value * 2`.

Source materials:

- [AMBOSS, `Error Handling` and `Maximum Efficiency With Semaphores`](https://medium.com/amboss/applying-modern-go-concurrency-patterns-to-data-pipelines-b3b5327908d4#50e3)
- [`go-resiliency` README patterns](https://github.com/eapache/go-resiliency)
- [Bryan C. Mills, `Rethinking Classical Concurrency Patterns`](https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view)

Expected behavior:

- Successful jobs publish doubled results.
- The first error cancels further work.

Bonus:

- Replace the token channel with `golang.org/x/sync/semaphore`.
- Add retries for transient failures.
