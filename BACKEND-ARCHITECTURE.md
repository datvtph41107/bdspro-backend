# BDSPro Backend Architecture

Tài liệu này mô tả source truth đang chạy. `bdspro-golang-microservices-r2f0-i21`
là source được phát triển; `used/bdspro-golang-microservices` chỉ là bằng chứng
production để kiểm tra compatibility.

## 1. System và capability

- BDSPro là system/product boundary.
- QHPro là capability quy hoạch do TQD Service sở hữu.
- React Admin và React Native không sở hữu giá, entitlement, quota hoặc usage.
- Gateway là HTTP edge; business authority luôn nằm trong service owner.

## 2. Canonical owner map

| Owner | Trách nhiệm |
| --- | --- |
| Gateway | HTTP ingress, JWT verification, route auth, metadata, error/response mapping |
| User | Profile, public Auth/OAuth, IAM, catalog, checkout intent, subscription, entitlement |
| Auth | `AuthInternal` compatibility và permission snapshot; không nhận business capability mới |
| Organization | Organization, membership và quyền hành động trong organization |
| Payment | Order, attempt, provider evidence, settlement, fulfillment và Outbox |
| TQD | QHPro, planning, quota, durable usage, report acceptance/job/rendering |
| File | File metadata/content và signed access |
| Notification | Inbox dedupe, logical notification, delivery intent/effect |
| Hub | Location/system configuration/API key/version capabilities hiện hữu |
| Assistant | AI provider boundary dùng bởi TQD |

## 3. Runtime flow đã được chứng minh

```text
React Native / Admin
        |
        v
     Gateway
        |
        +--> User --------> Organization (organization action authority)
        |      |
        |      +--> Payment --> provider webhook
        |                        |
        |                        +--> User settlement/Subscription
        |                        +--> Payment Outbox
        |                                  |
        |                                  v
        |                              RabbitMQ
        |                                  |
        |                                  v
        |                        Notification Inbox/effect
        |
        +--> TQD --> User Entitlement
               |       |
               |       +--> subscription authority
               +--> Redis admission + PostgreSQL durable quota/usage
               +--> Report + UsageEvent + Job transaction
               +--> report actor --> File --> signed PDF
```

Payment Outbox publisher và fulfillment worker là component trong lifecycle
của Payment process. Consumer của Payment events và delivery worker là
component trong lifecycle của Notification process. TQD report actor cũng
thuộc TQD process. Chúng có thể còn binary compatibility trong source, nhưng
Compose mặc định không tạo standing container riêng.

## 4. Data rules

- Money dùng `int64` minor units.
- Provider webhook là boundary riêng, không dùng JWT Admin/User.
- Checkout/order/report command có idempotency identity và conflict semantics.
- Payment settlement, Outbox và fulfillment lưu durable state trước retry.
- Notification Inbox dùng event ID làm dedupe identity và chỉ ACK sau DB commit.
- Report acceptance ghi Report + UsageEvent + Job trong transaction.
- PostgreSQL là durable usage truth; Redis chỉ là admission/runtime projection.
- Migration versioned tại `<service>/migrate`; application schema mode phải tắt
  ở môi trường deploy.

## 5. Process rules

- Một service có một default process/container.
- Composition root mở DB/Redis/RPC/broker và sở hữu close order.
- Component chạy cùng cancellation domain với server.
- Log ra stdout/stderr; Docker/runtime platform sở hữu retention.
- Local và production dùng cùng binary, Dockerfile và config schema. Khác biệt
  chỉ đến từ runtime values/secrets/topology do operator inject.
- Development dùng một PostgreSQL/PostGIS, một Redis và một RabbitMQ. Logical
  database, Redis DB index và migration authority vẫn thuộc từng service.

## 6. Production comparison

Production Auth đăng ký Auth/OAuth/AuthInternal/Permission/Role/RoleGroup/
UserInfo trong `auth-service/cmd/grpc.go`. Source mới đăng ký public Auth/OAuth/
IAM tại `user-service/cmd/grpc/main.go` và chỉ đăng ký AuthInternal tại
`auth-service/main.go`. Gateway trỏ Auth public tới `user:8201`. Customer và
Admin runtime acceptance đã chứng minh cutover này hoạt động.

Appointment có canonical serving owner tại CRM
(`crm-service/cmd/grpc/main.go`). Các deployable production bị thiếu khỏi source
mới gồm `appointment-service`, `marketing-service`, `membership-service`,
`task-service`, `transaction-service`; chỉ Appointment đã có owner replacement
được chứng minh. Các service còn lại chưa được phép coi là retired nếu chưa có
endpoint/data/deployment reachability proof.

Chi tiết workflow, topology và deployment gap nằm tại
[`documents/OPERATIONAL-MAPS.md`](documents/OPERATIONAL-MAPS.md).
Ma trận parity và blocker có source evidence nằm tại
[`documents/PRODUCTION-COMPATIBILITY.md`](documents/PRODUCTION-COMPATIBILITY.md).
