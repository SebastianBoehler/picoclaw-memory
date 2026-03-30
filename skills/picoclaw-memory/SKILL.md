---
name: picoclaw-memory
description: Board-first guidance for building, extending, and operating the PicoClaw Memory service. Use when working in this repository on the Go API, SQLite-backed memory storage, ARM64 Docker deployment, memory semantics such as save or recall or forget or profiles, or when another agent needs repo-specific instructions for changing and verifying the service.
---

# PicoClaw Memory

## Overview

Use this skill to work inside the `picoclaw-memory` repo without re-deriving the project shape.

The service is intentionally small and board-first:

- Go service
- SQLite in WAL mode
- `linux/arm64` friendly Docker image
- memory APIs before dashboards or connectors

Preserve that bias when making changes.

## Key Files

Read these first when the task is non-trivial:

- `README.md`
- `docs/architecture.md`
- `docs/roadmap.md`
- `internal/httpapi/router.go`
- `internal/memory/service.go`
- `internal/storage/sqlite/store.go`
- `Dockerfile`
- `Taskfile.yml`

Use the docs for intent and the `internal/` packages for the current implementation.

## Repo Workflow

When making changes:

1. Read the relevant docs and packages before editing.
2. Keep files small and modular.
3. Preserve explicit errors over hidden fallback behavior.
4. Prefer single-process, SQLite-first solutions unless there is a strong reason not to.
5. Update docs when API shape or runtime behavior changes.

## Service Boundaries

Current HTTP surface:

- `GET /healthz`
- `POST /v1/memories`
- `GET /v1/memories`
- `POST /v1/recall`

Current domain bias:

- container-tag scoped memory
- lexical recall
- SQLite persistence
- ARM64 container deployment

If adding features, extend in this order unless the task explicitly requires otherwise:

1. memory semantics such as forget, expiry, profile, contradiction handling
2. search quality such as FTS or hybrid recall
3. integration surfaces such as MCP
4. document ingestion and connectors

Do not jump straight to a cloud-heavy architecture.

## Implementation Rules

- Target Go `1.21`.
- Keep ARM64 Docker builds working.
- Prefer SQLite WAL mode and local SSD persistence.
- Avoid introducing Postgres, Redis, Kafka, or local model inference by default.
- Do not add mock data or silent fallback behavior unless explicitly requested.
- Keep new files under roughly 300 lines where practical.

## Verification

Use these checks after meaningful changes:

- `GOTOOLCHAIN=local go build ./...`
- `docker build -t picoclaw-memory .`

When HTTP behavior changes, run a container smoke test:

```bash
docker run --rm -d -p 18080:8080 --name picoclaw-memory-smoke picoclaw-memory
curl -fsS http://localhost:18080/healthz
curl -fsS -X POST http://localhost:18080/v1/memories \
  -H 'Content-Type: application/json' \
  -d '{"containerTag":"user_123","content":"User runs ARM64 PicoClaw boards with low RAM.","source":"smoke-test"}'
curl -fsS -X POST http://localhost:18080/v1/recall \
  -H 'Content-Type: application/json' \
  -d '{"containerTag":"user_123","query":"ARM64 PicoClaw","limit":5}'
docker rm -f picoclaw-memory-smoke
```

## Common Tasks

### Add a new memory endpoint

- update `internal/httpapi/router.go`
- extend `internal/memory/service.go`
- add store support in `internal/storage/sqlite/store.go`
- update `README.md` if the public API changed

### Improve search quality

- preserve the existing storage model first
- prefer SQLite FTS before adding a heavier retrieval layer
- keep RAM impact explicit in the change summary

### Prepare board deployment changes

- keep the image multi-stage and small
- preserve `linux/arm64` compatibility
- avoid runtime assumptions that require large memory headroom
