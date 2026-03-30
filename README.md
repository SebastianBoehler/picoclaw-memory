# PicoClaw Memory

PicoClaw Memory is a board-first memory service inspired by Supermemory and designed for small ARM64 Linux machines.

The project focuses on the core value of agent memory without dragging in a cloud-heavy stack:

- persistent memory writes
- container or project scoping
- recall over stored context
- profile-ready APIs
- Docker deployment on low-RAM boards

The design target is a single Go service with SQLite on local SSD, external model providers for expensive extraction or embeddings, and a runtime profile that works well on Orange Pi class hardware.

## Status

This project is early but functional.

Implemented now:

- `GET /healthz`
- `POST /v1/memories`
- `GET /v1/memories`
- `POST /v1/recall`
- SQLite persistence with WAL mode
- container-tag scoped storage
- lexical recall
- `linux/arm64` and `linux/amd64` Docker builds

Planned next:

- `forget` and expiry support
- profile synthesis
- SQLite FTS-backed recall
- contradiction and update handling
- embeddings and hybrid recall
- MCP adapter

## Why Go

Go is the default implementation language here because it gives us:

- lower delivery risk on ARM64 than a full Rust rewrite
- a simple static binary for boards and containers
- low operational overhead
- good enough performance for board-local memory workloads

Rust is still a valid future option for specific hot paths. The first version optimizes architecture before language choice.

## Design Principles

- one process, one local database
- no Postgres requirement
- no local model inference by default
- external providers for heavy extraction or embeddings
- keep RAM usage conservative
- prefer simple deployment and recovery

## Architecture

High-level architecture and scope decisions live in:

- [docs/architecture.md](./docs/architecture.md)
- [docs/roadmap.md](./docs/roadmap.md)

Current package layout:

- `cmd/server`: binary entrypoint
- `internal/config`: env-driven configuration
- `internal/httpapi`: HTTP routes and transport concerns
- `internal/memory`: domain types and service layer
- `internal/storage/sqlite`: SQLite-backed persistence

## Quick Start

### Requirements

- Go `1.21+`
- Docker, if you want containerized execution
- [`task`](https://taskfile.dev/) for the convenience commands

### Run locally

```bash
task tidy
task run
```

Default config:

- listen addr: `:8080`
- sqlite db: `./var/memory.db`

### Example write

```bash
curl -X POST http://localhost:8080/v1/memories \
  -H 'Content-Type: application/json' \
  -d '{
    "containerTag": "user_123",
    "content": "User prefers concise answers and runs ARM64 boards.",
    "source": "chat"
  }'
```

### Example recall

```bash
curl -X POST http://localhost:8080/v1/recall \
  -H 'Content-Type: application/json' \
  -d '{
    "containerTag": "user_123",
    "query": "ARM64 boards",
    "limit": 5
  }'
```

### Example list

```bash
curl 'http://localhost:8080/v1/memories?containerTag=user_123&limit=10'
```

## Docker

Build and run:

```bash
docker build -t picoclaw-memory .
docker run --rm -p 8080:8080 -v "$(pwd)/data:/data" picoclaw-memory
```

The image is intended to work on `linux/arm64` and `linux/amd64`.

## API

### `GET /healthz`

Returns a simple health payload.

### `POST /v1/memories`

Create a memory record.

Request body:

```json
{
  "containerTag": "user_123",
  "content": "User runs ARM64 boards with low RAM.",
  "source": "chat",
  "expiresAt": "2026-04-30T00:00:00Z"
}
```

### `GET /v1/memories`

List recent memories for one container.

Query params:

- `containerTag`
- `limit`

### `POST /v1/recall`

Recall memories by query within one container.

Request body:

```json
{
  "containerTag": "user_123",
  "query": "ARM64 low RAM",
  "limit": 5
}
```

## Verification

Current verification performed on this repo:

- `go build ./...`
- `docker build -t picoclaw-memory .`
- container smoke test for `healthz`, memory write, and recall

## Contributing

Contributions are welcome. Start with [CONTRIBUTING.md](./CONTRIBUTING.md).

For larger changes, open an issue or a draft PR first so the design direction stays aligned with the board-first scope.

## License

MIT. See [LICENSE](./LICENSE).
