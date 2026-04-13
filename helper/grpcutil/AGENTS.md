<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# grpcutil

## Purpose
gRPC server interceptors for request timing and rate limiting used by carbonserver's gRPC endpoint.

## Key Files

| File | Description |
|------|-------------|
| `grpcutil.go` | `UnaryServerTimeHandler` — interceptor that logs request payload, peer, and duration. `StreamServerTimeHandler` — same for streaming RPCs. `NewRateLimiter` — token-bucket rate limiter for gRPC calls. |

## For AI Agents

### Working In This Directory
- Interceptors are wired into carbonserver's gRPC server setup.
- Rate limiter uses atomic operations for lock-free concurrency.

## Dependencies

### External
- `google.golang.org/grpc` — gRPC framework

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
