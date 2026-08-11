# Target source status

Baseline backend: `bdspro-golang-microservices-20260811-091814`.

Mục tiêu của archive này là đọc và kiểm toàn bộ flow trước khi migration/cutover.
Target flow mặc định chưa bật trong runtime cũ.

## Đã có source

```text
Frontend HTTP command
  ↓
Gateway request identity + trusted Actor/Caller
  ↓ signed gRPC
TQD business operation
  ↓
Durable idempotency check
  ↓ gRPC
user-service GetAccess
  ↓
ReserveQuota
  ↓
PostgreSQL transaction:
  user_reported
  quota_usage_events
  report_jobs
  ↓
CommitQuota
  ↓
response
```

Background:

```text
expired reservation -> usage event có/không -> commit/cancel
saved usage         -> compare-and-set Redis used -> fix when different
stale running job   -> release -> pending -> claim lại
pending report job  -> worker -> ready/failed khi có Generator thật
```

## Source đã kiểm tra tĩnh

`bash scripts/validate-reconstruction-target.sh`:

- target files tồn tại;
- target business code không import HTTP/Gin;
- access/quota/businessoperation không phụ thuộc DB/Redis/gRPC;
- internal user/file proto không có HTTP annotation;
- generated internal proto không có grpc-gateway handler;
- Gateway target không forward raw Authorization/API key/plan fields;
- Gateway chỉ dùng target metadata ở đúng report route khi flag bật, RPC khác giữ legacy bridge;
- CreateGeneratedReport chỉ bỏ qua legacy identity/operation/quota path khi flag bật;
- TQD target file client dùng gRPC, không dùng HTTP.
- file-service internal gRPC được nối runtime sau flag mặc định tắt.

Toàn bộ Go files của reconstruction đang `gofmt` sạch.
Frontend target files parse sạch bằng TypeScript 5.8.3 `transpileModule`.

## Core logic tests đã chạy

Các package target được test trực tiếp bằng Go 1.25 của repo.

Pass:

```text
common/businessoperation
common/access
common/requestlog
user/internal/access
tqd/internal/quota
tqd/internal/report
tqd/internal/reportjob
tqd/internal/quotacleanup
tqd/internal/usagecheck
```

Các case quan trọng đã có test:

- reserve -> commit -> retry không charge lại;
- 99/100 concurrent request chỉ một request được reserve;
- cancel trả lại reserved quota;
- canceled reservation có thể retry;
- operation không có quota không tạo reservation;
- report accepted chỉ charge một lần;
- retry cùng Idempotency-Key sang kỳ quota mới không charge lại;
- race cùng command qua ranh giới kỳ quota chỉ finalize reservation có durable usage event tương ứng;
- cùng Idempotency-Key nhưng report input khác -> conflict;
- DB fail trước accepted -> quota về 0/0;
- expired reservation + usage event -> commit;
- expired reservation không có usage event -> cancel;
- saved usage khác Redis used -> sửa Redis used;
- compare-and-set không ghi đè quota commit đồng thời;
- usage checker chỉ quay vòng các kỳ còn nằm trong cửa sổ Redis retention, không tái tạo counter lịch sử vô hạn;
- running job stale được release về pending;
- feature flag tắt giữ nguyên legacy bridge;
- client cũ chỉ gửi `clientRequestId` trong proto vẫn được bind làm Idempotency-Key; nếu body/header khác nhau request bị reject;
- request log không ghi raw Idempotency-Key.

## Chưa được gọi là runtime-proven

1. Migration V3 `000009_report_quota_target` đã có up/down nhưng chưa chạy trên DB thật.
2. Chưa xác nhận dữ liệu catalog V2 ở môi trường dev có đúng subscription/plan/access của report.
3. Chưa cấu hình trusted metadata secret/key cho các service target.
4. Chưa chạy Redis/PostgreSQL integration thật.
5. Baseline không có renderer generated-report thật; `report_jobs` có thể pending nhưng worker render chưa nên bật.
6. File-service gRPC target được nối sau `QHPRO_FILE_GRPC_ENABLED`, nhưng chưa có renderer target sử dụng và chưa thay runtime HTTP cũ.

## Một rule cố ý fail thay vì đoán

Catalog có `day`, `calendar_month`, `lifetime`, nhưng baseline chưa chỉ ra đủ boundary để dùng an toàn.
Target hiện chỉ dùng `subscription_cycle` cho report và sẽ trả lỗi nếu gặp các period kia.
Migration review phải chốt timezone cho day/month và mốc lifetime trước khi hỗ trợ.

Operation report đã khai `UsageCode` nhưng DB thiếu usage allowance cũng trả lỗi, không tự hiểu thành unlimited.

## Runtime flag

`QHPRO_REPORT_TARGET_ENABLED` mặc định tắt.

`QHPRO_FILE_GRPC_ENABLED` mặc định tắt; `QHPRO_FILE_GRPC_PORT` mặc định `8203`.

Không bật flag trước khi hoàn thành migration checklist trong `04-migration-gates.md`.
