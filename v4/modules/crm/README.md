# modules/crm

CRM ownership for the v4 modular layout (ported from `crm-service`).

```
controller/     HTTP entry (gateway path parity for /v2/crm/*)
services/       Usecase layer (from crm-service/internal/usecase)
dto/            Request/response shapes
entity/         Domain models (from crm-service/internal/domain)
repository/     Store ports
postgres/       Store implementations (from crm-service/infra/postgre)
enums/          Domain catalogs
mapper/         Domain ↔ DTO/protobuf mappers
provider/       Outbound provider ports
```

## Wired `/v2/crm` surface (core CRM)

| Group | Status |
|-------|--------|
| `admin/sales` (incl. `contacts`) | wired |
| `contact` | wired |
| `friend` / `friend-group` / `follow` / `block` | wired |
| `tag` / `pipeline` / `stage` / `lead` / `rule` | wired |
| `sharing` / `invitation` | wired |
| `admin/support` | wired |

SEO / public content stays under `modules/seo` (`/v2/crm/admin/seo*`, `/v2/crm/public/*` content).

Schema: CRM-owned tables on `database.crm.dsn` (`db_crm` in prod).
