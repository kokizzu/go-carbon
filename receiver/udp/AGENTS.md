<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# udp

## Purpose
UDP receiver for Graphite plaintext protocol. Receives one metric per line from UDP datagrams.

## Key Files

| File | Description |
|------|-------------|
| `udp.go` | `UDP` struct — reads datagrams, splits by newline, parses via `parse.Plain()`. Optional internal buffer channel. |
| `udp_test.go` | Integration tests |
| `stop_test.go` | Graceful shutdown tests |

## For AI Agents

### Working In This Directory
- Simplest receiver: reads UDP packets, splits by `\n`, parses each line as plaintext metric.
- Default listen port: `:2003` (same as TCP — both can run simultaneously).
- `BufferSize` option adds an internal channel buffer between packet reading and cache storage.

### Testing Requirements
- `go test ./receiver/udp/`

## Dependencies

### Internal
- `receiver/` — plugin registry
- `receiver/parse/` — `Plain()` parser
- `points/`, `helper/`

### External
- `prometheus/client_golang`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
