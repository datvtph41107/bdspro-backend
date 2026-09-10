# organization-service

## Role

Organization, membership và organization action rights.

## Runtime status

Canonical default Compose service.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc` -> service composition.

`make server` chạy `go run . grpc`.

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

Root có thể gọi `make migrate service=organization-service`.

## Dependencies / durable state

PostgreSQL `organization` và RPC clients. Canonical migration: `migrate/`.

## Read source from here

`cmd/`, `internal/`, `infra/`, `config/`, `migrate/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
