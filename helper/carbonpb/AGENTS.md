<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-04-13 | Updated: 2026-04-13 -->

# carbonpb

## Purpose
Generated protobuf definitions for the carbon gRPC protocol. Defines the `CarbonLink` service and `Payload` message used by the gRPC API and protobuf receiver.

## Key Files

| File | Description |
|------|-------------|
| `carbon.proto` | Protobuf schema: `Payload` (metrics batch), `Metric`, `Point` messages, `CarbonLink` service |
| `carbon.pb.go` | Generated Go code from `carbon.proto` |
| `gen.go` | `go:generate` directive for protoc code generation |

## For AI Agents

### Working In This Directory
- **Do not edit `carbon.pb.go` directly** — it is generated. Edit `carbon.proto` and regenerate.
- Regenerate with: `go generate ./helper/carbonpb/`
- Uses `gogo/protobuf` (not standard `google.golang.org/protobuf`) for generation.

### Common Patterns
- The `Payload` message is a batch of `Metric` messages, each with a name and list of `Point` (value + timestamp).

## Dependencies

### External
- `github.com/gogo/protobuf` — protobuf runtime and code generation

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
