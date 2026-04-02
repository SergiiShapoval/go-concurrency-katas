# Heartbeats And Supervision

Write a worker that periodically emits:

- heartbeat pulses to show it is alive
- results on a slower cadence

Requirements:

- The worker returns two channels: heartbeat and results.
- Heartbeats should be non-blocking.
- The supervisor should stop if no pulse or result arrives before a timeout.

Source materials:

- [kat-co, `fig-interval-heartbeat.go`](https://github.com/kat-co/concurrency-in-go-src/blob/master/concurrency-at-scale/heartbeats/fig-interval-heartbeat.go)
- [luk4z7, `HeartBeats`](https://github.com/luk4z7/go-concurrency-guide)
- [Sameer Ajmani, `Advanced Go Concurrency Patterns`](https://go.dev/talks/2013/advconc.slide?utm_source=golangweekly&utm_medium=email#43)

Expected behavior:

- The program prints a mix of `pulse` and `result`.
- It exits cleanly on timeout or cancellation.

Bonus:

- Add a deliberately misbehaving worker and detect it.
