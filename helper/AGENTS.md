<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# helper

## Purpose
Shared utilities and infrastructure used across the go-carbon codebase. Provides lifecycle management for long-running components (`Stoppable`), rate limiting (`ThrottleTicker`), atomic stat helpers, platform-specific file stat types, generated protobuf code, gRPC utilities, atomic file writes, and test helpers.

## Key Files
| File | Description |
|------|-------------|
| `stoppable.go` | `Stoppable` struct — base type for long-running components. Provides `Start()`, `Stop()`, `StartFunc()`, `StopFunc()`, and goroutine launcher `Go`. Embedded by most subsystems in go-carbon. |
| `throttle.go` | `ThrottleTicker` — rate-limiting ticker with soft (`NewThrottleTicker`) and hard (`NewHardThrottleTicker`) modes. Used to cap max-updates-per-second in the persister. |
| `throttle_test.go` | Unit tests for `ThrottleTicker`. |
| `atomic_utils.go` | Atomic stat helpers: `SendAndSubstractUint32/64`, `SendUint32/64`, `SendAndZeroIfNotUpdatedUint32`. Defines `StatCallback` type `func(string, float64)`. |
| `hash.go` | Hash utility functions. |
| `hash_test.go` | Tests for hash utilities. |

## Subdirectories
| Directory | Purpose |
|-----------|---------|
| `atomicfiles/` | See [atomicfiles/AGENTS.md](atomicfiles/AGENTS.md) — atomic file write via temp-file + rename |
| `carbonpb/` | See [carbonpb/AGENTS.md](carbonpb/AGENTS.md) — generated protobuf/gRPC for the carbon cache-query protocol |
| `carbonzipperpb/` | See [carbonzipperpb/AGENTS.md](carbonzipperpb/AGENTS.md) — generated protobuf for the carbonzipper protocol |
| `grpcutil/` | See [grpcutil/AGENTS.md](grpcutil/AGENTS.md) — gRPC server interceptors for timing and status metrics |
| `qa/` | See [qa/AGENTS.md](qa/AGENTS.md) — test helper for creating and cleaning up temp directories |
| `stat/` | See [stat/AGENTS.md](stat/AGENTS.md) — platform-specific file stat types (`FileStats` with Size, RealSize, ATime, CTime, MTime) |

## For AI Agents

### Working In This Directory
- `Stoppable` uses an internal `exit chan bool` (closed on stop) and a `sync.WaitGroup`. All goroutines launched with `s.Go(func(exit chan bool))` must return when `exit` is closed.
- `ThrottleTicker.C` is a `chan bool`. Soft mode: sends `true` at the target rate. Hard mode: sends `true` for the first N per second, then nothing — callers must handle backpressure themselves.
- `StatCallback` is the universal stats reporting type; all packages use `helper.SendAndSubstractUint32` (counter-style, drains after read) or `helper.SendUint32` (gauge-style).
- Do not add business logic here. This package has no knowledge of metrics storage or network protocols.
- Protobuf files in `carbonpb/` and `carbonzipperpb/` are generated — edit `.proto` files and run `go generate` via the `gen.go` files, do not hand-edit `.pb.go` files.

### Testing Requirements
```
go test ./helper/...
```
Tests exist for `throttle.go` and `hash.go`. The `qa/` subpackage is a test helper (imported only in `_test.go` files).

### Common Patterns
- Embed `helper.Stoppable` in any struct that needs start/stop lifecycle.
- Launch goroutines with `s.Go(func(exit chan bool) { ... select { case <-exit: return } ... })`.
- Report stats with `helper.SendAndSubstractUint32("metricName", &s.stat.counter, send)` inside a `Stat(send helper.StatCallback)` method.
- Atomic file writes: write to temp file in same directory, sync, close, then `os.Rename` — guaranteed by `atomicfiles.WriteFile`.

## Dependencies

### Internal
None — `helper` has no imports from other go-carbon packages. It is a leaf dependency.

### External
- `sync`, `sync/atomic` — standard library concurrency primitives
- `time` — standard library (used in throttle)
- `os`, `path` — standard library (used in atomicfiles, qa)
- `google.golang.org/grpc` — gRPC framework (used in grpcutil)
- `github.com/gogo/protobuf` / `github.com/golang/protobuf` — protobuf runtime (used in generated pb files)

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
