# Migration Guide

## Authorities

Mỗi service Makefile là owner của `MIGRATION_DIR`; danh sách dưới đây chỉ là developer-facing projection của owner đó. Không lặp lại version/count ở đây; `.github/workflows/acceptance-source-integrity.yml` kiểm inventory thật để tránh tạo authority thứ hai.

```text
user-service/database/migrations
organization-service/migrate
payment-service/database/migrations
tqd-service/database/migrations
notification-service/database/migrations
file-service/database/migrations
hub-service/database/migrations
bdspro-service/database/migrations
crm-service/database/migrations
```

Mỗi version có đúng một cặp cùng tên:

```text
NNNNNN_name.up.sql
NNNNNN_name.down.sql
```

`organization-service/migrate` là authority hiện hữu của Organization và không thuộc batch canonicalization này.

## Commands

```bash
make verify-migrations
make migrate service=payment-service
make migrate
make -C payment-service migration name=<schema_change>
make -C payment-service migrate
make -C payment-service migration-version
make -C payment-service rollback
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
