# Broken Double-Checked Locking, Then Fix It

Create two versions of lazy initialization:

- `BadGet`: uses an unsynchronized boolean fast-path
- `GoodGet`: uses proper synchronization to publish exactly one initialized instance

Requirements:

- Keep both implementations in the file so you can compare them.
- The solution should explain why `BadGet` is logically incorrect under the memory model.

Interview angle:

- This is exactly the kind of thing interviewers ask when they want to probe whether you understand publication, not just syntax.

Source materials:

- [Go memory model: incorrect synchronization](https://go.dev/ref/mem)
- [Go memory model: double-checked locking example](https://go.dev/ref/mem)
- [Go memory model: `sync.Once`](https://go.dev/ref/mem)
