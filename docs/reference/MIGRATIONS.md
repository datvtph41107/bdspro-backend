# Migration Guide

## Authorities

Mỗi service Makefile là owner của `MIGRATION_DIR`; danh sách dưới đây chỉ là developer-facing projection của owner đó. Không lặp lại version/count ở đây; `.github/workflows/acceptance-source-integrity.yml` kiểm inventory thật để tránh tạo authority thứ hai.

```text
internal/user/database/migrations
internal/organization-legacy/migrate
internal/payment/database/migrations
internal/planning/database/migrations
internal/notification/database/migrations
internal/file/database/migrations
internal/hub/database/migrations
internal/property/database/migrations
internal/crm/database/migrations
```

Mỗi version có đúng một cặp cùng tên:

```text
NNNNNN_name.up.sql
NNNNNN_name.down.sql
```

`internal/organization-legacy/migrate` là authority hiện hữu của Organization và không thuộc batch canonicalization này.

## Commands

```bash
make verify-migrations
make migrate service=payment-service
make migrate
make -C internal/payment migration name=<schema_change>
make -C internal/payment migrate
make -C internal/payment migration-version
make -C internal/payment rollback
```

`make migrate` chạy mọi source-level authority. BDSPro/CRM không phải standing
service trong default Compose, nhưng local PostgreSQL init vẫn tạo database rỗng
`db_bdspro` và `db_crm` để developer có thể chạy history của chúng.

## Rules

- không renumber migration history;
- không để ad-hoc SQL trong canonical version directory;
- recovery/hotfix SQL lịch sử phải nằm trong docs/history owner và chạy có chủ đích;
- không sửa migration đã được dùng ở môi trường cần compatibility nếu chưa có
  compatibility proof;
- serving process không được biến migration thành hidden side effect mặc định.
