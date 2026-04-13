<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# stat

## Purpose
Platform-specific file stat types for accessing file creation time (Ctim) used by carbonserver's file scanning.

## Key Files

| File | Description |
|------|-------------|
| `types.go` | Defines `StatCallback` type alias used throughout the project for metrics collection |
| `stat.go` | Default (non-Linux) stat helper |
| `stat_linux.go` | Linux-specific stat with `Ctim` (creation time) access via `syscall.Stat_t` |

## For AI Agents

### Working In This Directory
- Platform-specific files use build tags (filename convention `_linux.go`).
- `StatCallback` is `func(metric string, value float64)` — the universal stat reporting signature.

## Dependencies

### External
- Standard library (`syscall`)

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
