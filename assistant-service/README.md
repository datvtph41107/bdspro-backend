# assistant-service

## Role

AI provider boundary dùng bởi TQD/other internal callers; provider mode được chọn từ runtime config/ENV.

## Runtime status

Canonical default Compose Go service.

## Entrypoint / lifecycle

`main.go` -> `cmd/grpc.RunGRPCServer()`.

`make server`.

## Developer commands

```bash
make help
make test
make build
make server
```

## Dependencies / durable state

BDSPro/Hub RPC và external AI providers khi provider mode không phải stub; không có canonical DB migration.

## Read source from here

`cmd/grpc/`, `infra/client/`, `internal/usecases/`, `config/`, `wire/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
