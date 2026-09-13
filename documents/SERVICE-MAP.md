# BDSPro Service Map

## Canonical default local runtime

| Service | Owner | Entrypoint | Default mode | DB migration |
| --- | --- | --- | --- | --- |
| `gateway-service` | HTTP ingress/JWT/routing/error mapping | `main.go` | `http` | none |
| `user-service` | profile, Auth/OAuth/IAM, catalog, subscription, entitlement | `cmd/grpc/main.go` | gRPC | `database/migrations/` |
| `auth-service` | `AuthInternal` compatibility/permission snapshot | `main.go` | gRPC | none |
| `organization-service` | organization/membership/action rights | `main.go` | `grpc` | `migrate/` |
| `payment-service` | order/attempt/settlement/fulfillment/outbox | `main.go` | `grpc` | `database/migrations/` |
| `tqd-service` | QHPro planning/quota/usage/report | `main.go` | `-server=grpc` | `database/migrations/` |
| `notification-service` | inbox/logical notification/delivery | `main.go` | `grpc` | `database/migrations/` |
| `file-service` | file content/metadata/signed access | `main.go` | HTTP | `database/migrations/` |
| `hub-service` | location/system/config/API-key capabilities | `main.go` | `grpc` | `database/migrations/` |
| `assistant-service` | AI provider boundary used by TQD | `main.go` | gRPC | none |

## Repository modules outside default Compose

Các module dưới đây vẫn được source gate kiểm tra. Việc không có standing
container trong `compose.yaml` **không phải retirement proof**.

| Module | Current source role | Entry mode |
| --- | --- | --- |
| `bdspro-service` | BDS/property source module | Cobra `grpc`/`http`; Make defaults `grpc` |
| `crm-service` | CRM + public content/SEO; payment event consumer compatibility | Cobra `grpc` + separate consumer binary |
| `chat-service` | chat source module | Cobra `grpc`/`http` |
| `chat-v1-service` | older chat source module | Cobra |
| `map-service` | map/location HTTP source module | `main.go` HTTP |
| `relay-service` | websocket relay source module | `main.go` |
| `search-service` | Elasticsearch-backed search source module | `main.go` HTTP |
| `social-service` | social source module | Cobra `grpc`/`http` |
| `ai-service` | Python AI runtime with separate operational model | Python/Docker files owned inside module |

## Reading rule

Nhận task -> tìm owner ở bảng -> mở README service -> Makefile -> entrypoint ->
config -> business flow -> durable store -> test.
