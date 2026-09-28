# modules/crm

CRM contact ownership for the v4 modular layout.

```
controller/   HTTP entry (gateway path parity)
services/     Local contact list (+ SyncUtil, profile enrichment)
dto/          Request shapes for HTTP
entity/       GORM models for CRM-owned tables
repository/   Postgres stores on database.crm.dsn
```

Phase-1 routes (local on CRM DSN, no crm-service TCP for these reads):

- `GET /v2/crm/contact/list/user/{ownerId}` — contacts of the authenticated user
  (`ownerId=0` → JWT profile; size capped at 20; Redis `SyncKeyContactMe`)

Schema: CRM-owned `tb_contact` / `friend` / `tb_contact_product` / `origin_profile`
on `database.crm.dsn` (crm-service `db_crm`). Distinct from the thinner
`tb_contact` stub on the v4 user DSN.

Create/update/delete, org/group owner lists, leads, and pipelines remain out of
scope for phase-1.
