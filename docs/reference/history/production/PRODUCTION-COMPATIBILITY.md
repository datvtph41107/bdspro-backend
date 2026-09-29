# BDSPro Production Compatibility Matrix

Ngày audit: 2026-09-03. Production evidence chỉ được đọc từ
`../used/bdspro-golang-microservices`; source phát triển duy nhất là repository
hiện tại. Tài liệu này không xác nhận deployment production và không thay thế
runtime proof.

## Kết luận hiện tại

```text
QHPro commercial/quota core: runtime proven
Legacy production parity:    not proven
Whole backend:                NOT CLOSED
```

Source production có 23 thư mục `*-service`; source hiện tại có 19, trong đó
`chat-v1-service` là deployable mới. Việc tên service biến mất hoặc code compile
không phải bằng chứng retirement.

## Ma trận deployable

| Production deployable | Owner/source hiện tại | Evidence hiện tại | Trạng thái |
| --- | --- | --- | --- |
| `gateway-service` | Gateway | Canonical Compose + smoke + core E2E | runtime-proven cho core |
| `user-service` | User | Auth public, IAM, catalog, checkout, subscription, entitlement; core E2E | runtime-proven cho core |
| `auth-service` | User + Auth compatibility | Public Auth/OAuth/IAM ở User; Auth chỉ đăng ký `AuthInternal` | core cutover proven; caller chưa zero |
| `organization-service` | Organization | Checkout dùng authority; vẫn mở Transaction client `:8215` | core path proven; toàn service chưa proven |
| `payment-service` | Payment | Order/attempt/webhook/settlement/outbox/fulfillment E2E + recovery | runtime-proven cho core |
| `tqd-service` | TQD | Quota/usage/report/file E2E + Redis/process recovery | runtime-proven cho core |
| `notification-service` | Notification | Payment-event Inbox/effect E2E + broker/process recovery | runtime-proven cho core |
| `file-service` | File | PDF, response-loss convergence và signed read | runtime-proven cho core |
| `hub-service` | Hub | Compose healthy; API key/file dependencies exercised | runtime-proven cho core only |
| `assistant-service` | Assistant | Compose healthy; full capability E2E chưa có | runtime-present, partial proof |
| `bdspro-service` | BDSPro | test/vet/race/build; không trong standing Compose | source-healthy only |
| `crm-service` | CRM | test/vet/race/build; không trong standing Compose | source-healthy only |
| `chat-service` | Chat | test/vet/race/build; Gateway trỏ logical `chat` | source-healthy only |
| `map-service` | Map | test/vet/race/build; không trong standing Compose | source-healthy only |
| `relay-service` | Relay | test/vet/race/build; không trong standing Compose | source-healthy only |
| `search-service` | Search | test/vet/race/build; không trong standing Compose | source-healthy only |
| `social-service` | Social | test/vet/race/build; Gateway còn routes | source-healthy only |
| `ai-service` | chưa chốt | Source còn nhưng không thuộc root module/Compose acceptance | unclassified |
| `appointment-service` | CRM candidate | CRM đăng ký `AppointmentService`; chưa chứng minh data/deploy cutover | replacement candidate |
| `membership-service` | User candidate | Commercial owner ở User; BDSPro/Organization còn legacy plan code | unresolved/blocking |
| `marketing-service` | CRM candidate | CRM còn handler/use cases nhưng không register gRPC | unresolved/blocking |
| `task-service` | chưa có | Production chỉ có stub service + ticker; current không có owner | retirement candidate, proof thiếu |
| `transaction-service` | BDSPro candidate | BDSPro có APIs; Organization vẫn gọi logical `transaction:8215` | unresolved/blocking |
| `chat-v1-service` (mới) | chưa chốt | Trùng Chat/BackgroundImage contract; không có canonical Compose owner | deployment ambiguity |

Production `feedback` selector không tương ứng một thư mục production riêng;
capability feedback hiện nằm trong CRM. Đây là selector compatibility cần audit,
không phải service thứ 24.

## Blocker cụ thể

### Membership

- `bdspro-service/infra/handler/member_plan_handler.go` tải plan từ URL/IP
  production viết cứng khi constructor chạy.
- `bdspro-service/infra/providers/member_plan_provider.go` còn một URL/IP legacy
  khác.
- `organization-service/infrastructure/client/membership_client.go` trả plan
  `Free` và branch limit `10` viết cứng.

User Catalog là owner của commercial/QHPro mới nhưng chưa thay thế toàn bộ
semantics membership cũ.

### Transaction

- Organization vẫn bind `TransactionGrpcClient` và yêu cầu
  `ORGANIZATION_TRANSACTION_GRPC_ADDRESS`.
- Canonical Compose gán `transaction:8215` nhưng không có container đó.

gRPC dial lazy cho phép process healthy, nhưng code path dùng transaction có thể
fail. Health không phải dependency parity.

### Marketing, Chat và AI

CRM có Marketing code nhưng không có serving registration. Hai Chat
implementation cùng contract chưa có traffic decision. AI chưa được phân loại
deployable/asset/dead. Không được xóa trước reachability/deployment proof.

## Deployment compatibility

`shared/code/deploy.sh` giữ đúng operator intent `deploy <one service>`, nhưng còn
hard-coded SSH host, copy development config/secret, thiếu immutable version,
readiness/smoke và rollback. Production source cũng không có root Compose được
script tham chiếu, nên topology thật trên host là external evidence chưa có.

Không tự động chạm host production trong workstream này. Modernization cần
production Compose/config inventory, selector mapping và quyền operator riêng.

## Delete gate

Một legacy owner chỉ được retire khi endpoint caller, internal RPC caller,
table/data owner, event, job, Admin/Mobile client, deploy selector/runbook và
production traffic dependency đều bằng zero.

Ưu tiên tiếp: Membership → Transaction → Marketing → Task → Chat/Chat-v1 → AI;
sau đó runtime-proof BDSPro/CRM/Map/Relay/Search/Social và chốt whole backend.
