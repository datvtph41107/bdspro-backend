# Group 2 — Task 2.3 Implementation Report

## 1. Kết luận

Task 2.3 đã triển khai caller classification xuyên Gateway, unary gRPC và streaming gRPC.

```text
JWT user
API-key application
Gateway anonymous
Internal service
```

Caller classification chưa được dùng làm quota subject, permission hay billing identity.

## 2. Kiến trúc sau thay đổi

```text
HTTP request
→ JWT/API-key authentication
→ callercontext.Caller
→ grpc-gateway metadata
→ signed transport assertion
→ server verification
→ service caller context
→ signed downstream propagation
```

Hai context được giữ riêng:

```text
callercontext.Caller
= nguồn logic của request

transporttrust.Caller
= service vận chuyển metadata đã được xác minh
```

## 3. Caller matrix

| Kind               | Điều kiện                        | Actor | API-key app identity |
| ------------------ | -------------------------------- | ----: | -------------------: |
| `user`             | JWT hợp lệ                       |    Có |            Có thể có |
| `api_key`          | Hub verify thành công, không JWT | Không |                   Có |
| `anonymous`        | Gateway public route             | Không |                Không |
| `internal_service` | Service tự phát sinh call        | Không |                Không |

JWT và API key có thể cùng xuất hiện. Khi đó caller vẫn là `user`; API key chỉ bổ sung identity của ứng dụng gọi.

## 4. Thay đổi chính

### Caller source of truth

Tạo package:

```text
shared/common/callercontext
```

Package chỉ lưu caller kind và API-key identity không nhạy cảm.

### API-key verification

Tạo contract tại:

```text
shared/common/jwt/api_key.go
```

Gateway dùng:

```text
gateway-service/service/hub_api_key_verifier.go
```

để gọi `hub-service.VerifyApiKey`.

Mapping lỗi:

```text
Header sai       → 400
Key invalid      → 401
Hub unavailable  → 503
```

### Public route matching

Plain path chuyển từ arbitrary prefix sang path-segment boundary:

```text
/v2/auth/login/key-malicious
```

không còn match route public:

```text
/v2/auth/login/key
```

### Trusted metadata

Các field sau được bind vào HMAC assertion:

```text
x-qhpro-caller-kind
x-qhpro-api-key-id
x-qhpro-api-key-app
x-api-key
```

### Streaming gRPC

Bổ sung:

```text
StreamClientInterceptor
ParseGrpcMetadataContextStreamMiddleware
```

Tất cả `grpc.NewServer` được quét đều có unary và stream metadata boundary. Tất cả client sử dụng shared unary interceptor đều có shared stream interceptor.

### Legacy server helper

`shared/common/setup.go` cũng được đưa qua common unary/stream metadata boundary để không còn một helper tạo server bypass contract.

## 5. Security decisions

### Anonymous enforce policy

```text
Signed Gateway anonymous → allow
Unsigned direct gRPC     → reject
```

Anonymous là caller không có user/API-key identity, không phải transport không đáng tin.

### Raw API key

Raw key không được lưu trong `callercontext`. Chỉ `key_id` và `app_name` được sử dụng làm caller identity.

Raw key vẫn được forward trong signed metadata cho interceptor legacy cần tái xác minh. Đây là compatibility bridge.

## 6. Tests

Đã thêm regression tests cho:

- Caller validation và metadata round-trip.
- API-key header normalization.
- Conflicting hoặc unsafe API-key header.
- Verified API key trên private route.
- Invalid API key trên public route.
- API-key authority unavailable.
- Public-route prefix bypass.
- Caller metadata tampering.
- Signed anonymous caller.
- Invalid caller/actor combination.
- Stream server context propagation.
- Stream client signing.
- Hub verifier response mapping.

## 7. Validation result

### PASS

```text
Static architecture gate
All changed Go files gofmt
Callercontext isolated tests
JWT/API-key isolated tests
Transport-trust isolated tests
Metadata middleware isolated tests
Gateway Hub verifier isolated tests
Auth middleware isolated tests
Patch apply check on Task 2.2 baseline
```

### Pending

Full repository module tests chưa chạy được trong sandbox:

```text
Sandbox Go:          1.23.2
shared/common needs: 1.24.1
gateway needs:       1.25.0
network:             blocked
```

Full gate cố tải Go 1.24.1 từ `proxy.golang.org` và thất bại do network bị chặn. Vì vậy Task chưa được đánh dấu full runtime PASS.

## 8. Trade-offs và debt

### Hub lookup latency

Mỗi API-key request hiện gọi Hub đồng bộ. Kiến trúc sau nên bổ sung cache ngắn hạn, có expiry và revocation/version strategy.

### API-key storage

Source hiện còn đường legacy lưu/hiển thị raw key. Cần migration sang hash hoặc fingerprint và chỉ hiển thị secret một lần.

### v1 HTTP proxy

Legacy `/v1/:service/*path` không tạo signed gRPC caller identity cho downstream HTTP service. Không dùng đường này làm Quota V2 identity source trước khi migrate hoặc retire.

### Client IP

`X-Forwarded-For` chưa có trusted-proxy contract. IP chỉ là telemetry, chưa phải anonymous quota subject.

### Shared secret

HMAC shared secret là bước chuyển tiếp. Đích dài hạn vẫn là workload identity/mTLS, network policy và per-service authorization.

## 9. Ranh giới Task 2.3

Task không triển khai:

```text
Quota subject resolution
Organization membership verification
Caller permission matrix
Anonymous rate limit
Usage meter
Billing
Audit ledger
```

## 10. Next decision gate

Task kế tiếp hợp lý:

```text
Group 2 — Task 2.4
Quota Subject Resolution Rules
```

Cần chốt user, organization, API key, project, subscription và anonymous traffic được map thành quota subject như thế nào. Không bắt đầu task đó cho tới khi full Task 2.3 gate PASS trên môi trường local.
