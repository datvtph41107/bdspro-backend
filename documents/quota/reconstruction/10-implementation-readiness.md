# Implementation readiness

## Trạng thái source

Reconstruction bundle đã được ghép vào backend và frontend. Các lỗi compile và các khoảng trống migration phát hiện trong lúc review đã được sửa trong source:

- feature flag report được snapshot khi TQD khởi động; flag tắt giữ nguyên identity, operation và quota legacy;
- Gateway chỉ dùng canonical metadata adapter cho `POST /v2/tqd/map-workspace/reports` khi flag bật;
- internal protobuf có template generate riêng, không tạo HTTP gateway handler;
- `clientRequestId` trong proto được giữ làm fallback cho client cũ, nhưng phải khớp Idempotency-Key nếu cả hai cùng có;
- usage reconciliation dùng compare-and-set và quay vòng qua các page scope;
- quota cleanup/finalization đối chiếu đúng `reservation_id`, nên cùng command ở kỳ quota khác không thể dùng nhầm durable evidence;
- report runner release job `running` stale về `pending`;
- access query bỏ qua subscription đã hết hạn và từ chối usage entitlement sai amount/period;
- internal file gRPC được nối sau flag riêng, giới hạn 64 MiB, kiểm secure path/Actor/file access và cleanup upload lỗi;
- migration V3 `000009_report_quota_target` có cả up/down.
- local lab có smoke generator opt-in để kiểm chứng trọn vòng job -> file gRPC -> PDF; đây không phải renderer nghiệp vụ production.

## Flags

```text
QHPRO_REPORT_TARGET_ENABLED=false   # default
QHPRO_FILE_GRPC_ENABLED=false       # default
QHPRO_FILE_GRPC_PORT=8203           # default khi file gRPC bật
QHPRO_REPORT_GENERATOR_MODE=none    # dùng smoke chỉ trong local lab
```

Không bật report target nếu Gateway và TQD không nhận cùng giá trị flag.

## Kiểm tra source đã chạy

```bash
bash scripts/generate-internal-protobuf.sh
bash scripts/validate-reconstruction-target.sh
```

Targeted Go tests đã pass cho shared canonical identity, Gateway, user access, TQD report/quota/job/reconciliation và file gRPC. Frontend Jest cho request identity/report idempotency pass 20 tests; ESLint các file frontend target không có lỗi/warning.

Full React Native typecheck hiện chưa xanh vì baseline có nhiều lỗi ngoài reconstruction. Hai HTTP integration test của shared common không chạy được trong filesystem sandbox vì `httptest` không được phép mở loopback port; các package tương ứng vẫn compile và các unit test không dùng listener đã pass.

## Chưa runtime-proven

- Chưa chạy migration trên PostgreSQL thật.
- Chưa chạy Redis Lua/PostgreSQL concurrency integration.
- Chưa cấu hình trusted metadata key/secret trên môi trường triển khai.
- Chưa xác minh catalog/subscription dữ liệu thật cho `workspace.report.generate`.
- Chưa có renderer generated-report nghiệp vụ thật; local lab có thể bật generator `smoke`, production để `none` cho tới khi renderer thật được inject.
- Chưa có evidence end-to-end frontend -> Gateway -> TQD -> user -> Redis/PostgreSQL.

Các mục này là deployment gates, không được thay bằng `AutoMigrate` hoặc tự bật flag từ application code.
