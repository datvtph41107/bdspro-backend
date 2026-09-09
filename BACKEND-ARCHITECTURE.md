# BDSPro Backend Architecture

## Architectural intent

BDSPro là nhiều application/service độc lập về business nhưng thống nhất về developer experience.

> **Độc lập ở business. Thống nhất ở cách developer làm việc.**

## Five ownership zones

```text
ROOT
  public workspace UX and system orchestration

SHARED/CODE
  repository engineering application

SERVICE
  product/application owner

SHARED RUNTIME
  reusable technical mechanics

CONTRACTS
  cross-owner communication truth
```

### Root

Root cố ý ít file có semantics toàn workspace: `README.md`, `Makefile`, `compose.yaml`, `go.work`, `.vscode/`, docs index. Root Makefile là front door và include implementation từ `shared/code/development/root.mk`.

### `shared/code`

Sở hữu engineering workflows: pinned toolchain, local config generation, native dev mechanics, protobuf/Wire generation, verification, integration orchestration helpers, deployment/release compatibility. Không sở hữu product business.

### Services

Mỗi service sở hữu business, runtime configuration semantics, database migration history, composition, API implementation, workers/actors và tests của mình.

User Service sở hữu identity/IAM nên cũng sở hữu one-time root-operator
bootstrap. Đây là một explicit CLI role của cùng service, không phải seed/startup
side effect hay service mới. Sau bootstrap, mọi administrator/role mutation đi
qua authenticated IAM API; `QHPRO_SYSTEM_ROOT` không thể gán, sửa hoặc xóa qua
API và được cấp toàn bộ permission identity hiện hành trong User catalog.

### Shared runtime

`shared/common`/`shared/base` chỉ nhận logic generic thực sự reusable. Không dùng shared để né ownership decision.

### Contracts

`shared/protobuf/schema` là communication source of truth; generated files là artifact.

## Source vs process vs container

Ba topology không đồng nghĩa:

```text
SOURCE OWNERSHIP ≠ PROCESS TOPOLOGY ≠ CONTAINER TOPOLOGY
```

Một Payment service có thể chạy API + outbox/fulfillment actors trong một hay nhiều process mà vẫn là một business owner. Migration container là one-shot process, không phải service.

## Development topology

Default local runtime:

```text
Postgres/Redis/RabbitMQ (Docker)
           ↑ host ports
all application services (native binaries)
           ↑ focused edit loop
one selected service (native Air in IDE terminal)
```

Image/package integration (explicit only):

```text
make integration-up → Compose full topology
→ health
→ smoke
→ E2E
```

Production packaging có thể dùng Docker images mà không buộc daily edit loop phải build image.

## Dependency rule

Business dependency phải hướng vào owner-local abstractions hoặc explicit contracts. Root/engineering không được trở thành global service locator/business kernel.

Boundary test: “Nếu Payment biến mất, file này còn có ý nghĩa không?”

- không → Payment;
- có và là generic runtime mechanic → shared runtime;
- có và thao tác/generate/build/verify repo → `shared/code`;
- là communication surface → contracts;
- là whole-workspace UX/topology → root.

## Correctness that must survive simplification

Không đánh đổi các invariant sau để lấy source đẹp hơn: int64 minor-unit money, DB transaction boundaries, idempotency/unique keys, webhook verification, durable outbox/inbox/jobs, quota correctness, fail-closed auth/authz, graceful shutdown, versioned migrations.
