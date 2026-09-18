# Quy chuẩn log và response backend

Ngày rà soát: 14/09/2026. Tài liệu này phân biệt cơ chế đã triển khai với các call-site cũ còn phải chuyển đổi. `make verify-backend` chỉ kiểm tra một số ràng buộc cấu trúc; PASS không chứng minh mọi business rule đúng.

## 1. Bài toán thực tế và nguyên nhân

`POST /v2/user/admin/catalog/plan-versions` với `draft.tier_rank = 0` từng trả HTTP 500, `code: 13`, `internal server error`.

Chuỗi lỗi trước đây:

1. `admin.validateDraftCommand` trả `errors.New("tier rank must be positive")`.
2. `planServiceError` dò chuỗi `required`, `invalid`, `page size`; câu trên không khớp.
3. Handler thay nguyên nhân bằng `codes.Internal`, `plan version administration failed`.
4. Gateway ẩn thông tin Internal — đúng trách nhiệm bảo vệ lỗi hạ tầng, nhưng loại lỗi đã bị sai từ bước 2.
5. Không có log nguyên nhân tại owner trước khi mapping nên dev phải thêm `fmt.Println`.

Một nhánh khác còn trả HTTP 200 khi nhận `shared.ErrorResponse`; ngoài ra Organization dùng các số 3001–5002 làm gRPC status. Đó là ba lỗi khác nhau: mất loại lỗi, HTTP thành công giả, trộn namespace mã.

## 2. Owner và nơi đặt code

| Thành phần | Owner | Trách nhiệm |
|---|---|---|
| Điều kiện `tier_rank > 0`, giá, entitlement | User domain/usecase | Quy tắc và lỗi có kiểu; không import HTTP/gRPC/logger |
| `common/validation` | Shared | Kiểu lỗi field/reason/message, ghi vị trí tạo lỗi để debug |
| `common/errors` | Shared | Bảng HTTP ↔ gRPC và encode BadRequest/ErrorInfo |
| Handler của từng service | Service | Map sentinel nghiệp vụ, log cause tại nơi còn biết operation/context |
| `gateway/internal/httperror` | Gateway | Đổi transport, giữ reason/field/IDs, che chi tiết 5xx |
| `common/logging` | Shared | Format, phân file, lọc level/category, rotation và RPC log |
| `logging.yml` | Cấu hình vận hành backend | Một file defaults + override theo service; không chứa bí mật |
| `common/db/sql_logger.go` | Shared DB adapter | Nối Trace/ParamsFilter của GORM vào logger; core logger không import ORM |
| `common/logging/fabric` | Adapter tương thích | Chuyển FabricLogger cũ một chiều vào slog |
| `main` mỗi process | Process owner | Khởi tạo logger một lần, đóng tài nguyên |
| `common/requestlog` | Shared | Gắn request_id/operation_id/actor đã có; không tạo danh tính/quyền |
| `native-stack.sh` | Supervisor | Giữ PID và stdout/stderr thô để thấy lỗi trước khi logger khởi tạo |

Config global này chỉ sở hữu logging, không trở thành một application config khổng lồ chứa business rule, địa chỉ DB, JWT hay quota. Các config `logging` cũ không có consumer (kể cả khóa Spring trong Go) đã được bỏ khỏi Auth/User/Assistant.

## 3. Cấu trúc log và cách cấu hình

Native `make up`, `make native-restart`, `make -C <service> server/dev` truyền `QHPRO_LOG_CONFIG` tới `logging.yml`. Khi chạy binary trực tiếp, truyền biến này tường minh. Không có biến này: JSON stdout; container/platform quản lý thu thập và retention.

```text
.tmp/development/logs/
  user/
    runs/20260914T...Z-<pid>/
      .writer.lock
      info/lifecycle.jsonl
      warn/application.jsonl
      error/application.jsonl
      error/dependency.jsonl
      debug/database.jsonl
  tqd/runs/<run_id>/...
  gateway/runs/<run_id>/...
  user.log, gateway.log, ...   # stdout/stderr thô của supervisor hiện có
```

Chỉ tạo file khi category/level có sự kiện. ERROR không ghi lặp vào INFO. `request_id`, tên table, method, feature, organization là **field**; không làm tên file/thư mục. Nếu tạo file theo các ID này, số lượng file sẽ tăng vô hạn.

File nghiệp vụ chính là JSONL, không có ANSI. `format: text|json` chỉ đổi console. Dùng viewer để có màu; JSON giữ nguyên khả năng tìm kiếm/index về sau. Supervisor capture không phải một error owner và không phân loại lại message. Không xóa log cũ của supervisor trong lần thay đổi này.

```yaml
version: 1
defaults:
  level: info
  console: true
  format: text
  add_source: true
  files: true
  root: .tmp/development/logs
  mode: session
  max_size_mb: 20
  max_backups: 3
  max_age_days: 7
  keep_runs: 10
  categories:
    application: info
    grpc: info
    dependency: info
    database: info
  sql:
    enabled: false
    slow_threshold: 1s
services:
  tqd:
    categories:
      database: debug
      grpc: debug
    sql:
      enabled: true
```

Đường dẫn `root` tương đối với file config, không theo cwd của service. Thứ tự: defaults trong code (stdout an toàn) → defaults YAML → services.<tên>. Ghi đè `false` hoạt động. Khóa sai, category không hợp lệ, version lạ, rotation không giới hạn đều bị từ chối ở startup. Muốn đổi cấu hình: restart process; chưa có hot reload để tránh một request có hai chính sách log.

| Mode | Hành vi | Dùng khi |
|---|---|---|
| `session` | Thư mục mới mỗi lần start, có `run_id` | So sánh hai phiên debug; mặc định native |
| `append` | Ghi tiếp active file, rotate khi chạm dung lượng | Cần lịch sử qua restart trong cùng thư mục |
| `truncate` | Xóa nội dung **active file đã biết ngay khi start**, giữ archive cũ | Dev chủ động chỉ muốn lần chạy hiện tại |

Có khóa file chống hai process dùng cùng thư mục append/truncate. Khóa tự giải phóng khi process crash. Trong session mode, startup dọn run cũ theo `keep_runs`/`max_age_days`, chỉ dọn tên run do logger tạo và không khóa được run đang chạy thì bỏ qua. Vì vậy quota dung lượng không phải hard cap khi nhiều process còn hoạt động. Không xóa thư mục tùy ý. Cơ chế file lock hiện phục vụ Linux (native WSL và Docker backend).

Mỗi file rotate theo dung lượng bằng Lumberjack; `max_age_days` là retention, **không phải** lệnh rotate lúc nửa đêm. Chưa có lịch calendar rotation. Lỗi ghi file được báo stderr, không đổi kết quả giao dịch; sự cố disk phải được supervisor theo dõi. Chưa triển khai hàng đợi async drop log: ghi đồng bộ, không mất âm thầm khi queue đầy.

### Xem log

Tại root backend:

```bash
make log-view SERVICE=user
make log-view SERVICE=tqd LOG_CATEGORY=database
make log-view SERVICE=gateway LOG_LEVEL=error
make log-view REQUEST_ID=req_9770ca1e2a97738baee1447ecd49df9f
```

`REQUEST_ID` nhận cả request ID và operation ID. Bỏ SERVICE sẽ gom các service về một timeline đọc, không bắt các process cùng ghi một file. `LOG_CONFIG=/đường/dẫn/logging.yml` chọn config khác. Viewer giới hạn 100 sự kiện mới nhất, màu tự bật khi terminal, tôn trọng `NO_COLOR`. Dùng `go run ./shared/common/cmd/logs -help` trong workspace để xem `-limit`, `-color`, `-stdin`; `tail -F <file.jsonl> | ... -stdin` phục vụ theo dõi liên tục.

## 4. Viết log có giá trị vận hành

| Category | Nội dung | Level thông thường |
|---|---|---|
| lifecycle | process ready/stopped, config đã chọn | INFO |
| application | quyết định nghiệp vụ, validation, nguyên nhân lỗi | INFO / WARN / ERROR |
| http | status và duration tại gateway | INFO; 4xx WARN; 5xx ERROR |
| grpc | kết quả RPC inbound, method, status | Thành công DEBUG, lỗi WARN/ERROR |
| dependency | RPC outbound, peer_service, method, duration | Thành công DEBUG, lỗi WARN/ERROR |
| database | query, rows, SQL verb, source, duration | DEBUG khi bật; chậm WARN; lỗi ERROR |
| worker | claim, retry, success/failure, job ID | INFO/WARN/ERROR theo outcome |

`logging.UnaryServer` đặt sau binder context và trước handler. `common/rpc.NewClient` gắn outbound logger. Không serialize nguyên request/response. Hai interceptor cũ của Payment/Organization từng dump request đã bị loại bỏ.

```go
requestlog.FromContext(ctx, nil).WarnContext(ctx,
    "Dữ liệu gói không hợp lệ",
    "category", "application",
    "event_name", "catalog.validation.failed",
    "error_reason", "PLAN_TIER_RANK_INVALID",
    "field", "draft.tier_rank",
    "error_source", validationError.Source,
)
```

Không log một lỗi nguyên nhân ở mọi tầng. Owner log cause; client interceptor log chuyến gọi thất bại; gateway log response cuối. Các event khác nhau nối qua IDs. HTTP có `request_failed`; `error_count` cũ chỉ đếm `gin.Context.Errors`, không phải tổng mọi lỗi downstream. `source` chỉ nơi ghi log; `error_source` chỉ nơi tạo lỗi validation; `query_source` chỉ repository gọi GORM. Đây không phải full distributed trace/span. Nếu request bị từ chối trước binder, không được log danh tính chưa xác thực như đã được tin cậy.

`log.Printf` cũ được bridge vào slog; FabricLogger bridge giữ severity/name. Không đoán severity từ chữ "error" trong text. Các `fmt.Print`, request ngoài common RPC hoặc GORM không `WithContext(ctx)` chưa tự có correlation; cần chuyển tại call-site. Không ghi token, mật khẩu, Authorization, cookie, DSN. Redaction field là lớp cuối, không thể cứu một secret đã bị ghép vào message bằng fmt.

SQL đi qua `common/db.Open` dùng logger chung. Chỉ gọi callback dựng SQL khi query cần ghi. `ParamsFilter` và hook `RecorderParamsFilter` của đường `Raw().Scan()` giữ placeholder `$1`…; SQL literal đã tự nội suy trong raw query vẫn phải sửa về tham số ở repository. Bật SQL thường cần **cả** `sql.enabled: true` và category `database: debug`. Log slow/error vẫn hoạt động khi query thường tắt. DB transaction và quyền truy cập không thuộc logger.

## 5. Response và quy tắc xử lý lỗi

| Tình huống | gRPC | HTTP | Client xử lý |
|---|---|---|---|
| Dữ liệu sai | InvalidArgument (3) | 400 | Hiển thị message + lỗi trường |
| Chưa xác thực | Unauthenticated (16) | 401 | Đăng nhập/làm mới phiên theo policy |
| Thiếu quyền | PermissionDenied (7) | 403 | Hiển thị thiếu quyền |
| Không tìm thấy | NotFound (5) | 404 | Không có tài nguyên |
| Trùng/đồng thời | AlreadyExists (6) / Aborted (10) | 409 | Làm mới hoặc đổi dữ liệu |
| Sai trạng thái | FailedPrecondition (9) | 412 | Hướng dẫn thao tác hợp lệ |
| Hết quota | ResourceExhausted (8) | 429 | Hiển thị usage/gói theo reason |
| Dependency không sẵn sàng | Unavailable (14) | 503 | Thử lại có kiểm soát |
| Timeout | DeadlineExceeded (4) | 504 | Truy vấn trạng thái trước retry mutation |
| Lỗi không dự kiến | Internal (13) | 500 | Thông báo an toàn + request ID |

Không dùng HTTP 200 để vận chuyển lỗi. `code` trong response lỗi gateway là số **gRPC**; `status` là số **HTTP**; `reason` là mã nghiệp vụ ổn định. `legacy_code` giữ mã cũ chỉ cho caller đang migrate. Response thành công đang có `code: 0` thì đúng; không đổi success envelope trong đợt này. Client phải dùng HTTP status/reason, không dò text hoặc trộn code HTTP với gRPC.

Response kỳ vọng cho case thực tế:

```json
{
  "status": 400,
  "code": 3,
  "message": "tier rank must be positive",
  "reason": "PLAN_TIER_RANK_INVALID",
  "domain": "user.catalog",
  "request_id": "req_...",
  "operation_id": "op_...",
  "errors": [{"field":"draft.tier_rank","error":"tier rank must be positive"}]
}
```

HTTP response không trả stack, đường dẫn source, câu SQL hoặc thông tin kết nối. Với 5xx, gateway bỏ detail/metadata; chi tiết tra trong log owner qua request ID. Đối với 4xx, service phải chỉ đưa dữ liệu công khai vào BadRequest/ErrorInfo.

### Dev thêm quy tắc mới

```go
// Trong usecase/domain: quy tắc có owner, không phụ thuộc transport.
if command.TierRank <= 0 {
    return DraftRecord{}, validation.New(
        "draft.tier_rank", "PLAN_TIER_RANK_INVALID", "tier rank must be positive",
    )
}
```

Handler dùng `errors.Is` cho sentinel lifecycle; `errors.As`/`ValidationStatus` cho lỗi trường; dùng `%w` để thêm context mà giữ loại lỗi. Không dùng `strings.Contains(err.Error(), ...)`. Không bọc mọi lỗi repository thành validation: "invalid database connection" vẫn phải là Internal. Kiểm tra SQLSTATE 23505 bằng kiểu driver (`pgconn.PgError`/`pq.Error`), không dò chữ duplicate key.

Khi đổi HTTP 200 → 4xx, các client xử lý lỗi trong success callback phải chuyển sang error callback. Đây là thay đổi contract có chủ ý, không coi chỉ là đổi log. Các interceptor axios hiện có phải giữ response body để hiển thị field/reason. Notification không được đánh dấu đã giao thành công nếu adapter trả lỗi; Payment không được retry mù settlement khi HTTP timeout.

## 6. Lựa chọn thư viện và chi phí

- [Go slog](https://go.dev/blog/slog): API chuẩn cho structured log, handler thay thế được; repository đã dùng ở gateway/requestlog nên chọn làm bề mặt chung.
- [Uber Zap](https://github.com/uber-go/zap): khoảng 24.7k stars tại thời điểm tra cứu; phù hợp throughput cao, hiện là dependency gián tiếp của Fabric. Benchmark upstream không phải benchmark workload QHPro. Không thay toàn bộ code sang API Zap khi chưa có bằng chứng CPU log là bottleneck.
- [Lumberjack](https://github.com/natefinch/lumberjack): khoảng 5.5k stars; chọn version **v2.2.1** ghim trong go.mod/go.sum, chỉ phụ trách rolling file, một writer/process. Shared giữ config/schema, tên sự kiện và routing; không fork thuật toán rotation.
- [gRPC errors](https://grpc.io/docs/guides/error/): dùng status + richer details (BadRequest, ErrorInfo), không tự chế wire protocol.
- [GORM logger](https://gorm.io/docs/logger.html): dùng interface Trace/ParamsFilter để nối SQL context và giữ tham số ngoài log.

Stars chỉ là tín hiệu cộng đồng, không chứng minh tính phù hợp. Đổi thư viện sink sau này không đổi field/reason và lời gọi slog trong service. Chưa có tuyên bố cải thiện throughput: đo disabled-path, enabled-path và I/O khi có tải tương ứng trước khi bật DEBUG toàn production.

## 7. Phạm vi đã nối và phần còn cần rà call-site

Đã nối main của 10 service native: auth, user, tqd, notification, payment, gateway, organization, file, hub, assistant; cả hai executable notification worker. Các entrypoint bdspro, crm, chat, chat-v1, social, search, map, relay, bootstrap admin, import location và CRM payment consumer cũng khởi tạo logger chung. Payment/Chat-v1/CRM consumer dùng SQL logger chung; PMTiles chuyển logger của thư viện vào category dependency. Shared DB/RPC áp dụng cho caller của constructor chung.

TQD chuyển lỗi của layer, region, label, family, authority, import, legend, land use, workspace và discovery theo kiểu lỗi tại owner. `infra/handler/errorresponse` là adapter chung của TQD: validation có field/reason/source; sentinel qua errors.Is; lỗi hạ tầng được log cùng request/method/handler_source rồi trả message an toàn. Repository authority nhận đúng SQLSTATE và tên unique index để xử lý cuộc đua tạo code. Build layer/family không còn gộp lỗi DB thành “not found”. HTTP import/discovery giữ tương thích tên trường response hiện có nhưng phân loại theo kiểu.

Catalog validation/publication dùng lỗi có kiểu; payment admin pagination cũng vậy. Organization có reason riêng; KYC/PriceTable chuyển số HTTP qua bảng shared thay vì cast thành gRPC. `strings.Contains` trong test chống lộ dữ liệu, tìm văn bản SEO hoặc phân tích codec vẫn có vai trò hợp lệ.

Giới hạn: nhiều log ứng dụng cũ vẫn đi qua adapter log.Printf/Fabric với category application; fmt.Print không tự thành structured log, và câu SQL thiếu WithContext không tự có request ID. Chưa tuyên bố mọi call-site đã có đầy đủ ngữ nghĩa và mọi tính năng đã được nghiệm thu runtime. Chuyển thêm phải dựa trên hợp đồng nghiệp vụ và test hành vi, không đổi hàng loạt mọi return err thành một status tùy đoán.

## 8. Kiểm chứng và bằng chứng

Chạy từ root repository:

```bash
make verify-backend
make test-backend
make native-build
make deps
make native-up
make native-status
make log-view SERVICE=tqd LOG_LEVEL=warn
make log-view REQUEST_ID=<request_id nhận trong response>
```

`verify-backend` kiểm tra cấu trúc/config, `test-backend` kiểm tra core services, `native-build` tạo binary; ba lệnh này không chứng minh dependency đã chạy. Chỉ xác nhận runtime khi native-status READY và request qua gateway trả đúng contract.

Các kết quả đợt đồng bộ được lưu tại `.tmp/development/code-review/`:

- `tqd-error-contract.log`: usecase, handler, repository, gồm lỗi DB có message giống lỗi nghiệp vụ.
- `backend-sync-tests.log`: core backend tests.
- `extra-services-sync-tests.log`, `chat-v1-sync-tests.log`: service ngoài native, chat-v1 chạy standalone vì trùng tên module chat.
- `sync-race-tests.log`: logging/SQL/errors/RPC và TQD error contract.
- `sync-native-build.log`: build native.

Probe local `check-tqd-errors.py` sử dụng đăng nhập dev của probe catalog hiện có; chỉ gửi dữ liệu validation sai hoặc ID không tồn tại. Khi dependency đang chạy:

```bash
python3 .tmp/development/code-review/check-tqd-errors.py
```

Kết quả mong đợi: catalog tier_rank=0 → 400/PLAN_TIER_RANK_INVALID; family name trắng → 400/LAYER_FAMILY_NAME_INVALID; authority code trắng → 400/AUTHORITY_CODE_INVALID; family hoặc import error không tồn tại → 404 với reason tương ứng. Request/operation ID trong body phải khớp header. File `tqd-runtime-errors.json` chỉ được tạo sau khi mọi assertion pass; không sử dụng bằng chứng phiên cũ để kết luận binary hiện tại đã chạy.

Nếu Docker Desktop/WSL integration chưa hoạt động, dependency và runtime acceptance chưa đạt. Giữ trạng thái đó rõ ràng trong bàn giao; không thay kiểm chứng runtime bằng kết quả build.

Kiểm thử bảo vệ thay đổi: slog handler contract (WithGroup/WithAttrs), concurrent writers, redaction, parameterized SQL/correlation, typo/override config, rotate/truncate/lock; handler rank=0; DB error có chữ invalid; gateway HTTP status/details; compile từng service cả workspace và standalone. Bằng chứng chạy nằm trong `.tmp/development/code-review/logging-*`.
