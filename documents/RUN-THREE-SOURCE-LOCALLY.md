# Chạy ba source QHPRO trên máy local

Runbook này là đường chạy tay chuẩn cho Backend, React Admin và React Native
Web. Không có `demo runtime`, không ghi thẳng trạng thái nghiệp vụ và không chạy
application service trong Docker. Docker chỉ cung cấp PostgreSQL/PostGIS,
Redis và RabbitMQ; Go, Admin và Client chạy trực tiếp để developer debug được
source đang sửa.

## 1. Cấu trúc workspace

Ba repository phải nằm cạnh nhau:

```text
/home/bop/projects/
├── bdspro-golang-microservices-r2f0-i21/
├── react-js-admin-bds/
└── react-native-quy-hoach/
```

Trên Windows, mở Docker Desktop, vào **Settings → Resources → WSL
Integration**, bật distro đang dùng rồi chọn **Apply & restart**. Kiểm tra từ
terminal WSL:

```bash
docker version
docker compose version
```

Hai lệnh phải hiện cả Docker client/server và không còn thông báo “command not
found” hoặc “WSL integration”.

## 2. Chuẩn bị lần đầu

### Terminal 1 — Backend

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make setup
make doctor
make config-check
```

`make doctor` phải báo Docker, Docker Compose, Go được pin, Buf, Wire, migrate,
Air, FFmpeg và WeasyPrint ở trạng thái `READY`.

### Terminal 2 — Admin

```bash
cd /home/bop/projects/react-js-admin-bds
npm ci
npm run accept:commercial-admin
```

### Terminal 3 — Client

```bash
cd /home/bop/projects/react-native-quy-hoach
yarn install --frozen-lockfile
yarn verify:environment-config
yarn verify:commercial-flow
```

Các bước cài dependency chỉ cần chạy lại khi lockfile/toolchain thay đổi.

## 3. Chọn giữ hay làm sạch dữ liệu

Mặc định `make up` giữ database của lần làm việc trước. Đây là lựa chọn phù hợp
cho daily development.

Nếu cần một lần nghiệm thu từ dữ liệu hoàn toàn sạch:

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make down
make reset
```

Nhập chính xác `reset-development-data` khi được hỏi. Lệnh này xóa volume local
của Postgres, Redis, RabbitMQ và file development; source code không bị xóa.
Không chạy `make reset` nếu còn dữ liệu local cần giữ.

## 4. Khởi động runtime

### Terminal 1 — Backend và hạ tầng

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make up
make status
make smoke
```

`make up` thực hiện theo thứ tự:

```text
PostgreSQL/PostGIS + Redis + RabbitMQ trong Docker
  → versioned migrations
  → development login identities
  → build Go binaries
  → chạy Go services trực tiếp trong WSL
```

Sau khi lệnh trả về, process Go vẫn được quản lý bằng PID/log trong
`.tmp/development/`.

Theo dõi log ở terminal thứ tư nếu cần:

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make logs service=gateway-service
make logs service=payment-service
make logs service=tqd-service
```

Nhấn `Ctrl+C` chỉ dừng việc theo dõi log, không dừng server. Xem một lần không
follow:

```bash
FOLLOW=0 TAIL=100 make logs service=tqd-service
```

### Terminal 2 — Admin

```bash
cd /home/bop/projects/react-js-admin-bds
npm run dev -- --host 127.0.0.1 --port 5173 --strictPort
```

Mở `http://127.0.0.1:5173` và đăng nhập:

```text
Tài khoản: admin
Mật khẩu: admin123
```

### Terminal 3 — React Native Web Client

```bash
cd /home/bop/projects/react-native-quy-hoach
QHPRO_WEB_DEV_PORT=3001 yarn web:local
```

Mở `http://localhost:3001/dashboard/pricing`. Dùng hostname `localhost` này
xuyên suốt checkout để browser return URL và session không đổi origin.

Chọn **Đăng nhập → Đăng nhập bằng tài khoản và mật khẩu**:

```text
Số điện thoại hoặc tài khoản: 0900000000 hoặc client
Mật khẩu: client123
```

Các boundary cần nhìn thấy:

| Thành phần | Địa chỉ |
| --- | --- |
| Gateway | `http://127.0.0.1:8000` |
| File | `http://127.0.0.1:8002` |
| Admin | `http://127.0.0.1:5173` |
| Client | `http://localhost:3001` |
| RabbitMQ UI | `http://127.0.0.1:15672` |

## 5. Tạo ba gói từ Admin

Mở `http://127.0.0.1:5173/admin/commercial/plans`:

1. Chọn **Thiết kế gói mới**.
2. Chọn một trong ba mẫu: Khởi đầu, Tiêu chuẩn hoặc Chuyên nghiệp.
3. Kiểm tra tên, giá, thời hạn và đối tượng **Khách hàng cá nhân**.
4. Kiểm tra riêng ba quyền lợi: phân tích khu vực, hồ sơ thửa chi tiết và báo
   cáo PDF. Mỗi quyền lợi có hạn mức hoặc `Không giới hạn` riêng.
5. Chọn **Lưu bản nháp**.
6. Chọn **Kiểm tra điều khoản**.
7. Chọn **Xuất bản cho Client**.
8. Lặp lại cho các mẫu còn lại.

Version đã xuất bản là bất biến. Nếu mã gói/version đã tồn tại, tạo version mới
như `1.0.1`; không sửa database. Client chỉ hiển thị version đang `active`, có
giá recurring và còn hiệu lực.

Quay lại tab Client. Chuyển focus sang tab hoặc refresh
`/dashboard/pricing`; catalog sẽ được đọc lại từ User Service. Gói miễn phí là
trạng thái không có subscription; ba gói trả phí là catalog do Admin vừa xuất
bản.

## 6. Mua gói và kích hoạt hoàn toàn bằng UI local

1. Trên Client chọn một gói trả phí và bấm **Đăng ký**.
2. Nhập họ tên, email và chọn **Chuyển khoản ngân hàng**.
3. Xác nhận. Client tạo Order rồi Payment Attempt và mở màn tiến trình bốn
   bước. Ghi lại mã tham chiếu của đơn.
4. Không cần chuyển tiền thật cho lần proof local đầu tiên. Mở Admin tại
   `/admin/commercial/payments`, tìm theo mã tham chiếu và mở chi tiết.
5. Chỉ trong vai trò operator đang đối soát, bấm **Xác nhận đã nhận tiền**,
   kiểm tra gói/số tiền, nhập lý do và chọn **Xác nhận và kích hoạt gói**.

Đây không phải thao tác ghi thẳng `paid`. Nó gọi command có permission,
`Idempotency-Key` và audit evidence của Payment Service; Payment tạo Settlement
và Fulfillment, rồi User Service kích hoạt Subscription/Entitlement. Client
poll order khoảng ba giây một lần và phải chuyển qua:

```text
Đang chờ thanh toán
  → Đã nhận tiền, đang kích hoạt
  → Gói đã được kích hoạt
```

Nếu fulfillment cần kiểm tra, mở `/admin/commercial/operations` và chỉ dùng
**Thử kích hoạt lại** sau khi dependency đã hoạt động. Không xác nhận tiền lần
hai và không sửa bảng database.

### Kiểm tra QR tại bước 2 với SePay Webhooks Test mode

1. Tạo tài khoản ngân hàng trong **SePay → Test mode → Tài khoản ngân hàng**.
2. Chạy task VS Code **BDSPro · Payment · Configure Webhook Test**, nhập mã
   BIN và số tài khoản vừa tạo. Task đặt `SEPAY_BANK_TRANSFER_FLOW=vietqr`,
   `SEPAY_GATEWAY_ENVIRONMENT=sandbox`, `SEPAY_WEBHOOK_ENVIRONMENT=test`.
3. Khởi động lại Payment và Gateway. Trong SePay Test mode, tạo webhook JSON
   cho tiền vào tới `https://<tunnel-gateway>/v2/payment/sepay/webhook`, xác thực
   API Key bằng giá trị `SEPAY_API_KEY` trong `payment-service/.env`.
4. Tạo phiên đăng ký mới. Bước 2 hiển thị QR và mã chuyển khoản ngay trong app.
   Phiên cũ vẫn giữ hành động thanh toán ban đầu để bảo đảm idempotency.
5. Trên SePay Test mode, mô phỏng tiền vào đúng số tiền và mã đơn. Không chuyển
   tiền thật bằng QR thử nghiệm và không bấm xác nhận thủ công trong Admin.
6. Webhook trả `200` với `{"success":true}`; app tự cập nhật đến khi gói được
   kích hoạt. Tải lại trang khi đang chờ sẽ khôi phục hành động của cùng phiên.

Tài khoản Test mode và tunnel phải được cấu hình trên SePay; task local không
tự tạo chúng. Xem [hướng dẫn Payment](../payment-service/README.md#inline-qr-with-bank-webhooks-test-mode).

Nếu chưa có URL công khai, chạy task **BDSPro · Payment · Webhook Test Tunnel**.
Task dùng `cloudflared` tại `.tmp/development/tools/cloudflared` và mở proxy
local ở `127.0.0.1:18080`. Chỉ `POST /v2/payment/sepay/webhook` được chuyển
đến Gateway `127.0.0.1:8000`. Lấy URL `https://….trycloudflare.com` trong
terminal rồi thêm `/v2/payment/sepay/webhook` để cấu hình trên SePay.
Giữ task chạy trong lúc thử; mỗi lần tạo tunnel mới có thể nhận URL khác.
Không chạy task thứ hai nếu tunnel nền đang giữ cổng `18080`.

Cloudflared được tải từ [bản phát hành chính thức](https://github.com/cloudflare/cloudflared/releases)
và kiểm tra SHA-256 trước khi chạy. [Quick Tunnel](https://developers.cloudflare.com/tunnel/setup/)
dành cho development; khi Gateway có HTTPS công khai thì dùng trực tiếp URL đó.

### Kiểm tra cổng thanh toán SePay sandbox (luồng mở trang ngoài)

Khi cần chứng minh provider integration, giữ
`SEPAY_GATEWAY_ENVIRONMENT=sandbox`, `SEPAY_BANK_TRANSFER_FLOW=hosted_checkout`
và credential trong `payment-service/.env`. Cấu hình IPN trên SePay thành:

```text
https://<HTTPS-tunnel-của-Gateway-local>/v2/payment/sepay/ipn
```

Tunnel phải chuyển đến `http://127.0.0.1:8000`. Khi dùng IPN sandbox, không bấm
xác nhận thủ công trong Admin. Browser redirect chỉ là UX; chỉ IPN đã xác thực
hoặc command đối soát thủ công có audit mới tạo settlement. Thanh toán thẻ giữ
tắt cho đến khi merchant được SePay cấp quyền CARD.

## 7. Kiểm chứng ba quota

Sau khi Client hiện gói `ACTIVE`, mở
`http://localhost:3001/dashboard/pricing`, bấm gói hiện tại và kiểm tra ba hàng
hạn mức độc lập.

Mở `http://localhost:3001/ban-do`:

1. Chọn công cụ **Phân tích**, khoanh một khu vực hợp lệ, xem preview rồi bấm
   **Phân tích vùng này**. Chỉ response thành công tăng một lượt phân tích.
2. Chọn một thửa đất và mở hồ sơ chi tiết. Chỉ hồ sơ chi tiết tải thành công
   tăng một lượt xem thửa.
3. Trong hồ sơ chọn **Xuất PDF**. Report được chấp nhận, worker tạo PDF và File
   Service trả URL tải về.
4. Quay lại `/dashboard/pricing`; ba số liệu phải phản ánh riêng từng tác vụ.
5. Retry ngay cùng khu vực/thửa hoặc HTTP retry không được trừ thêm lượt.

PMTiles chỉ là lớp hình ảnh. API phân tích và hồ sơ cần dữ liệu PostGIS tương
ứng. Nếu bản đồ hiện nhưng không chọn được thửa, hãy import development dataset
qua TQD owner workflow và dùng PMTiles sinh từ cùng version dataset; không trỏ
riêng API nghiệp vụ sang production để che thiếu dữ liệu local.

## 8. Đối chiếu trên Admin

Mở lần lượt:

```text
/admin/commercial/payments       order, attempt, settlement, evidence
/admin/commercial/operations     fulfillment và redrive
/admin/commercial/subscriptions  subscription đang hiệu lực
/admin/commercial/usage          ba durable meter và Redis projection
/admin/commercial/reports        report job và PDF/File evidence
/admin/user                      user, gói và quota tổng hợp
```

Kết quả đúng phải có cùng profile/order/plan xuyên suốt. Payment thành công
không được dùng làm bằng chứng rằng notification đã gửi; Notification có inbox
và trạng thái effect riêng.

Cuối cùng, trên Client mở gói hiện tại, chọn **Hủy gói trả phí**, nhập lý do và
xác nhận. Client phải trở lại gói Miễn phí; subscription/payment/report cũ vẫn
được giữ làm lịch sử.

## 9. Xử lý lỗi thường gặp

### `docker: command not found`

Bật Docker Desktop WSL Integration cho đúng distro, restart Docker Desktop rồi
mở terminal WSL mới. Không cài một Docker daemon thứ hai chỉ để né lỗi này.

### Gateway `8000` reset/refused

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make status
FOLLOW=0 make logs service=gateway-service
```

### Upstream `8219` refused

```bash
FOLLOW=0 make logs service=tqd-service
make native-restart service=tqd-service
make status
```

### Port Admin `5173` đã được dùng

```bash
lsof -nP -iTCP:5173 -sTCP:LISTEN
kill <PID-vừa-tìm-thấy>
```

Sau đó chạy lại Admin với `--strictPort`. Không chuyển sang 5174 vì proxy và
đường kiểm thử chuẩn đang dùng 5173.

### Client báo `No refresh token` hoặc `Authorization header missing`

Đăng xuất, mở DevTools → Application → Storage của `localhost:3001`, xóa session
cũ rồi đăng nhập lại. Không chép JWT từ Admin sang Client.

### Gói Admin tạo chưa xuất hiện trên Client

Kiểm tra gói đã **Xuất bản**, ở tab **Đang bán**, có giá recurring và
`subjectScope=profile`. Sau đó focus/refresh Client. Bản nháp không bao giờ xuất
hiện trên catalog công khai.

## 10. Dừng và vòng làm việc hằng ngày

Nhấn `Ctrl+C` tại terminal Admin và Client. Dừng Backend/hạ tầng nhưng giữ dữ
liệu:

```bash
cd /home/bop/projects/bdspro-golang-microservices-r2f0-i21
make down
```

Ngày làm việc tiếp theo chỉ cần:

```text
Terminal 1: make up
Terminal 2: npm run dev -- --host 127.0.0.1 --port 5173 --strictPort
Terminal 3: QHPRO_WEB_DEV_PORT=3001 yarn web:local
```

Khi sửa một Go service, dùng `make dev service=tqd-service` trong terminal riêng
để Air rebuild/restart service đó; các service còn lại và hạ tầng vẫn chạy.
