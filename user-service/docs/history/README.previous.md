# User Service

User Service is an independently buildable QHPRO backend service. It owns User
business/data state and consumes Shared only for reusable technical mechanisms
and cross-service contracts.

## Requirements

- Go 1.25+
- Docker + Docker Compose for local dependencies / image verification
- `golang-migrate` for schema operations
- Buf and Wire through the repository Shared generation workflow

## Local configuration

Committed YAML profiles contain non-secret defaults only. User runtime dùng
`user-service/.env`; root `.env` chỉ dành cho full Compose:

```bash
make -C .. env-check
```

The serving process accepts canonical resource variables such as
`USER_DATABASE_URL`, `GRPC_ADDRESS`, `USER_REDIS_ADDRESS` and
`JWT_KEY_GENERATE`.
Historical YAML/env keys remain compatibility inputs while old providers are
being retired.

## Daily development

Keep backing services in containers and run the active Go process on the host:

```bash
make deps-up
make migrate-up
make test
make server
```

Stop local dependencies with:

```bash
make deps-down
```

## Schema

`migrate/` is the only production schema-evolution authority for User Service.
Historical replacement/backfill/cutover SQL that is useful for an existing
deployment is recovered from source-control history; it is never part of the serving
startup path and never competes with fresh-schema migration history.

```bash
make migrate-up
make migrate-down            # one version only
make migrate-version
make migrate-create name=add_example_state
```

Serving runtime does not execute schema mutation. Schema changes are authored
as ordered SQL under `migrate/`; model/schema inspection tooling must remain
operator/developer-only and outside the serving process.

## Generation

User provides the local command vocabulary; Shared remains generation
authority:

```bash
make buf
make wire
make generate
```

## Build and container verification

```bash
make build
make docker-build
```

The Dockerfile builds from repository-root context because `go.mod` uses local
Shared module replacements.

For a bounded production-shaped lab:

```bash
make compose-up
make compose-logs
make compose-down
```

The Compose flow waits for PostgreSQL/Redis health, runs migration as a one-shot
process, then starts the actual User image. Optional background actors are off
in this local image lab so that it tests the serving artifact rather than
unrelated scheduled work.

## Verification

Fast service-local gate:

```bash
make architecture-check
make acceptance-static
```

Full closure additionally requires a real PostgreSQL migration up/down/up lab,
Gateway -> User integration, User <-> Payment commercial flow, Docker image
execution, graceful shutdown proof, and retirement of competing legacy
commercial authority.
