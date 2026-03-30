# Contributing

Thanks for contributing to PicoClaw Memory.

The project is intentionally opinionated: it prioritizes small-board deployments, low operational complexity, and clear memory semantics over broad platform coverage.

## Before You Start

- Check the current direction in [docs/architecture.md](./docs/architecture.md) and [docs/roadmap.md](./docs/roadmap.md).
- Prefer opening an issue or a draft PR for larger changes before investing heavily.
- Keep the board-first constraints in mind. A simpler design that fits ARM64 boards is usually preferred over a more ambitious cloud-style design.

## Development Setup

Requirements:

- Go `1.21+`
- Docker
- [`task`](https://taskfile.dev/) for convenience commands

Recommended workflow:

```bash
task tidy
task build
task test
```

Run locally:

```bash
task run
```

## Project Expectations

- Keep files reasonably small and modular.
- Avoid mock or fallback behavior unless it is explicitly part of the feature.
- Prefer explicit errors over silent degradation.
- Keep API behavior predictable and easy to reason about.
- Favor SQLite and single-process designs unless there is a strong reason not to.

## Pull Request Guidelines

- Keep PRs scoped.
- Include the motivation, implementation summary, and verification steps.
- Add or update docs when behavior changes.
- If you change API shape, include example requests and responses.
- If you add dependencies, explain why they are justified for board deployments.

## Coding Notes

- Use `gofmt`.
- Keep the public surface small and clear.
- Prefer straightforward code over abstraction-heavy patterns.
- Add tests for non-trivial behavior when practical.

## Reporting Issues

When filing a bug, include:

- expected behavior
- actual behavior
- reproduction steps
- board or host architecture
- Docker or local runtime details

## Security

Do not open public issues for sensitive vulnerabilities. Use the process in [SECURITY.md](./SECURITY.md).
