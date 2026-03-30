# Roadmap

## Phase 0

Goal: working board-safe skeleton.

- Go service with modular packages
- SQLite persistence
- health, write, list, recall endpoints
- Docker image
- ARM64-friendly defaults

## Phase 1

Goal: real memory semantics.

- `forget` support
- expiry timestamps
- profile summary endpoint
- SQLite FTS-backed lexical search
- structured metadata on memories

## Phase 2

Goal: closer Supermemory behavior.

- update and contradiction detection
- memory relationship edges
- background ingestion jobs
- external embedding integration
- hybrid recall and reranking

## Phase 3

Goal: agent integration.

- MCP adapter
- container-tag aware tools
- prompt-context export
- project discovery and listing

## Phase 4

Goal: document and connector ingestion.

- file upload
- extraction pipeline
- Gmail, Drive, GitHub, Notion style connectors
- webhook or poll-based sync

## Suggested Order

Do not start with connectors or web UI.

The right order is:

1. validate board runtime and persistence
2. validate memory semantics
3. validate agent integration
4. only then expand ingestion surface
