# CCSDS Spike — Notes

— Gabriela, 2026-04-12 (#tetelestai-core)

Tried wrapping the existing `frame.Parse` to also accept CCSDS Space Packets.
Idea was a single ingest path with format autodetection on the first byte.

Killed it because the version field collision made detection ambiguous (our
H&S frames don't have a discriminator byte where CCSDS expects its version)
and the spike branched into "let's redesign the frame format" which is way
out of scope for this quarter.

If someone re-attempts: treat them as separate ingest sources, not a single
polymorphic parser. Two `<-chan Frame` paths into the pipeline is fine; the
parser doesn't need to know about both formats.

Branch was `hfi/ccsds-spike`, deleted 2026-04-12.
