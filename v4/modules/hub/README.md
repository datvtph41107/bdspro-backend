# modules/hub

Location, user guide, FAQ, app-version, and request-stats ownership for the v4 modular layout.

```
controller/   HTTP entry (gateway path parity)
services/     Local FAQ / Location / LocationV2 / UserGuide / Version / RequestStats
dto/          Request shapes for HTTP
entity/       GORM models for hub tables on the v4 DSN
repository/   Postgres stores on the v4 DSN
```

Phase-1 routes (local, no hub-service TCP):

- `GET /v2/hub/provinces`, `/districts`, `/wards`, `/location/search`
- `GET /v2/hub/v2/provinces`, `/province/{code}`, `/wards`, `/ward/{code}`, `/location/search`
- `GET /v2/hub/user-guides`, `/user-guides/{id}`, `/user-guides/by-key`, `/user-guides/by-key/{key}`, `/user-guides/groups`
- `GET /v2/hub/faqs`, `/faqs/{id}`
- `GET /v2/hub/app-bundle/update-check/{appName}`
- `POST /v2/hub/app-bundle/version` (requires `api-key` / `x-api-key`)
- `POST /v2/hub/request-stats/track`
- `GET /v2/hub/request-stats/dashboard`, `/daily`, `/users/dashboard`, `/users/daily`, `/audit`

Schema: `v4/migrate/000030_hub_core.*.sql` on the same v4 Postgres DSN.

In-process lookups (`HubProvider` / profile enrichment) also read these tables locally.

Admin CRUD, event-queue, system-config, applinks, and error-log remain out of scope.
