# Group 2 — Task 2.2: Trusted Metadata Boundary

## Mục tiêu

Không coi các metadata sau là đáng tin chỉ vì chúng xuất hiện trong gRPC request:

```text
authid
profileid
originid
organizationid
session
role
type
planid
planfrom
```

Direct gRPC client có thể tự gửi các giá trị này. Task 2.2 tạo bằng chứng rằng metadata được phát ra từ Gateway hoặc internal service đã cấu hình.

## Ranh giới trách nhiệm

Task này xác minh **nguồn vận chuyển metadata**. Nó chưa:

- Xác định quota subject.
- Xác nhận organization membership.
- Xác nhận plan/entitlement còn hiệu lực.
- Phân quyền RPC theo caller service.
- Thay thế mTLS hoặc network policy.
- Trừ quota, ghi usage hoặc billing.

## Contract

Outgoing gRPC call thêm assertion:

```text
x-qhpro-internal-version
x-qhpro-internal-caller
x-qhpro-internal-issued-at
x-qhpro-internal-signature
```

Signature dùng HMAC-SHA256 và bind với:

- Full gRPC method.
- Caller service ID.
- Issued time.
- Actor metadata.
- Plan assertion.
- `request_id`.
- `operation_id`.
- `idempotency_key` nếu có.

Thay đổi một trường hoặc replay sang RPC khác làm signature không hợp lệ.

## Chế độ vận hành

### `audit`

```text
QHPRO_TRUSTED_METADATA_MODE=audit
```

- Verified assertion được đánh dấu `verified`.
- Metadata cũ chưa ký vẫn chạy để tương thích.
- Metadata cũ được đánh dấu `legacy_unverified`.
- Quota V2 không được dùng actor `legacy_unverified` làm quota subject.

Đây là chế độ migration, không phải production enforcement cuối cùng.

### `enforce`

```text
QHPRO_TRUSTED_METADATA_MODE=enforce
```

- Privileged metadata thiếu, sai hoặc hết hạn signature bị reject bằng `Unauthenticated`.
- Outbound client thiếu signer config bị fail trước khi gửi RPC.
- Anonymous request không mang privileged metadata vẫn có thể đi tiếp; method authorization là task riêng.

## Environment

```text
QHPRO_SERVICE_ID=gateway-service
QHPRO_INTERNAL_METADATA_SECRET=<secret-manager-injected-value>
QHPRO_TRUSTED_METADATA_MODE=audit|enforce
QHPRO_INTERNAL_METADATA_MAX_AGE=30s
QHPRO_INTERNAL_METADATA_CLOCK_SKEW=5s
```

Secret không được commit vào YAML, source hoặc database.

Mọi Gateway và service gọi gRPC phải có `QHPRO_SERVICE_ID` riêng, nhưng dùng secret phù hợp với policy triển khai hiện tại.

## Luồng

```text
Verified JWT at Gateway
→ canonical actor metadata
→ UnaryClientInterceptor signs metadata
→ service verifies signature
→ transporttrust.Caller{verified}
→ actor is parsed into identitycontext
→ downstream call is re-signed by current service
```

## One Source of Truth

```text
Actor data          → identitycontext.Actor
Transport trust     → transporttrust.Caller
Signing policy      → environment / secret manager
Quota subject       → chưa xác định trong Task 2.2
```

Actor và transport trust là hai concern khác nhau:

```text
Actor exists
≠ actor came from a verified internal caller
```

Code cần quyết định tài chính hoặc quota sau này phải kiểm tra:

```go
transporttrust.IsVerified(ctx)
```

## Failure modes

| Tình huống | Audit | Enforce |
|---|---|---|
| Không có actor metadata | Allow anonymous | Allow anonymous |
| Actor metadata chưa ký | Allow + unverified | Reject |
| Signature sai | Allow + unverified | Reject |
| Signature hết hạn | Allow + unverified | Reject |
| Method bị thay | Allow + unverified | Reject |
| Client signer thiếu secret | Gửi legacy | Fail outbound |

## Trade-off

HMAC assertion là lớp integrity/authentication phù hợp để đóng lỗ hổng hiện tại mà chưa thay toàn bộ transport.

Giới hạn:

- Shared secret làm tăng blast radius nếu một service bị compromise.
- Assertion có replay window ngắn; chưa có nonce ledger.
- Kết nối hiện vẫn có thể dùng insecure transport.
- mTLS helper hiện chưa được wire đầy đủ và CA verification còn chưa hoàn thiện.

Đích enterprise vẫn là:

```text
mTLS workload identity
+ network policy
+ per-service authorization
+ signed metadata binding
```

## Rollout

1. Deploy code với mode `audit`.
2. Cấu hình service ID và secret cho Gateway cùng internal services.
3. Theo dõi log `Legacy unverified metadata`.
4. Loại bỏ mọi gRPC client bypass interceptor.
5. Chuyển staging sang `enforce`.
6. Chạy smoke test public, authenticated và downstream calls.
7. Chuyển production sang `enforce` sau khi không còn legacy traffic.

## Điều kiện PASS

- Gateway ký mọi gRPC call qua shared interceptor.
- Internal clients chủ động truyền actor đều dùng shared interceptor.
- Server phân biệt `verified` và `legacy_unverified`.
- Enforce reject forged actor metadata.
- Signature bind full method và hết hạn theo policy.
- Secret chỉ đọc từ environment.
- Unit test và compile gate PASS.
