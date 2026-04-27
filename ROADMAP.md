# Roadmap

Stuff we owe ourselves on `hns-alertd`. Cross items off in the same commit
as the change — this file is the source of truth for scope.

- [ ] Severity escalation by repeat frequency. Gabriela had this on her plate
      since sprint 48 and it kept getting bumped. INFO → WARN → CRITICAL based
      on occurrences inside some window — probably *not* the dedup window,
      those are different concerns.
- [ ] Replace `chan any` in `internal/actor/subsystem.go` with typed dispatch.
      Cosmetic, but bugs me every time I open that file.
- [ ] Persistent alert store. SQLite probably. Want to survive a restart
      without losing in-flight dedup state.
- [ ] Backpressure observability. At minimum, drop counter on the ingest
      channel and queue depth gauge. Right now we have no idea if we're
      shedding load under burst.
- [ ] Figure out what happened to the CCSDS spike. Branch was
      `hfi/ccsds-spike`, Gabriela deleted it after a week, left a note in
      `.archive/`. Worth a re-read in #tetelestai-core before someone
      else attempts it.
- [ ] Per-spacecraft tenancy. Dedup window is global today, which is fine
      for one bird and wrong the moment we onboard a second.
