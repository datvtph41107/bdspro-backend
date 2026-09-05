# BDSPro Backend

BDSPro Backend là workspace Go chứa các service nghiệp vụ của BDSPro. QHPro là
capability quy hoạch thuộc TQD Service, không phải tên của toàn hệ thống.

Admin và ứng dụng mobile chỉ gửi command hoặc đọc dữ liệu qua Gateway. Giá,
thanh toán, subscription, entitlement, quota và usage luôn do backend quyết
định.

## Bắt đầu

Yêu cầu: Linux/WSL2, Git, Make và Docker Engine có Docker Compose.

```bash
make setup
make up
make status
make smoke
```

- `setup` tạo config local còn thiếu và cài đúng toolchain của repository.
- `up` build image, chạy migration và khởi động backend.
- `status` cho biết container nào đang chạy và health hiện tại.
- `smoke` gọi public API thật qua Gateway.

Xem toàn bộ lệnh dành cho developer:

```bash
make help
```

## Công việc hằng ngày

```bash
make up
make logs service=payment-service
make test service=payment-service
make test-e2e
make down
```

`down` giữ dữ liệu. `reset` xóa database/queue/file local và luôn yêu cầu xác
nhận.

Các lệnh dài cũ như `dev-up`, `test-backend` và
`runtime-value-chain-acceptance` hiện chỉ là compatibility/internal targets.
Developer không cần ghi nhớ chúng.

## Tìm code theo nghiệp vụ

| Cần sửa | Owner |
| --- | --- |
| HTTP public, JWT, routing, error mapping | `gateway-service` |
| Profile, login, IAM, catalog, subscription, entitlement | `user-service` |
| Thành viên và quyền trong tổ chức | `organization-service` |
| Order, payment attempt, webhook, settlement, outbox | `payment-service` |
| QHPro, planning, quota, usage, report | `tqd-service` |
| File metadata, nội dung và signed access | `file-service` |
| Notification inbox và delivery | `notification-service` |

`auth-service` đang là compatibility boundary trong lúc chuyển authorization
về User. Không thêm nghiệp vụ mới vào Auth.

Khi đọc một service, bắt đầu từ `main.go` hoặc `cmd/<server>/main.go`, sau đó đi
theo constructor tới package nghiệp vụ và repository tương ứng. Không sửa file
generated `*.pb.go` hoặc `wire_gen.go` bằng tay.

## Các lệnh chính

```bash
make dev service=payment-service
make test
make test service=tqd-service
make migrate
make migrate service=user-service
make generate
make logs service=notification-service
```

`make dev` rebuild và restart riêng image/container của service trong topology
đang chạy. Full integration dùng `make up`. Native debugger cần host/container
wiring hai chiều nên đang được chuẩn hóa lần lượt, bắt đầu từ Payment; README
không tuyên bố mode đó đã sẵn sàng trước khi được kiểm chứng.

## Local và production

Source, binary và Dockerfile của application là một. Hai môi trường chỉ khác
endpoint, secret, tài nguyên và số instance được inject từ bên ngoài.

- `compose.yaml` là topology development.
- `<service>/config/runtime.yml` chứa default kỹ thuật không nhạy cảm.
- `.env` bị ignore và chỉ cung cấp runtime values.
- `<service>/migrate/` là migration authority của service.
- Application log ra stdout/stderr.

Production hiện vẫn dùng semantics của `shared/code/deploy.sh`: operator chọn
một service, build/chuyển artifact và restart service đó. Script này đang được
giữ tương thích trong khi production Compose, config và selector được audit;
không được dùng file `.env` development làm production config.

## Trạng thái chuyển đổi

Correctness của commercial/quota/report vẫn được giữ: integer money, database
transaction, idempotency, webhook verification, outbox/inbox, durable jobs,
quota durability và graceful shutdown.

Topology development dùng một PostgreSQL/PostGIS, một Redis và một RabbitMQ;
database/Redis index/migration authority vẫn tách theo owner. Migration là
one-shot container. Payment, Notification và TQD chạy actor bền vững trong
lifecycle của default process, không tạo standing worker container riêng.

Fresh workspace, fresh database, recovery labs và commercial/quota E2E đã PASS.
Whole backend vẫn chưa CLOSED vì Membership/Transaction/Marketing/Task,
Chat/Chat-v1, AI, non-core runtime và production deployment parity còn gate.

Ba bản đồ current/target và deployment gap nằm tại
[`documents/OPERATIONAL-MAPS.md`](documents/OPERATIONAL-MAPS.md).

Đối chiếu 23 production deployable với source hiện tại nằm tại
[`documents/PRODUCTION-COMPATIBILITY.md`](documents/PRODUCTION-COMPATIBILITY.md).

Kiến trúc owner hiện tại nằm tại
[`BACKEND-ARCHITECTURE.md`](BACKEND-ARCHITECTURE.md). Kết quả kiểm chứng và các
gate chưa đóng được ghi tại [`BACKEND-ACCEPTANCE.md`](BACKEND-ACCEPTANCE.md).

Đọc [`CONTRIBUTING.md`](CONTRIBUTING.md) trước khi đổi public contract,
transaction, migration, payment, entitlement, quota hoặc event liên service.
