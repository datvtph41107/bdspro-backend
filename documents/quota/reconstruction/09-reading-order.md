# Reading order - đọc source theo đúng request flow

Use case chuẩn: `POST /v2/tqd/map-workspace/reports` -> `CreateGeneratedReport`.

Không cần mở toàn bộ repo cùng lúc. Đọc theo thứ tự dưới đây.

## 1. Frontend tạo command

```text
frontend/src/Services/request/identity.ts
frontend/src/Services/request/client.ts
frontend/src/Services/tqd/report.ts
frontend/src/Services/tqd-workspace-map-service.ts
```

Câu hỏi cần trả lời:
- Operation-ID được giữ qua retry thế nào?
- Idempotency-Key được tạo ở action nào?
- frontend có gửi quota/plan/business operation không? -> không.

## 2. Gateway nhận HTTP

```text
gateway-service/internal/httpmiddleware/request_id.go
gateway-service/internal/httpmiddleware/operation_id.go
gateway-service/internal/httpmiddleware/idempotency_key.go
shared/common/jwt/middleware.go
gateway-service/internal/grpcmetadata/from_http.go
shared/common/grpcclient/interceptor.go
```

Đọc để thấy:

```text
raw HTTP
 -> canonical request IDs
 -> verified Actor/Caller
 -> gRPC metadata
 -> sign
```

HTTP dừng ở đây.

## 3. TQD nhận trusted gRPC

```text
shared/common/requestgrpc/interceptor.go
shared/common/identitygrpc/interceptor.go
tqd-service/cmd/target_ingress.go
tqd-service/cmd/business_operation_target.go
tqd-service/infra/grpcoperation/interceptor.go
```

Kết quả context trước handler:

```text
Request-ID
Operation-ID
Idempotency-Key
Actor
Caller
BusinessOperation = workspace.report.generate
```

## 4. gRPC handler chỉ chuyển dữ liệu

```text
tqd-service/infra/handler/grpc/workspace_map_grpc_handler.go
tqd-service/infra/reportgrpc/adapter.go
```

Không có quota/DB logic ở đây.

## 5. Report flow chính

```text
tqd-service/internal/report/service.go
```

Đây là file nên đọc kỹ nhất.

Flow:

```text
Actor + request identity + operation
 -> durable command check
 -> build report data
 -> GetReportAccess
 -> ReserveQuota
 -> AcceptReport transaction
 -> CommitQuota
 -> outcome
```

## 6. Durable command check

```text
tqd-service/internal/report/request_hash.go
tqd-service/infra/postgres/report/save_store.go
```

Rule:
- same key + same input -> return existing;
- same key + different input -> conflict;
- existing command không được charge lại, kể cả retry sang kỳ quota khác.

## 7. TQD hỏi user-service về access

```text
tqd-service/infra/useraccess/client.go
shared/protobuf/schema/user/access_internal.proto
user-service/infra/handler/accessgrpc/handler.go
user-service/internal/access/service.go
user-service/infra/postgres/access/subscription_store.go
user-service/infra/postgres/access/plan_access_store.go
```

TQD không gửi profile ID.
user-service lấy profile từ Actor đã trusted.

G4 dừng ở:

```text
allowed
limit
period
subscription evidence
```

Không đọc Redis ở đây.

## 8. Reserve quota

```text
tqd-service/internal/quota/service.go
tqd-service/internal/quota/types.go
tqd-service/infra/redisquota/reserve.go
```

Atomic check:

```text
used + reserved + amount <= limit
```

## 9. Business accepted

```text
tqd-service/infra/postgres/report/save_store.go
tqd-service/infra/postgres/usageevent/store.go
tqd-service/infra/postgres/reportjob/store.go
```

Một PostgreSQL transaction lưu:

```text
user_reported
quota_usage_events
report_jobs
```

Transaction commit là mốc business accepted.

## 10. Finalize quota

```text
tqd-service/infra/redisquota/finalize.go
```

```text
reserved -> committed
```

DB đã commit thì Redis/transport error sau đó không rollback report.

## 11. Background recovery

```text
tqd-service/internal/quotacleanup/service.go
tqd-service/internal/quotacleanup/runner.go
tqd-service/internal/usagecheck/service.go
tqd-service/internal/usagecheck/runner.go
```

Hai việc khác nhau:

```text
reservation quá hạn
 -> nhìn usage event
 -> commit hoặc cancel

DB usage khác Redis used
 -> compare-and-set Redis used

running report job stale
 -> release về pending
```

## 12. Async report job

```text
tqd-service/internal/reportjob/worker.go
tqd-service/internal/reportjob/runner.go
tqd-service/infra/postgres/reportjob/store.go
```

Quota đã charge ở lúc job accepted, không đợi render xong.
Renderer thật chưa có trong baseline nên runner render target chưa bật.

## 13. Observability

```text
shared/common/requestlog/logger.go
documents/quota/reconstruction/05-observability.md
```

Dùng Operation-ID để nối toàn bộ request.
Không log raw Idempotency-Key.

## 14. Đọc ngược

Sau khi đọc forward, mở:

```text
documents/quota/reconstruction/06-reverse-review.md
```

Bắt đầu từ `quota_usage_events` rồi đi ngược về frontend command.
Nếu một bước không tìm được key để nối, source/evidence vẫn chưa đủ rõ.

Sau cùng đọc `10-implementation-readiness.md` để biết phần nào đã được kiểm bằng source test và phần nào vẫn là deployment gate.
