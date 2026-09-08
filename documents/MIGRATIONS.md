# Migration Guide

## Authorities

```text
user-service/migrate                         000001..000009
organization-service/migrate                 000001..000003
payment-service/migrate                      000001..000007
tqd-service/migrate                          000001..000044
notification-service/migrate                 000001..000002
file-service/migrate                         000001..000004
hub-service/migrate                          000001..000003
bdspro-service/infra/db/migrations/v2        000001..000003
crm-service/infra/db/migrate_v2              000001..000025
```

Mỗi version có đúng một cặp cùng tên:

```text
NNNNNN_name.up.sql
NNNNNN_name.down.sql
```

## Commands

```bash
make verify-migrations
make migrate service=payment-service
make migrate
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
