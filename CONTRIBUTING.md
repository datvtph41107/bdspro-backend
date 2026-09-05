# Contributing to QHPRO Backend

## Một workflow cho toàn đội

Mọi developer dùng cùng public commands:

```bash
make configure
make doctor
make dev-up
make dev-status
make dev-smoke
```

IDE tasks phải gọi các command này, không sao chép một workflow riêng vào IDE.

## Trước khi sửa code

Trả lời các câu hỏi phù hợp:

1. User outcome là gì?
2. Service nào sở hữu business state?
3. Invariant nào phải luôn đúng?
4. Actor và subject là ai?
5. Durable identity/idempotency key là gì?
6. Transaction kết thúc ở đâu?
7. Response loss, retry, duplicate và concurrency sẽ ra sao?
8. Client/server cũ có còn tương thích không?
9. Log, metric và test nào chứng minh kết quả?

Không bắt đầu bằng việc chọn folder hoặc tạo abstraction.

## Dependency direction

```text
handler -> usecase -> domain
adapter -> consumer-owned interface
composition root -> concrete dependencies
```

Không cho phép:

- domain import transport, SQL, Redis hoặc broker;
- service đọc database của service khác;
- Gateway quyết định business policy;
- Shared sở hữu commercial, payment, quota hoặc report policy;
- constructor tự đọc environment hoặc mở hidden connection;
- detached goroutine, deep `log.Fatal` hoặc `os.Exit`.

## Contract và migration

- Sửa source `.proto`, sau đó chạy `make generate-backend`.
- Không sửa generated file.
- Thay đổi contract theo hướng additive.
- Không tái sử dụng protobuf field number.
- Mỗi schema chỉ có một migration authority.
- Migration dữ liệu đang dùng theo expand, backfill, switch, constrain, retire.
- Serving process không chạy schema mutation mặc định.

## Proof ladder

```text
format/compiler
-> package unit test
-> service integration test
-> runtime flow smoke
-> backend verification
-> acceptance checkpoint
```

Thay đổi money, quota, queue/worker hoặc shared state phải xem xét concurrency,
idempotency và restart recovery.

## Definition of Done

- owner và invariant rõ ràng;
- code nằm đúng responsibility;
- contract compatible;
- database constraint/transaction phù hợp;
- retry và duplicate an toàn;
- context/deadline được truyền;
- resource được đóng;
- log không lộ secret;
- test và smoke liên quan pass;
- tài liệu/runbook được cập nhật khi behavior vận hành thay đổi;
- `git diff --check` pass.

Không `git add`, commit hoặc push thay cho người sở hữu worktree.
