# modules/bdspro

Product, property, and deal ownership for the v4 modular layout.

```
controller/   HTTP entry (gateway path parity)
services/     Local Sync / Product / Property / Deal / Public (v4 DSN + SyncUtil)
dto/          Request/response shapes for HTTP
entity/       GORM models for phase-1 tables
repository/   Postgres stores on the v4 DSN
enums/        In-process enum catalogs (public/enums)
```

Phase-1 routes (local, no bdspro/crm/notification TCP for these reads):

- `GET /v2/bdspro/ts` — app sync timestamps (Redis / SyncUtil)
- `GET /v2.1/org/deals/user` — deals of current user
- `GET /v2/bdspro/v2/product/me`, detail, new, delete, `:id/timestamp`, `:id/deals`
- `GET /v2/bdspro/product-notes/all/:productId`
- `GET /v2.1/appointment/list/by-product/:id`
- `GET /v2/notification/property-history/:id`
- `GET /v2/bdspro/v2/property/me`, detail
- `GET /v2/bdspro/public/enums`, `/enums/{text}`

Schema: `000029_bdspro_core` + `000031_bdspro_product_tab` on the v4 Postgres DSN.
