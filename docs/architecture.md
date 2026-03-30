# Architecture

## Goal

Build a board-native memory service that captures the useful parts of Supermemory without inheriting its cloud-heavy deployment model.

The first version optimizes for:

- small RAM usage
- ARM64 Docker deploys
- low operational complexity
- acceptable recall quality for personal and agent memory

## Non-Goals

Not in v0:

- full Supermemory parity
- Cloudflare Workers compatibility
- web dashboard
- realtime SaaS auth flows
- local model inference on the board

## Runtime Shape

One process, one local database.

```text
client or agent
    |
    v
HTTP API / MCP adapter
    |
    v
memory service
    |
    +-- SQLite metadata + text + future embeddings
    |
    +-- optional async job runner
    |
    +-- optional external embedding / extraction providers
```

## Storage Choice

Use SQLite in WAL mode on SSD.

Why:

- trivial deploy and backup story
- low idle memory
- strong enough concurrency for board-local use
- no separate Postgres or vector service

Initial tables:

- `memories`: stored memory units
- future `memory_edges`: update, extends, derives links
- future `profiles`: precomputed profile summaries
- future `ingest_jobs`: background extraction state

## Recall Strategy

### v0

Use lexical search first.

- direct content storage
- container filtering
- recency ordering
- substring scoring

This is enough to validate API shape and board behavior.

### v1

Add hybrid recall.

- lexical candidate generation with SQLite FTS
- optional embeddings from external provider
- rerank top lexical candidates
- avoid full ANN infrastructure until corpus size justifies it

### v2

If corpus size grows materially:

- add quantized embeddings
- add lightweight ANN or `sqlite-vec`
- keep lexical prefiltering to control RAM

## Memory Semantics

The important product behavior to preserve from Supermemory is not the exact implementation, but the semantics:

- stable facts
- recent activity
- project or container scoping
- updates replacing stale facts
- forgetting expired facts

Planned evolution:

1. raw memory records
2. explicit expiry support
3. update detection for contradictions
4. derived profile summaries
5. relationship graph between memories

## Board Constraints

Design assumptions:

- ARM64 host
- limited RAM
- enough SSD for local state
- Docker available
- external network access for embeddings or extraction when needed

Operational choices:

- no local embedding model by default
- no Postgres
- no Kafka, Redis, or sidecar queue
- keep indexes small and explicit
- persist only what is useful

## API Plan

### v0 API

- `GET /healthz`
- `POST /v1/memories`
- `GET /v1/memories`
- `POST /v1/recall`

### v1 API

- `POST /v1/forget`
- `POST /v1/profile`
- `POST /v1/documents`
- `GET /v1/jobs/:id`

### v2 API

- MCP endpoints
- connector sync endpoints
- admin and compaction endpoints

## Docker Strategy

Use a multi-stage Go build that emits a small final image.

- build once for `linux/amd64` and `linux/arm64`
- mount `/data` for SQLite persistence
- keep config entirely env-driven

## Why This Is Close Enough

If we preserve these, the system will feel close to Supermemory for your use case:

- same mental model of memory per user or project
- write once, recall later
- profile-ready context
- future contradiction handling
- future hybrid recall
- future MCP integration for agents

The cloud platform, SaaS dashboard, and broad connector surface are secondary for the board deployment.
