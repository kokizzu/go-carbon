<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# http

## Purpose
HTTP receiver for Graphite plaintext protocol over HTTP POST requests.

## Key Files

| File | Description |
|------|-------------|
| `http.go` | `HTTP` struct — accepts POST requests with plaintext metrics in the body. Configurable max message size (default 64MB). |
| `http_test.go` | Integration tests |
| `stop_test.go` | Graceful shutdown tests |

## For AI Agents

### Working In This Directory
- Default listen port: `:2007`.
- Reads the full request body, splits by newline, parses each line via `parse.Plain()`.
- `MaxMessageSize` limits body size to prevent memory exhaustion.

### Testing Requirements
- `go test ./receiver/http/`

## Dependencies

### Internal
- `receiver/`, `receiver/parse/`, `points/`, `helper/`

### External
- `prometheus/client_golang`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
