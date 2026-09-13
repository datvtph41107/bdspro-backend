# hub-service

## Role

Location/system configuration/API-key/version support capabilities hiện hữu.

## Runtime status

Canonical default Compose service.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc` -> `cmd/grpc`.

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
make migrate
make rollback
make migration name=<schema_change>
make migration-version
```

Root có thể gọi `make migrate service=hub-service`.

## Dependencies / durable state

PostgreSQL `hub_service`, Redis DB 3, User/Notification RPC. Canonical migration: `database/migrations/`.

## Read source from here

`cmd/grpc/`, `internal/`, `infra/`, `config/`, `database/migrations/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
