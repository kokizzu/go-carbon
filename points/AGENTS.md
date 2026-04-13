<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# points

## Purpose
Foundational data types for go-carbon. Defines the `Point` (value + timestamp) and `Points` (metric name + slice of points) structs used everywhere in the pipeline. Also provides text and binary serialization, deserialization, and the `Glue` batching helper used by the collector for remote metric forwarding. Every other package imports this one — it must stay dependency-free of other go-carbon packages.

## Key Files
| File | Description |
|------|-------------|
| `points.go` | `Point{Value float64, Timestamp int64}` and `Points{Metric string, Data []Point}`; constructors `New()`, `OnePoint()`, `NowPoint()`; `Add()`, `Append()`, `Copy()`, `Eq()`; `WriteTo()` (plain text), `WriteBinaryTo()` (delta-varint encoding); `ParseText()` for the `"metric value timestamp\n"` line protocol |
| `glue.go` | `Glue(exit, in, chunkSize, chunkTimeout, callback)` — reads `*Points` from a channel, serializes to plain-text lines, batches into chunks by size or timeout, invokes callback; used by `Collector` for remote TCP/UDP metric forwarding |
| `reader.go` | `ReadPlain()`, `ReadBinary()`, `ReadFromFile()` — stream readers for dump files; `ReadFromFile` auto-detects binary vs plain by `.bin` suffix |
| `points_test.go` | Tests for constructors, `ParseText`, `Eq`, `WriteBinaryTo`/`ReadBinary` round-trips |
| `reader_test.go` | Tests for `ReadPlain`, `ReadBinary`, and `ReadFromFile` |

## For AI Agents

### Working In This Directory
- This package must have zero imports from other go-carbon packages. Before adding any import, check that it does not create a cycle. The only allowed imports are from the Go standard library.
- Binary encoding in `WriteBinaryTo()` uses delta-varint compression: the first point stores absolute value and timestamp, subsequent points store deltas. `ReadBinary()` in `reader.go` reverses this with a running accumulator. Any change to the encoding must update both sides simultaneously.
- `ParseText()` parses the Graphite plaintext protocol: `"<metric> <value> <timestamp>\n"`. It accepts float timestamps (truncates to int64) and rejects NaN values. The commented-out range checks (1980–2100) are intentionally disabled — do not re-enable without discussion.
- `Glue()` is used exclusively by `carbon.Collector` for the remote-endpoint code path. It is not used for normal cache operations.
- `MetaMetric` at the bottom of `points.go` is a stub for a planned feature (realtime metric size updates). Do not remove it, but do not build on it without verifying current status.
- `ReadFromFile()` dispatches to binary or plain reader based on filename suffix `.bin` (case-insensitive). Dump files written by `cache.DumpBinary()` always use `.bin`; xlog files written during graceful stop use plain text.

### Testing Requirements
```
go test ./points/
```
All tests are self-contained with no filesystem or network dependencies.

### Common Patterns
- Constructors return `*Points`. Single-point construction: `points.OnePoint(metric, value, timestamp)`. Current-time shorthand: `points.NowPoint(metric, value)`.
- Chained builder style: `p.Add(value, timestamp).Add(value2, timestamp2)` — `Add()` returns `*Points`.
- Equality check via `Eq()` compares metric name and all point values/timestamps element-by-element. Used in tests.
- `WriteTo(w io.Writer)` writes one `"metric value timestamp\n"` line per data point — the standard Graphite plaintext format.

## Dependencies

### Internal
None. This package has no imports from other go-carbon packages by design.

### External
Standard library only:
- `encoding/binary` — varint encoding for binary dump format
- `bufio`, `io`, `os` — buffered reading for dump files
- `math` — `Float64bits` / `Float64frombits` for lossless float encoding
- `strconv`, `strings`, `fmt`, `time` — text parsing and formatting

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
