<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# api

## Purpose
Provides a gRPC API for cache queries, implementing a CarbonLink-compatible interface over gRPC. Allows external tools (e.g. carbonapi) to query in-memory metric data points from the write cache before they are persisted to disk.

## Key Files
| File | Description |
|------|-------------|
| `api.go` | `Api` struct: wraps a gRPC server, implements `carbonpb.CarbonServer`. Methods: `New()`, `Listen()`, `Stop()`, `Addr()`, `Stat()`, `CacheQuery()`. Tracks per-request counters atomically. |
| `sample/cache-query/cache-query.go` | Standalone CLI client for manual testing of the gRPC cache query API. Connects to a running go-carbon instance and queries named metrics. |

## Subdirectories
| Directory | Purpose |
|-----------|---------|
| `sample/cache-query/` | Example gRPC client binary (`cache-query`); not imported by any other package |

## For AI Agents

### Working In This Directory
- `Api` embeds `stop.Struct` (from `lomik/stop`), not `helper.Stoppable` — use `StartFunc`/`Go` patterns from that package.
- The gRPC server is registered with `carbonpb.RegisterCarbonServer` and gRPC reflection. The only RPC method is `CacheQuery`.
- Stats counters (`cacheRequests`, `cacheRequestMetrics`, `cacheResponseMetrics`, `cacheResponsePoints`) use `sync/atomic` and are drained (subtract-and-send) via `helper.SendAndSubstractUint32`.
- The default port for the gRPC API is `7003` (see `sample/cache-query/cache-query.go` default flag).
- Do not add new RPC methods without updating `helper/carbonpb/carbon.proto` and regenerating the pb file.

### Testing Requirements
```
go test ./api/
```
No test files currently exist in `api/`; the `sample/` binary can be used for manual integration testing against a live instance.

### Common Patterns
- Goroutine lifecycle via `api.Go(func(exit chan struct{}) { ... })` — goroutines watch the `exit` channel for shutdown.
- All stat increments use `atomic.AddUint32`.
- gRPC server is stopped on exit via a dedicated goroutine that waits on the exit channel.

## Dependencies

### Internal
- `cache/` — `*cache.Cache` is the only data source; queried via `cache.Get(metric)`
- `helper/` — `helper.StatCallback`, `helper.SendAndSubstractUint32`
- `helper/carbonpb/` — generated gRPC service and message types

### External
- `github.com/lomik/stop` — provides `stop.Struct` with `StartFunc`/`Go` lifecycle
- `google.golang.org/grpc` — gRPC server
- `google.golang.org/grpc/reflection` — gRPC server reflection (for tooling like grpcurl)
- `golang.org/x/net/context` — context package

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
