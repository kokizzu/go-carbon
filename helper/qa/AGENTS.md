<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# qa

## Purpose
Test helpers shared across packages.

## Key Files

| File | Description |
|------|-------------|
| `tmp_dir.go` | `Root()` — creates a temporary directory for a test, calls cleanup on defer |

## For AI Agents

### Working In This Directory
- Single-function package used by test files throughout the project.
- `Root(t, func(dir string))` pattern ensures cleanup even on test failure.

## Dependencies

### External
- Standard library only (`os`, `testing`)

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
