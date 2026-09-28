# modules/tqd

TQD location / planning ownership for the v4 modular layout.

```
controller/   HTTP entry (gateway path parity)
services/     Local location catalog (+ SyncUtil)
dto/          Request shapes for HTTP
entity/       GORM models for TQD tables (province_v2, ward_v2, …)
repository/   Postgres stores on database.tqd.dsn
enums/        Location unit / data-source catalogs
```

Phase-1 routes (local on TQD DSN, no tqd-service TCP for these reads):

- `GET /v2/tqd/client/provinces` — guest map province catalog + catalog timestamp
- `GET /v2/tqd/client/wards` — national ward catalog (paginated) + shared catalog timestamp

Schema: TQD-owned `province_v2` / `ward_v2` on `database.tqd.dsn` (separate from hub’s
`province_v2` on the v4 DSN). Catalog sync key: `tqd:loc:provinces-wards`.

Admin CRUD, parcels, planning, layers, and PMTiles remain out of scope for phase-1.
