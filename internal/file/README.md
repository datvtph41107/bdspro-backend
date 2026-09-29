# file-service

## Role

File metadata/content, signed access, media/version artifact HTTP boundary.

## Runtime status

Canonical default Compose service.

## Entrypoint / lifecycle

`main.go` trực tiếp compose config, Auth/Hub clients, DB, file services, HTTP router và graceful shutdown.

`make server`.

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

Root có thể gọi `make migrate service=file-service`.

## Dependencies / durable state

PostgreSQL `file_service`, Auth/Hub RPC và local file volume. Canonical migration: `database/migrations/`.

## Read source from here

`main.go`, `internal/file*`, `infra/handler/filehttp/`, `repositories/`, `database/migrations/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../../docs/reference/SERVICE-MAP.md`](../../docs/reference/SERVICE-MAP.md) và [`../../docs/reference/DEVELOPER-OPERATING-GUIDE.md`](../../docs/reference/DEVELOPER-OPERATING-GUIDE.md).
