# TQD Service

TQD is the QHPRO Planning Platform service. Canonical business owners include Planning, Spatial, Discovery, Related Entity, Projection, Workspace, Generated Report, Usage and Quota.

## Durable authority

- Planning/Report/Usage/Report Job: PostgreSQL/PostGIS.
- Quota used/reserved runtime projection: Redis.
- Generated Report acceptance keeps `Report + UsageEvent + Job` in one PostgreSQL transaction.
- Redis loss must not erase durable Usage.

## Runtime ownership

The gRPC process loads typed runtime config once, opens one PostgreSQL pool, one Redis pool and its outbound gRPC connections, then injects them into Wire. Business packages do not open process resources. Generated Report reuses the same process DB/Redis/User Access resources and owns only its report actors/adapters.

## Local operation

```bash
# Makefile tự nạp tqd-service/.env, không dùng .env.local.
make doctor
make deps-up
make migrate-up
make generate
make build
make run-grpc
```

`QHPRO_REPORT_GENERATOR_MODE=smoke` kiểm luồng report/quota/File với renderer
xác định; `TQD_CLASSIFY_ENABLED=false` giữ AI classify ngoài minimal boot.

## Migrations

`migrate/` is the only production schema evolution authority. Historical SQL
is recovered from source-control history or the legacy repository when an
adoption investigation needs it; the service does not keep a duplicate
`migrate/` là schema owner duy nhất; không còn cây script migration song song.

## Developer gates

```bash
make check
make test-race
make acceptance
```

`make acceptance` là entrypoint kiểm tra code cho vet, unit test và build. Hạ tầng PostgreSQL/Redis do root `compose.yaml` sở hữu; dùng `make deps-up` từ service hoặc `make dev-dependencies SERVICE=tqd` từ repository root.
