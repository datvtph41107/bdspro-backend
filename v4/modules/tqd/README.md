# modules/tqd

TQD ownership for the v4 modular layout (ported from `tqd-service`).

```
controller/     HTTP entry (gateway path parity for /v2/tqd/client + /v2/tqd/public)
services/       Usecase layer (from tqd-service/internal/usecase)
dto/            Request/response shapes
entity/         Domain models (from tqd-service/internal/domain)
repository/     Store ports
postgres/       Store implementations (from tqd-service/infra/postgres)
enums/          Domain catalogs
mapper/         Domain ↔ DTO/protobuf mappers
provider/       Outbound provider ports
access/         Access-subject helpers
geometry/       Geometry helpers
config/         Limits + layer resolver helpers
utils/          Small shared helpers (adm join, document ext)
```

## Migration status

| Source (`tqd-service`) | Target (`v4/modules/tqd`) | Status |
|------------------------|---------------------------|--------|
| `internal/domain/**` | `entity/**` | Migrated |
| `internal/usecase/**` | `services/**` | Migrated |
| `internal/enums/**` | `enums/**` | Migrated |
| `internal/dto/**` | `dto/**` | Migrated |
| `internal/interface/repo/**` | `repository/**` | Ports migrated |
| `internal/interface/provider/**` | `provider/**` | Ports migrated |
| `infra/mapper/**` | `mapper/**` | Migrated |
| Client-needed `infra/postgres/**` | `postgres/**` | Migrated + wired |
| Admin / discovery / tile workers HTTP | — | Not mounted on v4 HTTP yet |

## `/v2/tqd/client` + public (wired on `:8288`)

All proto client paths are registered via `controller/*` + `providers.go` / wire. Guest JWT allowlist is in `cmd/http/auth_middleware.go`.

| Group | Status |
|-------|--------|
| Location catalog (`provinces`, `wards`) | wired |
| Administrative units (+ followed units) | wired |
| Layers / labels / land-use / legends / build-tile | wired |
| Layer families | wired |
| Regions (`list`, `by-point`, `by-bbox`) | wired |
| Parcels (`search`) | wired |
| Planning client (projects/docs/events/follow/…) | wired |
| `nearby-areas`, `planning-documents/by-adm` | 501 (same Unimplemented stubs as tqd-service) |
| `GET /v2/tqd/public/enums` | wired |

`planning-documents/{relation-layer\|view\|{id}}/{…}` share one ServeMux pattern to avoid Go 1.22 wildcard conflicts; public URLs stay unchanged.

## Package mapping notes

- Root domain `package domain` → `package entity`
- `qh_domain` → `entity/qh` (`package qh`)
- Nested discovery/relatedentity/… model packages → `entity/<capability>`
- Usecase nested trees keep the same layout under `services/`
- Import prefix: `userv4/modules/tqd/...`

Schema: TQD-owned tables on `database.tqd.dsn`. Canonical SQL lives in `tqd-service/migrate/`.
`NewTqdDB` runs pending golang-migrate files on process start (`TQD_MIGRATION_DIR` override supported).
Catalog sync key: `tqd:loc:provinces-wards`.

Integration tests that still need `tqd-service` infra workers are tagged `//go:build tqd_service_infra`.
