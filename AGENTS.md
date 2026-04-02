# Go Concurrency Exercise Authoring Guide

This file documents the current exercise pattern for `go_concurrency`.

It is written for any automation agent or contributor that needs to add,
rewrite, renumber, or clean up tasks in this subproject. It is intentionally
tool-agnostic: follow the contract, not a specific assistant workflow.

## Goal

Each exercise should:

- teach one concrete concurrency idea
- point to real source material
- give the learner a starter implementation in `main.go`
- let tests teach the edge cases
- keep a working reference implementation in `solution/solution.go`

The subproject is not a dump of random snippets. It is a curated practice set
ordered from easier to harder tasks.

## Directory Contract

Each exercise lives in a directory with this shape:

```text
stable_slug/
  README.md
  main.go
  check_test.go
  solution/
    solution.go
    check_test.go
```

Rules:

- The directory name is a stable slug.
- The slug should be short, descriptive, and stable.
- The learner-facing file is `main.go`.
- The reference implementation is `solution/solution.go`.
- `check_test.go` in the exercise root and `solution/check_test.go` must match exactly.

## Stable Identity And Ordering

The storage path is not the learning order.

Ordering is defined in:

- [README.md](README.md) for humans browsing the repository
- [curricula/core_easy_to_hard.json](curricula/core_easy_to_hard.json) for automation

Rules:

- Keep directory names stable once created.
- Order exercises by learning complexity, not by discovery date.
- If a new easier task is inserted early, update curriculum manifests and the
  ordered list in [README.md](README.md), not the directory names.
- Use the stable directory slug itself as the task ID.

## What A Good Exercise Looks Like

A good task has one dominant lesson.

Prefer:

- one core concept per task
- a realistic but compact API
- tests that expose common mistakes
- source materials that justify the exercise design

Avoid:

- mixing several unrelated synchronization lessons in one task
- README text that gives away the exact implementation
- tests that depend on hidden behavior not stated in the README
- artificial APIs that exist only to force one library feature

## README Rules

Exercise `README.md` files should describe:

- what the learner must implement
- the public API or relevant types
- behavior and constraints
- why the pattern matters
- source materials

The README should not be a step-by-step implementation recipe.

Prefer:

- contract language
- behavioral constraints
- references that let the learner research the idea

Avoid:

- “use X in line Y”
- naming that leaks the hidden trick when that ruins the exercise
- internal TODO commentary

If tests require behavior that is not obvious, the README must say so.

## Starter Code Rules

`main.go` is the learner starter.

Rules:

- keep exported and tested signatures aligned with `solution/solution.go`
- leave TODOs only where the learner should implement logic
- keep enough scaffolding so the package builds once the TODOs are filled
- include a small `main()` demonstration when it helps orient the learner

The starter may omit implementation details, but it should not omit required
types, helper signatures, or tested surface area.

## Solution Rules

`solution/solution.go` should be a clean reference implementation.

Rules:

- keep the same public signatures and type names as `main.go`
- do not add extra learner-visible API surface unless the starter also has it
- helper types or functions used only internally may exist in the solution, but
  prefer mirroring them in the starter if they are part of the conceptual shape
- keep the solution readable; it is reference code, not code golf

## Test Rules

Tests are the main teacher of mistakes.

Rules:

- `check_test.go` and `solution/check_test.go` must stay identical
- tests should cover the contract and the most likely wrong implementations
- tests should be deterministic
- tests should prefer behavior over implementation details
- benchmarks belong in tests when performance comparison is part of the lesson

Prefer tests that reveal:

- missing synchronization
- cancellation leaks
- incorrect close discipline
- wrong ordering guarantees
- failure to preserve invariants under concurrency
- publication bugs
- panic semantics when relevant

Avoid:

- timing-sensitive assertions without control points
- tests that only restate the happy path
- comments in one test file but not the mirrored one

## Source Material Rules

Each task should have clear provenance inside its own README.

Rules:

- prefer primary materials: official docs, talks, repo examples, original articles
- use source materials to shape the contract, not to copy code mechanically
- if the exercise intentionally diverges from the source, make that explicit in
  the README or design notes

## Link Rules

For markdown files inside this subproject:

- use relative links for internal files
- use normal external URLs for outside references

Do not use absolute local filesystem links in repository markdown.

## New Task Workflow

When adding a new exercise:

1. Pick the concept and source materials.
2. Decide the learner-facing contract first.
3. Create a new directory `stable_slug`.
4. Write `README.md` with contract, constraints, and sources.
5. Add starter code in `main.go`.
6. Write `check_test.go` for the contract and common mistakes.
7. Copy the same test file to `solution/check_test.go`.
8. Implement `solution/solution.go`.
9. Insert the task slug into [curricula/core_easy_to_hard.json](curricula/core_easy_to_hard.json).
10. Add the task to [README.md](README.md) in the correct difficulty position.
11. Run solution-package tests.
12. Verify starter and solution signatures still match.

## Cleanup Workflow

When modifying an existing exercise:

1. Check whether README, starter, tests, and solution still describe the same contract.
2. If you change tests, mirror the same edits into `solution/check_test.go`.
3. If you change public signatures or type names, update both `main.go` and `solution/solution.go`.
4. If you rename a task slug, update curricula and documentation links.
5. Run the affected `solution` tests.

## Verification Checklist

Before finishing changes, verify:

- all internal markdown links are relative
- ordered curricula are updated
- task and solution tests are identical for every exercise
- starter and solution signatures match
- tasks appear in the correct order in [README.md](README.md)
- `go test ./.../solution` passes for the affected packages

## Practical Heuristics

Use these heuristics when designing tasks:

- If a README makes the exact implementation obvious, it is probably too explicit.
- If a learner can satisfy the README but still fail the tests, the README is incomplete.
- If a task mainly exists to show a benchmark difference, make the benchmark part of the exercise contract.
- If a pattern is already taught elsewhere in the set, avoid duplicating it unless the new task highlights a different failure mode.
- If one concept is naturally solved by a simpler standard primitive, redesign the task so the intended lesson is still meaningful.
