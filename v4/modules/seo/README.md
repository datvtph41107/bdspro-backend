# modules/seo

CRM SEO capabilities owned by this module on the v4 DSN.
crm-service is no longer required for these endpoints.

```
controller/   HTTP entry (CRM path parity)
services/     Local news + seo-domain lifecycle services
provider/     Local CrmSeoProvider for modules/tqd (replaces noop/RPC)
dto/          Request/response shapes matching CRM public/admin JSON
entity/       GORM models for SEO tables on the v4 DSN
repository/   Postgres stores on the v4 DSN
```

## Schema

Migration: `v4/migrate/000037_seo_core.up.sql` / `.down.sql`

Tables: `seo_domain`, `seo_links`, `seo_relative`, `seo_generation_log`,
`seo_count`, `seo_count_log`, `seo_category` (+seed), `user_saved_planning_news`.

Also registered in `db/migrate.go` AutoMigrate.

## Registered HTTP Routes

### Public news / content
- `GET /v2/crm/public/seo/news`, `/{id}`
- `POST /v2/crm/public/seo/click`
- `GET /v2/crm/public/content/categories`
- `GET /v2/crm/public/content/search`
- `GET /v2/crm/public/seo/documents/{module}/{slug}`
- `GET /v2/crm/public/content/page/{path...}` (HTML)
- `GET /v2/crm/public/content/sitemap.xml`, `/v2/crm/sitemap.xml`, `/v2/seo/sitemap.xml`
- `GET /v2/crm/public/seos/{refType}/{refId}/list`
- `GET /v2/crm/internal/seos/{refType}/{refId}/count`

### Public SEO stats
- `GET /v2/crm/public/seo-domains/top-clicks`
- `GET /v2/crm/public/seo-domains/click-daily`
- `GET /v2/crm/public/seo-domains/count`

### Client (authenticated)
- `GET /v2/crm/client/content/planning-news/saved`
- `POST|DELETE /v2/crm/client/content/planning-news/{id}/saved`

### Admin categories
- `GET|POST /v2/crm/admin/content/categories`
- `PUT|DELETE /v2/crm/admin/content/categories/{id}`
- `PUT /v2/crm/admin/content/planning-news/{id}/category`

### Admin SEO domains
- `GET|POST /v2/crm/admin/seo-domains`
- `POST /v2/crm/admin/seo-pages/drafts`
- `GET|PUT|DELETE /v2/crm/admin/seo-domains/{id}`
- `PUT /v2/crm/admin/seo-domains/batch`
- `POST /v2/crm/admin/seo-domains/{id}/publish|unpublish|archive|restore`
- `POST /v2/crm/admin/seo-pages/{id}/publish|unpublish|archive|restore|preview|generate`
- `POST /v2/crm/admin/seo-pages/preview`
- `POST /v2/crm/admin/seo-domains/touch-by-ref`
- `POST /v2/crm/admin/seo-domains/generate-from-ref`
- `POST /v2/crm/admin/seo-domains/parcel`
- `POST /v2/crm/admin/seo-domains/parcels|layers|planning-projects|planning-documents|administrative-units/generate`
- `GET /v2/crm/admin/seo/ref-types`, `/sources`
- Links / relatives under `/v2/crm/admin/seo-domains/...`

## TQD integration

`provider.NewLocalCrmSeoProvider` implements `modules/tqd/provider.CrmSeoProvider`
and is wired in place of the previous noop provider.
