# BDSPro Service Map

## Canonical default local runtime

Runtime identity is stable protocol/process identity. Source path is repository topology.
They are intentionally separate so source can move without renaming service assertions,
logs, images, ports, or wire contracts.

| Service identity | Source path | Owner | Entrypoint | Default mode | DB migration |
| --- | --- | --- | --- | --- | --- |
| `gateway-service` | `api/gateway/` | HTTP ingress/JWT/routing/error mapping | `main.go` | `http` | none |
| `user-service` | `internal/user/` | profile, Auth/OAuth/IAM, catalog, subscription, entitlement | `cmd/grpc/main.go` | gRPC | `database/migrations/` |
| `auth-service` | `internal/authentication/` | `AuthInternal` compatibility/permission snapshot | `main.go` | gRPC | none |
| `organization-service` | `internal/organization-legacy/` | legacy organization/membership/action-right compatibility | `main.go` | `grpc` | `migrate/` |
| `payment-service` | `internal/payment/` | order/attempt/settlement/fulfillment/outbox | `main.go` | `grpc` | `database/migrations/` |
| `tqd-service` | `internal/planning/` | QHPro planning/quota/usage/report | `main.go` | `-server=grpc` | `database/migrations/` |
| `notification-service` | `internal/notification/` | inbox/logical notification/delivery | `main.go` | `grpc` | `database/migrations/` |
| `file-service` | `internal/file/` | file content/metadata/signed access | `main.go` | HTTP | `database/migrations/` |
| `hub-service` | `internal/hub/` | location/system/config/API-key capabilities | `main.go` | `grpc` | `database/migrations/` |
| `assistant-service` | `internal/assistant/` | AI provider boundary used by planning | `main.go` | gRPC | none |

## Repository modules outside default Compose

Các module dưới đây vẫn được source gate kiểm tra. Việc không có standing
container trong `compose.yaml` **không phải retirement proof**.

| Runtime/module identity | Source path | Current source role | Entry mode |
| --- | --- | --- | --- |
| `bdspro-service` | `internal/property/` | BDS/property implementation area | Cobra `grpc`/`http`; Make defaults `grpc` |
| `crm-service` | `internal/crm/` | CRM + public content/SEO; payment event consumer compatibility | Cobra `grpc` + separate consumer binary |
| `chat-service` | `internal/chat/` | current chat source module | Cobra `grpc`/`http` |
| `chat-v1-service` | `internal/chat-legacy/` | older chat compatibility source | Cobra |
| `map-service` | `internal/map-legacy/` | legacy map/location HTTP source | `main.go` HTTP |
| `relay-service` | `internal/realtime-relay/` | websocket relay source module | `main.go` |
| `search-service` | `internal/search/` | Elasticsearch-backed search source module | `main.go` HTTP |
| `social-service` | `internal/social/` | social source module | Cobra `grpc`/`http` |
| `ai-service` | `internal/ai/` | Python AI runtime with separate operational model | Python/Docker files owned inside module |

## Shared technical/runtime areas

| Source path | Role |
| --- | --- |
| `infrastructure/runtime/` | generic runtime mechanics and accepted compatibility debt; not a business authority |
| `infrastructure/base/` | reusable base mechanics |
| `proto/` | protobuf/wire contracts and generated code |
| `tools/development/` | repository engineering application |
| `tools/scripts/` | repository/support scripts |
| `docs/reference/` | canonical repository documentation |

## Reading rule

Nhận task -> tìm semantic owner/runtime identity -> resolve source path ->
README -> Makefile -> entrypoint -> config -> business flow -> durable store -> test.

Source path does not by itself prove business ownership, process boundary, or network boundary.
