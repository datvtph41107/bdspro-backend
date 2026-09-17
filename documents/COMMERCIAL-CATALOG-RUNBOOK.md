# Khai báo gói, chức năng và quota

Catalog thương mại thuộc User Service. Admin tạo và xuất bản version; Client
chỉ đọc catalog và trạng thái quota. Không khai báo hạn mức trong React,
Payment hoặc TQD.

## Mô hình khai báo

Một `plan version` có định danh (`productCode`, `planCode`, `version`,
`displayName`, `subjectScope`), thương mại (`tierRank`,
`subscriptionTermDays`, `prices[]`), entitlement `feature_access` và
entitlement `usage_allowance`. Mỗi chức năng phải có cùng một `featureCode`,
`meterCode` và operation binding; `unitsPerAction` là số đơn vị bị tiêu thụ khi
thao tác được chấp nhận.

```json
{
  "featureCode": "workspace.report.generate",
  "meterCode": "workspace.report_generation.accepted",
  "entitlements": [
    {"code":"workspace.report.feature","kind":"feature_access","featureCode":"workspace.report.generate","period":"none"},
    {"code":"workspace.report.allowance","kind":"usage_allowance","meterCode":"workspace.report_generation.accepted","amount":15,"period":"subscription_cycle"}
  ],
  "operations": [
    {"code":"workspace.report.generate","featureCode":"workspace.report.generate","meterCode":"workspace.report_generation.accepted","unitsPerAction":1}
  ]
}
```

Các operation hiện được thực thi là `workspace.area.analyze`,
`workspace.parcel.detail.view` và `workspace.report.generate`. Bản đồ nền, PBF,
tìm kiếm và xem trước công khai không tiêu hao quota. Thêm operation mới phải
có handler TQD, điểm acceptance/usage, projection quota, client và test; thêm
dòng catalog một mình sẽ tạo quyền không có thực thi.

## Quy trình Admin trên browser

1. Mở **Thương mại → Gói**, tạo hoặc mở bản nháp.
2. Kiểm tra `planCode`, version, display name, scope, `tierRank`, thời hạn và
   giá VND.
3. Trong **Quyền lợi trong gói**, bật từng chức năng. Chọn **Giới hạn theo chu
   kỳ** để nhập lượt hoặc **Không giới hạn**. Tắt chức năng sẽ bỏ cả quyền,
   allowance và operation tương ứng.
4. Đối chiếu **Xem trước trên Client**, bấm **Lưu bản nháp** rồi **Kiểm tra điều
   khoản**. Backend tính `termsChecksum`.
5. Bấm **Xuất bản cho Client**. Draft chuyển active và terms bị khóa. Khi đổi
   giá/quyền, tạo version mới; retire chỉ ngăn mua mới và giữ evidence cũ.

API tương ứng:

```text
POST /v2/user/admin/catalog/plan-versions
POST /v2/user/admin/catalog/plan-versions/{id}:validate
POST /v2/user/admin/catalog/plan-versions/{id}:publish
POST /v2/user/admin/catalog/plan-versions/{id}:retire
```

`tierRank` là thứ tự nghiệp vụ để phân loại upgrade, không suy ra từ giá hoặc
tên hiển thị.

## Từ checkout đến quota

1. Client đọc `GET /v2/user/commercial/plans`; chỉ version active trong cửa sổ
   thời gian và scope phù hợp được trả về.
2. `POST /v2/user/checkout` với idempotency key. User lưu snapshot terms trước
   khi gọi Payment; Payment giữ `amountMinor`, reference và QR.
3. Webhook SePay xác thực giao dịch. Đúng amount/reference mới settlement; sai
   credential trả `401`, sai amount không kích hoạt.
4. User tạo hoặc nâng cấp subscription. Upgrade giữ period và usage; cùng tier
   bị từ chối; downgrade schedule chưa được triển khai.
5. TQD kiểm tra entitlement. Khi acceptance thành công, report/usage/job được
   ghi cùng command identity; retry cùng key không trừ lần hai. Worker tạo file,
   File Service sở hữu nội dung.
6. Client đọc profile, usage và report status; không tự cộng quota. `ready` có
   PDF URL thật thì hiển thị **Mở PDF**.

Hủy thu hồi quyền trả phí ngay, giữ evidence và không tự hoàn tiền. Mua lại tạo
   chu kỳ mới. Subscription quá hạn được materialize thành `expired` trước
   settlement tiếp theo; entitlement lọc theo `access_until`.

## Checklist nghiệm thu

Admin draft/validate/publish; Client hiển thị đúng quyền và giá; checkout replay
cùng order; webhook sai credential/sai amount không cấp quyền; từng operation trả
đúng `allowed`, `limit`, `durable_used`, `remaining`; exhausted/replay không tính
trùng; PDF ready trả HTTP 200 `%PDF`; upgrade giữ period/usage; cancel trả free;
repurchase reset usage; Admin đối chiếu orders, fulfillments, subscriptions,
usage và audit events.

Chạy gate:

```bash
make runtime-value-chain-acceptance
make runtime-admin-acceptance
```

Gate dùng webhook mô phỏng local và không chèn order, subscription, report hoặc
usage bằng SQL.
