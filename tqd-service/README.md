# tqd-service

## Role

QHPro/planning capability, quota, durable usage và report acceptance/job/rendering.

## Runtime status

Canonical default Compose service. Report actor thuộc owner process; PMTiles là special mode có chủ đích.

## Entrypoint / lifecycle

`main.go`; mode mặc định `-server=grpc`, optional `-server=pmtiles`.

`make server`; `make pmtiles` cho tile server.

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

Root có thể gọi `make migrate service=tqd-service`.

## Dependencies / durable state

PostgreSQL `qhpro_tqd`, Redis DB 1, User/Auth/Assistant/File. Canonical migration: `migrate/`.

Generated Report chấp nhận Report + UsageEvent + Job trước khi worker render.
“Tạo lại” đưa chính job đã chấp nhận về hàng đợi, không chỉ sửa trạng thái hiển
thị và không ghi thêm usage. Renderer đọc effect đã commit từ File trước khi
render/upload để retry không ghi đè PDF có byte khác. Khi rollout, triển khai
endpoint lookup của File trước TQD sử dụng contract này.

Metadata thửa thuộc `qh_parcel_info`, geometry/ID thuộc `parcels`; overview và
detail dùng projection chung trong `infra/postgres/parcel_info_projection.go`.
Không thêm cột metadata song song vào `parcels` để chữa query sai owner.

Root `make runtime-value-chain-acceptance` kiểm tra thêm vòng đời gói:
QR, webhook/replay, quota exhaustion, recovery/share PDF, upgrade giữ kỳ/usage,
cancel và repurchase. Output JSON/PDF ở `.tmp/development/acceptance-current/`.
Đây là webhook mô phỏng local; kết quả GIS thành công cần dữ liệu thửa/vùng thật.

## Read source from here

`cmd/`, `internal/`, `infra/`, `config/`, `migrate/`, `template/`.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).

Chuẩn triển khai operation và comment nghiệp vụ: [`../documents/QUOTA-DEVELOPMENT-STANDARD.md`](../documents/QUOTA-DEVELOPMENT-STANDARD.md).
