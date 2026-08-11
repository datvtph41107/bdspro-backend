# QHPRO target source map

Baseline backend: `bdspro-golang-microservices-20260811-091814`.

Giai đoạn này viết source target để đọc và kiểm logic trước. Chưa chạy DB migration/cutover.

## Đường request chính

```text
Frontend HTTP
  ↓
gateway-service
  ↓ gRPC
TQD
  ↓ gRPC
user-service GetAccess
  ↓
TQD ReserveQuota
  ↓
PostgreSQL AcceptReport
  ↓
Redis CommitQuota
  ↓
HTTP response qua gateway
```

## Vị trí source

```text
shared/common/requestcontext      request identity đã có
shared/common/requestgrpc         request identity qua gRPC
shared/common/identitycontext     Actor
shared/common/callercontext       Caller
shared/common/identitygrpc        Actor/Caller qua trusted gRPC
shared/common/businessoperation   business action
shared/common/access              kết quả quyền sử dụng

user-service/internal/access      đọc subscription + plan access
user-service/infra/postgres/access
user-service/infra/handler/accessgrpc

tqd-service/internal/report       create report flow
tqd-service/internal/quota        reserve/commit/cancel
tqd-service/infra/redisquota      Redis atomic state
tqd-service/internal/usage        durable usage event
tqd-service/infra/postgres/usageevent
tqd-service/internal/reportjob    background report job
tqd-service/infra/postgres/reportjob
tqd-service/internal/quotacleanup xử lý reservation quá hạn
tqd-service/internal/usagecheck   so DB usage với Redis used

gateway-service/internal/grpcmetadata
                                HTTP context -> gRPC metadata
```

## Rule đọc folder

- `shared/common`: dữ liệu hoặc rule nhỏ cần dùng qua nhiều service.
- `<service>/internal/<capability>`: logic service sở hữu.
- `<service>/infra/<technology-or-boundary>`: DB, Redis, gRPC adapter.
- `<service>/cmd`: chỉ lắp các phần lại với nhau.
