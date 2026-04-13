<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# persister

## Purpose
Writes cached metrics to Whisper (`.wsp`) files on disk. Manages a worker pool
that drains the cache's WriteoutQueue, creating or updating whisper files
according to retention schemas, aggregation rules, and optional per-namespace
quotas. Also supports online migration of whisper file configuration without
downtime.

## Key Files
| File | Description |
|------|-------------|
| `whisper.go` | Core `Whisper` struct and worker pool. Workers shard by `crc32(metricName) % workersCount`. Handles sparse file creation, `flock`, compressed whisper, hash filenames for tagged metrics, and online migration. Uses 32768 (`storeMutexCount`) sharded mutexes for per-file locking. |
| `whisper_schema.go` | Parses `storage-schemas.conf`. Maps metric name regex patterns to `whisper.Retentions`. `WhisperSchemas.Match()` returns the first matching schema. |
| `whisper_schema_test.go` | Tests for schema parsing and matching. |
| `whisper_aggregation.go` | Parses `storage-aggregation.conf`. Maps metric name patterns to aggregation method (`average`, `sum`, `min`, `max`, `last`) and `xFilesFactor`. Default: `average` / `0.5`. |
| `whisper_quota.go` | Parses `storage-quotas.conf` (go-carbon-specific). Defines per-namespace limits on metric count, logical/physical size, data points, and throughput with a configurable dropping policy. |
| `ini.go` | INI file parser used by schemas, aggregation, and quotas config files. Implements the standard Graphite INI format (not TOML). |

## Subdirectories
None.

## For AI Agents

### Working In This Directory
- The `Whisper` struct embeds `helper.Stoppable`; use `Stop()` / wait patterns from that helper for graceful shutdown.
- File-level locking uses `storeMutex[crc32(path) % storeMutexCount]` — 32768 mutexes sharded by path hash. Never hold multiple mutexes simultaneously to avoid deadlock.
- `hashFilenames` mode hashes tagged metric names to avoid filesystem path length limits (`maxPathLength = 4095`, `maxFilenameLength = 255`). Tagged metric detection lives in the `tags` package.
- Online migration (`onlineMigration`) live-updates retention, `xFilesFactor`, and aggregation method in existing whisper files. It is rate-limited via `helper.ThrottleTicker`. Controlled globally by `persister.Whisper` flags but can be overridden per-schema or per-aggregation entry.
- Config files (`storage-schemas.conf`, `storage-aggregation.conf`, `storage-quotas.conf`) use INI format parsed by `ini.go`, not TOML. Do not introduce TOML parsing here.
- `WhisperSchemas` and `WhisperAggregation` are matched in order of descending `Priority`. Ensure new schemas/aggregation entries set a meaningful priority.

### Testing Requirements
```
go test ./persister/
```
Includes `whisper_schema_test.go` and `whisper_stop_test.go`. The stop test verifies graceful shutdown of the worker pool.

### Common Patterns
- Sharded mutex pattern: `storeMutex[crc32.ChecksumIEEE([]byte(path)) % storeMutexCount]`.
- Rate limiting via `helper.ThrottleTicker` for both `maxUpdatesPerSecond` and online migration rate.
- Prometheus counters (`created`, `updateOperations`, `committedPoints`, `oooDiscardedPoints`) are updated via `atomic.AddUint32`.
- Structured logging via `go.uber.org/zap` (field-based, not `fmt.Sprintf`). Separate `logger` and `createLogger` fields for normal vs. file-creation log lines.

## Dependencies

### Internal
- `github.com/go-graphite/go-carbon/helper` — `Stoppable`, `ThrottleTicker`, `StatCallback`
- `github.com/go-graphite/go-carbon/points` — `Points` struct from the WriteoutQueue
- `github.com/go-graphite/go-carbon/tags` — Tagged metric detection and filename hashing

### External
- `github.com/go-graphite/go-whisper` — Whisper file read/write library
- `go.uber.org/zap` — Structured logging
- `github.com/prometheus/client_golang/prometheus` — Prometheus metrics
- `github.com/lomik/zapwriter` — zap writer initialization

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
