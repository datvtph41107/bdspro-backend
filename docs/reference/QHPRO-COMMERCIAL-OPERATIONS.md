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
- Client: `http://127.0.0.1:3001/dashboard/pricing`.

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
Authorization: <provider credential>
```

Local development muốn kiểm tra thanh toán phải dùng sandbox/provider callback
đi vào đúng endpoint trên. Không có lệnh Admin đổi thẳng order thành `paid` và
không có settlement shortcut trong public Make API.

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

Response loss hoặc callback lặp lại không được tạo settlement/subscription thứ
hai. Admin chỉ quan sát evidence và redrive fulfillment theo permission.

## 5. Quota, usage và báo cáo

Client tạo report bằng command identity ổn định. TQD kiểm tra entitlement và
quota của subject trước khi chấp nhận lệnh.

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
| Payment provider | `SEPAY_API_KEY`, callback HTTPS | callback hợp lệ, callback lặp không tạo settlement thứ hai |
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
publish tại `http://127.0.0.1:3001/dashboard/pricing`. Có thể đối chiếu public
payload:

```bash
curl -sS 'http://127.0.0.1:8000/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile'
```

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
  -H "Authorization: $PROVIDER_KEY" \
  -H 'Content-Type: application/json' \
  -d "{\"gateway\":\"LOCAL_PROVIDER\",\"transactionDate\":\"2026-09-06 12:00:00\",\"accountNumber\":\"LOCAL\",\"content\":\"$PAYMENT_REFERENCE\",\"transferType\":\"in\",\"description\":\"QHPRO purchase\",\"transferAmount\":$AMOUNT_MINOR,\"referenceCode\":\"$PAYMENT_REFERENCE\",\"id\":$PROVIDER_TRANSACTION_ID}"
```

Quan sát đến khi `progressState` thành `ACTIVE`, rồi đọc subscription/quota:

```bash
curl -sS "http://127.0.0.1:8000/v2/payment/commercial/orders/$ORDER_ID" \
  -H "Authorization: Bearer $USER_TOKEN"
curl -sS http://127.0.0.1:8000/v2/user/commercial/profile \
  -H "Authorization: Bearer $USER_TOKEN"
curl -sS http://127.0.0.1:8000/v2/tqd/commercial/generated-report-usage \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 8.4 Tiêu thụ quota và PDF

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

Phải thấy cùng order reference, subject profile, subscription active, usage
tăng đúng một đơn vị, report/PDF và notification thanh toán. Nếu fulfillment ở
`requires_review`, Admin chỉ được redrive command; settlement evidence không bị
viết lại.
