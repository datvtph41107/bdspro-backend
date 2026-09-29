# Change Summary — 2026-09-04

This candidate keeps the accepted commercial theorem and also fixes runtime
gaps exposed by acceptance rather than hiding them in fixtures or UI behavior.

## Structural/operating changes

- Root `Makefile` reduced to a stable front door; implementation moved to `shared/code/development/root.mk`.
- Root `tools/` removed. Development, legacy and manual engineering utilities now have one physical owner under `shared/code/`.
- All service Makefiles consume the engineering contract from `shared/code/development/`.
- `make dev service=X` now means native/Air daily development, not Docker image rebuild.
- Added `make deps [service=X]` / `deps-down` to make backing services explicit.
- `make rebuild service=X` is the explicit integration-image refresh; `make rebuild` remains full candidate rebuild.
- VS Code tasks restore the familiar `Run Task → Dev · Service` workflow.
- Notification generated env now explicitly includes the published business-event exchange contract discovered by real local acceptance.
- Canonical docs rewritten around owner → capability → flow → proof and daily-vs-integration separation.

## Runtime corrections proved on 2026-09-04

- User owns Admin User credentials/profile lifecycle. Password material is
  write-only at the provider boundary, preventing credential-hash disclosure
  and accidental double hashing during profile linkage.
- Optional CRM origin and history recording have bounded login budgets; their
  outage cannot consume the whole authentication request deadline.
- User Admin IAM is checked at User ingress with durable `USER_ADMIN_VIEW` /
  `USER_ADMIN_MANAGE` permissions.
- Notification owns a reconnecting Payment-event consumer supervisor. A real
  RabbitMQ restart kept the Notification API healthy, reconnected the consumer
  and was followed by a passing full value-chain.
- Admin User CRUD uses the canonical User contract and commercial User detail
  reads Subscription, Payment, Usage/Quota, Report/Job/File evidence from their
  respective owners.

## Intentionally unchanged

- Product business behavior.
- Service ownership boundaries.
- Protobuf schemas/contracts other than regeneration mechanics.
- Migration histories.
- Dockerfiles/runtime image semantics.
- Existing full Compose/E2E acceptance logic.
