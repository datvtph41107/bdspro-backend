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
- `GET /v4/deals/me` — deals of current user (was `/v2.1/org/deals/user`)
- `GET /v4/deal/{id}`, `{id}/ts`, `{id}/products`, `{id}/summary`, `{id}/bank-account`, `{id}/member` — deal detail / timestamps / products / summary / bank / members (was `/v2.1/org/deal/...`)
- `GET /v4/deal/invitation/list/{dealId}` — deal invitations (was `/v2.1/org/deal-invitation/list/{id}`)
- `GET|POST|PUT|DELETE /v2.1/org/investment...` and `/v4/org/investment...` — investment CRUD, approve/reject/confirm, member histories
- `GET /v2.1/org/investments/list/{dealId}` and `/v4/deal/investments/list/{dealId}` — investment list
- `GET /v2.1/transaction/contract/deal/list/{dealId}` (+ `/v4/.../{id}/list` legacy) — contracts list
- `GET /v2.1/transaction/contract/deal/statistics/{dealId}` (+ `/v4/.../{id}/statistic`) — cost/payment statistics
- `GET /v2.1/transaction/contract/deal/{id}/process` — contract process timeline
- `POST /v2.1/transaction/contract/deal/payment`, `POST .../payment/action`, `GET .../detail/{id}`, `DELETE .../{id}`, `GET .../me` — contract deal writes/reads (approve + financial summary not implemented)
- `GET|POST|PUT|DELETE /v2.1/transaction/deal-cost/...` and `/v4/transaction/deal-cost/...` — deal cost CRUD + list
- `GET|POST|PUT|DELETE /v2.1/transaction/cost-type` and `/v4/transaction/cost-type` — cost types
- `GET /v2/auth/permission/deal/{id}` — stays on auth (not migrated here)
- `GET /v2/bdspro/v2/product/me`, detail, new, delete, `:id/timestamp`, `:id/deals`
- `GET /v2/bdspro/product-notes/all/:productId`
- `GET /v4/appointment/me` — appointments of current user (was `/v2.1/appointment`)
- `GET /v4/appointment/{id}` — appointment detail (was `/v2/appointment/{id}` / `/v2.1/appointment/{id}`)
- `GET /v2.1/appointment/list/by-product/:id`
- `GET /v2/notification/property-history/:id`
- `GET /v2/bdspro/v2/property/me`, detail, `:id/timestamp`
- `GET /v2/bdspro/public/enums`, `/enums/{text}`

Schema: `000029_bdspro_core` + `000031_bdspro_product_tab` + `000033_bdspro_deal_contract` + `000034_bdspro_deal_lists` + `000035_bdspro_bank_accounts` on the v4 Postgres DSN.
