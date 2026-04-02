# Go Concurrency Katas

`go-concurrency-katas` is a learning-by-doing collection of Go concurrency exercises.

[![CI](https://github.com/SergiiShapoval/go-concurrency-katas/actions/workflows/ci.yml/badge.svg)](https://github.com/SergiiShapoval/go-concurrency-katas/actions/workflows/ci.yml)

It is designed for engineers who already know basic Go syntax and want hands-on practice with:

- pipelines
- fan-out / fan-in
- cancellation
- channel signaling
- `sync` primitives
- atomics
- common concurrency bugs and tradeoffs

Curriculum files:

- [curricula/core_easy_to_hard.json](curricula/core_easy_to_hard.json)

Suggested learning path, easiest first:

- [`001` Graceful Pipeline Shutdown](graceful_pipeline_shutdown/README.md)
- [`002` Fan-Out / Fan-In](fan_out_fan_in/README.md)
- [`003` Merge With `WaitGroup` And Close Discipline](merge_with_waitgroup_and_close_discipline/README.md)
- [`004` Startup Gate With `sync.Once` And `sync.Cond`](startup_gate_with_sync_once_and_sync_cond/README.md)
- [`005` Resource Lifecycle Evolution To `sync.Pool`](resource_lifecycle_evolution_to_sync_pool/README.md)
- [`006` Concurrent Cache With `sync.Map`](concurrent_cache_with_sync_map/README.md)
- [`007` Ordered Unlimited-Buffer Adapter](ordered_unlimited_buffer_adapter/README.md)
- [`008` Cancellation-Safe `orDone` And `tee`](cancellation_safe_or_done_and_tee/README.md)
- [`009` Bridge Channel](bridge_channel/README.md)
- [`010` Heartbeats And Supervision`](heartbeats_and_supervision/README.md)
- [`011` Bounded Parallel Map](bounded_parallel_map/README.md)
- [`012` Reusable Timer Batcher](reusable_timer_batcher/README.md)
- [`013` Replicated Requests](replicated_requests/README.md)
- [`014` Retry With Deadline And Backoff](retry_with_deadline_and_backoff/README.md)
- [`015` Singleflight-Style Duplicate Suppression](singleflight_style_duplicate_suppression/README.md)
- [`016` Sneaky Race Conditions And Correct Transfers](sneaky_race_conditions_and_correct_transfers/README.md)
- [`017` Error Propagation With Bounded Concurrency](error_propagation_with_bounded_concurrency/README.md)
- [`018` Publication Signals](publication_signals/README.md)
- [`019` `or` Channel](or_channel/README.md)
- [`020` Broken Double-Checked Locking, Then Fix It](broken_double_checked_locking_then_fix_it/README.md)
- [`021` Count Work Across Background Workers With Atomics](count_work_across_background_workers_with_atomics/README.md)
- [`022` Build Your Own `sync.Once`](build_your_own_sync_once/README.md)

Each exercise follows the same layout:

- `README.md`: task description, constraints, and source materials
- `main.go`: starter code with TODOs
- `*_test.go`: behavior checks for the candidate solution
- `solution/solution.go`: one reference solution

Quick start:

```bash
go test ./...
go test ./.../solution
go test -run '^$' -bench . -benchmem ./.../solution
```

CI runs the reference implementations only:

```bash
go test ./.../solution
```

This keeps the default GitHub Actions check green while learner starter files in `main.go` still contain TODOs or intentionally incomplete implementations.

Benchmark tracking:

- GitHub Actions also runs `go test -run '^$' -bench . -benchmem ./.../solution` on `main`.
- Benchmark history is stored in the `gh-pages` branch and rendered as a GitHub Pages dashboard:
  [Benchmark Dashboard](https://sergiishapoval.github.io/go-concurrency-katas/dev/bench/).
- Before the first benchmark run, create an empty `gh-pages` branch and set GitHub Pages to publish from that branch.
- The benchmark workflow compares current results with previous history and leaves a workflow summary on each run. Alert comments are enabled for large regressions.
- Pull requests also run the same benchmark command on both the PR branch and its base branch, then compare the two results on the same runner.
- The PR benchmark workflow comments on regressions above `130%` and fails above `150%` to reduce noise from GitHub-hosted runner variance.

How to use this repository:

1. Pick a task from the ordered list above.
2. Read the task `README.md`.
3. Implement the missing logic in `main.go`.
4. Run the tests for that task.
5. Use `solution/solution.go` only after attempting the exercise yourself.

Notes:

- Stable directory slugs are used as task IDs.
- The global curriculum order lives in [curricula/core_easy_to_hard.json](curricula/core_easy_to_hard.json).
