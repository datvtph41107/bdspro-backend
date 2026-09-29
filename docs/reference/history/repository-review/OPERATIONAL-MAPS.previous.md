# BDSPro Operational Maps

Tài liệu này ghi current truth ngắn gọn để việc đơn giản hóa không làm mất
runtime hoặc deployment compatibility.

## 1. Developer workflow

### Public workflow

```text
clone
  -> make setup
  -> make up
  -> make status
  -> make smoke
  -> make logs service=<name>-service
  -> make test service=<name>-service
  -> make test-e2e
  -> make down
```

Nguồn sự thật:

| Concern | Owner |
| --- | --- |
| Public command names | root `Makefile` |
| Development topology | root `compose.yaml` |
| Service build/run/test | `<service>/Makefile` trong thời gian chuyển đổi |
| IDE commands | `.vscode/tasks.json`, chỉ gọi root Makefile |
| Tool versions | `shared/code/development/toolchain.versions` |
| Local runtime values | ignored root/service `.env` |

Canonical local infrastructure hiện là một PostgreSQL/PostGIS, một Redis và
một RabbitMQ. PostgreSQL vẫn có database riêng theo owner; Redis dùng DB index
riêng (User 0, TQD 1, Notification 2, Hub 3). Đây là gom physical dependency,
không gom data/business ownership.

Các target dài còn lại là implementation/compatibility, không phải vocabulary
mà developer phải học.

## 2. Runtime map

### Business request flow

```text
Admin / Mobile
    -> Gateway
        -> User: identity, IAM, catalog, subscription, entitlement
        -> Organization: organization membership and action authorization
        -> Payment: order, attempt, webhook, settlement
        -> TQD: QHPro, quota, usage, report
        -> File: report artifact and signed access
        -> Notification: customer-visible notification
```

### Durable commercial flow

```text
Login
  -> Catalog
  -> Checkout
  -> Payment attempt
  -> Provider webhook
  -> Settlement
  -> Subscription
  -> Entitlement
  -> Quota reserve/consume
  -> Report job
  -> PDF/File
  -> Usage
  -> Notification
```

### Current process exceptions

| Owner | Current processes | Target until contrary evidence exists |
| --- | --- | --- |
| Payment | API + outbox publisher + fulfillment actor trong cùng process | giữ một standing Payment container |
| Notification | API + event consumer trong cùng process; delivery cùng process khi Firebase có cấu hình | giữ một standing Notification container |
| TQD | API and report actor in service lifecycle | keep together; reassess PDF isolation only from measured failures/resources |
| Auth | permission snapshot/compatibility proxy | move authorization capability to User after callers reach zero |

Migration and acceptance fixture containers are one-shot operations, not
business services. They must disappear from the steady-state mental model.

## 3. Production deployment map

Current production semantics discovered in `shared/code/deploy.sh`:

```text
operator selector
  -> map selector to service directory / Compose key
  -> build Linux binary locally
  -> copy binary + config + Dockerfile to a fixed SSH host/path
  -> remote docker compose build
  -> remote docker compose up selected service
```

The following are compatibility facts, not approved target design:

- Operators deploy one named service at a time.
- Historical selectors use short names such as `payment`, `org` and `tqd`.
- Compose service keys and deployable directory names are not consistently the
  same.
- Some owner services currently have multiple runtime components.

Deployment gaps that must be resolved before production cutover:

| Gap | Evidence | Required outcome |
| --- | --- | --- |
| Server identity in source | `REMOTE_HOST` is hard-coded | host/user/path supplied by operator environment |
| Stale config | script copies deleted `config/develop.yml` | runtime config stays on deployment host; safe defaults come from image |
| Secret transfer | Firebase credential filename is copied | secret provisioned out of repository and never copied by generic deploy |
| Compose ambiguity | script references `docker-compose.yml`; canonical development uses `compose.yaml` | production Compose path is explicit and independently owned |
| Missing deployables | selectors reference membership/appointment/marketing/transaction etc. absent from current source | compatibility matrix proves new owner or keeps legacy deployment |
| Version trace | binary/image has no immutable source version | running image/artifact reports a commit-derived version |
| Health | restart success is treated as deploy success | readiness and business smoke must pass |
| Rollback | server rebuild overwrites the prior artifact | previous version remains selectable |

The existing deploy script is intentionally not rewritten into a remote
migration/recreate tool until the production Compose and environment contract
are supplied and tested. Adding a more powerful SSH script without those facts
would increase production risk.

## 4. Migration order

1. Keep old deployment selector semantics.
2. Simplify README, root Make commands and IDE tasks.
3. Make Payment the readable reference service.
4. Co-locate lightweight Payment/Notification workers with their owner.
5. Local infrastructure đã được simplify; giữ logical database ownership và old volumes cho rollback.
6. Move authorization callers from Auth to User; retire Auth only at zero.
7. Repeat the convention for TQD, File, Organization and Gateway.
8. Build the production compatibility matrix for endpoints, tables, events,
   clients and deployments.
9. Run the complete commercial E2E before each runtime cutover.

Fresh workspace, fresh database và recovery proof đã PASS ngày 2026-09-03.
Production parity/legacy retirement vẫn là gate riêng; xem
[`PRODUCTION-COMPATIBILITY.md`](PRODUCTION-COMPATIBILITY.md).

## 5. Invariants that simplification cannot remove

- Money remains integer minor units.
- Provider webhook evidence is verified and durable.
- Commands are idempotent and conflicts are explicit.
- Settlement, outbox and fulfillment are recoverable.
- Quota and usage have one durable owner.
- Report acceptance and job creation remain transactional.
- Inbox prevents duplicate notification effects.
- Authorization fails closed when its authority is unavailable.
- Database evolution uses versioned migrations.
- Every process has bounded cancellation and resource close ownership.
