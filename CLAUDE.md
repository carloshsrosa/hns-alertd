# CLAUDE.md

## Project
Go 1.22 service that ingests spacecraft H&S frames, detects out-of-range
readings, and emits deduplicated alerts. Module: `github.com/carloshsrosa/hns-alertd`.

## Key Files
- `cmd/alertd/main.go` — entrypoint, wires the pipeline; no business logic here.
- `internal/actor/subsystem.go` — per-subsystem actor; owns per-subsystem state.
- `internal/alert/store.go` — alert dedup with sliding window.
- `internal/limits/limits.go` — limit rule evaluation; pure, no I/O.
- `internal/clock/clock.go` — `Clock` interface; `clock/fake.go` is the test fake.
- `internal/pipeline/pipeline.go` — frame intake → subsystem actors → store → output.

## Conventions
1. Time reads in business logic go through `clock.Clock`. The only direct
   `time.Now()` calls live in `internal/clock/clock.go` and `cmd/alertd`.
   We added this after a week of flaky tests — don't reintroduce direct
   `time.Now()` in `internal/`, even "just for logging".
2. Per-subsystem state (counters, history, in-flight tracking) lives inside the
   subsystem actor in `internal/actor`. Don't spread it across other packages
   with extra locks — if it's per-subsystem, it goes in the actor.
3. Errors crossing package boundaries are wrapped: `fmt.Errorf("context: %w", err)`.
4. Channel buffer size is 64 unless a benchmark in the same PR justifies otherwise.
5. The `alert.Alert` struct is immutable after construction. Build via
   `alert.New`; add fields to the constructor, not setters. Gabriela got
   bitten by a setter pattern early on and we're not going back.
6. Test names follow `TestType_Behavior` (e.g. `TestStore_DedupesWithinWindow`).

## Do Not
1. Bypass pre-commit hooks. `git commit --no-verify` is forbidden — fix the
   failure instead. The hook runs `go test -race`; if it's catching something,
   the something is real.
2. Use `time.Sleep` to wait on time-dependent behavior in tests. Use
   `clock.Fake.Advance`. A bounded `time.After` is fine only for asserting
   absence of async output.
3. Add a dependency without recording it in `go.mod` and explaining the choice
   in the PR body.
4. Use `interface{}` / `any` for alert payloads or rule values. Define typed
   structs. (Yes, the `chan any` in `actor` is on the roadmap. Don't add more.)
5. Log from `internal/`. Return errors and let `cmd/alertd` decide.

## Tests
- `make test` runs `go test -race ./...`. Race detector is mandatory.
- Tests live next to the code they cover (`foo.go` + `foo_test.go`).
- Time-dependent tests inject `clock.NewFake` and call `Advance`.

## Build & Lint
- `make build` produces `./bin/alertd`.
- `make lint` runs `golangci-lint run ./...`. CI fails on any finding.
- `make pre-commit` runs `fmt + lint + test`. The git hook in
  `.githooks/pre-commit` runs the same target; install it with `make setup`.
