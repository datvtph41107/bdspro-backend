# crm-service

## Role

CRM + public content/SEO source module; có PaymentCompleted consumer compatibility.

## Runtime status

Repository source-health module, không có standing service trong canonical default Compose. Không suy diễn retirement/cutover completion.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc`; compatibility consumer tại `cmd/payment-completed-consumer/`.

`make server`; `make payment-completed-consumer` khi cần actor riêng.

## Developer commands

```bash
make help
make test
make build
make server
```

Nếu thay schema:

```bash
make migrateup
make migrate-version
make new_migration name=<schema_change>
```

Root có thể gọi `make migrate service=crm-service`.

## Dependencies / durable state

Local DB fallback `db_crm`; canonical migration authority `infra/db/migrate_v2/`.

## Read source from here

`cmd/grpc/`, `internal/`, `infra/`, `config/`, `infra/db/migrate_v2/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
