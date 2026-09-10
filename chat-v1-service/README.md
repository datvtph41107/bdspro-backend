# chat-v1-service

## Role

Older Chat source module còn trong repository.

## Runtime status

Repository source-health/compatibility module, không nằm default Compose. Không xóa nếu chưa có caller/runtime retirement proof.

## Entrypoint / lifecycle

`main.go` -> Cobra commands.

`make server` mặc định `grpc`; `make dev` khi dùng Air.

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

Root có thể gọi `make migrate service=chat-v1-service`.

## Dependencies / durable state

Xem `config/runtime.yml`; không có canonical migration authority trong current source.

## Read source from here

`main.go`, `cmd/`, `internal/`, `infra/`, `config/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
