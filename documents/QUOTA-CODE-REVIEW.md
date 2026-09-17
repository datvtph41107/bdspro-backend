# Review code thương mại và quota — 2026-09-14

## Phạm vi và mức độ bằng chứng

Review này đọc code và sửa lỗi trên trục User → Payment → TQD → File, kèm
Gateway và caller frontend/admin. Không phải chứng nhận mọi file của repository
đã được review. Kiểm kê riêng bảy service liên quan có 1.389 file Go; việc tìm
kiếm toàn source khác với đọc và kiểm chứng từng đường thực thi.

| Phần | Đã đọc/kiểm tra | Mức bằng chứng |
| --- | --- | --- |
| TQD quota | action accounting, reserve/commit/cancel, Redis scripts, usage PostgreSQL, cleanup, projection, registry, wiring | Sửa logic + regression test + compile cmd/wire |
| TQD report | acceptance, request replay, transaction boundary, lease context | Dùng chung lease policy + test application/processing/rendering |
| User | access result, entitlement, actor/subject, checkout snapshot, current subscription và expiry | Siết validation + test internal, checkout/subscription PostgreSQL packages, seed |
| Payment | SePay handler, order/settlement/fulfillment, wallet deposit, money conversion, transaction query | Sửa validation, locking và replay + unit test + SQL DryRun |
| Gateway | route dùng credential provider, chuyển metadata, webhook response | Đọc code + test các package tương ứng |
| File | owned file identity và lookup | Đọc code + test filecontent/HTTP handler |
| React Native | identity, parcel/detail caller, payment polling, checkout/report | Test 6 suite, 26 case và commercial-flow guard |
| Admin | UsagePage, PlanEditor capability mapping, catalog API | Đọc source; chưa chạy browser Admin ở lượt này |

Các log kiểm tra được giữ ở `.tmp/development/code-review/`. Browser acceptance
và giao dịch thật từ các lượt trước không được tính là kiểm chứng runtime cho
bản sửa này. Chưa deploy, chạy migration production hay chuyển tiền trong lượt review này.

## Các lỗi đã sửa

### 1. P1 — hủy quota của command đang được retry

**Trước:** `ActionService.Run` hủy reservation ngay khi callback lỗi hoặc
`SaveUsageEvent` báo lỗi. Hai request cùng command có thể dùng cùng reservation;
request lỗi có thể hủy quota của request sắp commit. Lỗi mạng database cũng
không chứng minh transaction chưa commit.

**Sau:** giữ lease có thời hạn. Cleanup đọc bằng chứng usage theo reservation
để commit/cancel; handler và store dùng deadline trước khi lease hết hạn.
Policy deadline ở `quota/acceptance_context.go` được dùng chung bởi report và
lượt đọc. Callback lỗi không ghi usage mới.

**Đánh đổi vận hành:** request lỗi có thể giữ tạm một lượt đến khi lease hết
(mặc định 5 phút) và cleanup xử lý. Đây là reserved, không phải used. Cần theo
dõi runner cleanup; không cộng quota thủ công để chữa phần đang giữ tạm.

### 2. P1 — đã tính lượt nhưng client nhận lỗi Redis

**Trước:** database đã ghi usage, Redis commit lỗi thì response thành thất bại.

**Sau:** giữ response thành công; log `read accepted; quota projection pending
repair`. Reconciliation sửa projection. Không bảo đảm client nhận được byte nếu
kết nối mạng đã mất; idempotency vẫn cần cho retry.

### 3. P1 — commit reservation mới dù usage thuộc kỳ cũ

**Trước:** bỏ qua `created=false` từ insert chống trùng, rồi commit reservation
hiện tại. Retry qua kỳ mới có thể tăng runtime used mà ledger không có usage mới.

**Sau:** interface `ReadActionUsageStore` buộc kiểm tra usage có thuộc reservation
hiện tại hay không. Reservation dư chờ cleanup, không tăng used. Interface thuộc
consumer; PostgreSQL và test store cùng thực hiện. Đã regenerate Wire.

### 4. P1 — entitlement sai trở thành quyền không tính quota

**Trước:** `Allowed=true`, có meter, `Limit=0`, `Period=none` vẫn hợp lệ;
`UsesQuota=false` mở đường bỏ qua accounting. Read adapter/projection còn bỏ qua
validation khi cho rằng quyền không metered.

**Sau:** cả User và TQD từ chối allowance hữu hạn bằng 0; không giới hạn phải
khai báo `Unlimited`. Read adapter và usage projection xác minh decision trước
khi dùng. Không diễn giải feature không tính quota là feature miễn phí.

### 5. P1 — webhook tiền ra đi vào nhánh nạp ví

**Trước:** chỉ nhánh order QHP kiểm tra `transferType=in`. Nhánh DEP nạp ví có
thể nhận giao dịch tiền ra; request nil còn có đường panic.

**Sau:** kiểm tra request, chiều tiền và số tiền trước khi chọn order/ví.
Order không có settlement dependency trả unavailable, không rơi xuống ví.
Mã không nhận diện được trả lỗi; không tự tạo quyền hoặc nạp tiền.

### 6. P1 — hai webhook đồng thời cộng ví hai lần

**Trước:** đọc trạng thái transaction không khóa, sau đó mới khóa wallet.
Hai request có thể cùng đọc pending, lần lượt cộng cùng khoản vào wallet.
Replay completed cũng trả lỗi khiến provider tiếp tục retry.

**Sau:** `GetTransactionByCodeForUpdate` dùng `SELECT ... FOR UPDATE` trong
transaction của caller trước khi đọc trạng thái. Chỉ pending deposit được
chuyển completed; replay cùng số tiền trả lại transaction cũ. Canceled hoặc
transaction loại khác bị từ chối. Kiểm tra overflow trước khi cộng số dư.

Test usecase mô phỏng hai webhook đồng thời và xác minh một lần cập nhật/audit;
test SQL DryRun xác nhận query adapter thật có khóa. Chưa có bằng chứng chạy
race này trên PostgreSQL thật trong lượt này.

### 7. P1 — float64 ở biên int64 đổi dấu số tiền

`float64(MaxInt64)` làm tròn thành 2^63. So sánh `>` trước kia cho lọt chính
ngưỡng đó. `MinorFromWire` nay loại `>=` ngưỡng dương trước khi cast. Webhook
dùng lại hàm chuyển đổi tiền hiện có, không thêm validator tiền song song.

### 8. P2 — test/registry/comment tạo cảm giác an toàn sai

- `ProfileOperations` chuyển từ slice public có thể sửa thành hàm trả bản sao.
  Danh sách chỉ sở hữu projection của TQD; User vẫn sở hữu quyền theo gói.
- Test đối chiếu protobuf annotation với projection; bỏ test chỉ tự chép lại
  ba constant và khóa cứng thứ tự.
- Memory usage store trả lỗi với event sai/context hủy như adapter thực,
  không giả thành replay thành công.
- Sửa comment Service ở đúng type, giải thích bằng tiếng Việt owner, lý do giữ
  lease, điểm commit và recovery.
- Sửa tài liệu: report tính tại acceptance, không chờ client tải PDF; usage key
  không chứa period; FK không thể đi trực tiếp qua database User/TQD. Retire
  catalog không có nghĩa xóa version lịch sử.

## Các vấn đề còn mở — không coi là production-ready

| Mức | Bằng chứng | Công việc tiếp theo / điều kiện hoàn tất |
| --- | --- | --- |
| DONE | Read command đã có `quota_read_commands`, hash deterministic của protobuf request và conflict khác payload | Migration 000045; TQD handler bind trước action. Legacy row không có hash fail-closed; cần chạy migration khi deploy |
| DONE | Redis command key v2 thêm operation; key v1 được đọc chuyển tiếp chỉ khi reservation lưu đúng operation | Đã sửa Lua/memory store; cần rollout Redis thật và theo dõi key v1 hết TTL |
| P1 | `quota_pools` chỉ xuất hiện ở migration 000038, không có Go reader/writer | Inventory dữ liệu môi trường đích; xác định owner/retention rồi mới lập migration bỏ bảng. Không xóa migration lịch sử |
| P1 | Admin SQL lọc profile; UsagePage cũng chặn reconcile organization | Mở API/projection theo subject và kiểm thử membership, org dùng chung quota, quyền đọc report trước khi mở organization billing |
| P1 | Frontend tạo key đọc theo bucket thời gian và hash ngắn | Tách identity của một lần tương tác khỏi retry; test retry qua ranh giới bucket, reload và collision. Quy tắc chống tính trùng phải do backend bảo vệ |
| P1 | Legacy ví chưa lưu provider transaction id trong flow DEP | Bổ sung sổ bằng chứng ngân hàng và đối soát giao dịch trùng/đổi reference; khóa transaction hiện sửa double credit cùng deposit, chưa thay thế bank ledger |
| P2 | PlanEditor có capability mapping riêng; protobuf test chỉ đối chiếu TQD RPC/projection | Contract check giữa catalog capability, Admin và frontend; không đặt giá/allowance vào registry mã operation |
| P2 | GIS legacy queries, dữ liệu local khác production | Test thửa/vùng thật và so sánh output bản đồ; không suy luận từ test quota rằng GIS đúng |

## Chuẩn phát triển từ các lỗi này

1. Bắt đầu bằng tình huống thất bại và kết quả người dùng nhìn thấy. Viết test
   cho tình huống đó trước khi gọi một thay đổi là cải tiến.
2. Mỗi dữ liệu có owner rõ. Không gom tất cả service thành một owner; thống
   nhất contract giữa các owner và transaction trong phạm vi từng database.
3. Comment nêu bất biến và lý do: ai sở hữu lease, có thể retry đồng thời không,
   commit nào là bằng chứng và việc nào thuộc recovery.
4. Tên phương thức thể hiện yêu cầu quan trọng: `ForUpdate` báo cho caller rằng
   phải có transaction. Không giữ tên/alias cũ nếu không có caller cần nó.
5. Không lấy việc test pass làm chứng nhận toàn source. Ghi package, loại test,
   giới hạn mock/DryRun và phần runtime chưa chạy.
6. Khi đổi policy ảnh hưởng quota đang giữ hoặc key đang dùng, có rollout và
   đối soát. Không deploy chỉ vì formatter và unit test xanh.

## Kết quả kiểm tra của bản sửa

| Gate | Kết quả | Log |
| --- | --- | --- |
| TQD access/quota/report + cmd/wire | 9 package có test đạt; 5 package compile không có test | `tqd-tests.log` |
| Memory usage store | 1 package đạt | `memory-usage-tests.log` |
| User internal + checkout/subscription stores + seed | 19 package có test đạt | `user-tests.log` |
| Payment internal/provider/handler/wallet store | 10 package có test đạt | `payment-tests.log` |
| Gateway auth/metadata/response | 4 package đạt | `gateway-tests.log` |
| File content/HTTP | 2 package đạt | `file-tests.log` |
| Frontend payment/report/request identity | 6 suite, 26 test đạt | `frontend-tests.log` |
| Frontend commercial architecture guard | Đạt | `frontend-contract.log` |
| TQD quota race detector | Đạt | `tqd-race.log` |
| Payment wallet race detector | Đạt | `payment-race.log` |
| TQD read command fingerprint/replay | Đạt | `read-command-tests.log` và `infra/db` |

Log nằm trong `.tmp/development/code-review/`; bảng này chỉ tính package có test
thật, không gọi `[no test files]` là có coverage. Các test dùng mock/DryRun không
thay cho nghiệm thu PostgreSQL/Redis thật hoặc browser.

Để chạy lại backend, tại thư mục service dùng `go test -p 1 -timeout 90s` với
package đã ghi trong log/phạm vi bên trên. Máy dev ít RAM nên hạn chế số compiler
đồng thời (`GOMAXPROCS=2`). `wire/wire_gen.go` phải generate từ `wire/wire.go`,
không sửa tay. Với thay đổi về race, chạy thêm `go test -race -p 1` cho
`tqd-service/internal/usecase/quota` và `payment-service/internal/usecase/wallet`.

## Verification note — 2026-09-14

Sau khi thêm read command binding, các gate trực tiếp đạt: `go test` cho
TQD quota/usage/Redis adapter/command/migration contracts, `go test -race` cho
TQD quota và Payment wallet, các package User/Payment/Gateway/File, cùng 7 suite
frontend (37 tests) và commercial architecture guard. `go test ./...` toàn TQD
không thể gọi là xanh trong container hiện tại vì một test `httptest.NewServer`
bị môi trường cấm bind IPv6 loopback `[::1]:0`; đây là giới hạn môi trường, không
phải assertion nghiệp vụ thất bại. TypeScript web typecheck đầy đủ vẫn có lỗi
baseline ở nhiều module legacy; không quy các lỗi đó cho thay đổi quota khi chưa
có diff tương ứng.
