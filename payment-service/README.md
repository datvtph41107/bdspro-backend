# Payment Service

Payment Service sở hữu đơn hàng, hướng dẫn thanh toán, bằng chứng nhận tiền,
đối soát và tác vụ kích hoạt quyền lợi. Luồng đăng ký gói bằng chuyển khoản là:
**Payment sinh chuỗi VietQR → Client render QR → SePay gửi Bank Webhook →
Payment đối soát → User kích hoạt gói → Client đọc trạng thái mới**.

README này mô tả code và quy trình vận hành; không xác nhận rằng một deployment
hoặc giao dịch cụ thể đã chạy thành công. Bằng chứng hoàn tất phải lấy từ
settlement, fulfillment và subscription của đúng môi trường.

## 1. Ownership và nguồn dữ liệu chuẩn

| Thành phần | Trách nhiệm |
| --- | --- |
| User Service | Catalog/version gói, checkout của người dùng, subscription và entitlement |
| Payment Service | Order với điều khoản đã chốt, Payment Attempt, Settlement, Fulfillment, command audit và completion Outbox |
| Gateway Service | HTTP boundary, chuyển credential vào metadata RPC được bảo vệ, phản hồi theo contract SePay |
| Client | Chọn gói, gửi checkout command, render payload QR và đọc tiến trình |
| TQD Service | Kiểm tra quyền và ghi nhận sử dụng quota cho phân tích, hồ sơ thửa, báo cáo |
| Notification Service | Xử lý thông báo; trạng thái gửi thông báo có lifecycle riêng |

Payment còn có capability ngân hàng/ví trong `internal/domain/bank`,
`internal/domain/wallet` và các usecase tương ứng. Mã nạp ví `DEP...` thuộc
capability ví; đăng ký gói `QHP-...` đi qua commercial settlement. Không dùng
lịch sử ví để suy ra subscription đã được kích hoạt.

Các nguyên tắc khi sửa code:

- Giá và quyền lợi đến từ catalog của User, được đóng băng trong Order;
  không lấy số tiền do Client tự khai làm giá thanh toán.
- Định dạng mã đơn và chuẩn hóa nội dung ngân hàng có một owner:
  [`internal/domain/payment/reference.go`](internal/domain/payment/reference.go).
- Provider adapter tạo hướng dẫn thanh toán và chuyển dữ liệu provider vào
  domain; quyết định nhận tiền nằm ở settlement usecase.
- Webhook và command đối soát của operator dùng cùng settlement usecase.
- Chỉ User thay đổi subscription/entitlement; Payment gọi qua RPC.
- Chỉ TQD sở hữu bộ đếm sử dụng quota. Payment không ghi vào database User/TQD.
- QR, browser redirect và nút làm mới không phải bằng chứng đã nhận tiền.

## 2. Luồng đăng ký gói

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gateway
    participant U as User Service
    participant P as Payment Service
    participant S as SePay
    C->>G: Checkout và tạo Payment Attempt
    G->>U: Xác thực người dùng, chọn điều khoản gói
    U->>P: Tạo Order và Payment Attempt
    P-->>C: Qua User/Gateway: reference + VietQR EMV
    Note over C,S: Client render QR; người dùng chuyển khoản
    S->>G: POST /v2/payment/sepay/webhook
    G->>P: Payload + credential digest
    P->>P: Đối soát, lưu Settlement và Fulfillment
    P-->>S: Qua Gateway: 200 + success:true
    P->>U: Fulfillment gọi ApplySettlement
    U-->>P: Kết quả kích hoạt subscription/entitlement
    P->>P: Hoàn tất Fulfillment và ghi completion Outbox
    C->>G: Poll tiến trình đơn
    G->>P: GetMyPaymentOrder
    P-->>C: Qua Gateway: ACTIVE khi fulfillment hoàn tất
```

Client poll mỗi 3 giây chỉ khi người dùng đang mở modal thanh toán, màn hình
Pricing còn focus, tab/app đang active và đơn ở `AWAITING_PAYMENT` hoặc
`ACTIVATING`. Modal đóng và polling dừng khi người dùng rời luồng, đóng modal,
hoặc đơn sang `ACTIVE`/`ACTION_REQUIRED`. Thanh toán vẫn được backend xử lý khi
Client không poll; phiên còn chờ được hiển thị bằng nút tiếp tục thay vì tự mở.

### Các trạng thái cần phân biệt

| Dữ liệu | Trạng thái/ý nghĩa |
| --- | --- |
| Order | `pending_funds`, `funds_confirmed`, `requires_review` |
| Payment Attempt | Hướng dẫn và lần thử thanh toán; `pending_action` không có nghĩa đã nhận tiền |
| Settlement | `funds_confirmed`, `unmatched`, `requires_review`, `late` |
| Fulfillment | `pending`, `running`, `retry`, `completed`, `requires_review` |
| Client progress | `AWAITING_PAYMENT`, `ACTIVATING`, `ACTIVE`, `ACTION_REQUIRED` |

`ACTIVE` trong tiến trình Payment được suy ra từ fulfillment `completed`.
Để kiểm tra gói còn hiệu lực và quota hiện tại, đọc subscription/entitlement
của User và usage của TQD.

HTTP `200` với `{"success":true}` chỉ xác nhận Bank Webhook đã được xử lý.
Settlement sai số tiền hoặc không tìm thấy đơn vẫn có thể được lưu và trả ACK;
ACK không bảo đảm có fulfillment hoặc gói đã active. Khi đối soát, đọc riêng
`settlement_status`, `review_reason` và `fulfillment_status`.

## 3. Contract VietQR và mã chuyển khoản

Với `SEPAY_BANK_TRANSFER_FLOW=vietqr`, provider gọi `vietqr.Generate(...)`
từ BIN, số tài khoản, số tiền VND và mã đơn. `nextAction.kind=display_qr` mang
`qrPayload` là chuỗi EMV đầy đủ, bắt đầu bằng `000201`, có checksum CRC.

Client chuyển nguyên chuỗi vào QR renderer. Không tự ghép BIN/số tiền/mã đơn,
không sinh mã đơn và không hỗ trợ nhánh ảnh PNG, Base64 hay URL ảnh QR.
`qrPayload` được lưu trong Payment Attempt để replay cùng hướng dẫn.

Contract nội bộ ở
[`commerce_internal.proto`](../shared/protobuf/schema/payment/commerce_internal.proto);
contract checkout công khai ở
[`checkout.proto`](../shared/protobuf/schema/user/checkout.proto).

Mã chuẩn có dạng `QHP-` + 16 ký tự hex viết hoa. Ví dụ minh họa:

```text
QHP-0123456789ABCDEF
```

Infrastructure sinh 8 byte ngẫu nhiên; domain format mã chuẩn. Khi webhook
đến, domain nhận diện mã trong `content`, sau đó `referenceCode` nếu cần,
và chuẩn hóa chữ hoa cùng dấu phân cách:

```text
BankAPINotify QHP0123456789ABCDEF → QHP-0123456789ABCDEF
```

Việc ngân hàng bỏ dấu `-` là chuẩn hóa dữ liệu đầu vào, không tạo định dạng
lưu trữ hoặc parser thứ hai. `referenceCode` của ngân hàng như `FT...` không
phải mã đơn QHPro. Trường `code` trên SePay có thể null; handler hiện nhận diện
mã từ `content`/`referenceCode`.

## 4. Runtime và phụ thuộc

Entrypoint: `main.go` → Cobra `grpc` →
[`cmd/grpc/runtime.go`](cmd/grpc/runtime.go).

Process mặc định chạy gRPC, fulfillment worker, completion Outbox publisher
và các job ví/dashboard. Các actor dùng cùng lifecycle của Payment; không cần
khởi động một fulfillment service riêng.

| Cấu hình | Mục đích / mặc định |
| --- | --- |
| `PAYMENT_DATABASE_URL` | DSN PostgreSQL của Payment; bắt buộc có cấu hình hợp lệ |
| `PAYMENT_RABBIT_URL` | RabbitMQ; bắt buộc khi khởi động runtime |
| `PAYMENT_GRPC_ADDRESS` | Override địa chỉ gRPC; mặc định host/port là `0.0.0.0:8205` |
| `QHPRO_USER_GRPC_ADDRESS` | User RPC; `localhost:8201` |
| `QHPRO_AUTH_GRPC_ADDRESS` | Auth RPC; `localhost:8216` |
| `QHPRO_NOTIFICATION_GRPC_ADDRESS` | Notification RPC; `localhost:8204` |
| `PAYMENT_ORDER_TTL` | Thời hạn đơn mới; `30m` |
| `PAYMENT_FULFILLMENT_POLL_INTERVAL` | Chu kỳ lấy việc; `1s` |
| `PAYMENT_FULFILLMENT_CLAIM_LEASE` | Lease giữ việc; `30s` |
| `PAYMENT_FULFILLMENT_RETRY_DELAY` | Khoảng chờ retry; `5s` |
| `PAYMENT_EVENT_EXCHANGE` | Rabbit exchange; `qhpro.payment.events` |
| `PAYMENT_OUTBOX_POLL_INTERVAL` | Chu kỳ publisher; `1s` |
| `PAYMENT_OUTBOX_CLAIM_LEASE` | Lease publisher; `30s` |
| `PAYMENT_OUTBOX_RETRY_DELAY` | Khoảng chờ publish lại; `5s` |

Defaults ở [`config/runtime.yml`](config/runtime.yml) và
[`config/config.go`](config/config.go); ENV override giá trị runtime.
`CONFIG_FILE` chỉ dùng khi chủ động mount cấu hình khác. Native Make commands
nạp `.env` root trước rồi `.env` của service. Compose truyền ENV cho container;
sửa file local không tự thay cấu hình deployment từ xa.

Fulfillment dùng lease/claim version và effect key của Order để retry qua User.
Khi fulfillment hoàn tất, completion event được ghi vào Outbox cùng transaction;
publisher chuyển event sang RabbitMQ. Gói đã active không đồng nghĩa thông báo
đã được gửi thành công.

## 5. Chọn chế độ SePay và đúng loại khóa

| Chế độ | Cấu hình lựa chọn | Endpoint xác nhận |
| --- | --- | --- |
| VietQR + mô phỏng ngân hàng | `vietqr`, `sandbox`, `test` | `/v2/payment/sepay/webhook` |
| VietQR + tiền thật | `vietqr`, `production`, `live` | `/v2/payment/sepay/webhook` |
| Hosted Checkout | `hosted_checkout`, `sandbox` hoặc `production` | `/v2/payment/sepay/ipn` |

Ba giá trị trong hai dòng VietQR lần lượt là `SEPAY_BANK_TRANSFER_FLOW`,
`SEPAY_GATEWAY_ENVIRONMENT`, `SEPAY_WEBHOOK_ENVIRONMENT`. Dù tên biến có chữ
`GATEWAY`, client hiện cũng dùng biến này để nhận biết môi trường thanh toán.
Code chỉ chấp nhận cặp `sandbox + test` hoặc `production + live` cho VietQR.
Giá trị flow mặc định trong code là `hosted_checkout`; luôn đặt rõ `vietqr`
khi vận hành QR trong Client.

| Biến | Dùng cho |
| --- | --- |
| `SEPAY_API_KEY` | Xác thực Bank Webhook qua `Authorization: Apikey <key>` |
| `SEPAY_MERCHANT_ID` | Định danh merchant Hosted Checkout |
| `SEPAY_MERCHANT_SECRET_KEY` | Ký dữ liệu gửi Hosted Checkout |
| `SEPAY_IPN_SECRET_KEY` | Xác thực Hosted IPN qua `X-Secret-Key` |

**Sửa merchant ID/secret không sửa được lỗi API key của Bank Webhook.**
Giá trị `SEPAY_API_KEY` không chứa tiền tố `Apikey `. Gateway băm credential
và chuyển digest qua metadata RPC được bảo vệ; Payment kiểm tra digest với
khóa cấu hình. Digest này không phải chữ ký HMAC của body.

Adapter Bank Webhook hiện hỗ trợ API Key. Không chọn HMAC/OAuth trên SePay
khi backend chưa có adapter tương ứng. Cách gửi header được mô tả trong
[tài liệu xác thực SePay](https://developer.sepay.vn/vi/sepay-webhooks/xac-thuc).

### VietQR với Test mode

Tạo tài khoản trong SePay Webhooks Test mode, rồi đặt trong `.env` của Payment:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=sandbox
SEPAY_BANK_TRANSFER_FLOW=vietqr
SEPAY_WEBHOOK_ENVIRONMENT=test
SEPAY_BANK_TRANSFER_ENABLED=true
SEPAY_CARD_ENABLED=false
SEPAY_QR_BANK_BIN=<BIN_6_CHU_SO>
SEPAY_QR_ACCOUNT_NUMBER=<TAI_KHOAN_TEST>
SEPAY_API_KEY=<KHOA_WEBHOOK_TEST>
```

Các placeholder phải được thay bằng giá trị thực của môi trường test.
Có thể chạy task VS Code **BDSPro · Payment · Configure Webhook Test** hoặc,
từ thư mục service:

```bash
python3 ops/configure-webhook-test.py --bank-bin <BIN> --account-number <TAI_KHOAN_TEST>
```

Script chỉ sửa `.env` local, giữ khóa đang có hoặc sinh khóa nếu thiếu; không
tạo tài khoản/webhook trên SePay. Dùng chức năng mô phỏng tiền vào của SePay,
đúng số tiền và mã của Order. Không dùng app ngân hàng chuyển tiền thật vào QR
Test mode. Xem [hướng dẫn Test mode](https://developer.sepay.vn/vi/sepay-webhooks/test-mode/bat-dau-nhanh).

### VietQR với Live

Liên kết tài khoản nhận tiền thật trên SePay Live và cấu hình:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=production
SEPAY_BANK_TRANSFER_FLOW=vietqr
SEPAY_WEBHOOK_ENVIRONMENT=live
SEPAY_BANK_TRANSFER_ENABLED=true
SEPAY_CARD_ENABLED=false
SEPAY_QR_BANK_BIN=<BIN_6_CHU_SO>
SEPAY_QR_ACCOUNT_NUMBER=<TAI_KHOAN_LIVE>
SEPAY_API_KEY=<KHOA_WEBHOOK_LIVE>
```

Live QR tạo chỉ dẫn chuyển tiền thật dù browser chạy ở localhost. BIN là mã
ngân hàng 6 chữ số, không phải số tài khoản; số tài khoản cấu hình dài 5–19 chữ
số. Kiểm tra tên người nhận, số tiền và nội dung trong app ngân hàng trước khi
xác nhận.

Test và Live cần database cùng secret riêng. Payload Bank Webhook không mang
cờ môi trường để backend tự phân biệt. Không đổi một database đang chứa đơn
Test sang Live rồi coi toàn bộ dữ liệu là giao dịch thật.

### Hosted Checkout khi có nhu cầu sử dụng

Đây là adapter đang có trong code; nó đưa bằng chứng về cùng settlement owner.
Để dùng nó, cấu hình đủ ba credential và chọn flow rõ ràng:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=sandbox
SEPAY_BANK_TRANSFER_FLOW=hosted_checkout
SEPAY_BANK_TRANSFER_ENABLED=true
SEPAY_CARD_ENABLED=false
SEPAY_MERCHANT_ID=<MERCHANT_ID>
SEPAY_MERCHANT_SECRET_KEY=<MERCHANT_SECRET>
SEPAY_IPN_SECRET_KEY=<IPN_SECRET>
SEPAY_CHECKOUT_BRIDGE_URL=http://127.0.0.1:8000/v2/payment/sepay/checkout
```

Payment ký các trường checkout, Gateway bridge tạo HTML form POST tới host
SePay được allowlist. `SEPAY_CHECKOUT_URL` có thể bỏ trống để runtime chọn
host theo môi trường. Production yêu cầu bridge HTTPS. Các biến
`SEPAY_SUCCESS_URL`, `SEPAY_ERROR_URL`, `SEPAY_CANCEL_URL` là URL quay lại Client,
có thể bỏ trống; chúng không kích hoạt gói.

Chỉ bật `SEPAY_CARD_ENABLED=true` khi merchant được cấp phương thức thẻ và
cấu hình hosted đầy đủ. IPN chấp nhận sự kiện `ORDER_PAID`, order `CAPTURED`,
transaction `APPROVED`/`PAYMENT`, số tiền và currency phù hợp. Merchant secret
không được trả về Client hay lưu trong payload công khai.

## 6. Chạy local và nối webhook

Các lệnh sau chạy từ root backend:

```bash
make setup
make doctor
make up
make status
make smoke
```

Luồng development mặc định chạy application native; Docker chạy hạ tầng.
Để sửa riêng Payment trong terminal có auto-reload:

```bash
make deps service=payment-service
make dev service=payment-service
```

Sau khi sửa ENV/source mà đang dùng binary native, build và restart Payment:

```bash
make native-restart service=payment-service
make logs service=payment-service
```

Khởi động Client ở terminal khác:

```bash
cd ../react-native-quy-hoach
QHPRO_WEB_DEV_PORT=3001 yarn web:local
```

Mở `http://localhost:3001/dashboard/pricing`, đăng nhập và tạo checkout.
Hướng dẫn đầy đủ ba source và tài khoản development ở
[runbook local](../documents/RUN-THREE-SOURCE-LOCALLY.md).

### Endpoint công khai

SePay cần gọi Gateway của **cùng backend/database mà browser đang sử dụng**:

```text
https://<public-gateway>/v2/payment/sepay/webhook
```

Nếu browser dùng backend local, webhook phải đi tới Gateway local qua tunnel.
Nếu webhook dùng `https://api-test.qhpro.vn/...`, phải cập nhật code và ENV của
backend phục vụ domain đó, trừ khi domain đã được cấu hình trỏ về local.
Tên domain không chứng minh nó đang dùng database nào.

Task **BDSPro · Payment · Webhook Test Tunnel** chạy
[`sepay-webhook-tunnel.py`](../gateway-service/ops/sepay-webhook-tunnel.py).
Script cần binary `cloudflared` ở `.tmp/development/tools/cloudflared`; xem
[runbook local](../documents/RUN-THREE-SOURCE-LOCALLY.md) để chuẩn bị tool.

Proxy chỉ chuyển `POST /v2/payment/sepay/webhook` từ `127.0.0.1:18080` tới
Gateway `127.0.0.1:8000`; không phục vụ hosted IPN. Khi tạo tunnel, terminal
in `COPY THIS URL INTO SEPAY`. Đọc lại URL từ root:

```bash
cat .tmp/development/sepay-webhook-url.txt
```

Giữ tunnel chạy và dùng URL của lần chạy hiện tại. Script xóa file URL khi
process dừng để không ai sao chép nhầm endpoint đã hết hiệu lực. Vận hành lâu
dài cần endpoint HTTPS ổn định thay vì phụ thuộc URL tạm.

### Cấu hình webhook trên SePay

1. Chọn đúng môi trường Test hoặc Live và đúng tài khoản nhận tiền.
2. Vào **Tích hợp → Webhooks → Thêm webhook/Sửa**.
3. Chọn sự kiện tiền vào, `application/json`, bật tự động gửi lại.
4. Điền endpoint công khai của Gateway tương ứng.
5. Chọn API Key, nhập giá trị trùng `SEPAY_API_KEY` của Payment đang chạy.
6. Chỉ bật bộ lọc mã khi đã kiểm tra SePay nhận diện đúng nội dung `QHP...`;
   tùy chọn bỏ qua giao dịch không có mã có thể chặn webhook khi `code=null`.

Các trường trên dashboard và bộ lọc được mô tả trong
[hướng dẫn tạo webhook SePay](https://developer.sepay.vn/vi/sepay-webhooks/tao-webhook).
Không đưa khóa thật vào README, commit hoặc log chia sẻ. Nếu khóa đã lộ, thay
khóa ở cả SePay và đúng deployment, restart Payment rồi kiểm tra lại delivery.

## 7. Đối soát và xử lý sự cố

Thứ tự kiểm tra: **giao dịch ngân hàng → nhật ký gửi SePay → Gateway →
Settlement → Fulfillment → Subscription → Client**. Ghi lại transaction ID,
reference, số tiền, thời gian giao dịch và môi trường; không ghi credential.

### Điều kiện nhận tiền cho đơn

- Bank Webhook vượt qua kiểm tra credential, là tiền vào và có số tiền VND
  nguyên dương; thời gian được parse theo UTC+7.
- Domain tìm được mã đơn; settlement lookup Order bằng reference.
- Số tiền/currency khớp chính xác điều khoản đã chốt.
- Thời điểm giao dịch không sau `Order.expires_at`.
- Provider transaction ID được dedupe; retry cùng giao dịch không tạo thêm
  settlement/fulfillment. Một giao dịch mới cho đơn đã phân loại được đưa vào review.

Thời hạn mặc định 30 phút được chốt khi tạo Order. Webhook tới muộn vẫn có thể
được chấp nhận nếu **thời gian giao dịch gốc** nằm trong hạn. Thay ENV TTL không
kéo dài các Order đã tạo.

| Triệu chứng | Kiểm tra / hướng xử lý |
| --- | --- |
| `401 {"message":"Invalid or expired token"}` | Gateway đang phân loại `Authorization: Apikey ...` thành JWT. Deploy/restart Gateway có `POST /v2/payment/sepay/webhook` trong `ProviderCredential`, dùng đúng path không có dấu `/` cuối và kiểm tra reverse proxy giữ nguyên method/path |
| Webhook `401` với lỗi provider credential | So sánh API key SePay với `SEPAY_API_KEY` của đúng Payment đang chạy; restart Payment sau khi sửa ENV và kiểm tra header qua proxy |
| Timeout hoặc `502` | Endpoint/tunnel, Gateway `8000`, kết nối Gateway → Payment và log process |
| Không có delivery | Tài khoản, môi trường, webhook bật/tắt và bộ lọc trên SePay |
| Mã không khớp | Đối chiếu `content` với Order; domain hỗ trợ bỏ dấu `-` và đổi hoa/thường, không suy diễn mã từ số tiền |
| `unmatched` | Reference có cú pháp nhận diện được nhưng không có Order trong database nhận webhook |
| `requires_review` | Đọc `review_reason`; thường sai số tiền/currency hoặc có giao dịch bổ sung |
| `late` | Giao dịch xảy ra sau hạn; không tự cấp quyền |
| `200` nhưng chưa active | Đọc Settlement và Fulfillment; ACK không chứng minh nhận đủ tiền |
| Fulfillment `retry` | Kiểm tra User RPC và lỗi tạm thời; worker tự thử lại |
| Fulfillment `requires_review` | Kiểm tra điều khoản/effect và lỗi User; chỉ redrive sau khi xử lý nguyên nhân |
| Backend completed nhưng UI chờ | Kiểm tra API tiến trình, user/order hiện tại, môi trường Client và request polling |
| Attempt hết hạn sau migration | Xem mục migration; replay khóa cũ không tái tạo QR |

Nếu đã chuyển tiền, **không chuyển thêm chỉ để thử lại**. Sửa nguyên nhân rồi
gửi lại webhook của chính giao dịch đó từ SePay, giữ transaction ID, số tiền,
nội dung và thời gian gốc. Không dựng webhook giả để xác nhận tiền Live.

Backend không cộng dồn các khoản thiếu thành một thanh toán đủ. Chuyển `5.000`
cho đơn `99.000` sẽ cần đối soát; không tự kích hoạt gói.

### Công cụ operator

Admin UI cung cấp:

| Trang | Dữ liệu |
| --- | --- |
| `/admin/commercial/payments` | Order, attempt, settlement và evidence |
| `/admin/commercial/operations` | Fulfillment, lỗi và redrive |
| `/admin/commercial/subscriptions` | Subscription do User sở hữu |
| `/admin/commercial/usage` | Usage/quota |
| `/admin/commercial/reports` | Report job và file |

**Xác nhận đã nhận tiền** là command có permission, idempotency key, actor và
lý do; nó đưa evidence vào cùng settlement usecase. Chỉ sử dụng sau khi operator
đối chiếu chứng từ. Command hiện yêu cầu đơn đang chờ, đã có attempt và chưa hết
hạn; không phải công cụ mở khóa mọi đơn late/review. Không ghi SQL `paid`, sửa
giá đơn hoặc tự cộng quota để giải quyết lỗi webhook.

**Thử kích hoạt lại** xử lý fulfillment, không xác nhận tiền lần hai. Khi một
giao dịch đã được lưu thành settlement review/late, gửi lại cùng ID không tự
phân loại lại kết quả; cần quy trình đối soát theo trạng thái thực tế.

## 8. Database và migration

[`migrate/`](migrate/) là authority schema. Runtime mặc định dùng
`QHPRO_PAYMENT_DB_SCHEMA_MODE=sql`, không tự sửa schema khi phục vụ request.
Chế độ ORM `automigrate` còn có cho bootstrap development; không thay thế SQL
migration, đặc biệt với data migration và constraint.

Các bảng commercial chính:

- `commerce_orders`: điều khoản, reference, trạng thái và hạn của đơn.
- `commerce_payment_attempts`: command, provider reference và next action.
- `commerce_settlements`: evidence đã đối soát, unique provider/transaction ID.
- `commerce_fulfillments`: việc cấp quyền, lease và kết quả.
- `commerce_command_effects`: audit/idempotency của command vận hành.
- Completion Outbox: schema tại [`000003`](migrate/000003_payment_completed_outbox.up.sql).

Từ root, chạy migration owner:

```bash
make migrate service=payment-service
make -C payment-service migrate-version
```

Trước khi chạy ngoài local, kiểm tra `MIGRATION_URL` trỏ đúng database Payment.
Makefile migration hiện ưu tiên `MIGRATION_URL`, nếu thiếu dùng `DB_URL` mặc định
local; không tự suy ra `PAYMENT_DATABASE_URL` của runtime.

### Chuyển sang payload EMV — migration 000008

[`000008`](migrate/000008_canonical_vietqr_payload.up.sql) làm hai việc:

1. Với attempt `display_qr` có payload không bắt đầu bằng `000201`: đặt status
   `expired`, next action `none`, xóa payload và ghi `legacy_qr_payload_retired`.
2. Thêm CHECK: `display_qr` phải có payload bắt đầu `000201` và redirect rỗng;
   các next action khác phải có `qr_payload` rỗng.

CHECK chỉ chặn sai representation ở mức tiền tố/field; không phải bộ kiểm tra
EMV/CRC đầy đủ. Payload chuẩn được provider tạo từ dữ liệu đã kiểm tra.
Migration áp dụng cho mọi attempt khớp điều kiện, kể cả attempt lịch sử; nó
không đổi Order, Settlement hoặc quyền lợi User. Down migration chỉ gỡ CHECK,
không phục hồi ảnh đã xóa hay status cũ. Sao lưu trước khi triển khai nếu cần
giữ presentation data lịch sử.

Khi rollout, phối hợp backend, migration và Client theo cùng release; dừng tạo
attempt trong thời gian chuyển đổi để không trộn binary sinh PNG với constraint
mới hoặc Client chỉ render EMV với response cũ.

Lưu ý phục hồi: cùng command key trả lại attempt đã lưu, kể cả đã hết hạn.
Schema còn unique `(provider, provider_reference)`, nên chỉ đổi command key
trên cùng Order không bảo đảm tạo được attempt thứ hai. Không xóa lịch sử để
vượt constraint. Với đơn chưa thanh toán, quay lại quy trình checkout khi được
User cho phép tạo Order mới; với đơn đã chuyển tiền, tiếp tục đối soát trên mã
Order gốc. Migration này không triển khai command tái phát hành QR cho cùng đơn.

## 9. Kiểm thử và bằng chứng hoàn tất

Từ root backend:

```bash
make test service=payment-service
make smoke
```

Khi cần kiểm tra đúng các boundary vừa sửa, từ thư mục service:

```bash
go test ./internal/domain/payment ./infra/reference ./infra/provider/sepay ./infra/handler/grpc/wallet
```

Nếu sửa contract/config/schema dùng các gate trong
[`CONTRIBUTING.md`](../CONTRIBUTING.md), gồm `make verify`; dùng `make accept`
cho candidate toàn hệ thống. Các suite tự tạo giao dịch chỉ được chạy trong
môi trường test riêng. `make smoke` chứng minh boundary hoạt động, không chứng
minh ngân hàng đã chuyển tiền hoặc fulfillment đã hoàn tất.

Client checks, từ repository `react-native-quy-hoach`:

```bash
yarn test src/Screens/DashboardPricingScreen/__tests__/PaymentProgressModal.test.tsx --runInBand
yarn verify:commercial-flow
yarn web:build:local
```

Một lần kiểm chứng end-to-end cần có đầy đủ:

1. Đơn mới có reference và QR chứa đúng tài khoản/số tiền/nội dung.
2. Giao dịch tương ứng xuất hiện trong SePay, delivery đúng môi trường được ACK.
3. Settlement `funds_confirmed`, fulfillment `completed` trên cùng Order.
4. Subscription/entitlement đúng user và version gói; Client đọc được `ACTIVE`.
5. Thực hiện phân tích vùng, mở hồ sơ thửa và tạo PDF; kiểm tra usage tương ứng
   tại TQD, theo [runbook quota](../documents/RUN-THREE-SOURCE-LOCALLY.md).
6. Retry cùng webhook không cấp lại quyền. Trạng thái gửi thông báo được kiểm
   tra riêng ở Notification.

## 10. Bản đồ source

| Câu hỏi cần trả lời | Source owner |
| --- | --- |
| Runtime cấu thành thế nào? | [`cmd/grpc/runtime.go`](cmd/grpc/runtime.go) |
| Biến ENV và validation ở đâu? | [`config/config.go`](config/config.go) |
| Mã chuyển khoản chuẩn là gì? | [`internal/domain/payment/reference.go`](internal/domain/payment/reference.go) |
| Order được chốt ra sao? | [`internal/usecase/order/service.go`](internal/usecase/order/service.go) |
| Replay attempt/phương thức thanh toán? | [`internal/usecase/attempt/service.go`](internal/usecase/attempt/service.go) |
| Chuỗi VietQR và hosted signature? | [`infra/provider/sepay/provider.go`](infra/provider/sepay/provider.go) |
| Ingress webhook/IPN và credential? | [`infra/handler/grpc/wallet/handler.go`](infra/handler/grpc/wallet/handler.go) |
| Ai quyết định nhận tiền? | [`internal/usecase/settlement/service.go`](internal/usecase/settlement/service.go) |
| Ai gọi User kích hoạt gói? | [`internal/usecase/fulfillment/service.go`](internal/usecase/fulfillment/service.go) |
| Client progress được suy ra thế nào? | [`internal/usecase/commercialprofile/service.go`](internal/usecase/commercialprofile/service.go) |
| Operator command và redrive? | [`internal/usecase/admin/service.go`](internal/usecase/admin/service.go) |
| Transaction/constraint lưu trạng thái? | [`infra/postgres/payment/`](infra/postgres/payment/), [`migrate/`](migrate/) |
| Gateway tạo ACK SePay? | [`httpresponse/sepay.go`](../gateway-service/internal/httpresponse/sepay.go) |

Xem thêm [schema authority](migrate/README.md),
[backend architecture](../BACKEND-ARCHITECTURE.md) và
[commercial operations](../documents/QHPRO-COMMERCIAL-OPERATIONS.md).


### Bất biến webhook nạp ví và mua gói

Webhook ngân hàng kiểm tra credential, chiều tiền vào và số tiền nguyên hợp lệ
trước khi chọn nhánh QHP (order) hoặc DEP (ví). Nhánh DEP khóa dòng transaction
bằng `GetTransactionByCodeForUpdate` trong cùng transaction với cập nhật wallet.
Webhook lặp cùng deposit đã completed trả lại kết quả cũ; canceled và giao dịch
khác loại không được cộng số dư. Nhánh ví chưa thay thế sổ bằng chứng ngân hàng
của commerce settlement. Xem [review code và giới hạn](../documents/QUOTA-CODE-REVIEW.md).
