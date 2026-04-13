<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# carbonzipperpb

## Purpose
Generated protobuf definitions for the carbonzipper protocol. Shared types used by carbonserver for carbonapi v2 gRPC compatibility.

## Key Files

| File | Description |
|------|-------------|
| `carbonzipper.proto` | Protobuf schema for carbonzipper request/response types |
| `carbonzipper.pb.go` | Generated Go code from `carbonzipper.proto` |
| `gen.go` | `go:generate` directive for protoc code generation |

## For AI Agents

### Working In This Directory
- **Do not edit `carbonzipper.pb.go` directly** — regenerate from the proto file.
- Uses `gogo/protobuf` for generation.
- These types are largely superseded by `go-graphite/protocol` carbonapi v3 types but still used for backward compatibility.

## Dependencies

### External
- `github.com/gogo/protobuf`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
