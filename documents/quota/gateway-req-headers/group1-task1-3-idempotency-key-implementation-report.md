# QHPRO Group 1 — Task 1.3 Implementation Report

## Kết luận

Đã triển khai input contract và propagation cho optional `Idempotency-Key`.

Gateway không tạo fallback. Key sai bị reject để tránh phá semantic retry.

## Luồng

```text
Client-owned Idempotency-Key
→ HTTP request config
→ Gateway validation
→ gRPC metadata
→ service context
→ outbound gRPC metadata
→ downstream service
```

## Tích hợp đầu tiên

Generated report dùng cùng key cho:

```text
Idempotency-Key
clientRequestId
```

Đây là compatibility bridge, chưa phải idempotency enforcement chung.

## Không triển khai

- Redis record;
- lock;
- payload fingerprint;
- result replay;
- quota charge deduplication;
- billing deduplication.

## Validation trong sandbox

- `gofmt`: PASS.
- requestcontext unit tests bằng Go 1.23 compatibility copy: PASS.
- isolated TypeScript check cho request identity/client: PASS.
- isolated TypeScript check cho workspace report service: PASS.
- full middleware Go tests: BLOCKED vì sandbox thiếu dependency cache/internet và repository yêu cầu Go 1.24.1.
- Gateway full tests: BLOCKED vì gateway module yêu cầu Go 1.25.0.
- Jest: BLOCKED vì source ZIP không chứa `node_modules`.
