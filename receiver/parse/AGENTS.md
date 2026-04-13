<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# parse

## Purpose
Metric parsing functions for all supported message formats. Used by all receivers to deserialize incoming metric data into `points.Points`.

## Key Files

| File | Description |
|------|-------------|
| `plain.go` | `Plain()` — parses Graphite plaintext format: `metric.name value timestamp\n` |
| `pickle.go` | `Pickle()` — parses Python pickle-encoded metric batches |
| `protobuf.go` | `Protobuf()` — parses `carbonpb.Payload` protobuf messages |
| `msgpack.go` | `Msgpack()` — parses msgpack-encoded metrics (carbon-relay-ng format) |
| `plain_test.go` | Plaintext parser tests |
| `pickle_test.go` | Pickle parser tests |
| `protobuf_test.go` | Protobuf parser tests |
| `msgpack_test.go` | Msgpack parser tests |
| `bench_test.go` | Benchmarks for all parse formats |
| `common_test.go` | Shared test helpers |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `data/` | Test fixture data files for parser tests |

## For AI Agents

### Working In This Directory
- Each parser is a single function: `func([]byte) ([]*points.Points, error)`.
- Plaintext format: `metric.name value timestamp\n` — the standard Graphite line protocol.
- Pickle uses `lomik/graphite-pickle` library for Python pickle deserialization.
- Protobuf uses `helper/carbonpb` generated types.
- Msgpack is specifically for `carbon-relay-ng` compatibility.
- All parsers return a slice of `*points.Points` (one per unique metric name in the batch).

### Testing Requirements
- `go test ./receiver/parse/`
- `go test -bench=. ./receiver/parse/` for performance benchmarks.

## Dependencies

### Internal
- `points/` — output data type
- `helper/carbonpb/` — protobuf message types

### External
- `lomik/graphite-pickle` — pickle decoder
- `gogo/protobuf` — protobuf runtime
- `vmihailenco/msgpack/v5` — msgpack decoder

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
