<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# tags

## Purpose
Handles Graphite tagged metrics: normalizes tag strings into canonical form and manages a persistent queue for sending tagged metric names to an external TagDB service. This implements the tagged metrics support described in [Graphite tag documentation](https://graphite.readthedocs.io/en/latest/tags.html).

## Key Files
| File | Description |
|------|-------------|
| `normalize.go` | Tag normalization. `Normalize(s)` — public function, sorts tags by key, deduplicates, returns canonical `metric;tag1=val1;tag2=val2` form. `normalizeOriginal(s)` — mirrors Python carbon's normalization. `FilePath(root, s, hashOnly)` — maps a tagged metric name to an on-disk path under `_tagged/`, using SHA256 hash prefix sharding. |
| `tags.go` | Package-level wiring (top-level type definitions or init logic). |
| `queue.go` | `Queue` struct — LevelDB-backed persistent queue. `NewQueue()`, `Add(metric)`, `Stop()`, `Lag()`. Background worker (`sendWorker`) drains the queue in configurable chunk sizes via a user-supplied `send func([]string) error` callback. Handles corrupted LevelDB by moving and recreating. |
| `normalize_test.go` | Tests for `Normalize` and `FilePath`. |
| `queue_test.go` | Tests for `Queue` (add/send lifecycle). |

## For AI Agents

### Working In This Directory
- Tagged metric format: `metricname;key1=val1;key2=val2`. Metrics without `;` are skipped by `Queue.Add` and are treated as untagged.
- `Normalize` uses a stable sort (`sort.Stable`) on the tag segment slice — this preserves ordering of equal keys before deduplication.
- `FilePath` with `hashOnly=false` uses the literal metric string (dots replaced with `_DOT_`) as the filename, sharded into `_tagged/<h0:3>/<h3:6>/`. With `hashOnly=true` only the SHA256 hash is used as filename — safer for very long metric names.
- LevelDB keys are `uint64 timestamp (big-endian) + metric string`. This ensures iteration order equals insertion order. The `Lag()` method reads the oldest key's timestamp to compute queue depth.
- The `send` callback is injected at construction time — it should call the TagDB HTTP endpoint. If it returns an error, the batch is retried after 100ms.
- `Queue` embeds `helper.Stoppable`; call `q.Stop()` for clean shutdown (flushes LevelDB close).

### Testing Requirements
```
go test ./tags/
```
Both `normalize_test.go` and `queue_test.go` are present. Queue tests use `helper/qa` for temp directories.

### Common Patterns
- Error handling in `Queue`: LevelDB errors increment atomic counters (`putErrors`, `deleteErrors`, `sendFail`) and are logged via `zapwriter`/`zap` — they do not stop the queue.
- Stat fields use `sync/atomic` increments, consistent with the rest of go-carbon.
- The `changed` channel (`cap=1`, non-blocking send) is used as a "work available" signal to the send worker goroutine, avoiding polling latency.

## Dependencies

### Internal
- `helper/` — `helper.Stoppable` (embedded in `Queue`), `zapwriter` (logging)

### External
- `github.com/syndtr/goleveldb/leveldb` — persistent queue storage
- `go.uber.org/zap` — structured logging
- `github.com/lomik/zapwriter` — logger factory (`zapwriter.Logger("tags")`)
- `crypto/sha256`, `encoding/binary`, `sort`, `strings` — standard library

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
