# tqd-service

## Role

QHPro/planning capability, quota, durable usage và report acceptance/job/rendering.

## Runtime status

Canonical default Compose service. Report actor thuộc owner process; PMTiles là special mode có chủ đích.

## Entrypoint / lifecycle

`main.go`; mode mặc định `-server=grpc`, optional `-server=pmtiles`.

`make server`; `make pmtiles` cho tile server.

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

Root có thể gọi `make migrate service=tqd-service`.

## Dependencies / durable state

PostgreSQL `qhpro_tqd`, Redis DB 1, User/Auth/Assistant/File. Canonical migration: `migrate/`.

## Read source from here

`cmd/`, `internal/`, `infra/`, `config/`, `migrate/`, `template/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
