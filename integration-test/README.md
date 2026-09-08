# Integration Test

`integration-test` là module kiểm chứng behavior xuyên nhiều service. Nó **không sở hữu business capability** và không phải một service runtime.

## Vai trò

- Chứng minh các flow cần nhiều owner cùng hoạt động.
- Giữ fixture/acceptance test ở một nơi thay vì nhét vào từng service.
- Chỉ chạy sau khi local stack và các migration cần thiết đã sẵn sàng.

## Khi nào dùng

Trong inner loop, ưu tiên test owner service trước:

```bash
make test service=payment-service
```

Khi thay đổi ảnh hưởng contract hoặc flow xuyên service, từ root chạy:

```bash
make test-e2e
```

Full local candidate:

```bash
make accept
```

## Boundary

Business rule vẫn nằm trong owner service. Module này chỉ **quan sát và chứng minh behavior tích hợp**; không chứa implementation thay thế cho User, Payment, TQD, Notification hay owner khác.

Fixtures nằm trong `fixtures/` và chỉ phục vụ acceptance/integration testing.
