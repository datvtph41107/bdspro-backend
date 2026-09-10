# BDSPro Backend Acceptance

Acceptance được phân tầng để developer không phải trả chi phí full-system cho mọi edit.

## 1. Service proof

```bash
make test service=payment-service
```

Trả lời: code/business của owner có đúng trong phạm vi service không?

## 2. Repository proof

```bash
make verify
```

Bao gồm generation consistency, config/layout/migration guards, docs/source checks, Go vet/test/race/build và whitespace guard. Gate này không được hiểu là runtime acceptance.

## 3. Runtime smoke

```bash
make up
make smoke
```

Trả lời: hạ tầng Docker và toàn bộ application process native sống, public
boundary cơ bản hoạt động?

Trên fresh database, trước lần đăng nhập Admin đầu tiên, operator phải chạy
`make bootstrap-admin` với input `QHPRO_BOOTSTRAP_*` như tài liệu onboarding.
Proof hợp lệ phải cho thấy: lần đầu tạo đủ profile/credential/root role trong
một transaction, đăng nhập qua Gateway thành công, và lần gọi bootstrap thứ hai
bị từ chối. Không thay proof này bằng acceptance fixture.

## 4. Cross-service E2E

```bash
make test-e2e
```

Trả lời: customer/admin business value chain chạy qua các owner thật?

## 5. Final candidate

```bash
make accept
```

`accept` là composite gate: source proof + native runtime verification + smoke,
RabbitMQ/Notification recovery và customer/admin E2E. Image/package proof có
target riêng `make docker-backend` hoặc `make integration-up`; nó không làm thay
đổi cách developer chạy business code.
Chỉ khi command trở về shell với không `FAIL`/`Error` mới được đánh dấu local
acceptance CLOSED.

Recovery gate cố tình restart RabbitMQ và chỉ PASS khi cùng Notification native
PID vẫn sống và consumer mở lại named AMQP connection. Value-chain chạy ngay
sau đó chứng minh event durable tiếp tục đi qua Inbox và notification effect.

## Migration rule

Default native runtime chạy `migrate` CLI đã pin trực tiếp từ service owner trước
khi start. Các `<service>-migrate` trong Compose chỉ là one-shot executor cho
image/package integration. Canonical history vẫn nằm ở service owner và được
guard bằng:

```bash
make verify-migrations
```

## Artifact rule

Mọi artifact ZIP final phải có `SOURCE-MANIFEST.sha256` được regenerate sau cùng và checksum ZIP bên ngoài. Nếu source đổi sau khi đóng gói, checksum/manifest cũ không còn đại diện candidate mới.
