# BDSPro Three-Source Onboarding

Tài liệu này mô tả một developer mới đi từ ba repository vừa clone đến runtime
local và proof QHPRO quota. Development workflow chỉ tạo hai identity đăng
nhập cố định; không tạo ngầm plan, order, payment, subscription, quota, report
hoặc dữ liệu nghiệp vụ giả.

## 1. Workspace

```text
workspace/
├── bdspro-golang-microservices-r2f0-i21/
├── react-js-admin-bds/
└── react-native-quy-hoach/
```

Mở file `bdspro-golang-microservices-r2f0-i21/bdspro.code-workspace` để IDE nhận
đúng ba source và loại `node_modules`, generated protobuf, build artifact khỏi
search/watch.

## 2. Prerequisites

- Linux hoặc WSL2.
- Git, Make, Python 3, FFmpeg và Docker Engine/Compose.
- Docker Desktop phải bật WSL integration nếu dùng Windows.
- Node.js 22 đáp ứng Admin và cũng tương thích Client.
- npm đi cùng Node; Yarn `1.22.22` cho Client.

Backend tự bootstrap Go, Buf, Wire, migrate, Air và WeasyPrint đã pin; developer
không tự cài một phiên bản Go khác để chạy repository.

## 3. Setup một lần

Backend:

```bash
cd workspace/bdspro-golang-microservices-r2f0-i21
make setup
make doctor
make config-check
```

Admin:

```bash
cd workspace/react-js-admin-bds
npm ci
npm run verify:commercial-admin
```

Client:

```bash
cd workspace/react-native-quy-hoach
yarn install --frozen-lockfile
yarn verify:environment-config
yarn verify:commercial-flow
```

## 4. Runtime local

Terminal Backend:

```bash
cd workspace/bdspro-golang-microservices-r2f0-i21
make up
make status
make smoke
```

Khi `QHPRO_ENVIRONMENT=development`, `make up` khởi tạo idempotent hai tài khoản
trong database `user_service`:

| Giao diện | Tài khoản | Mật khẩu |
| --- | --- | --- |
| Admin | `admin` | `admin123` |
| Client | `0900000000` hoặc username `client` | `client123` |

Đây là credential đã hash trong đúng Auth/Profile/IAM tables và mọi lần đăng
nhập vẫn đi qua session/JWT/permission production path. Mỗi lần provision ở
development đưa hai mật khẩu về đúng giá trị công bố để fresh clone và database
đã dùng có cùng hành vi. Seed nằm ngoài thư mục migration và bị từ chối nếu
environment không phải `development`.

Với một installation không phải development, operator có quyền truy cập
database tạo đúng một quản trị viên gốc bằng input riêng:

```bash
export QHPRO_BOOTSTRAP_CONFIRM=CREATE_ROOT_OPERATOR
export QHPRO_BOOTSTRAP_ADMIN_USERNAME='root.operator'
export QHPRO_BOOTSTRAP_ADMIN_PASSWORD='replace-with-a-secret-of-at-least-12-characters'
export QHPRO_BOOTSTRAP_ADMIN_FULL_NAME='Quản trị hệ thống'
export QHPRO_BOOTSTRAP_ADMIN_EMAIL='root@example.com'
export QHPRO_BOOTSTRAP_ADMIN_PHONE=''
make bootstrap-admin
unset QHPRO_BOOTSTRAP_CONFIRM QHPRO_BOOTSTRAP_ADMIN_USERNAME \
  QHPRO_BOOTSTRAP_ADMIN_PASSWORD QHPRO_BOOTSTRAP_ADMIN_FULL_NAME \
  QHPRO_BOOTSTRAP_ADMIN_EMAIL QHPRO_BOOTSTRAP_ADMIN_PHONE
```

Command thuộc User Service, ghi profile + credential bcrypt + SYSTEM root role
trong một transaction và dùng advisory lock chống hai operator chạy đồng thời.
Nếu đã có root, command từ chối thay vì reset password hoặc nâng quyền ngầm.
Không ghi các biến bootstrap vào `.env` hay source control. Từ admin thứ hai trở
đi phải tạo qua Admin/IAM API.

Terminal Admin:

```bash
cd workspace/react-js-admin-bds
npm run dev -- --host 127.0.0.1 --port 5173 --strictPort
```

Terminal Client:

```bash
cd workspace/react-native-quy-hoach
QHPRO_WEB_DEV_PORT=3001 yarn web:local
```

Expected boundaries:

| Boundary | URL |
| --- | --- |
| Gateway | `http://127.0.0.1:8000` |
| File | `http://127.0.0.1:8002` |
| Admin | `http://127.0.0.1:5173` |
| Client pricing | `http://127.0.0.1:3001/dashboard/pricing` |
| RabbitMQ operator UI | `http://127.0.0.1:15672` |

## 5. Dữ liệu cần có owner

Source/migrations tạo schema, permission identity, reference data và catalog
baseline đã phê duyệt. Development workflow chỉ tạo hai login identity đã ghi
rõ ở trên; source không tự tạo payment, subscription, quota, report hoặc dữ
liệu quy hoạch giả.

Ba input phải được cung cấp theo đúng operational owner:

1. Production root operator: `make bootstrap-admin` đúng một lần với input do
   operator sở hữu; sau đó Admin tạo operator khác qua User/IAM API.
2. Payment provider: cấu hình credential và callback vào
   `POST /v2/payment/sepay/webhook`.
3. Planning dataset: import dữ liệu quy hoạch và cung cấp PMTiles; không commit
   production dataset vào application source.

Thiếu một trong ba input này không được che bằng code path `local` hoặc thao tác
ghi thẳng business table.

## 6. Thao tác flow thật

1. Operator đăng nhập Admin.
2. Tạo plan version draft, khai báo giá minor units, entitlement, meter và
   operation policy; validate rồi publish.
3. User đăng ký/đăng nhập Client, làm mới catalog và checkout plan đã publish.
4. Provider xác nhận giao dịch bằng webhook; Payment ghi evidence và settlement.
5. Fulfillment kích hoạt subscription trong User; Client đọc entitlement/quota.
6. User chọn dữ liệu quy hoạch đã import và yêu cầu report.
7. TQD ghi `Report + UsageEvent + Job` cùng transaction; worker tạo PDF và File
   trả artifact URL.
8. Admin đối chiếu Payment, Subscription, Usage, Report và Notification.

Các request `curl`, payload và màn Admin cần quan sát được ghi tại
[`QHPRO-COMMERCIAL-OPERATIONS.md`](QHPRO-COMMERCIAL-OPERATIONS.md). Đây là public
API thật mà Client đang dùng, không phải script cấp state trực tiếp.

## 7. Hai tầng bằng chứng

Kiến trúc/backend runtime:

```bash
cd workspace/bdspro-golang-microservices-r2f0-i21
make status
make smoke
make verify
```

Quota/value-chain tự động, chỉ trong development/test/CI:

```bash
make test-e2e
```

Admin và Client:

```bash
cd workspace/react-js-admin-bds
npm run accept:commercial-admin

cd workspace/react-native-quy-hoach
yarn verify:commercial-flow
yarn web:build:local
yarn web:verify:local
```

`test-e2e` dùng test-only identity trong `integration-test/.env`; production
guard từ chối fixture. Đây là proof của cùng public API/state machine, không
phải một cách chạy sản phẩm.

## 8. Tiêu chí hoàn thành

- Infrastructure healthy và mọi native service mở đúng port.
- Gateway liveness/readiness và public catalog PASS.
- Order có frozen terms và exact money.
- Provider evidence tạo đúng một settlement/subscription.
- Entitlement cho phép report đúng subject.
- Một accepted report tạo đúng một durable usage event/job.
- Retry cùng command không tiêu thụ quota lần hai.
- PDF mở qua File URL.
- Admin và Client đọc cùng durable state từ owner API.
