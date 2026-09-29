# Organization và Map owner cutover

Trạng thái được ghi từ working tree ngày 2026-09-06. Tài liệu này phân biệt
owner canonical với source compatibility để không xóa dữ liệu hoặc contract chỉ
vì một service không nằm trong Compose mặc định.

## Quyết định

| Capability | Canonical owner | Quyết định |
| --- | --- | --- |
| Identity, đăng nhập, profile, IAM toàn hệ thống | User | Không chuyển sang CRM hoặc Auth |
| Organization identity, membership, role và permission trong tổ chức | User | Đang chuyển theo từng contract; Admin directory và commercial checkout đã cut over |
| Deal, investment và group collaboration hiện hữu | CRM | Chuyển nguyên trạng theorem theo từng runtime path trước khi retire Organization |
| Branch/access thuộc danh tính tổ chức | User | Không trộn với CRM sales pipeline |
| Authentication compatibility/snapshot | Auth | Không nhận organization membership |
| Contact, friend và public/SEO content | CRM | CRM có thể dùng membership projection từ User nhưng không quyết định quyền tổ chức |
| Named spatial map points | TQD | Đã cut over khỏi Map Service qua `MapPointService` |

`organization-service` là source migration/compatibility, không còn là owner
đích. Nó chưa đủ điều kiện xóa vật lý vì Gateway còn public RPC surface, CRM
còn gọi membership và database cũ còn hơn 20 bảng durable. Việc xóa chỉ diễn ra
sau khi từng capability đã có new owner, dữ liệu đã được đối soát và mọi caller
về zero; không xóa directory để tạo cảm giác kiến trúc đã đơn giản.

## Organization cutover đã có proof

Admin organization directory hiện đi theo một đường duy nhất:

```text
Admin /v2/user/admin/organizations
  -> Gateway trusted identity
  -> User AdminOrganizationService
  -> User organizations + organization_members
```

Các thao tác list, create, get, update và archive đã chạy thật qua Gateway.
Admin UI đã bỏ CRUD giả và chuyển adapter khỏi `/v2/org/organization/...`.
Commercial checkout cũng không còn gọi Organization RPC: User kiểm tra member
active và role `owner/admin/billing` trong chính durable owner của nó.
CRM friend và contact organization-member projection đã đổi sang signed internal
`User.InternalOrganizationMembershipService.CheckMembers`; response nội bộ mang
membership identity, canonical role và status để CRM chỉ map sang response legacy.
Organization client của CRM chỉ còn phục vụ deal/group compatibility chưa migrate.
Luồng login/switch organization của User cũng đọc trực tiếp membership active
từ User database; nó không còn nhờ Organization RPC quyết định có được phát
organization context trong token hay không.
`GET /v2/user/profile/accounts/me`, là projection để client chọn tài khoản cá
nhân/tổ chức, cũng đã chuyển từ `GetOrganizationsByUserId` RPC sang
`User.OrganizationService.ListForProfile`; lỗi authority không còn bị nuốt.

IAM role/permission cũng đã được khóa tại User. Database sạch seed
`IAM_ROLE_VIEW` và `IAM_ROLE_MANAGE`; toàn bộ Role handler dùng hai key này.
Contract `POST /v2/auth/role/permissions` đã bổ sung `body: "*"` tại protobuf
owner rồi regenerate, vì contract cũ khiến grpc-gateway bỏ request body và gửi
`roleId=0`. Nghiệm thu hiện đọc role cùng 15 permission và replace cùng tập ID
qua Gateway thành công.

Evidence ngày 2026-09-06:

```text
Organization CRUD proof:
create -> update -> get -> list -> archive -> list(0)  PASS

Admin browser path:
Vite :5173 -> Gateway :8000 -> User :8201
QHPRO_DEMO_ORGANIZATION, verified, total_members=1  PASS

Regression:
make verify-backend  PASS
make smoke           PASS
make test-e2e        PASS
profile=21, order=42, report=47, usage.used=1

Admin:
architecture guard   PASS
TypeScript boundary  PASS (33 owned entry files)
Vite production build PASS (10,551 modules)

IAM role permission matrix:
login -> list roles -> read 15 permission IDs -> replace through Gateway PASS
```

Phần denominator còn lại trước khi xóa service:

1. Public member/role/permission/branch API chuyển sang User.
2. Deal/investment/group API và bảng dữ liệu chuyển sang CRM.
3. Deal/group enrichment giữ tại CRM sau khi dữ liệu được chuyển; organization
   membership enrichment đã chuyển sang User.
4. Gateway không còn đăng ký Organization public handlers.
5. Mobile/Admin không còn `/v2/org`; production data count/checksum PASS.
6. Bỏ Organization khỏi Wire client graph, config, Compose/deploy và `go.work`.

## MapPoint runtime cutover

Contract trình duyệt được giữ nguyên:

```text
/v1/map/locations/*
        -> Gateway auth + generated grpc-gateway
        -> TQD MapPointService
        -> TQD map_points
```

Map không còn nằm trong danh sách HTTP reverse proxy của Gateway và không là
standing container trong `compose.yaml`. Admin sử dụng URL `/v1/map/locations`
để local Vite và production reverse proxy có cùng contract.

Runtime proof đã chạy qua Gateway với một record test và PASS theo thứ tự:

```text
create -> get -> nearby -> polygon -> update -> delete
```

Record test được xóa ngay cuối proof. Migration `000042_add_map_points` là owner
schema mới cho fresh database.

## Production data gate của Map Service

Không xóa thư mục `map-service/` trước khi operator hoàn tất cả các bước sau:

1. Dừng write vào Map Service cũ trong maintenance window.
2. Export `id`, `name`, geometry SRID 4326, `created_at`, `updated_at` từ
   `db_map.map_points`.
3. Import idempotent vào TQD `map_points`; reset sequence cao hơn `MAX(id)`.
4. So sánh row count, null/invalid geometry count và checksum theo `id`.
5. Chạy lại create/get/nearby/polygon/update/delete qua production Gateway.
6. Xác nhận Admin và mọi external caller không gọi port 8101 trực tiếp.
7. Bỏ selector `map` khỏi production deploy script/runbook, giữ image cũ trong
   rollback window đã thống nhất.
8. Sau rollback window, xóa source Map và entry tương ứng khỏi `go.work`,
   development tooling, manifest và deployment compatibility.

Cho tới khi gate trên hoàn tất, `map-service/` là rollback/data-migration source,
không phải runtime owner. Không được thêm feature mới vào đó.

## Architecture guards

`make verify-backend` khóa các invariant:

- TQD đăng ký `MapPointService` đúng một lần.
- Gateway đăng ký generated handler đúng một lần.
- Gateway không reverse-proxy tới Map Service cũ.
- Compose không chạy Map Service cũ như standing container.
- Admin Organization chỉ có một serving owner tại User và một Gateway handler.
- Commercial checkout dùng User-owned membership authority, không gọi
  Organization RPC legacy và không suy diễn từ global role trong token.
- Internal Organization contract cũ chỉ là compatibility denominator cho CRM
  và các capability chưa migrate; không được nhận feature mới.
