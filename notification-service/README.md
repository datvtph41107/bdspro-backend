# notification-service

## Role

Payment event Inbox dedupe, logical notification và delivery intent/effect.

RabbitMQ connection/channel thuộc payment-consumer supervisor. Broker restart
không được kéo gRPC API xuống; supervisor reconnect với bounded backoff, còn
durable queue + Inbox giữ recovery authority.

## Runtime status

Canonical default Compose service. Payment consumer/delivery actor thuộc owner process; compatibility binaries vẫn build được nhưng không có standing worker containers.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc` -> `cmd/grpc_server.go`.

`make server`; compatibility: `make event-worker`, `make delivery-worker` khi điều tra/recovery.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

Canonical database commands:

```bash
make migration name=<schema_change>
make migrate
make migration-version
make rollback
```

Root có thể gọi `make migrate service=notification-service`.

## Dependencies / durable state

PostgreSQL `qhpro_notification`, Redis DB 2, RabbitMQ, User RPC. Canonical migration: `database/migrations/`.

## Read source from here

`cmd/grpc_server.go`, `infra/postgres/eventing/`, `infra/worker/`, `internal/usecase/`, `database/migrations/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
