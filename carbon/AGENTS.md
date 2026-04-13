<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# carbon

## Purpose
Application orchestrator for go-carbon. Owns all top-level components and controls their lifecycle: startup, shutdown, config reload, and graceful dump/restore. This is the entry point package — `main` constructs an `App` and calls `Start()`, `Loop()`, and `Stop()`.

## Key Files
| File | Description |
|------|-------------|
| `app.go` | `App` struct owning all components; `New()`, `Start()`, `Stop()`, `Loop()`, `ReloadConfig()`, `configure()`, `startPersister()` |
| `config.go` | `Config` struct with TOML tags; `ReadConfig()`, `NewConfig()`, `PrintDefaultConfig()`; `Duration` wrapper for TOML time parsing |
| `collector.go` | `Collector` gathers stats from all components at `MetricInterval`; implements `statModule` interface via `Stat(send StatCallback)`; sends internal metrics either to local cache or remote TCP/UDP endpoint |
| `grace.go` | `DumpStop()` for SIGUSR2 graceful restart — diverts cache writes to xlog, dumps binary cache to disk, stops listeners; `Restore()` / `RestoreFromDir()` / `RestoreFromFile()` for reading dump files back into cache on startup |
| `app_stop_test.go` | Tests for clean shutdown ordering |
| `restore_test.go` | Tests for dump/restore round-trip |

## For AI Agents

### Working In This Directory
- `App` is the single authoritative owner of all component pointers. Always acquire `app.Lock()` before reading/writing component fields outside of `Start()`.
- `ReloadConfig()` only recreates `Persister`, `Tags`, and `Collector` — it does not restart receivers or carbonserver. Do not broaden this scope without understanding the startup ordering in `Start()`.
- `startPersister()` is called both from `Start()` and `ReloadConfig()`. Keep it free of side effects that must not run twice.
- `Collector` must be re-created on every config reload (comment in `App` struct is authoritative).
- The `DumpStop()` flow is order-sensitive: stop persister → divert cache to xlog → dump cache → stop listeners → flush xlog. Do not reorder.
- `grace.go` sorts dump files by nanotimestamp with `cache.*` before `input.*` at the same timestamp to ensure correct restore order.
- `Duration` in `config.go` is a TOML-aware wrapper around `time.Duration`; always call `.Value()` to get the underlying duration.
- Config sections are unexported structs (`commonConfig`, `whisperConfig`, etc.) embedded in the exported `Config`.

### Testing Requirements
```
go test ./carbon/
```
Tests require a writable temp directory (created via `carbon.TestConfig(rootDir)`).

### Common Patterns
- Component lifecycle: check `if app.X != nil { app.X.Stop(); app.X = nil }` before re-creating.
- All components that expose stats implement `statModule` interface: `Stat(send helper.StatCallback)`.
- Receivers are registered once via `registerPluginsOnce.Do(...)` in `New()`.
- Config-driven feature flags: each major subsystem has an `Enabled bool` TOML field checked before initialization.
- Logging via `zapwriter.Logger("component-name")`, structured fields with `zap.String`, `zap.Error`, `zap.Duration`.

## Dependencies

### Internal
- `github.com/go-graphite/go-carbon/api` — gRPC API server
- `github.com/go-graphite/go-carbon/cache` — in-memory metric buffer
- `github.com/go-graphite/go-carbon/carbonserver` — HTTP/gRPC read server
- `github.com/go-graphite/go-carbon/helper` — `Stoppable`, `StatCallback`, throttle ticker
- `github.com/go-graphite/go-carbon/persister` — whisper write workers
- `github.com/go-graphite/go-carbon/points` — core data types
- `github.com/go-graphite/go-carbon/receiver` — receiver plugin registry and built-in receivers
- `github.com/go-graphite/go-carbon/tags` — tag normalization and TagDB sync

### External
- `github.com/BurntSushi/toml` — config file parsing
- `github.com/prometheus/client_golang/prometheus` — Prometheus metrics registry
- `github.com/lomik/zapwriter` — structured logger configuration
- `go.uber.org/zap` — structured logging

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
