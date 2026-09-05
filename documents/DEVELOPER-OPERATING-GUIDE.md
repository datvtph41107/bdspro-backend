# QHPRO Developer Operating Guide

## Workspace

Ba repository là ba Git source độc lập nhưng được mở trong một IDE workspace:

```text
QHPRO
├── BACKEND
├── ADMIN
└── CLIENT
```

Backend phải ready trước khi chạy hai UI.

## First-day flow

```bash
git clone <backend-repository> bdspro-golang-microservices-r2f0-i21
git clone <admin-repository> react-js-admin-bds
git clone <client-repository> react-native-quy-hoach

cd bdspro-golang-microservices-r2f0-i21
make bootstrap
make dev-up
make dev-status
make dev-smoke
```

`configure` là idempotent và không ghi đè file hiện có. `bootstrap` dùng
version/checksum trong `tools/development/toolchain.versions`.

## Manual IDE flow

```text
Ctrl+Shift+P
-> Tasks: Run Task
-> QHPRO · Configure
-> QHPRO · Doctor
-> Backend · Up
-> Backend · Status
-> Backend · Smoke
```

Các terminal được tách thành:

```text
CONTROL       lệnh ngắn, status, smoke, verify
BACKEND       build/start output
BACKEND LOGS  log theo dõi dài hạn
ADMIN         Vite
CLIENT        React Native/Web
```

## Quan sát runtime

```bash
make dev-status
make dev-logs
make dev-logs SERVICE=payment ROLE=api
make dev-logs SERVICE=payment ROLE=outbox-publisher SINCE=10m
make dev-logs SERVICE=notification ROLE=payment-consumer FOLLOW=0
```

Docker Desktop hiển thị một project `qhpro-backend`. Image được đặt theo owner
(`qhpro/payment`, `qhpro/notification`); container được đặt theo process role
(`payment`, `payment-publisher`, `notification-payment-worker`).
Compose giới hạn mỗi container ở 5 file log, tối đa 10 MiB/file; ứng dụng vẫn
chỉ ghi stdout/stderr để image không phụ thuộc môi trường chạy.

| Business service | Role | Compose process |
| --- | --- | --- |
| Payment | `api` | `payment` |
| Payment | `outbox-publisher` | `payment-publisher` |
| Notification | `api` | `notification` |
| Notification | `payment-consumer` | `notification-payment-worker` |
| Notification | `delivery-worker` | `notification-delivery-worker` |

Runtime log phải có service, role, operation, outcome và correlation identity
khi request có cung cấp. Không log password, token, signing key hoặc DSN chứa
credential.

## Daily loop

```bash
make dev-up
make dev-status
make -C user-service test
make dev-smoke
make verify-backend
make dev-down
```

`dev-down` giữ dữ liệu. Reset là thao tác riêng:

```bash
make dev-reset
```

Non-interactive automation phải xác nhận rõ:

```bash
make dev-reset CONFIRM=reset-development-data
```

## Production-like boundary

Development và production dùng cùng domain, use case, contract, migration,
database constraint, handler, adapter và process lifecycle. Những input được
phép khác là endpoint, credential, replica count, resource và telemetry sink.

Code và default `config/runtime.yml` không đổi theo môi trường. Host/container
khác nhau ở biến wiring và secret được inject; root `compose.yaml` là một cách
cung cấp các input đó cho development, không phải một phiên bản business khác.

Không đặt development account, fake settlement, quota grant hoặc report result
trong product UI/serving code. Test provisioning thuộc integration tooling và
phải đi qua public business boundary nơi cần chứng minh behavior.

## Troubleshooting order

1. `make doctor`
2. `make config-check`
3. `make dev-status`
4. Xem log đúng service/role.
5. Lấy request ID, command key hoặc aggregate ID.
6. Đi theo owner thay vì restart toàn stack.
7. Sau fix, thêm regression proof hoặc runbook.
