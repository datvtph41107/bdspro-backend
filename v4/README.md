# v4 (mdn4)

Sandbox clone of `user-service` for side-by-side experiments.

## Module

Go module: `mdn4` (folder `v4/` at repo root). Do not conflate with owner `user-service`.

## Entrypoints

| Entrypoint | Default listen                               | Notes                                 |
| ---------- | -------------------------------------------- | ------------------------------------- |
| `cmd/grpc` | `:8211` (`GRPC_ADDRESS` / `server.tcp_port`) | Legacy gRPC (same as user-service)    |
| `cmd/http` | `:8288` (`HTTP_ADDRESS` / `server.port`)     | **Real HTTP (Gin)** — modular surface |

## HTTP routes (`cmd/http`)

Owned by `modules/auth` (`LoginController.RegisterHTTP`). Same handlers under **`/v2`**, **`/v2.1`**, **`/v3`** (gateway-compatible):

```
GET  /health
GET  {prefix}/auth/phone/check/:phone
POST {prefix}/auth/otp/request
POST {prefix}/auth/otp/verify
POST {prefix}/auth/otp/resend
POST {prefix}/auth/password/login
POST {prefix}/auth/login/key
POST {prefix}/auth/admin/login
POST {prefix}/auth/token/refresh
```

## Developer commands

```bash
make build
make server          # cmd/grpc → :8211
make http            # cmd/http → :8288
make test
```

## Config layering

`LoadProperties` loads:

1. `config/runtime.yml` — base (CORS `http.allowed_origins`, ports, defaults)
2. `config/<ENV>.yml` — environment overlay (`local` / `develop` / …) for DSN/secrets
3. optional `QHPRO_CONFIG_FILE` / `CONFIG_FILE` — final merge

HTTP CORS for `cmd/http` comes from `runtime.HTTP.AllowedOrigins`.
