# modules/hub

All hub-service HTTP capabilities are now owned by this module on the v4 DSN.
hub-service is no longer required for these endpoints.

```
controller/   HTTP entry (gateway path parity)
services/     Local services for all hub capabilities
dto/          Request shapes for HTTP
entity/       GORM models for hub tables on the v4 DSN
repository/   Postgres stores on the v4 DSN
```

## Registered HTTP Routes

### Location (public)
- `GET /v2/hub/provinces`, `/districts`, `/wards`, `/location/search`
- `GET /v2/hub/v2/provinces`, `/province/{code}`, `/wards`, `/ward/{code}`, `/location/search`

### User Guide (public + admin)
- `GET /v2/hub/user-guides`, `/{id}`, `/by-key`, `/by-key/{key}`, `/groups`
- `GET /v2/hub/admin/user-guides`
- `POST /v2/hub/user-guides`, `PUT /{id}`, `DELETE /{id}`

### FAQ (public + admin)
- `GET /v2/hub/faqs`, `/{id}`
- `GET /v2/hub/admin/faqs`, `/{id}`
- `POST /v2/hub/admin/faqs`, `PUT /{id}`, `DELETE /{id}`

### App Version
- `GET /v2/hub/app-bundle/update-check/{appName}`
- `POST /v2/hub/app-bundle/version` (requires `api-key` / `x-api-key`)
- `GET /v2/hub/admin/versions`, `/{id}`
- `POST /v2/hub/admin/versions`, `PUT /{id}`, `DELETE /{id}`

### Request Stats
- `POST /v2/hub/request-stats/track`
- `GET /v2/hub/request-stats/dashboard`, `/daily`, `/users/dashboard`, `/users/daily`, `/audit`

### Error Log (public)
- `POST /v2/hub/error/log`

### Applinks
- `POST /v2/hub/applinks`
- `GET /v2/hub/applinks/{code}` (public)

### API Keys (admin)
- `POST /v2/hub/admin/api-keys`, `GET`, `DELETE /{id}`, `POST /verify`

### System Config (public key/group, admin bulk-upsert)
- `GET /v2/hub/system-config/by-group/{groupKey}` (public)
- `GET /v2/hub/system-config/key/{key}` (public)
- `GET /v2/hub/settings/{groupKey}` (public)
- `POST /v2/hub/system-config/bulk-upsert`
- `POST /v2/hub/system-config/reset-default/{groupKey}`

### Event Queue (admin)
- `GET /v2/hub/event-queue/list`, `/detail/{id}`
- `POST /v2/hub/event-queue`, `PUT /edit/{id}`, `DELETE /{id}`, `PUT /status/{id}`

### Interactive Events
- `POST /v2/hub/v2/interactive-event`

### Click URLs
- `POST /v2/hub/v2/click-url`

### Update Data Sync
- `GET /v2/hub/update/{resource}/sync/status`
- `GET /v2/hub/update/{resource}/sync/status/{id}`
- `GET /v2/hub/update/{resource}/refresh/ids`
- `POST /v2/hub/update/{resource}/flush/ids`

### Keyword Search Stats
- `GET /v2/hub/keyword-search-stats`

### Search Latency
- `GET /v2/hub/search-latency/dashboard`
- `GET /v2/hub/search-latency/daily`

## Schema
- Phase 1 tables: `v4/migrate/000030_hub_core.*.sql`
- Remaining tables: `v4/migrate/000036_hub_remaining.*.sql`
