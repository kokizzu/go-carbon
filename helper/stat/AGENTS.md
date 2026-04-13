<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# stat

## Purpose
Platform-specific file stat helpers for accessing file metadata (size, real size, atime, ctime, mtime). Used by carbonserver's file scanning to get accurate file statistics across platforms.

## Key Files

| File | Description |
|------|-------------|
| `types.go` | Defines `FileStats` struct: Size, RealSize (blocks*512), ATime, CTime, MTime with nanosecond variants |
| `stat.go` | Non-Linux `GetStat()` — extracts `FileStats` from `os.FileInfo`, falls back to MTime for CTime |
| `stat_linux.go` | Linux `GetStat()` — uses `syscall.Stat_t` to provide real ATime and CTime via `Atim`/`Ctim` |

## For AI Agents

### Working In This Directory
- Platform-specific files use build tag convention (`_linux.go` suffix).
- `GetStat(os.FileInfo) FileStats` is the single public function — same signature on all platforms, different implementations.
- On non-Linux, CTime falls back to MTime since `syscall.Stat_t` doesn't expose `Ctim` on all platforms.
- `RealSize` is calculated from `Blocks * 512` to account for sparse files.

## Dependencies

### External
- Standard library only (`os`, `syscall`)

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
