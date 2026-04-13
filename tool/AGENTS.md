<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# tool

## Purpose
Standalone CLI tools for go-carbon operations and administration. These are compiled binaries, not importable packages. Currently contains one tool for pre-deployment validation of persister configuration changes.

## Key Files
| File | Description |
|------|-------------|
| `persister_configs_differ/main.go` | CLI tool to diff persister configs. Reads old and new `storage-schemas.conf` and `storage-aggregation.conf`, iterates over a carbonserver file-list cache (`.gzip`), and reports how many metrics would change schema retention or aggregation method. Outputs a summary of change counts grouped by rule transition (e.g. `60s:7d->10s:30d`). Supports `--print-metrics` to emit each affected metric path. |

## Subdirectories
| Directory | Purpose |
|-----------|---------|
| `persister_configs_differ/` | Source for the `persister_configs_differ` binary |

## For AI Agents

### Working In This Directory
- This is a `package main` binary, not a library. Do not import it from other packages.
- The tool reads a carbonserver file-list cache (produced by go-carbon's carbonserver when `file-list-cache` is configured). Pass the `.gzip` cache file via `-file-list-cache`.
- All four config flags (`-old-schema`, `-new-schema`, `-old-aggregation`, `-new-aggregation`) are required for a meaningful diff — the tool will panic on missing/invalid files.
- There is a `TODO` comment in the source suggesting this could become a go-carbon sub-command; do not remove it without addressing the underlying design question.
- Commented-out debug logging (`log.Printf`, counter increments) is intentionally left in the source as scaffolding — do not remove without discussion.

### Testing Requirements
```
go build ./tool/...
```
No unit tests exist. Validation is done by running the binary against real config files and a file-list cache snapshot.

### Common Patterns
- Uses `flag` package for CLI argument parsing (standard go-carbon tooling pattern).
- Panics on errors from config parsing and cache reading — this is intentional for a CLI tool where errors are fatal.
- Output is plain text to stdout: one section per config type (`schema-changes`, `aggregation-changes`), each line is `old->new count`.

## Dependencies

### Internal
- `carbonserver/` — `carbonserver.NewFileListCache`, `carbonserver.FLCVersionUnspecified` for reading the metric file-list cache
- `persister/` — `persister.ReadWhisperSchemas`, `persister.ReadWhisperAggregation` for parsing config files

### External
- `flag`, `fmt`, `io`, `errors`, `strings` — standard library only

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
