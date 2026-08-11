# QHPro generated report: đọc code và chạy end-to-end trên local

Tài liệu này là một lab thực hành. Mỗi chặng trả lời bốn câu hỏi:

1. Nghiệp vụ đang cần bảo vệ điều gì?
2. File/hàm nào chịu trách nhiệm?
3. Input biến thành output như thế nào?
4. Kiểm chứng bằng HTTP, PostgreSQL hoặc Redis ra sao?

Lab chỉ dùng database/Redis local, profile `42` và quota `3 report/subscription cycle`.
Không dùng credential trong `config/local.yml` và không dùng `make migrateup` của TQD.

## 1. Bức tranh tổng thể

Một lần bấm “tạo báo cáo” đi theo đường sau:

```text
React Native/Web
    |
    | POST /v2/tqd/map-workspace/reports
    | JWT + X-Operation-ID + Idempotency-Key
    v
Gateway HTTP :8000
    | 1. verify JWT, tạo Actor(profile=42)
    | 2. bind request identity
    | 3. chuyển sang gRPC metadata và ký HMAC
    v
TQD gRPC :8219
    | 4. verify chữ ký, restore Actor/request identity
    | 5. map RPC -> workspace.report.generate
    | 6. build report input
    |
    +------ signed gRPC ------> user-service :8201
    |                          đọc subscription + entitlement
    |                          trả allowed, limit, period
    |
    | 7. Redis Lua reserve quota
    | 8. PostgreSQL transaction:
    |      user_reported + quota_usage_events + report_jobs
    | 9. Redis Lua commit: reserved -> used
    v
HTTP response: report accepted / processing
    |
    | background report worker
    v
Smoke generator -> signed file gRPC :8203 -> file-service DB/filesystem
    |
    v
report_jobs=completed, user_reported.status=ready, pdf_url có giá trị
```

Ba loại identity không được nhầm:

| Identity | Ý nghĩa | Ví dụ |
|---|---|---|
| Request-ID | trace một HTTP request vật lý | retry có Request-ID khác |
| Operation-ID | một lần thực hiện logic từ đầu đến cuối | `op-report-001` |
| Idempotency-Key | identity của command do client sở hữu | `idem-report-001` |

Nếu có Idempotency-Key, nó là `commandKey`. Nếu thiếu, Operation-ID được dùng làm
fallback. Retry cùng command không được tạo report hoặc trừ quota lần hai.

## 2. Chặng 1: Frontend tạo command

File chính:

- `react-native-quy-hoach/src/Services/tqd/report.ts`
- `react-native-quy-hoach/src/Services/tqd-workspace-map-service.ts`
- `react-native-quy-hoach/src/Services/request.ts`

`createGeneratedReport()` không tự quyết định plan hoặc quota. Nó chỉ tạo payload,
chọn một `clientRequestId` ổn định và gửi cùng giá trị đó qua header
`Idempotency-Key`.

Ví dụ:

```json
{
  "reportType": 1,
  "entityType": 1,
  "entityId": 1,
  "parcelId": 1,
  "title": "QHPro local smoke",
  "clientRequestId": "idem-report-001"
}
```

- `reportType=1`: planning report, format mặc định là PDF.
- `entityType=1`: parcel; lab seed parcel PostGIS ID 1 để kiểm chứng cả source lookup.
- `clientRequestId`: compatibility cho client cũ; nếu header cũng có thì hai giá trị
  phải giống nhau.

Kiểm chứng frontend contract:

```bash
cd /home/bop/projects/react-native-quy-hoach
yarn test \
  src/Services/__tests__/request.operationId.test.ts \
  src/Services/__tests__/tqdWorkspaceReport.idempotency.test.ts \
  --runInBand
```

## 3. Chặng 2: Gateway xác thực và tạo context chuẩn

File chính:

- `shared/common/jwt/middleware.go`
- `gateway-service/internal/httpmiddleware/operation_id.go`
- `gateway-service/internal/httpmiddleware/idempotency_key.go`
- `shared/common/requestcontext/*_contract.go`

Thứ tự middleware tại Gateway:

1. Verify JWT bằng `jwt.key-generate`.
2. Tạo `Actor` từ claims, quan trọng nhất là `ProfileID`.
3. Nhận hoặc generate `X-Operation-ID`.
4. Validate `Idempotency-Key`; Gateway không tự tạo key này.

Các ID chỉ chấp nhận chữ, số, `-`, `_`, `.`, `:` và tối đa 128 ký tự. Header
trùng lặp hoặc chứa khoảng trắng bên trong bị trả HTTP 400.

Token lab có claim `profile=42`. Script tạo token không gọi login thật, nhưng vẫn
đi qua đúng code verify JWT của Gateway:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/token.sh
```

## 4. Chặng 3: Chỉ route report mới dùng metadata target

File chính:

- `gateway-service/internal/grpcmetadata/from_http.go`
- `gateway-service/cmd/http/main.go`

`DuringReportMigration()` chỉ chọn bridge mới khi đồng thời thỏa:

```text
QHPRO_REPORT_TARGET_ENABLED=true
method = POST
path = /v2/tqd/map-workspace/reports
```

Bridge mới chỉ chuyển canonical Actor, Caller, Request-ID, Operation-ID và
Idempotency-Key. Nó không chuyển raw JWT, API key hoặc `planId/planFrom` cũ.
Điều này buộc downstream phải hỏi user-service thay vì tin dữ liệu thương mại do
client gửi lên.

Sau đó `common/grpcclient.UnaryClientInterceptor` ký method và metadata bằng HMAC.
Trong lab, Gateway, TQD, user-service và file-service dùng cùng một secret local.

## 5. Chặng 4: TQD verify transport và gắn BusinessOperation

File chính:

- `tqd-service/cmd/target_ingress.go`
- `shared/common/requestgrpc/interceptor.go`
- `shared/common/identitygrpc/interceptor.go`
- `shared/common/transporttrust/trust.go`
- `tqd-service/cmd/business_operation_target.go`
- `tqd-service/infra/grpcoperation/interceptor.go`

TQD không tin metadata chỉ vì Gateway gửi tới. Nó:

1. Verify caller, timestamp, gRPC method và HMAC.
2. Restore request identity vào `context.Context`.
3. Restore Actor đã được Gateway xác thực.
4. Map method `/types.tqd.MapWorkspaceService/CreateGeneratedReport` thành
   `workspace.report.generate`.

Business operation là tên ổn định của hành động, không phải tên HTTP route và
không phải meter. Report service chỉ chạy nếu operation trong context đúng chính
xác.

## 6. Chặng 5: Adapter và idempotency compatibility

File chính:

- `tqd-service/infra/reportgrpc/adapter.go`
- `tqd-service/internal/report/request_hash.go`

Adapter chuyển protobuf thành `report.Input`. Nó cũng xử lý giai đoạn migration:

- Chỉ có header key: dùng header.
- Chỉ có `clientRequestId`: bind nó vào context.
- Có cả hai và giống nhau: hợp lệ.
- Có cả hai nhưng khác nhau: `InvalidArgument`, HTTP 400.

`BuildRequestHash()` SHA-256 toàn bộ input. Vì thế cùng key nhưng payload khác sẽ
trả conflict thay vì vô tình trả report cũ không đúng nội dung.

## 7. Chặng 6: Report service điều phối nghiệp vụ

File chính:

- `tqd-service/internal/report/service.go`
- `tqd-service/internal/report/build.go`

Đây là use case trung tâm. Thứ tự được cố ý giữ như sau:

1. Đọc Actor; report luôn thuộc `Actor.ProfileID`.
2. Kiểm tra business operation.
3. Chọn command key.
4. Tìm report cũ trước; retry có thể kết thúc sớm.
5. Validate/build report.
6. Hỏi access từ user-service.
7. Reserve quota.
8. Commit DB transaction.
9. Commit quota.

Không reserve trước khi validate input để request sai không giữ quota. Không commit
quota trước DB vì DB rollback sẽ làm người dùng mất quota dù report chưa accepted.

`build.go` có ba nhánh source:

- Parcel: bắt buộc parcel tồn tại, lấy location/spatial từ DB.
- Region: bắt buộc region tồn tại, lấy metadata vùng từ DB.
- Location snapshot: dùng location/spatial trong request.

Lab chọn nhánh Parcel. PostgreSQL có một polygon thật trong `parcels`; server tự
lấy address, bounds và centroid bằng PostGIS. Payload client không được quyền tự
quyết định các giá trị này.

## 8. Chặng 7: User-service quyết định access

File chính:

- `tqd-service/infra/useraccess/client.go`
- `user-service/infra/handler/accessgrpc/handler.go`
- `user-service/internal/access/service.go`
- `user-service/infra/postgres/access/subscription_store.go`
- `user-service/infra/postgres/access/plan_access_store.go`
- `user-service/infra/postgres/access/catalog_codes.go`

TQD không gửi profile ID trong body gRPC. User-service lấy profile từ Actor đã ký,
ngăn service gọi tự khai profile khác.

Mapping catalog:

```text
Business operation: workspace.report.generate
Feature grant:      workspace.report.generate
Usage meter:        workspace.report_generation.accepted
```

User-service tìm subscription còn hiệu lực, rồi tìm hai entitlement của plan
version. Kết quả lab:

```text
subject       = profile:42
allowed       = true
limit         = 3
period        = subscription_cycle
period_start  = subscription.current_period_start
period_end    = subscription.current_period_end
```

Access chỉ trả chính sách; nó không đọc `used/reserved`. Redis quota thuộc TQD.

## 9. Chặng 8: Redis reserve/commit/cancel

File chính:

- `tqd-service/internal/quota/service.go`
- `tqd-service/infra/redisquota/reserve.go`
- `tqd-service/infra/redisquota/finalize.go`
- `tqd-service/infra/redisquota/keys.go`

Usage bucket:

```text
quota:usage:profile:42:workspace.report.generate:<periodStart>:<periodEnd>
```

Hai số quan trọng:

- `used`: report đã được DB chấp nhận.
- `reserved`: request đang ở giữa Redis reserve và DB commit.

Lua reserve thực hiện nguyên tử:

```text
if used + reserved + amount > limit -> quota exceeded
else reserved += amount, tạo reservation và command mapping
```

Sau DB success, Lua commit chuyển `reserved -= 1`, `used += 1`. Nếu DB fail,
cancel chỉ giảm `reserved`. Gọi commit lại trên reservation đã committed không
tăng `used` lần nữa.

## 10. Chặng 9: PostgreSQL transaction là bằng chứng bền vững

File chính:

- `tqd-service/infra/postgres/report/save_store.go`
- `tqd-service/infra/postgres/usageevent/store.go`
- `tqd-service/infra/postgres/reportjob/store.go`
- `tqd-service/infra/db/migrations/v3/000009_report_quota_target.up.sql`

Một transaction ghi ba thứ:

1. `user_reported`: report mà người dùng nhìn thấy.
2. `quota_usage_events`: bằng chứng durable cho một lần tính usage.
3. `report_jobs`: durable queue để tạo file sau response.

Unique command `(user_id, client_request_id)` và `request_hash` bảo vệ retry/race.
Usage event chứa đúng `reservation_id`, nhờ vậy cleanup không thể lấy bằng chứng
của reservation khác.

## 11. Chặng 10: Background worker và output file

File chính:

- `tqd-service/internal/reportjob/worker.go`
- `tqd-service/internal/reportjob/runner.go`
- `tqd-service/infra/reportgenerator/smoke.go`
- `tqd-service/infra/filegrpc/client.go`
- `file-service/grpc_target.go`
- `file-service/internal/filecontent/service.go`

Worker claim job bằng `FOR UPDATE SKIP LOCKED`, đổi `pending -> running`, gọi
`Generator`, rồi cùng transaction đổi:

```text
report_jobs.status = completed
user_reported.status = 10 (ready)
user_reported.pdf_url = output URL
```

Smoke generator tạo một PDF nhỏ có report ID, user ID, operation và job ID, rồi
upload bằng signed gRPC. Nó chỉ chạy khi:

```text
QHPRO_REPORT_GENERATOR_MODE=smoke
```

Đây là renderer kiểm chứng kỹ thuật, không phải renderer bản đồ/quy hoạch thật.
Mặc định flag trống nên production behavior không đổi.

## 12. Setup local từ số 0

### 12.1 Yêu cầu

- Go theo `go.work` của repo.
- Node 20+ và Yarn.
- `curl`, `jq`.
- Docker Desktop với WSL integration, hoặc Docker Engine + Compose plugin.

Nếu WSL báo `docker: command not found`:

1. Mở Docker Desktop trên Windows.
2. Settings -> Resources -> WSL Integration.
3. Bật distro đang dùng.
4. Apply & Restart.
5. Mở terminal WSL mới và chạy `docker version`.

### 12.2 Tạo PostgreSQL/PostGIS và Redis

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/up.sh
```

Lab dùng port riêng:

```text
PostgreSQL/PostGIS  127.0.0.1:55432
Redis               127.0.0.1:56379
```

PostgreSQL init tự tạo:

- `qhpro_user_lab`: Catalog V2 + subscription profile 42.
- `qhpro_tqd_lab`: PostGIS parcel ID 1 + base `user_reported` + migration target.
- `qhpro_file_lab`: file-service sẽ AutoMigrate bảng file khi start.

Volume có tên cố định `qhpro_report_lab_*`, không đụng database/Redis khác.

Muốn xóa dữ liệu report/quota để chạy lại các case từ đầu:

```bash
CONFIRM_QHPRO_REPORT_LAB_RESET=yes scripts/report-lab/reset-state.sh
```

Script chỉ truncate ba bảng target trong `qhpro_tqd_lab`, bảng file trong
`qhpro_file_lab` và `FLUSHDB` Redis container riêng. Nó không xóa các file đã sinh
trong `file-service/files`; có thể giữ chúng làm evidence.

### 12.3 Generate proto và test code

```bash
cd /home/bop/projects/bdspro-golang-microservices
bash scripts/generate-internal-protobuf.sh
bash scripts/validate-reconstruction-target.sh

cd tqd-service
go test ./infra/reportgenerator ./internal/reportjob ./internal/report ./internal/quota
```

### 12.4 Start service đúng thứ tự

Mở bốn terminal.

Terminal 1, user-service:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/run-user.sh
```

Terminal 2, file-service:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/run-file.sh
```

Terminal 3, TQD:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/run-tqd.sh
```

Terminal 4, Gateway:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/run-gateway.sh
```

Các port:

| Process | Port | Protocol |
|---|---:|---|
| Gateway | 8000 | HTTP public |
| file-service | 8002 | HTTP file delivery |
| user-service | 8201 | internal gRPC |
| file-service | 8203 | internal gRPC |
| TQD | 8219 | internal gRPC |

## 13. Gọi API thực tế và tải output

Sau khi bốn service đã listen:

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/smoke.sh
```

Script thực hiện:

1. Tạo JWT profile 42.
2. POST report qua Gateway.
3. Lấy `reportId` từ response envelope.
4. Poll GET report tối đa 20 giây.
5. Chờ status `10` và `pdfUrl`.
6. Download file.
7. Kiểm tra magic bytes `%PDF-`.

Output cuối:

```text
PASS report_id=<id>
PASS pdf_url=http://127.0.0.1:8002/v1/file/load/<encoded-path>/local
PASS output=.../scripts/report-lab/qhpro-smoke-output.pdf
```

Gateway dùng envelope chung nên create response có dạng:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true,
    "message": "report generation started",
    "item": {
      "reportId": 1,
      "status": 20
    }
  }
}
```

Sau worker, GET trả `data.status=10`, `data.assets.pdfUrl` và
`data.reportMeta.fileSize`.

## 14. Quan sát bằng DB và Redis

```bash
cd /home/bop/projects/bdspro-golang-microservices
scripts/report-lab/inspect.sh
```

Mong đợi sau một command:

```text
user_reported:      1 row, status=10, pdf_url != ''
quota_usage_events: 1 row, amount=1
report_jobs:        1 row, status=completed, attempts=1
Redis usage hash:   used=1, reserved=0
```

## 15. Các tình huống cần tự kiểm chứng

### A. Retry đúng idempotency

Gửi lại cùng header, cùng `clientRequestId` và cùng body.

Mong đợi:

- HTTP 200.
- Cùng report ID.
- Business message `report request already accepted`.
- Vẫn chỉ một row ở cả ba bảng.
- Redis `used` không tăng.

### B. Cùng key nhưng body khác

Giữ `Idempotency-Key`, đổi title.

Mong đợi HTTP 409 vì `request_hash` khác. Không có usage/report/job mới.

### C. Không có subscription

Tạo token bằng profile khác:

```bash
QHPRO_LAB_PROFILE_ID=43 scripts/report-lab/token.sh
```

Gọi POST với token đó. Mong đợi HTTP 403; Redis và TQD DB không đổi.

### D. Hết quota

Fixture có limit 3. Gửi bốn command với bốn Idempotency-Key khác nhau.

- Ba command đầu accepted.
- Command thứ tư trả HTTP 429.
- DB có đúng ba usage event.
- Redis `used=3`, `reserved=0`.

### E. Header/body idempotency xung đột

Header là `idem-A`, body `clientRequestId` là `idem-B`.

Mong đợi HTTP 400 trước reserve quota.

### F. File-service tắt

Tắt file-service rồi tạo report.

- HTTP create vẫn accepted vì DB/quota không phụ thuộc renderer.
- Job retry tối đa ba lần.
- Sau retry, job/report chuyển failed.
- Usage vẫn được tính vì business action “request report” đã accepted.

Đây là khác biệt quan trọng giữa acceptance và render completion.

## 16. Troubleshooting theo checkpoint

| Hiện tượng | Checkpoint cần xem |
|---|---|
| HTTP 401 | JWT secret/claim profile ở Gateway |
| `trusted gRPC metadata is required` | service ID và trust secret giữa các process |
| HTTP 403 | subscription/feature entitlement của profile |
| HTTP 429 | Redis usage đã đạt limit |
| HTTP 409 | cùng command key nhưng request hash khác |
| HTTP 503 sau create | DB đã accepted nhưng Redis commit lỗi; retry cùng key |
| report mãi processing | generator disabled, file gRPC không chạy hoặc job retry |
| PDF URL có nhưng 404 | file-service HTTP không chạy hoặc filesystem file bị mất |

Lệnh xem log infra:

```bash
docker compose -f scripts/report-lab/compose.yml logs postgres redis
```

Không đặt `ENV_RUNTIME=local` khi chạy các script trực tiếp trên host. Trong source,
giá trị non-empty làm client tìm hostname Docker như `user`/`tqd`; lab chủ động đặt
`ENV_RUNTIME=` để dùng `localhost`.

Khi hoàn tất, dừng container nhưng giữ dữ liệu:

```bash
scripts/report-lab/down.sh
```

## 17. Điều lab chứng minh và chưa chứng minh

Lab chứng minh:

- HTTP -> signed gRPC -> TQD -> signed user access gRPC.
- Catalog subscription/entitlement.
- Redis reserve/commit và quota exceeded.
- PostgreSQL atomic acceptance + durable job.
- Background worker -> signed file gRPC -> file DB/filesystem.
- API trả report ready và tải được PDF hợp lệ.

Lab chưa chứng minh renderer nghiệp vụ thật. PDF smoke không có bản đồ, parcel,
planning layers hay template sản phẩm. Để production-ready cần thay
`SmokeGenerator` bằng renderer thật nhưng giữ nguyên interface
`reportjob.Generator`; quota, transaction và worker lifecycle không cần đổi.

## 18. Nối React Native vào lab

Nên chạy `smoke.sh` thành công trước. Khi đó lỗi còn lại ở ứng dụng chỉ thuộc lớp
mobile/session, không còn phải đoán Redis, database hay worker có hoạt động hay
không.

Các file cần đọc theo thứ tự:

1. `react-native-quy-hoach/config/application-environments.json`: nguồn URL theo
   environment. Profile `local` trỏ API tới port `8000`, file tới `8002`.
2. `react-native-quy-hoach/src/Env/index.ts`: trên Android local, hostname tự đổi
   từ `127.0.0.1` thành `10.0.2.2`; iOS simulator giữ `127.0.0.1`.
3. `react-native-quy-hoach/src/Services/request.ts`: gắn JWT,
   `X-Operation-ID` và `Idempotency-Key` vào HTTP request.
4. `react-native-quy-hoach/src/Services/tqd-workspace-map-service.ts`: tạo, lấy,
   liệt kê, xóa và chia sẻ report.
5. `react-native-quy-hoach/src/Services/tqd/report.ts`: facade ngắn cho action tạo
   và lấy report.
6. `PersonalMapPanel.tsx`/`usePersonalMapWorkspaceData.ts`: state và UI danh sách
   report của người dùng.

Kiểm chứng riêng contract phía app, không cần emulator:

```bash
cd /home/bop/projects/react-native-quy-hoach
yarn test src/Services/__tests__/request.operationId.test.ts \
  src/Services/__tests__/tqdWorkspaceReport.idempotency.test.ts --runInBand
```

Chạy Android emulator sau khi backend đã lên:

```bash
cd /home/bop/projects/react-native-quy-hoach
yarn install
yarn android:local
```

Android emulator dùng `10.0.2.2:8000` để tới host. Nếu dùng thiết bị thật, đặt
`QHPRO_LOCAL_HOST` trong một file env local riêng thành IP LAN của máy, rồi chạy
với `QHPRO_ENV_FILE=<đường-dẫn-file-env>`. Đồng thời start TQD bằng public file URL
mà điện thoại truy cập được:

```bash
QHPRO_FILE_PUBLIC_BASE_URL=http://<IP-LAN>:8002 \
  scripts/report-lab/run-tqd.sh
```

Lab seed profile `42` và `smoke.sh` tự tạo JWT cho profile đó. Ứng dụng mobile
không tự nhận token này: để thao tác bằng UI, phiên đăng nhập hiện tại của app
phải có profile `42`, hoặc cần một cơ chế debug inject token riêng. Không nên thêm
backdoor token vào code production. Vì auth/login đầy đủ kéo theo auth-service và
các dependency ngoài phạm vi report, API smoke là bằng chứng E2E chuẩn của feature
này; test service phía React Native là bằng chứng cho mapping header/body.
