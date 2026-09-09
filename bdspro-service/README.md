# bdspro-service

## Role

BDS/property source module hiện hữu.

## Runtime status

Repository source-health module, không có standing service trong canonical default Compose. Không suy diễn là retired.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc`/`http`; Make default `grpc`.

`make server`; `make dev` khi dùng Air.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

Nếu thay schema:

```bash
make migrateup
make migrate-version
make new_migration name=<schema_change>
```

Root có thể gọi `make migrate service=bdspro-service`.

## Dependencies / durable state

Local DB fallback `db_bdspro`. Canonical source migration authority: `infra/db/migrations/v2/`.

## Read source from here

`main.go`, `cmd/`, `internal/`, `infra/`, `config/`, `infra/db/migrations/v2/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
