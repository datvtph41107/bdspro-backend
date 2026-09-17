# Audit thiết kế quota và cách làm việc

Tài liệu này ghi nhận cách hệ thống đang vận hành, các câu hỏi phải trả lời
trước khi mở rộng và những điểm cần xử lý để tránh nợ kỹ thuật. Đây là tài liệu
ra quyết định; không dùng để thay thế contract của từng service.

## Kết luận hiện tại

Trục thương mại đang chạy được: catalog → checkout snapshot → Payment
settlement → User subscription/entitlement → TQD usage → File output. Ba
operation được khai báo trong contract là (không đồng nghĩa đã nghiệm thu
với dữ liệu GIS production):

```text
workspace.area.analyze
workspace.parcel.detail.view
workspace.report.generate
```

Tuy nhiên database và code vẫn có dấu hiệu chuyển tiếp cần xử lý trước khi gọi
là quota platform tổng quát.

## Phát hiện theo mức độ

### P0 — phải chốt trước khi mở thêm chức năng

**Hai đường accounting trong TQD.** `generatedreport/application.Service` có
transaction acceptance riêng; `quota/ActionService` bọc các RPC đọc parcel và
ghi usage sau callback. Hai đường này dùng chung concept reserve/commit nhưng
khác transaction boundary và recovery. Đây là hai *policy adapter* hợp lệ cho
hai loại hành động, operation mới phải chọn điểm ghi nhận phù hợp; chỉ thêm adapter khi có
transaction boundary khác được giải thích và kiểm thử. Mỗi operation mới phải
chọn đúng một adapter và ghi rõ lý do.

**Subject profile/organization chưa có policy đầy đủ.** Actor là profile,
subject có thể là organization. Cần chốt quota là pool chung hay quota riêng,
report thuộc profile hay organization, và thành viên nào được đọc output.

**Điểm tính giá trị chưa thống nhất cho mọi loại chức năng.** Report tính tại
acceptance bất đồng bộ; parcel read tính sau callback thành công. Chức năng mới
phải tuyên bố rõ “thành công” trước khi viết code.

### P1 — rủi ro dữ liệu và vận hành

**`quota_pools` đang là bảng mồ côi.** Migration tạo bảng có khóa
`subject_type/subject_id/meter_code/period_start`, nhưng không có store/usecase
nào đọc hoặc ghi bảng này. Runtime hiện dùng Redis cho reserve/commit và
`quota_usage_events` cho durable usage. Giữ bảng mà không owner sẽ gây hiểu nhầm
và khiến người bảo trì sửa nhầm nguồn. Cần quyết định: bỏ qua và đánh dấu legacy,
hay chuyển runtime sang bảng này bằng migration/benchmark/acceptance đầy đủ.
Không được tự xóa production table khi chưa inventory dữ liệu.

**`quota_usage_events` chưa FK tới catalog/subscription.** Bảng lưu snapshot
`subscription_id`, `plan_code`, `plan_version`, `policy_version`,
`limit_snapshot`, nhưng hiện chủ yếu dùng check/nullability. Đây là snapshot đi qua biên database của User và TQD; PostgreSQL không
đặt FK trực tiếp giữa hai database. Retire không yêu cầu xóa version đã bán.
Cần kiểm tra bằng chứng ở application và có đối soát; chưa có bằng chứng cho
phép xóa catalog lịch sử.

**Admin usage mới lọc `subject_type = profile`.** Nếu organization sử dụng quota,
dashboard này không hiển thị đầy đủ pool tổ chức. Phải mở projection theo subject
hoặc ghi rõ màn hình chỉ dành cho profile trước khi bật organization billing.

### P2 — khả năng bảo trì

- `ProfileOperations` là danh sách operation hiển thị; catalog database mới là
  authority. Thêm operation mà quên danh sách này sẽ khiến Client không hiển thị.
- Tên `ActionService` quá chung; nên đổi thành tên nêu rõ nó là adapter cho
  metered read actions, còn report acceptance là adapter bất đồng bộ.
- Một số query GIS legacy còn đọc metadata geometry trực tiếp thay vì projection
  metadata owner. Không thêm cột song song để chữa lỗi; refactor theo schema
  owner sau khi có dataset kiểm thử.
- `period = day/calendar_month/lifetime` chưa có timezone/reset contract; hiện
  fail closed là đúng, nhưng chưa phải tính năng đã bán được.

## Mô hình database cần giữ

### Catalog User Service

`plan version` lưu terms bất biến: giá, tier, thời hạn, feature access,
allowance và operation policy. Đây là nơi xác định *được phép làm gì* và *được
bao nhiêu*.

### Subscription User Service

Subscription lưu subject, plan version, period, access window và lifecycle.
Mọi entitlement query phải lọc theo thời gian. Checkout phải bỏ qua subscription
đã hết hạn; settlement phải materialize `expired` khi có giao dịch tiếp theo.

### Usage TQD

`quota_usage_events` là durable ledger theo:

```text
usage_key = hash(subject + operation + meter + command)
period_start/period_end = bằng chứng chu kỳ, không nằm trong usage_key
```

`usage_key` là idempotency authority của usage. Redis là admission projection
(`used`, `reserved`), không phải sổ cái. Reconciliation sửa Redis từ durable
usage, không tạo event mới.

### Resource/File

Report/job và file phải giữ command identity, subject/profile owner,
subscription evidence và output URL. Mở lại hoặc tạo lại output không tạo usage
mới.

## Câu hỏi bắt buộc trước khi thêm chức năng

1. Actor là ai? Subject trả tiền và chịu quota là profile hay organization?
2. Một “lần sử dụng” là command nào và thành công ở checkpoint nào?
3. Lỗi validation, timeout, worker failure và retry ảnh hưởng quota ra sao?
4. Kết quả thuộc profile hay organization? Ai được đọc/chia sẻ/xóa?
5. Meter có reset theo subscription cycle, ngày, tháng hay lifetime? Timezone nào?
6. Upgrade, cancel, expiry và repurchase tác động tới usage cũ thế nào?
7. Nếu Redis down sau durable acceptance, reconciliation chạy ở đâu và SLA bao lâu?
8. Admin xem được profile và organization usage bằng cùng projection nào?
9. Có cần refund/credit khi output không được tạo không? Ai quyết định?
10. Dữ liệu đầu vào có owner và retention policy nào?

Nếu chưa trả lời được một câu, chức năng chỉ nên tồn tại dưới dạng internal
feature flag/draft, chưa publish thành quyền trả phí.

## Quy trình phát triển chuẩn

```text
1. Viết business decision record và acceptance examples.
2. Chọn owner và adapter accounting duy nhất.
3. Định nghĩa operation/feature/meter/period.
4. Thêm catalog validation và Admin form/projection.
5. Thêm entitlement mapping ở User.
6. Thêm handler tại TQD với command identity.
7. Chọn transaction boundary: synchronous action hoặc async acceptance.
8. Ghi durable usage và recovery/reconciliation.
9. Thêm Client state/error/output.
10. Viết test invalid, denied, exhausted, retry, timeout, upgrade, cancel, expiry.
11. Chạy runtime acceptance bằng API thật; không chèn order/usage bằng SQL.
12. Chỉ sau khi gate pass mới publish plan version.
```

## Tiêu chí “mua giá trị thực”

Một giao dịch chỉ được coi là tạo ra giá trị khi có thể lần theo:

```text
payment evidence
  -> frozen order terms
  -> active subscription
  -> entitlement decision
  -> accepted command
  -> durable usage event
  -> report/job or synchronous result
  -> real output/resource
```

Thiếu một mắt xích phải trả trạng thái rõ ràng và giữ khả năng retry/reconcile;
không được làm xanh bằng cách cộng quota hoặc sửa database thủ công.

## Gate tối thiểu

```bash
make status smoke
make runtime-value-chain-acceptance
make runtime-admin-acceptance
```

Ngoài test backend, phải kiểm tra Admin browser nhìn thấy draft/version/usage,
Client nhìn thấy đúng quyền và giới hạn, và Network chỉ gọi endpoint thuộc owner
đã định nghĩa. Các kết quả nghiệm thu hiện tại nằm trong
`.tmp/development/acceptance-current/` và được mô tả ở
`react-native-quy-hoach/docs/quota-operational-acceptance.md`.


Kết quả review và sửa lỗi chi tiết: [QUOTA-CODE-REVIEW.md](QUOTA-CODE-REVIEW.md).
