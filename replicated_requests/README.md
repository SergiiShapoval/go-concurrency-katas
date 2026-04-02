# Replicated Requests

Start the same logical request against several providers and return the first successful result.

Requirements:

- Launch one goroutine per provider.
- The first success wins.
- Cancel or ignore slower responses.

Interview angle:

- Be ready to explain latency-vs-cost tradeoffs and avoiding leaks.

Source materials:

- [luk4z7, `Replicated Requests`](https://github.com/luk4z7/go-concurrency-guide)
- [Sameer Ajmani, `Advanced Go Concurrency Patterns`](https://go.dev/talks/2013/advconc.slide?utm_source=golangweekly&utm_medium=email#43)
