# map-service

## Role

Map/location HTTP source module.

## Runtime status

Repository source-health module, không nằm default Compose.

## Entrypoint / lifecycle

`main.go` -> Wire runtime/app -> Gin/router HTTP `:8101`.

`make server`.

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

Root có thể gọi `make migrate service=map-service`.

## Dependencies / durable state

Xem `config/runtime.yml` và Wire providers; không có canonical migration authority trong current source.

## Read source from here

`main.go`, `wire/`, `internal/`, `infra/`, `config/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
