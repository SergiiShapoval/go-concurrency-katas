# Sneaky Race Conditions And Correct Transfers

Model bank accounts and concurrent transfers.

Requirements:

- Keep each account balance protected by a mutex.
- Implement `Transfer(from, to *Account, amount int)` so the total balance never changes.
- Avoid deadlocks by locking accounts in a stable order.
- Prove the invariant with many concurrent transfers.

This exercise is about logical correctness, not just the absence of data races.

Source materials:

- [blogtitle, `Sneaky race conditions and granular locks`](https://blogtitle.github.io/sneaky-race-conditions-and-granular-locks/)
- [VictoriaMetrics, `sync.Mutex`](https://victoriametrics.com/blog/go-sync-mutex/)
- [Antonz, `Go Concurrency`](https://antonz.org/go-concurrency/)

Expected behavior:

- Final total balance matches the initial total.
- No account goes negative if the transfer rule forbids it.

Bonus:

- Add a deliberately broken `TransferWrong` that uses `Balance()` and `SetBalance()`.
