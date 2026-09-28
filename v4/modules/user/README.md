# modules/user

User profile ownership for the v4 modular layout.

```
controller/  ProfileController + Handler (gRPC + HTTP)
services/    ProfileService (+ legacy ProfileUsecase adapter during migration)
dto/         (reserved for request shapes as more routes move in)
```

Profile info me currently routes:

`GrpcProfileService` → `ProfileController` → `ProfileService` → legacy `ProfileUsecase` (+ KYC / hub / sync)

HTTP mounts (cmd/http):

- `GET /v2/user/profile/info/me`

Next steps: move GetInfoMe body off ProfileUsecase into `services/`,
then expand accounts/me, set-info, and related profile routes.
