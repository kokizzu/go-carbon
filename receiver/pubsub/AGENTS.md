<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# pubsub

## Purpose
Google Cloud Pub/Sub receiver for consuming metrics from a Pub/Sub subscription.

## Key Files

| File | Description |
|------|-------------|
| `pubsub.go` | `PubSub` struct — subscribes to a GCP Pub/Sub subscription, parses messages as plaintext metrics. Supports gzip-compressed messages. Configurable goroutine count and max outstanding messages/bytes. |
| `pubsub_test.go` | Tests using mock Pub/Sub server |

## For AI Agents

### Working In This Directory
- Uses `cloud.google.com/go/pubsub` client library.
- Automatically detects and decompresses gzip'd messages (pooled gzip readers for efficiency).
- Configured via TOML `[[receiver.pubsub]]` section with project, subscription, goroutine count.

### Testing Requirements
- `go test ./receiver/pubsub/` — uses a mock Pub/Sub server, no real GCP needed.

## Dependencies

### Internal
- `receiver/`, `receiver/parse/`, `points/`, `helper/`

### External
- `cloud.google.com/go/pubsub` — Google Cloud Pub/Sub client
- `prometheus/client_golang`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
