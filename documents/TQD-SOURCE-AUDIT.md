# TQD source audit

## Phạm vi và owner

TQD là owner canonical cho dữ liệu bản đồ, operation thương mại, quota và usage. Gateway chỉ chuyển tiếp request; User Service quyết định catalog, subscription và entitlement; Payment Service quyết định tiền và settlement; File Service sở hữu file đầu ra.

## Quy tắc cấu trúc

- `cmd/`: khởi động process, đọc cấu hình, mở dependency và lắp interceptor.
- `internal/usecase/`: điều phối nghiệp vụ, không phụ thuộc handler cụ thể.
- `internal/domain/`: invariant và kiểu nghiệp vụ.
- `infra/handler/`: chuyển đổi gRPC/HTTP sang usecase.
- `infra/postgres/`, `infra/redis/`, `infra/rpc/`: adapter bên ngoài.
- `migrate/`: thay đổi schema có thứ tự và kiểm tra contract.
- `wire/`: composition dependency injection.

`transport_ingress.go` và `operation_routing_test.go` đã được chuẩn hóa tên. `operation_target.go` vẫn là điểm cần tách tiếp thành routing, quota interceptor và status mapping; việc này không được làm bằng cách đưa nghiệp vụ vào `cmd` thêm nữa.

## Phát hiện theo mức độ

### P0 — chặn phát hành

- PostgreSQL Docker hiện lỗi bind mount `postgres-init.sql`, vì vậy chưa thể xác nhận migration và acceptance runtime.
- Download report cũ trả placeholder PDF; đã đổi sang lỗi rõ ràng cho đến khi file-service adapter canonical được nối.

### P1 — cần xử lý trong chu kỳ kế tiếp

- Một số handler/worker cũ vẫn có TODO nghiệp vụ; mỗi TODO phải được phân loại thành tính năng thật, loại bỏ hoặc trả `Unimplemented`, không để phản hồi thành công giả.
- Tách các logic tile khỏi `cmd` để composition root không trở thành business/infrastructure container.
- Thêm architecture test cấm handler tự truy cập repository hoặc tự tính quota.
- Chạy full integration lifecycle sau khi PostgreSQL hoạt động: checkout, webhook, activation, usage, exhaustion, upgrade, cancel và expiry.

### P2 — chuẩn hóa dài hạn

- Gom mapper/validator theo bounded context thay vì chỉ theo loại file.
- Chuẩn hóa error contract giữa gRPC và gateway.
- Cập nhật source manifest sau khi hoàn tất rename; không sửa hash thủ công.
- Bổ sung metric cho reservation, usage commit, webhook replay và job retry.

## Invariant bắt buộc

1. Một operation chỉ có một mã canonical trong protobuf/registry.
2. Một command key không được dùng lại với request hash khác.
3. Read action chỉ ghi usage sau khi action thành công.
4. Retry không được tạo usage thứ hai.
5. Webhook trùng là replay an toàn; webhook khác dữ liệu bị từ chối.
6. Entitlement hết hạn hoặc không hợp lệ phải fail closed.
7. Không trả dữ liệu thành công giả khi adapter downstream chưa sẵn sàng.

## Bằng chứng hiện có

- `go test ./cmd` và `go test ./infra/handler/grpc` đạt sau refactor.
- Quota, usage, payment wallet race và frontend commercial tests đã đạt ở các lần kiểm tra trước.
- `git diff --check` đạt.

Audit này là checklist vận hành cho các thay đổi tiếp theo; mọi PR sửa TQD phải nêu owner, invariant bị tác động và kiểm thử tương ứng.
