# QHPRO Commercial Operations

Đây là luồng vận hành duy nhất của QHPRO từ local development đến production.
Không có runtime `demo`, không seed user/order/subscription/quota/report và
không có API riêng để bỏ qua state machine nghiệp vụ.

## 1. Khởi động hệ thống

Lần đầu clone source:

```bash
make setup
make up
make status
make smoke
```

Docker chỉ chạy PostgreSQL/PostGIS, Redis và RabbitMQ. Các application service
chạy trực tiếp trên host. `make up` chạy versioned migrations, tạo hai identity
đăng nhập cố định `admin/admin123` và `0900000000/client123` khi
`QHPRO_ENVIRONMENT=development`, rồi khởi động runtime. Nó không tạo plan,
order, payment, subscription, quota, report hoặc dữ liệu quy hoạch giả.

Migrations được phép tạo schema, constraint, permission identity, reference
data và commercial catalog baseline đã được phê duyệt. Dữ liệu khách hàng và
giao dịch phải được tạo qua API/UI của owner tương ứng.

## 2. Khởi động Admin và Client

Giả sử ba repository được clone cạnh nhau trong cùng một workspace. Admin:

```bash
cd ../react-js-admin-bds
npm run dev -- --host 127.0.0.1 --port 5173 --strictPort
```

Client:

```bash
cd ../react-native-quy-hoach
QHPRO_WEB_DEV_PORT=3001 yarn web:local
```

Mở:

- Admin: `http://127.0.0.1:5173/admin/commercial/plans`.
- Client: `http://localhost:3001/dashboard/pricing` (dùng đúng hostname này
  xuyên suốt checkout để session browser không đổi origin khi SePay trả về).

Ở development, Admin đăng nhập bằng `admin/admin123`; Client chọn “Đăng nhập
bằng tài khoản và mật khẩu”, dùng `0900000000/client123` hoặc
`client/client123`. Cả hai đi qua contract Auth/User thật, session/JWT và IAM;
password được hash trong PostgreSQL User. Fresh installation ngoài development
dùng `make bootstrap-admin` đúng một lần với sáu biến `QHPRO_BOOTSTRAP_*` được
mô tả trong `THREE-SOURCE-ONBOARDING.md`. User Service thực hiện transaction và
từ chối khi root đã tồn tại; production không chạy development identity fixture
hoặc acceptance fixture.

## 3. Vòng đời gói

1. Admin có quyền `CATALOG_PLAN_MANAGE` tạo plan version ở trạng thái draft.
2. Admin cấu hình giá bằng minor units, entitlement, meter và operation policy.
3. Backend validate toàn bộ invariant trước khi publish.
4. Chỉ version đã publish và còn hiệu lực xuất hiện trên Client.
5. Checkout đóng băng plan version, giá, currency và terms checksum vào order.
6. Version đã publish không sửa/xóa; Admin chỉ có thể retire khỏi giao dịch mới.
7. Subscription đã mua giữ nguyên snapshot dù catalog thay đổi sau đó.

## 4. Thanh toán và kích hoạt

Client tạo order và payment attempt qua Gateway. Payment Service sinh reference
và sở hữu toàn bộ order/attempt/settlement/fulfillment state.

Trong production, ngân hàng/provider gửi callback đến:

```text
POST /v2/payment/sepay/webhook
Authorization: Apikey <provider credential>
```

Local development có hai boundary hợp lệ: sandbox/provider callback đi vào
đúng endpoint, hoặc operator dùng tác vụ Admin `ConfirmPaymentFunds` sau khi
đối chiếu mã tham chiếu/chứng từ. Tác vụ Admin có permission, idempotency và
durable audit evidence; nó không phải form ghi trực tiếp order status. Không có
settlement shortcut trong public Make API.

Sau khi callback được xác minh:

```text
Payment records provider evidence
  -> settlement durable
  -> fulfillment command
  -> User activates subscription
  -> entitlement becomes readable
  -> outbox publication
  -> Notification inbox/effect
```

Response loss, callback lặp lại hoặc command Admin retry không được tạo
settlement/subscription thứ hai. Admin quan sát evidence, xác nhận đối soát thủ
công khi có căn cứ và redrive fulfillment theo permission.

## 5. Quota, usage và báo cáo

Catalog QHPRO cấu hình ba quyền lợi được TQD thực thi thật:

| Quyền lợi trên Admin/Client | Operation owner | Durable meter |
| --- | --- | --- |
| Phân tích khu vực vẽ trên bản đồ | `workspace.area.analyze` | `workspace.area_analysis.accepted` |
| Xem hồ sơ thửa đất chi tiết | `workspace.parcel.detail.view` | `workspace.parcel_detail.viewed` |
| Xuất báo cáo quy hoạch PDF | `workspace.report.generate` | `workspace.report_generation.accepted` |

Gợi ý tìm kiếm, bản đồ nền/PMTiles và quick-view công khai không tiêu hao
quota. Hai thao tác đọc nâng cao chỉ ghi usage sau khi handler trả dữ liệu thành
công. Client gửi `Idempotency-Key` ổn định trong một cửa sổ thao tác, nên React
render lại, HTTP retry hoặc refresh ngay không trừ lần hai. Nếu handler lỗi,
reservation được hủy. TQD kiểm tra entitlement và quota của đúng subject trước
khi chạy action.

Riêng report là command tạo business state. Client tạo report bằng command
identity ổn định và TQD ghi:

```text
Report + UsageEvent + Job
```

được ghi trong một transaction. Redis chỉ làm admission/runtime projection;
PostgreSQL là durable truth. Worker tạo PDF, File Service lưu artifact và Client
poll report cho đến khi nhận URL PDF. Admin đọc subscription, usage, report,
payment và fulfillment projection qua API của owner, không query chéo database.

## 6. Test không phải runtime vận hành

```bash
make test-e2e
```

là integration proof. Nó được phép tạo identity có prefix `acceptance` từ
`integration-test/.env` và chỉ chạy khi `QHPRO_ENVIRONMENT` là
`development`, `test`, `ci` hoặc `acceptance`. Guard sẽ từ chối production.

Fixture test không nằm trong migration, không được application service load và
không phải cách developer tạo dữ liệu để làm việc hàng ngày.

### Kênh thông báo được nghiệm thu

Payment event được Notification nhận qua RabbitMQ và ghi Inbox, notification
hiển thị trong ứng dụng, cùng delivery intent trong một transaction. Đây là
bằng chứng bền vững mặc định của luồng thanh toán.

- In-app notification hoạt động không phụ thuộc Firebase.
- Push mobile chỉ được gửi khi operator cấu hình
  `NOTIFICATION_FIREBASE_CREDENTIAL_FILE` và thiết bị đã đăng ký token.
- Email/Gmail thanh toán chưa phải channel đã triển khai; không được ghi nhận
  là PASS chỉ vì Inbox hoặc push intent tồn tại.

Local không có Firebase vẫn phải quan sát được payment notification trên API/UI
và thấy log `delivery-worker status=disabled reason=firebase-not-configured`.

## 7. Production contract

Production sử dụng cùng binary, migrations, RPC/HTTP contracts và state machine
như development. Khác biệt chỉ nằm ở operator-owned config, secrets, database,
provider credentials và dữ liệu khách hàng thật.

Production startup không chạy `test-e2e`, không bật Compose profile
`acceptance` và không load `integration-test/.env`.

Các kết nối bên ngoài phải do operator cung cấp, không hard-code vào source:

| Boundary | Input vận hành | Bằng chứng tối thiểu |
| --- | --- | --- |
| Payment provider | SePay merchant/secret, IPN secret, callback HTTPS (`SEPAY_API_KEY` chỉ cho webhook ngân hàng cũ) | hosted QR/thẻ hoạt động; IPN hợp lệ và IPN lặp không tạo settlement thứ hai |
| Google OAuth | client ID/secret và redirect URI theo môi trường | login/callback tạo đúng một session |
| Mobile push | Firebase service credential ngoài repository | device thật nhận push, retry không tạo effect trùng |
| Planning data | import dataset và PMTiles do TQD sở hữu | map query và report dùng cùng version dữ liệu |

Giá trị development và production khác nhau, nhưng binary, API, migration và
state machine phải giống nhau.

## 8. Runbook browser + terminal trên local

Đây là cách quan sát cùng một giao dịch từ Client, terminal và Admin. Các lệnh
dưới đây gọi public Gateway giống browser; chúng không ghi database trực tiếp.
Đặt giá trị thật vừa nhận được vào biến shell, không lưu token/password vào
source.

### 8.1 Identity

Đăng nhập root operator trên Admin tại `http://127.0.0.1:5173/login`. Sau đó có
thể tạo user ở màn quản lý người dùng, hoặc gọi đúng Admin API:

```bash
ADMIN_TOKEN='paste-access-token-from-admin-login'
curl -sS -X POST http://127.0.0.1:8000/v2/user/admin/user/new \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"username":"customer.one","password":"replace-with-at-least-12-characters","fullName":"Khách hàng Một","email":"customer.one@example.com"}'
```

Đăng nhập user vừa tạo qua public contract và chép `accessToken`, `profileId`
từ response:

```bash
curl -sS -X POST http://127.0.0.1:8000/v2/auth/password/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"customer.one","password":"replace-with-at-least-12-characters","version":"0.1.0","platform":"WEB","os":"Linux","deviceName":"Local browser"}'

USER_TOKEN='paste-user-access-token'
PROFILE_ID='paste-profile-id'
```

### 8.2 Catalog, checkout và payment attempt

Admin tạo/sửa/validate/publish gói tại
`http://127.0.0.1:5173/admin/commercial/plans`. Client chỉ thấy version đã
publish tại `http://localhost:3001/dashboard/pricing`. Có thể đối chiếu public
payload:

```bash
curl -sS 'http://127.0.0.1:8000/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile'
```

Màn “Thiết kế gói mới” có ba mẫu vận hành để tạo nhanh nhưng không seed dữ liệu:

| Mẫu | Giá tham chiếu | Phân tích khu vực | Hồ sơ chi tiết | Báo cáo |
| --- | ---: | ---: | ---: | ---: |
| QHPRO Khởi đầu | 49.000đ / 30 ngày | 20 | 30 | 3 |
| QHPRO Tiêu chuẩn | 149.000đ / 30 ngày | 100 | 150 | 15 |
| QHPRO Chuyên nghiệp | 399.000đ / 30 ngày | Không giới hạn | Không giới hạn | 100 |

Mỗi mẫu vẫn đi qua draft → validate → publish. Admin có thể đổi giá và hạn mức
trước khi publish; Client không hard-code ba con số này.

Client gửi `planCode`, không gửi giá/quota có tính authority. Dùng một command
key mới cho mỗi ý định mua và giữ nguyên key khi retry do mất response:

```bash
CHECKOUT_KEY='checkout-customer-one-001'
curl -sS -X POST http://127.0.0.1:8000/v2/user/checkout \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Idempotency-Key: $CHECKOUT_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"planCode":"qhpro.basic","subjectKind":"profile"}'

ORDER_ID='paste-order-id'
ATTEMPT_KEY='payment-attempt-customer-one-001'
curl -sS -X POST http://127.0.0.1:8000/v2/user/checkout/payment-attempts \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Idempotency-Key: $ATTEMPT_KEY" \
  -H 'Content-Type: application/json' \
  -d "{\"orderId\":$ORDER_ID,\"subjectKind\":\"profile\",\"method\":\"bank_transfer\"}"
```

Chép `reference` và `amountMinor` từ order/progress. `amountMinor` là số nguyên
minor units, không phải số thực.

### 8.3 Provider callback và activation

Trong production callback do provider gửi. Trên local, operator mô phỏng đúng
provider boundary bằng credential trong `payment-service/.env`; không có Admin
API “đổi thành paid”. Chọn một transaction ID chưa từng dùng:

```bash
PROVIDER_KEY='paste-SEPAY_API_KEY-from-operator-config'
PAYMENT_REFERENCE='paste-order-reference'
AMOUNT_MINOR='paste-order-amount-minor'
PROVIDER_TRANSACTION_ID='202609060001'

curl -sS -X POST http://127.0.0.1:8000/v2/payment/sepay/webhook \
  -H "Authorization: Apikey $PROVIDER_KEY" \
  -H 'Content-Type: application/json' \
  -d "{\"gateway\":\"LOCAL_PROVIDER\",\"transactionDate\":\"2026-09-06 12:00:00\",\"accountNumber\":\"LOCAL\",\"content\":\"$PAYMENT_REFERENCE\",\"transferType\":\"in\",\"description\":\"QHPRO purchase\",\"transferAmount\":$AMOUNT_MINOR,\"referenceCode\":\"$PAYMENT_REFERENCE\",\"id\":$PROVIDER_TRANSACTION_ID}"
```

Để dùng checkout QR/thẻ sandbox thật, đăng ký tại `https://my.sepay.vn`, bật
**Cổng thanh toán → Sandbox** rồi đặt `SEPAY_MERCHANT_ID`,
`SEPAY_MERCHANT_SECRET_KEY` và `SEPAY_IPN_SECRET_KEY` vào file `.env` bị Git
bỏ qua. Phương thức hiển thị trên Client do Payment công bố, không được hard-code:

```dotenv
SEPAY_BANK_TRANSFER_ENABLED=true
SEPAY_CARD_ENABLED=true
SEPAY_SUCCESS_URL=http://localhost:3001/dashboard/pricing?payment=success
SEPAY_ERROR_URL=http://localhost:3001/dashboard/pricing?payment=error
SEPAY_CANCEL_URL=http://localhost:3001/dashboard/pricing?payment=cancel
```

Chỉ bật `SEPAY_CARD_ENABLED=true` sau khi merchant sandbox/production đã được
SePay cấp phương thức `CARD`. Có thể kiểm tra contract mà Client đang dùng tại
`GET /v2/user/checkout/payment-methods` bằng access token của user.

Cấu hình IPN trên dashboard là URL HTTPS công khai:

```text
https://<development-tunnel>/v2/payment/sepay/ipn
```

Localhost không nhận được callback từ SePay; dùng một tunnel HTTPS và giữ
Gateway chạy ở cổng 8000. Redirect về client chỉ phục vụ UX. Payment chỉ active
gói từ IPN, không active từ query `payment=success` trên browser.

Quan sát đến khi `progressState` thành `ACTIVE`, rồi đọc subscription/quota:

```bash
curl -sS "http://127.0.0.1:8000/v2/payment/commercial/orders/$ORDER_ID" \
  -H "Authorization: Bearer $USER_TOKEN"
curl -sS http://127.0.0.1:8000/v2/user/commercial/profile \
  -H "Authorization: Bearer $USER_TOKEN"
curl -sS http://127.0.0.1:8000/v2/tqd/commercial/generated-report-usage \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 8.4 Tiêu thụ ba quota và PDF

Sau khi gói `ACTIVE`, mở Client `http://localhost:3001/ban-do`:

1. Bật công cụ khoanh vùng, vẽ một polygon và chạy phân tích. Một kết quả thành
   công tăng `workspace.area_analysis.accepted` đúng một đơn vị.
2. Chọn một thửa và mở “Hồ sơ chi tiết”. Một response chi tiết thành công tăng
   `workspace.parcel_detail.viewed` đúng một đơn vị.
3. Trong hồ sơ thửa chọn xuất báo cáo. Report được chấp nhận tăng
   `workspace.report_generation.accepted` và worker tạo PDF.

Mở lại cùng vùng/thửa ngay hoặc để HTTP retry phải giữ nguyên số đã dùng. Kiểm
tra cả ba projection từ cùng endpoint Client đang dùng:

```bash
curl -sS http://127.0.0.1:8000/v2/tqd/commercial/usage \
  -H "Authorization: Bearer $USER_TOKEN"
```

Admin quan sát từng meter bằng tên nghiệp vụ tại
`http://127.0.0.1:5173/admin/commercial/usage`.

`REPORT_KEY` đại diện một ý định tạo báo cáo. Gửi lại cùng body + key phải trả
cùng report và không tăng `durableUsed` lần hai:

```bash
REPORT_KEY='report-customer-one-001'
curl -sS -X POST http://127.0.0.1:8000/v2/tqd/map-workspace/reports \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Idempotency-Key: $REPORT_KEY" \
  -H 'Content-Type: application/json' \
  -d "{\"reportType\":2,\"entityType\":2,\"entityId\":$PROFILE_ID,\"title\":\"Báo cáo quy hoạch của Khách hàng Một\",\"location\":{\"address\":\"Hà Nội\",\"province\":\"Hà Nội\"},\"spatial\":{\"centroid\":{\"lat\":21.0278,\"lon\":105.8342},\"bounds\":{\"minLon\":105.83,\"minLat\":21.02,\"maxLon\":105.84,\"maxLat\":21.03}}}"

REPORT_ID='paste-report-id'
curl -sS "http://127.0.0.1:8000/v2/tqd/map-workspace/reports/$REPORT_ID" \
  -H "Authorization: Bearer $USER_TOKEN"
curl -sS http://127.0.0.1:8000/v2/tqd/commercial/generated-report-usage \
  -H "Authorization: Bearer $USER_TOKEN"
```

Khi report trả `assets.pdfUrl`, mở URL đó để kiểm tra `%PDF-`/nội dung. Client
hiển thị lịch sử report và quota còn lại từ cùng API.

### 8.5 Đối chiếu Admin

Không sửa trạng thái bằng tay. Mở lần lượt:

```text
/admin/commercial/payments
/admin/commercial/operations
/admin/commercial/subscriptions
/admin/commercial/usage
/admin/commercial/reports
/admin/user/<PROFILE_ID>
```

Phải thấy cùng order reference, subject profile, subscription active, ba usage
meter tăng theo đúng thao tác thành công, report/PDF và notification thanh toán. Nếu fulfillment ở
`requires_review`, Admin chỉ được redrive command; settlement evidence không bị
viết lại.
