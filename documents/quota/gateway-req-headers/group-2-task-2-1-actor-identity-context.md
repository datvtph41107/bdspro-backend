# Group 2 — Task 2.1: Actor Identity Context

## Mục tiêu

Tạo một actor identity thống nhất cho mỗi request đã có JWT hợp lệ:

```text
JWT verified at Gateway
→ Actor assertion
→ gRPC metadata
→ canonical service context
→ downstream gRPC
```

Task này chưa xác định quota owner và chưa chứng minh metadata đến từ một internal caller đáng tin cậy.

## Vấn đề source hiện tại

Identity đang tồn tại dưới nhiều dạng song song:

- `jwt.Principal` tại Gateway.
- Các metadata rời: `profileId`, `originId`, `authId`, `organizationId`.
- Các context key rời trong service.
- `PrincipalFromContext` đọc lại raw metadata.
- Quota cũ tự chọn ID đầu tiên có giá trị.

Điều này làm mất một nguồn sự thật và không phân biệt actor với quota subject.

## Contract

Actor chứa:

```text
auth_id
profile_id
origin_id
organization_id assertion
session_id
role
token_type
```

Actor không chứa:

```text
plan_id
plan_from
quota subject
quota limit
operation code
billing information
```

`plan_id` thuộc entitlement/subscription, không thuộc identity.

## Ownership

```text
JWTAuthMiddleware
→ xác thực JWT một lần tại HTTP boundary
→ lưu verified Principal trong request context

InjectGrpcMetadataContextMiddleware
→ tạo Actor từ verified Principal
→ serialize Actor thành gRPC metadata

ParseGrpcMetadataContextMiddleware
→ validate metadata shape
→ lưu Actor vào identitycontext
→ project sang legacy context keys

GrpcClientMetadataFromContext
→ ưu tiên Actor canonical
→ fallback legacy keys để tương thích
```

## Trust boundary

Actor trong service hiện là một transported assertion.

Nó chưa tự chứng minh rằng caller là Gateway hoặc internal service hợp lệ. Direct gRPC vẫn cần network isolation, mTLS hoặc signed internal identity envelope ở Task 2.2.

Do đó chưa được dùng Actor này để:

- Chọn quota subject.
- Miễn quota.
- Quyết định billing.
- Tin organization membership.
- Tin API key ownership.

## Compatibility

Các context key cũ vẫn được project từ Actor để business handler hiện tại không bị phá:

```text
ProfileIDKey
OriginIDKey
OrganizationIDKey
AuthIDKey
SessionKey
RoleKey
TypeKey
```

Code mới phải đọc `identitycontext.FromContext`. Legacy helpers được giữ cho migration từng bước.

## Failure rules

- Không có actor metadata: xử lý như anonymous/public request.
- Public route có JWT sai hoặc hết hạn: không tạo Actor từ claims chưa xác thực.
- Actor metadata có nhiều giá trị cho cùng field: `Unauthenticated`.
- Actor numeric field sai định dạng: `Unauthenticated`.
- Actor string field có control character hoặc vượt 128 byte: `Unauthenticated`.

## Điều kiện PASS

- Gateway tái sử dụng Principal đã được JWT middleware xác thực.
- Actor có một package context duy nhất.
- Service lưu actor canonical trước business handler.
- Downstream ưu tiên actor canonical.
- Legacy context vẫn hoạt động.
- `planId` không nằm trong Actor.
- Quota middleware chưa sử dụng Actor.
