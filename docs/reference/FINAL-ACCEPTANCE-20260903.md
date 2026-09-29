# BDSPro Final Development Candidate Acceptance — 2026-09-03

## Scope

Candidate này chuẩn hóa **source + developer/local operation**. Production
deployment mechanics bị đóng băng và không được thay đổi.

Baseline archive:

`bdspro-backend-current-20260903-173907.zip`

SHA-256 baseline: `6c6804e875d43fca07a5cb6744ba3d022948df3f0c88dec2997e48de040b5943`

## Implemented

### Developer operating model

- Root public commands: `setup/up/rebuild/status/logs/dev/test/migrate/smoke/test-e2e/verify/accept/down/reset`.
- `make up` reuse image hiện có; `make dev service=...` rebuild đúng owner; `make accept` ép full rebuild.
- VS Code tasks/workspace chỉ gọi public Make commands và dùng tên BDSPro.
- Fresh-source config được tái tạo bằng `make setup`/`configure.sh`; archive final không cần mang secret local.

### Service contract

- 18 Go service modules có README/Make developer contract rõ; `integration-test` có README riêng như cross-service evidence module.
- Generic Go mechanics nằm ở `shared/code/development/service.mk`.
- Generic migration mechanics nằm ở `shared/code/development/migration.mk`.
- Service-specific buf/wire/pmtiles/compatibility actor vẫn ở owner Makefile.
- Notification compatibility worker binaries vẫn build được nhưng default Compose không tạo standing worker container riêng.

### Migration

Canonical source authorities:

`User 4 / Organization 3 / Payment 7 / TQD 41 / Notification 2 / File 4 / Hub 3 / BDSPro 3 / CRM 25`.

Corrections:

- BDSPro `000003`: up/down cùng tên, rollback không còn file rỗng.
- CRM Payment Inbox: canonicalized thành `000025_payment_event_inbox`.
- CRM ad-hoc cleanup SQL được chuyển khỏi versioned migration directory.
- CRM gRPC startup không còn gọi empty `AutoMigrate()`.
- Local PostgreSQL bootstrap có `db_bdspro` và `db_crm` cho source-health migration work.
- `verify-migrations.sh` enforce version continuity + exact up/down pair/name.

### Repository/docs cleanup

- Root chỉ giữ canonical control-plane files.
- Historical TQD/repository/deployment docs được giữ trong `documents/history/`.
- Legacy proto/manual scripts được giữ trong `shared/code/legacy/` và `shared/code/manual/`.
- Canonical docs: operating guide, architecture standard, service map, migration guide, acceptance.
- Old service README content được lưu dưới `<service>/docs/history/README.previous.md` khi có.
- Legacy Cursor rules được chuyển vào `documents/history/editor-rules/`; chỉ `/.cursor/rules/bdspro-development.mdc` còn active.

## Static proof rerun on rebuilt candidate

`PASS`:

- development shell syntax
- root/service env ownership
- runtime config isolation
- migration contract 9/9 authorities
- source layout guard
- backend owner/source architecture guard
- AI Python syntax: 16 files
- service Make parse: 18 Go services
- canonical markdown links: 31 docs
- YAML parse: 19 files, Compose: 22 entries
- workspace/tasks JSON
- Go parser syntax: 3653 files
- fresh-source config lab from 0 env files
- deployment frozen-boundary byte comparison

## Deployment frozen-boundary proof

The following files are byte-for-byte unchanged from baseline:

- `ai-service/deploy.sh`
- `ai-service/docker-compose.yml`
- `shared/code/deploy.sh`
- `shared/code/docker-compose.yml`

`shared/code/deploy.sh` SHA-256:

`80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`

## Runtime proof not claimed in this build environment

This execution environment has no Docker CLI/daemon and only system Go 1.23.2,
while the repository pins Go 1.25.10. Therefore this report does **not** claim
that `go test`, race, canonical build, Docker health, smoke or E2E were rerun
here after the final changes.

After import into WSL/local with Docker, run:

```bash
make setup
make verify
make rebuild
make status
make smoke
make test-e2e
make accept
```

Only when `make accept` passes on that machine should runtime acceptance be
marked CLOSED.
