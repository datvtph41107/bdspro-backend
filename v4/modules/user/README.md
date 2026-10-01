# modules/user

User profile ownership for the v4 modular layout. user-service is no longer required
for these HTTP endpoints once traffic points at v4 `cmd/http` (or v4 gRPC for dual-serve).

```
controller/  Profile + Catalog + Tag + KYC + Bookmark + AdminUser
services/    Local adapters over existing usecases + tile session (Redis)
dto/         Request shapes for HTTP
```

## Registered HTTP Routes

### Profile
- `GET  /v2|/v2.1|/v3/user/profile/info/me`
- `GET  /v2|/v2.1|/v3/user/profile/info/{id}`
- `POST /v2|/v2.1|/v3/user/profile/session` — tile map sessionK + AES key (Redis, public)
- `GET  /v2|/v2.1|/v3/user/profile/accounts/me`
- `PUT  /v2|/v2.1|/v3/user/profile/set-info`
- `PUT  /v2|/v2.1|/v3/user/profile/data`
- `PUT  /v2|/v2.1|/v3/user/profile/privacy`
- `POST /v2|/v2.1|/v3/user/profile/quick-setup`
- `GET  /v2|/v2.1|/v3/user/profile/search/public` (public)
- `GET  /v2|/v2.1|/v3/user/profile/find-by-phone` (public)
- `GET  /v2|/v2.1|/v3/user/profile/completion`
- `GET  /v2|/v2.1|/v3/user/status-online/{id}`
- `POST /v2|/v2.1|/v3/user/demo-mode`
- `GET  /v3/user/profile/{access}/{id}`

### Catalogs
- `GET/POST /v2/user/admin/main-area`, `PUT/DELETE /{id}`, `GET /v2/user/main-area`
- `GET/POST /v2/user/admin/purpose-use`, `PUT/DELETE /{id}`, `GET /v2/user/purpose-use`
- `GET/POST /v2/user/admin/profession`, `PUT/DELETE /{id}`, `GET /v2/user/professions`

### Tags
- `POST/PUT/DELETE/GET /v2/user/admin/tags`
- `GET /v2/user/profile/tags/{tagType}`
- `GET /v2/pub/user/profile/tags/{tagType}` (public)

### KYC
- `POST /v2/user/kyc/submit`
- `GET  /v2/user/admin/kycs`
- `POST /v2/user/admin/kyc/approve/{id}`
- `POST /v2/user/admin/kyc/reject/{id}`

### Bookmarks
- `GET/PUT /v2/user/bookmark-users`

### Admin users / staff profiles
- `GET/POST/PUT/DELETE /v2/user/admin/users...` (+ lock/unlock/approve/stats/dashboard)
- `POST/PUT/DELETE/GET /v2/user/admin-profile/...`

## Not in this module (yet)
- Commercial / checkout / catalog plan-versions / pack price-tables — still `infra/handler/grpc`
- `/v2/auth/*` — `modules/auth`
- Organization admin — org surface

## Cutover
1. Prove parity on v4 HTTP `:8288`
2. Point gateway `user` gRPC dial at v4 `cmd/grpc` (not HTTP `:8288`)
3. Stop `user-service` container
