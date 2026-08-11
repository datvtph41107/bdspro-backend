# Test map

Mục đích của bảng này là nhìn tên test và biết rule nào đang được chứng minh.

## G1 request identity

```text
shared/common/requestgrpc/metadata_test.go
```

Kiểm:
- outgoing chỉ truyền ID đã có;
- incoming Operation-ID/Idempotency-Key sai bị reject;
- context conflict không bị silent overwrite.

## G2 trusted identity

```text
shared/common/identitygrpc/metadata_test.go
```

Kiểm:
- Actor/Caller được restore từ metadata đã trusted;
- metadata duplicate/invalid bị reject;
- caller và actor phải khớp nhau.

## G3 business operation

```text
shared/common/businessoperation/*_test.go
tqd-service/infra/grpcoperation/*_test.go
tqd-service/cmd/target_migration_bridge_test.go
```

Kiểm:
- operation code hợp lệ;
- context không silent overwrite;
- đúng gRPC method mới chọn `workspace.report.generate`;
- method chưa migrate không tự có operation.
- target flag tắt không skip legacy operation/quota bridge.

## G4 access

```text
user-service/internal/access/service_test.go
```

Kiểm:
- active subscription + plan access -> allowed + limit;
- thiếu subscription -> deny;
- period chưa có boundary rõ -> error thay vì đoán.

## G5/G6 quota

```text
tqd-service/internal/quota/service_test.go
```

Kiểm:
- reserve/commit/retry;
- last quota concurrency 99/100;
- cancel;
- retry sau cancel;
- operation không cần quota.

Redis Lua cần integration test thật sau migration.
MemoryStore chỉ chứng minh service/state rules.

## G7/G8 report accepted + usage event

```text
tqd-service/internal/report/service_test.go
```

Kiểm:
- report accepted + usage một lần;
- durable command check trước quota;
- retry sang kỳ mới không charge lại;
- reuse key với input khác bị conflict;
- DB fail trước commit không để lại quota used/reserved.
- body `clientRequestId` cũ được bind khi thiếu header; body/header conflict bị reject.

## G9 background checks

```text
tqd-service/internal/quotacleanup/service_test.go
tqd-service/internal/usagecheck/service_test.go
```

Kiểm:
- reservation hết hạn có usage -> commit;
- không có usage -> cancel;
- DB usage khác Redis used -> sửa Redis used.
- Redis used thay đổi đồng thời -> compare-and-set không ghi đè commit mới.

## G10 request log

```text
shared/common/requestlog/logger_test.go
```

Kiểm:
- Request-ID, Operation-ID, Actor, BusinessOperation có trong trace fields;
- chỉ log `idempotency_key_set=true`;
- không ghi raw Idempotency-Key.

## Async report job

```text
tqd-service/internal/reportjob/worker_test.go
```

Kiểm:
- worker complete -> completed;
- temporary failure -> retry;
- quá số lần -> failed.
- stale running job -> pending để worker claim lại.

PostgreSQL `FOR UPDATE SKIP LOCKED` cần integration/concurrency test sau migration.
