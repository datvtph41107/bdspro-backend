# auth-service

## Role

`AuthInternal` compatibility/permission snapshot cho service-to-service authorization; không nhận business capability mới.

## Runtime status

Canonical default Compose compatibility service.

## Entrypoint / lifecycle

`main.go` trực tiếp load config, mở User RPC, đăng ký `AuthInternalService`, chạy actor permission service và graceful shutdown.

`make server`.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

## Dependencies / durable state

User RPC; không sở hữu DB migration.

## Read source from here

`main.go`, `infra/handler/`, `infra/services/`, `infra/rpc/`, `config/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
