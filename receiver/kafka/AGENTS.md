<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# kafka

## Purpose
Apache Kafka receiver for consuming metrics from Kafka topics. Supports multiple message formats and consumer group management.

## Key Files

| File | Description |
|------|-------------|
| `kafka.go` | `Kafka` struct — Sarama consumer group that reads metrics from Kafka topics. Supports plaintext, pickle, protobuf, and msgpack formats. Persists consumer offsets atomically to disk. Custom `Offset` type handles "newest"/"oldest" config values. |

## For AI Agents

### Working In This Directory
- Uses IBM/sarama Kafka client library.
- Message format is configurable: `plain`, `pickle`, `protobuf`, `msgpack`.
- Offsets can be persisted to a local file (via `helper/atomicfiles`) as backup beyond Kafka's consumer group offsets.
- TLS, SASL, and Kerberos authentication supported via config.
- Custom `Offset` type with `MarshalText`/`UnmarshalText` for TOML config parsing of "newest"/"oldest" strings.

### Testing Requirements
- `go test ./receiver/kafka/` — requires no real Kafka broker (uses mocks).

## Dependencies

### Internal
- `receiver/`, `receiver/parse/`, `points/`, `helper/`, `helper/atomicfiles/`

### External
- `github.com/IBM/sarama` — Apache Kafka client
- `prometheus/client_golang`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
