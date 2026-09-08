# payment-service

## Role

Order, payment attempt/provider evidence, settlement, fulfillment và durable Outbox.

## Runtime status

Canonical default Compose service. API + owned fulfillment/outbox actors chạy trong owner process mặc định.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc` -> `cmd/grpc/runtime.go`.

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

Root có thể gọi `make migrate service=payment-service`.

## Dependencies / durable state

PostgreSQL `payment_service`, RabbitMQ, User/Auth/Notification RPC. Canonical migration: `migrate/`.

## Read source from here

`cmd/grpc/runtime.go`, `internal/`, `infra/`, `config/`, `migrate/`. Theo flow settlement/fulfillment/outbox thay vì học layer trước.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
