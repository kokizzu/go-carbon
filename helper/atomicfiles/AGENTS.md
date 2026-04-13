<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# atomicfiles

## Purpose
Provides atomic file write operations — writes to a temp file then renames, ensuring no partial writes are visible to readers.

## Key Files

| File | Description |
|------|-------------|
| `atomicfiles.go` | `WriteFile()` — write-to-temp-then-rename pattern for crash-safe file updates |

## For AI Agents

### Working In This Directory
- Single-function package. The write pattern is: create temp file in same dir → write → sync → close → rename.
- Used by the kafka receiver to persist consumer offsets atomically.

### Testing Requirements
- No dedicated tests; tested implicitly through kafka receiver tests.

## Dependencies

### External
- Standard library only (`os`, `path`)

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
