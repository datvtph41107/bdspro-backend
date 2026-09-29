# BDSPro Backend

BDSPro Backend là workspace Go multi-service. QHPro chỉ là capability quy hoạch thuộc `tqd-service`; nó không phải namespace toàn hệ thống.

Mục tiêu của repository là tạo một **golden path** ổn định: developer mới clone source, chuẩn bị máy, chạy đúng service, thay đổi business và chứng minh thay đổi mà không cần hiểu toàn bộ Docker/topology trước.

## Quick start

Yêu cầu host: Linux/WSL2, Git, Make, Python 3, Ripgrep (`rg`), FFmpeg, WeasyPrint và Docker
Engine/Compose cho PostgreSQL/PostGIS, Redis và RabbitMQ. `make doctor` xác nhận
cả toolchain lẫn runtime dependency trước khi start.

```bash
make help
make setup
make doctor
make up
make status
make logs service=gateway-service
```

Development login sau `make up`: Admin `admin/admin123`; Client
`0900000000/client123` (hoặc username `client`). Mật khẩu được hash trong
database User; không có authentication bypass trong application runtime.

`make setup` cài WeasyPrint đã pin vào repository toolchain ở user cache, không
ghi vào Python hệ thống. Nếu `doctor` báo thiếu sau một clone cũ, chạy lại:

```bash
make setup
```

Daily loop cho Payment:

```bash
make deps service=payment-service
make dev service=payment-service
# sửa/save source; Air build + restart native process
make test service=payment-service
```

`make up` build rồi chạy toàn bộ application binary trực tiếp trên WSL/Linux;
PID và log nằm trong ignored `.tmp/development/`. `make dev` chuyển đúng service
đang sửa từ SUPERVISED sang DEV/Air trong terminal hiện tại. `make status` hiển
thị owner `SUPERVISED / DEV / FOREIGN / DOWN` và readiness riêng, nên port mở
không còn bị coi là bằng chứng ownership. Docker chỉ chạy PostgreSQL/PostGIS,
Redis và RabbitMQ trong golden path này.

Proof toàn hệ thống dùng chính các process native mà developer đang debug:

```bash
make up
make status
make smoke
make test-e2e
```

Chỉ khi cần kiểm tra image/package Compose:

```bash
make rebuild service=payment-service   # một service
make rebuild                           # toàn Compose candidate
make integration-up
make integration-status
```

Candidate cuối:

```bash
make verify
make accept
```

## Mental model

```text
TASK
  ↓
OWNER SERVICE
  ↓
BUSINESS CAPABILITY
  ↓
HANDLER → BUSINESS/USECASE → STORE/CLIENT → STATE/EFFECT
  ↓
SMALLEST SUFFICIENT PROOF
```

Developer không bắt đầu từ container. Developer bắt đầu từ owner của business.

## Source ownership

| Vùng                             | Sở hữu                                                                             |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| Root | workspace UX: `README`, public `Makefile`, `compose.yaml`, `go.work`, IDE config |
| `api/gateway/` | external HTTP/API edge source; không phải business authority |
| `internal/*` | current business/runtime implementation areas; source path không tự định nghĩa network boundary |
| `infrastructure/runtime/`, `infrastructure/base/` | reusable technical/runtime mechanics |
| `proto/` | communication contracts + generated contract code |
| `tools/development/` | repository engineering application: setup/dev/generate/build/verify/release |
| `integration-test/` | cross-owner runtime proof |

Nguyên tắc: **Root tổ chức hệ thống. Business/runtime source nằm dưới `internal/` hoặc `api/`. Infrastructure chỉ sở hữu mechanics. `proto/` sở hữu wire contracts. `tools/development/` sở hữu engineering workflow. Runtime service identity và source path là hai khái niệm độc lập.**

## Public command API

```text
make setup                         config local + pinned toolchain
make deps                          Postgres + Redis + RabbitMQ
make deps service=payment-service dependencies tối thiểu + migration của owner
make dev service=payment-service  native Air development
make test service=payment-service service proof
make migrate service=payment-service explicit schema operation
make generate                      protobuf/Wire generation

make up                            infra Docker + toàn bộ application native
make rebuild service=...          rebuild một integration image
make rebuild                       rebuild toàn integration stack
make status                        infra health + runtime owner/readiness
make logs service=...              native stdout/stderr trong .tmp/development
make smoke                         public boundary smoke
make test-e2e                      cross-service business proof
make integration-up                explicit full Compose image/package proof
make integration-down              dừng full Compose image/package topology

make verify                        source/code/repository proof
make accept                        rebuild + runtime + E2E acceptance
make down                          dừng native processes + infra, giữ data
make reset                         xóa local volumes có xác nhận
```

Các target dài bên trong là implementation/compatibility API và không phải vocabulary developer phải học thuộc.

## Cách đọc một service

Ví dụ nhận task Payment settlement:

```text
internal/payment/README.md
  ↓
main.go / cmd/<role>/main.go
  ↓
config/
  ↓
settlement/order/payment capability liên quan
  ↓
handler → usecase/business → store/client
  ↓
database/migrations/ nếu state schema đổi
  ↓
tests
```

Một service có thể sở hữu API, worker, outbox actor hoặc CLI role cùng lúc. **Process topology không tự động tạo thêm business/service owner.**

## Docker và migration

Development và local E2E:

```text
Docker: PostgreSQL/PostGIS + Redis + RabbitMQ
WSL/Linux: toàn bộ Go services; service đang sửa có thể chạy bằng Air
```

Packaging/image integration tường minh:

```text
make integration-up: full application Compose topology
```

Các `*-migrate` container trong Compose là **one-shot jobs**: chờ database, chạy migration của service owner, exit. Chúng không phải microservice riêng.

## Documentation map

- [`CONTRIBUTING.md`](CONTRIBUTING.md) — luật thay đổi source/PR.
- [`BACKEND-ARCHITECTURE.md`](BACKEND-ARCHITECTURE.md) — ownership và dependency boundaries.
- [`BACKEND-ACCEPTANCE.md`](BACKEND-ACCEPTANCE.md) — proof ladder và final gate.
- [`tools/development/README.md`](tools/development/README.md) — engineering application.
- [`docs/reference/DEVELOPER-OPERATING-GUIDE.md`](docs/reference/DEVELOPER-OPERATING-GUIDE.md) — clone → first change → integration.
- [`docs/reference/DEVELOPMENT-ARCHITECTURE-STANDARD.md`](docs/reference/DEVELOPMENT-ARCHITECTURE-STANDARD.md) — quy chuẩn dài hạn.
- [`docs/reference/SERVICE-MAP.md`](docs/reference/SERVICE-MAP.md) — service owners/entrypoints.
- [`docs/reference/MIGRATIONS.md`](docs/reference/MIGRATIONS.md) — migration authorities.
- [`docs/reference/THREE-SOURCE-ONBOARDING.md`](docs/reference/THREE-SOURCE-ONBOARDING.md) — fresh clone → config → ba runtime → quota proof.
- [`docs/reference/QHPRO-COMMERCIAL-OPERATIONS.md`](docs/reference/QHPRO-COMMERCIAL-OPERATIONS.md) — vận hành Admin/Client từ gói đến quota và report.
- [`docs/reference/FINAL-ACCEPTANCE-20260904.md`](docs/reference/FINAL-ACCEPTANCE-20260904.md) — trạng thái proof của artifact này.
