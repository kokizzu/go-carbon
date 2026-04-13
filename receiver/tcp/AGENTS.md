<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# tcp

## Purpose
TCP receiver for plaintext, pickle, and protobuf Graphite protocols over TCP connections. Registers three protocols: `tcp`, `pickle`, and `protobuf`.

## Key Files

| File | Description |
|------|-------------|
| `tcp.go` | `TCP` struct — plaintext line protocol receiver. `Register()` registers tcp/pickle/protobuf protocols. Supports compression (gzip, snappy). |
| `tcp_test.go` | Integration tests for TCP receiver |
| `pickle_test.go` | Tests for pickle protocol framing |
| `protobuf_test.go` | Tests for protobuf protocol framing |
| `stop_test.go` | Graceful shutdown tests |

## For AI Agents

### Working In This Directory
- `Register()` registers **three** protocols: `tcp` (plaintext), `pickle` (Python pickle framing), `protobuf` (length-prefixed protobuf).
- Pickle and protobuf use `FramingOptions` with length-prefixed message framing.
- TCP plaintext parses one metric per line via `parse.Plain()`.
- Compression support: gzip and snappy decompression on incoming connections.

### Testing Requirements
- `go test ./receiver/tcp/`
- Tests spin up real TCP listeners on random ports.

## Dependencies

### Internal
- `receiver/` — plugin registry
- `receiver/parse/` — metric parsing
- `points/` — data types
- `helper/` — `Stoppable` lifecycle

### External
- `klauspost/compress` — gzip and snappy decompression
- `lomik/graphite-pickle/framing` — pickle message framing
- `prometheus/client_golang` — metrics

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
