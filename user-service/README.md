# user-service

## Role

Profile, public Auth/OAuth, IAM, commercial catalog, checkout intent, subscription và entitlement.

## Runtime status

Canonical default Compose service; business owner của public auth/IAM và commercial entitlement.

## Entrypoint / lifecycle

`cmd/grpc/main.go`.

`make server` chạy native gRPC entrypoint.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

Nếu thay schema:

```bash
make migrateup
make migrate-version
make new_migration name=<schema_change>
```

Root có thể gọi `make migrate service=user-service`.

Database mới hoàn toàn cần đúng một root operator. Cấp input qua biến
`QHPRO_BOOTSTRAP_*` rồi chạy `make bootstrap-admin` ở repository root; command
ghi IAM atomically và từ chối chạy lại. Không đặt credential bootstrap trong
`.env`. Chi tiết nằm trong
[`../documents/THREE-SOURCE-ONBOARDING.md`](../documents/THREE-SOURCE-ONBOARDING.md).

Trong `QHPRO_ENVIRONMENT=development`, root `make up` còn khởi tạo idempotent
hai identity dùng để thao tác bằng UI thật: Admin `admin/admin123` và client
`0900000000/client123` (cũng có thể dùng username `client`). Credential được
hash trong PostgreSQL `user_service`; file khởi tạo nằm tại
`shared/code/development/identities.sql`, ngoài production migrations. Mỗi lần
provision ở development đưa hai mật khẩu về đúng giá trị công bố.

## Dependencies / durable state

PostgreSQL `user_service`, Redis DB 0, RPC clients. Canonical migration: `migrate/`.

## Read source from here

`cmd/grpc/`, `config/`, `infra/handler/grpc/`, `internal/`, `db/`, `migrate/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
