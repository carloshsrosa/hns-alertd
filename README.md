# hns-alertd

H&S alert daemon. Eats spacecraft health & status frames off stdin, runs them
past per-subsystem limit rules, and spits out deduplicated alerts.

Replaces the Monolito's `hsmon` daemon — same wire format, ~10x throughput,
predictable memory under burst. We've been running it on sabiá since sprint 47.

## Running locally

    make setup     # one-time: installs the git hook
    make test
    ./bin/alertd < testdata/synthetic_burst.bin

(`testdata/` isn't checked in — Gabriela had a generator script,
ping #tetelestai-core if you need the seed.)

## Layout

    cmd/alertd/        binary, no business logic
    internal/actor/    per-subsystem actor; one goroutine per subsystem
    internal/alert/    Alert type + dedup store
    internal/clock/    Clock interface — Fake for tests, System for prod
    internal/frame/    24-byte H&S wire format
    internal/limits/   limit rule evaluation, no I/O
    internal/pipeline/ wires intake → actors → store → out

See `CLAUDE.md` for conventions and the things that bit us.
`ROADMAP.md` is the source of truth for in-flight scope.
