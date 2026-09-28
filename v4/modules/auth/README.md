# modules/auth

Auth ownership for the v4 modular layout.

```
controller/  LoginController + PermissionController (gRPC + HTTP)
services/    LoginService / PermissionService (+ legacy adapters during migration)
dto/         Login and permission request/response types
entity/      AuthMethod sketch
repository/  AuthMethodRepository sketch
```

Login RPCs currently route:

`AuthHandler` → `LoginController` → `LoginService` → legacy `AuthUsecase`

Permission keys (`GET /v2/auth/permission/me`) currently route:

`PermissionHandler` → `PermissionController` → `PermissionService` → legacy `RoleUsecase`

HTTP mounts (cmd/http):

- login: `{/v2|/v2.1|/v3}/auth/...`
- permission me: `{/v2|/v2.1|/v3}/auth/permission/me`

Next steps: move VerifyOTP / LoginWithPassword bodies into `services/`,
then implement `repository/` against postgres without going through AuthUsecase.
Move permission-key resolution off `RoleUsecase` into this package.
