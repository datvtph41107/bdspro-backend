# Group 1 — Task 1.4: Retry and Cancellation Contract

## Mục tiêu

Chuẩn hóa retry và cancellation tại ranh giới Client → Gateway → gRPC → downstream service.

Task này không triển khai quota, usage ledger, billing hoặc idempotency storage.

## Retry contract

| Tình huống | Tự retry | Quy tắc |
|---|---:|---|
| HTTP 401 | Có, tối đa 1 lần | Refresh token rồi replay cùng request |
| Timeout | Không | Trả `TIMEOUT` cho caller |
| Network error | Không | Caller quyết định thử lại |
| HTTP 5xx | Không | Không suy đoán request có an toàn để chạy lại |
| Client cancellation | Không | Trả `CANCELED` |

Auth retry giữ nguyên:

```text
operation_id
idempotency_key nếu có
method
payload
caller metadata
```

Mỗi replay là một request attempt mới nên Gateway tạo `request_id` mới.

## Cancellation contract

Caller sở hữu `AbortController` và truyền `signal` vào request config.

```text
Canceled before dispatch
→ không gọi transport

Canceled while waiting token refresh
→ remove khỏi queue
→ không replay

Canceled during transport
→ normalize thành ApiError(CANCELED)
→ không retry
```

Cancellation của một request không hủy token refresh dùng chung và không hủy request khác trong queue.

## Server semantics

HTTP/gRPC context phải giữ nguyên:

```text
cancellation
deadline
request values
canonical metadata
```

`context.WithValue` và `metadata.NewOutgoingContext` phải dùng parent context hiện tại, không được thay bằng `context.Background()` cho request-scoped downstream calls.

## Unknown outcome

Client cancellation không đồng nghĩa nghiệp vụ đã rollback.

```text
Client canceled
≠ server chắc chắn chưa chạy
≠ database chắc chắn chưa commit
≠ async job chắc chắn đã dừng
```

Với command có tác động, retry an toàn tương lai phải dựa trên idempotency enforcement. Task 1.4 chỉ truyền key, chưa bảo đảm xử lý một lần.

## Async job

Sau khi server đã chấp nhận một job bền vững, client disconnect không tự hủy job. Hủy job cần một business command riêng và policy riêng.

## Trách nhiệm

### Client request layer

- nhận `AbortSignal`;
- không dispatch khi signal đã abort;
- loại request đã abort khỏi refresh queue;
- không auto-retry timeout, network, 5xx hoặc canceled request;
- chuẩn hóa cancellation thành `CANCELED`.

### Gateway và gRPC

- dùng request context hiện tại;
- truyền cancellation và deadline;
- không tạo background context cho request-scoped call.

### Business handler

- kiểm tra `ctx.Done()` khi tác vụ dài;
- quyết định rollback, refund hoặc tiếp tục theo business contract;
- không suy ra refund quota chỉ từ client cancellation.

## Ranh giới Task 1.4

Chưa triển khai:

- generic exponential retry;
- retry budget;
- circuit breaker;
- job cancellation API;
- idempotency record/replay;
- quota refund;
- billing reversal;
- operation success state machine.
