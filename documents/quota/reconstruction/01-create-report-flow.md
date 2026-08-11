# Create generated report - target flow

Use case chuẩn: `workspace.report.generate`.

## 1. Frontend

`src/Services/tqd-workspace-map-service.ts`

- tạo hoặc giữ `Idempotency-Key` cho command tạo report;
- cùng key được gửi ở HTTP header và `clientRequestId` để tương thích source cũ;
- request client giữ `Operation-ID` khi auth retry;
- frontend không gửi BusinessOperation và không tự tính quota.

## 2. Gateway

HTTP middleware tạo/nhận Request-ID, Operation-ID, Idempotency-Key.
JWT/API key boundary xác thực caller. JWT hợp lệ được bind thành Actor.

`gateway/internal/grpcmetadata.FromHTTPRequest` chỉ lấy canonical context và trace fields cần thiết.
Raw Authorization, API key, plan id không được forward bởi target adapter.

## 3. TQD gRPC ingress

Riêng `CreateGeneratedReport`:

```text
requestgrpc
  ↓
identitygrpc verify signed metadata
  ↓
grpcoperation -> workspace.report.generate
  ↓
legacy operation/quota: skip
```

RPC chưa migrate vẫn chạy đường cũ để cutover từng use case.

## 4. Durable idempotency check

Trước khi đọc access/quota, TQD tìm `user_reported` theo profile + command key.

- đã có report + cùng `request_hash` -> trả report cũ, không charge lại;
- đã có report nhưng input khác -> conflict;
- chưa có -> tiếp tục GetAccess/ReserveQuota.

Nhờ vậy retry cùng Idempotency-Key sang một subscription cycle khác cũng không bị charge lần hai.

## 5. GetAccess

TQD gọi `userpb.InternalAccessService/GetAccess`.
Request chỉ có `operation`; không gửi profile ID.

user-service lấy profile từ Actor đã được verify rồi:

```text
Actor profile
  ↓
active catalog_subscription
  ↓
plan version
  ↓
feature access
  ↓
usage allowance nếu có
  ↓
access.Result
```

Result dừng ở allowed/limit/period. Không đọc Redis usage.

## 6. Reserve quota

`quota.Service.ReserveQuota` chọn command key:

```text
Idempotency-Key nếu có
ngược lại Operation-ID
```

Redis Lua thực hiện atomically:

```text
used + reserved + amount <= limit
```

Cùng command đang reserved/committed không được giữ quota thêm lần nữa.
Reservation canceled được phép retry và tạo reservation mới.

## 7. Business acceptance

Một PostgreSQL transaction lưu:

```text
user_reported
quota_usage_events
report_jobs
```

Đây là điểm report được xem là accepted.

Nếu transaction fail -> CancelQuota.
Nếu transaction commit -> không rollback report vì lỗi transport/Redis sau đó.

## 8. Commit quota

Sau DB commit, `CommitQuota` chuyển:

```text
reserved -> used
```

Nếu Redis response fail sau khi DB đã commit, request trả lỗi retryable. Retry cùng command không tạo thêm report/usage event. Cleanup cũng có thể hoàn tất Redis bằng durable usage event.

## 9. Async job

`report_jobs` là queue PostgreSQL nhỏ, không thêm Kafka/RabbitMQ khi source chưa cần.
Worker claim bằng `FOR UPDATE SKIP LOCKED`.

Baseline chưa có renderer thật cho generated report. Vì vậy request flow có thể chạy với worker render tắt: report/job vẫn accepted và job ở `pending`. Khi có renderer thật, truyền `reportjob.Generator` vào composition rồi bật runner.
