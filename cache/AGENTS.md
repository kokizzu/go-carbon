<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# cache

## Purpose
Sharded in-memory buffer that sits between receivers and the persister. Receivers write points in via `Add()`; the persister reads them out via `WriteoutQueue` / `Pop()`; carbonserver reads live (unflushed) data via `Get()`. Designed for high write concurrency through CRC32-based sharding across 1024 independent `Shard` structs, each with its own `sync.RWMutex`.

## Key Files
| File | Description |
|------|-------------|
| `cache.go` | `Cache` struct and `Shard` struct; `Add()`, `Pop()`, `PopNotConfirmed()`, `Get()`, `Confirm()`; bloom filter for new-metric detection; throttle hook; `Stat()` for internal metrics |
| `writeout_queue.go` | `WriteoutQueue` orders metric keys for the persister; rebuilds lazily when drained; respects `WriteStrategy` ordering via `makeQueue()` |
| `queue.go` | `makeQueue()` — builds the ordered channel of metric keys based on write strategy (`MaximumLength`, `TimestampOrder`, `Noop`) |
| `carbonlink.go` | `CarbonlinkListener` — legacy pickle-protocol TCP server; allows graphite-web to query cached (not-yet-persisted) data points via `cache-query` requests |
| `dump.go` | `Dump()` (plain text) and `DumpBinary()` (varint-encoded) — iterate all shards including `notConfirmed` slice, write to `io.Writer` for graceful restart |
| `confirm.go` | `Confirm()` removes a `*points.Points` from the `notConfirmed` slice after successful persister write |
| `cache_test.go` | Unit tests for Add/Pop/Get and overflow behaviour |
| `dump_test.go` | Round-trip dump/restore tests |
| `confirm_test.go` | Tests for the confirm/notConfirmed flow |
| `carbonlink_test.go` | Tests for pickle request parsing and reply encoding |
| `xlog_test.go` | Tests for xlog (write-ahead log) diversion via `DivertToXlog()` |

## For AI Agents

### Working In This Directory
- Shard selection is `crc32(metricName) & 1023` — implemented in `helper.HashString`. Never assume a metric maps to a specific shard in tests.
- The `notConfirmed` mechanism exists to prevent data loss: `PopNotConfirmed()` removes a metric from `shard.items` but keeps the pointer in `shard.notConfirmed`. The persister calls `Confirm()` only after a successful whisper write. `Get()` scans both `items` and `notConfirmed` to return all live data.
- `Pop()` (used when confirmation is disabled) does a plain delete with no safety net. `PopNotConfirmed()` / `Confirm()` is the safe path used by the whisper persister.
- `DivertToXlog(w)` switches `Add()` to write plain-text lines to `w` instead of buffering in memory. Used by `DumpStop()` to capture incoming points while the cache is being serialized.
- Settings (`maxSize`, `xlog`, `tagsEnabled`) are stored in an `atomic.Value`-wrapped `cacheSettings` struct — copy-on-write, never mutate in place.
- Bloom filter (`newMetricCf`) tracks whether a metric has been seen before. When set, `Add()` only sends to `newMetricsChan` on a bloom miss, reducing channel pressure for existing metrics.
- `WriteoutQueue.Get(abort)` blocks until a metric key is available or `abort` is closed. It rebuilds the queue channel lazily with a 100ms debounce between rebuilds.
- `WriteStrategy` values map from TOML strings: `"max"` → `MaximumLength`, `"sorted"` → `TimestampOrder`, `"noop"` → `Noop`.
- The carbonlink protocol uses Python pickle encoding (protocols 2, 3, 4, 5). The parser in `carbonlink.go` handles all variants. Do not simplify it.

### Testing Requirements
```
go test ./cache/
```
No external dependencies required. Tests are self-contained with in-memory cache instances.

### Common Patterns
- Shard lock discipline: always `shard.mu.Lock()` for writes, `shard.mu.RLock()` for reads. Never hold the cache-level `c.mu` (used only for `writeStrategy` changes) while holding a shard lock.
- Atomic counters for stats: `atomic.AddUint32(&c.stat.overflowCnt, ...)`, `atomic.AddInt64(&c.stat.size, ...)`. Stats are read in `Stat()` via `helper.SendAndSubstractUint32` (delta reporting) or `helper.SendUint32` (cumulative).
- `helper.Stoppable` embedded in `CarbonlinkListener` for goroutine lifecycle management (`Go()`, `Stop()`, `StartFunc()`).

## Dependencies

### Internal
- `github.com/go-graphite/go-carbon/helper` — `Stoppable`, `StatCallback`, `HashString`, `SendAndSubstractUint32`
- `github.com/go-graphite/go-carbon/points` — `Points`, `Point` types; `WriteTo()`, `WriteBinaryTo()`
- `github.com/go-graphite/go-carbon/tags` — `tags.Normalize()` called in `Add()` when tags are enabled

### External
- `github.com/greatroar/blobloom` — blocked bloom filter for new-metric detection
- `github.com/lomik/graphite-pickle/framing` — framed TCP connection for carbonlink pickle protocol
- `go.uber.org/zap` — structured logging
- `github.com/lomik/zapwriter` — logger retrieval by name

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
