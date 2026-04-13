<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# receiver

## Purpose
Plugin-based metric ingestion layer. Provides a protocol registry that allows
multiple receiver types (TCP, UDP, HTTP, Kafka, Google Cloud Pub/Sub) to be
instantiated from TOML configuration and feed parsed metric points into the
cache via a shared `store` callback.

## Key Files
| File | Description |
|------|-------------|
| `receiver.go` | Plugin registry. Defines the `Receiver` interface (`Stop()`, `Stat()`, `InitPrometheus()`), the global `protocolMap`, `Register()` to add a protocol, `New()` to instantiate a receiver from a TOML options map, and `WithProtocol()` to serialize typed options into a map with a `protocol` key. |
| `tcp/` | TCP plaintext and pickle receiver |
| `udp/` | UDP plaintext receiver |
| `http/` | HTTP receiver (POST metrics) |
| `kafka/` | Apache Kafka consumer receiver |
| `pubsub/` | Google Cloud Pub/Sub receiver |
| `parse/plain.go` | Plaintext `metric value timestamp\n` parser |
| `parse/pickle.go` | Python pickle format parser |
| `parse/protobuf.go` | Protocol Buffers metric parser |
| `parse/msgpack.go` | MessagePack metric parser |

## Subdirectories
| Directory | Purpose |
|-----------|---------|
| `tcp/` | TCP receiver — plaintext and pickle protocols |
| `udp/` | UDP receiver — plaintext protocol |
| `http/` | HTTP receiver — accepts POSTed metrics |
| `kafka/` | Kafka consumer receiver |
| `pubsub/` | Google Cloud Pub/Sub receiver |
| `parse/` | Format parsers for all ingestion protocols |

## For AI Agents

### Working In This Directory
- The `Receiver` interface has exactly three methods: `Stop()`, `Stat(helper.StatCallback)`, and `InitPrometheus(prometheus.Registerer)`. All protocol implementations must satisfy this interface.
- Each protocol package (e.g. `tcp`, `udp`) calls `receiver.Register()` from its own `init()` or `Register()` function. These are wired up in `carbon/app.go` inside `registerPluginsOnce`.
- Configuration round-trips through TOML: `WithProtocol()` encodes a typed options struct to TOML and back to `map[string]interface{}`, appending a `"protocol"` key. `New()` reverses this to decode into the protocol's options type.
- To add a new protocol: create a subdirectory, define an options struct with TOML tags, implement `Receiver`, call `receiver.Register()` with a name, an options factory (`func() interface{}`), and a receiver factory, then call your `Register()` from `carbon/app.go`.
- The `store func(*points.Points)` callback passed to each receiver is the write path into the in-memory cache. Do not bypass it.

### Testing Requirements
```
go test ./receiver/...
```
Parser tests live in `receiver/parse/` and include benchmark (`bench_test.go`) and format-specific tests.

### Common Patterns
- Protocol registration via `receiver.Register()` at init time (plugin registry pattern).
- Options structs use `toml` struct tags and are round-tripped through TOML encoding for config decoding — do not use `json` tags here.
- `helper.StatCallback` is used for internal graphite self-monitoring stats; implement it in every new receiver.
- Prometheus metrics are registered via `InitPrometheus(prometheus.Registerer)` — use a `prometheus.Registerer` (not the default registry) so tests can register without conflicts.

## Dependencies

### Internal
- `github.com/go-graphite/go-carbon/helper` — `StatCallback`, `Stoppable`
- `github.com/go-graphite/go-carbon/points` — `Points` struct passed through the store callback

### External
- `github.com/BurntSushi/toml` — TOML encoding/decoding for options round-trip
- `github.com/prometheus/client_golang/prometheus` — Prometheus metrics registration

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
