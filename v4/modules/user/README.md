# modules/user

User profile ownership for the v4 modular layout.

```
controller/  ProfileController + Handler (gRPC + HTTP)
services/    ProfileService + DemoService (local seed/cleanup)
dto/         Request shapes for HTTP
```

Profile reads currently route:

`GrpcProfileService` → `ProfileController` → `ProfileService` / `DemoService`

HTTP mounts (cmd/http):

- `GET  /v2/user/profile/info/me`
- `GET  /v3/user/profile/{access}/{id}` — public profile (v4 `database.dsn`, not legacy `db_user`)
- `POST /v2/user/demo-mode` — seed/cleanup contacts, products, deals (+ bank/investment/cost/contract), appointments

Demo snapshots: `v4/data/demo/` (`contacts.json`, `products.json`, `medias.json`).
