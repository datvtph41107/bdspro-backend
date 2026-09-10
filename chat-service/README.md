# chat-service

## Role

Chat source module với gRPC và HTTP commands.

## Runtime status

Repository source-health module, không nằm default Compose; canonical-vs-compatibility classification cần proof riêng.

## Entrypoint / lifecycle

`main.go` -> Cobra; `cmd/grpc/` và `cmd/http/`.

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

Root có thể gọi `make migrate service=chat-service`.

## Dependencies / durable state

Xem `config/runtime.yml` và source clients; không có canonical migration authority trong current source.

## Read source from here

`cmd/`, `internal/`, `infra/`, `config/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
