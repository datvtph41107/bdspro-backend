# gateway-service

## Role

HTTP ingress của BDSPro: JWT verification, route auth, trusted metadata, HTTP/gRPC gateway/proxy, error/deadline/response mapping.

## Runtime status

Canonical default Compose service.

## Entrypoint / lifecycle

`main.go` -> Cobra `http` -> `cmd/http`.

`make server` chạy `go run . http`.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

## Dependencies / durable state

Các upstream gRPC service; không sở hữu database migration.

## Read source from here

`cmd/http/`, `internal/httpauth/`, `internal/grpcgateway/`, `internal/upstream/`, `internal/v1proxy/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
