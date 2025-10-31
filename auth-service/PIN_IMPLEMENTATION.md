# PIN Authentication Implementation

## Tổng quan
API xác thực bằng mã PIN 6 số cho auth-service. User có thể tạo, xác minh, và cập nhật mã PIN để đăng nhập nhanh chóng.

## Files đã tạo/cập nhật

### 1. Proto Definitions
**File**: `shared/protobuf/schema/auth/auth_method.proto`

Đã thêm:
- `rpc CreatePIN` - Tạo mã PIN mới
- `rpc VerifyPIN` - Xác minh và đăng nhập bằng PIN
- `rpc UpdatePIN` - Cập nhật mã PIN
- `rpc CheckPINExists` - Kiểm tra PIN có tồn tại

Messages:
- `CreatePINRequest/Response`
- `VerifyPINRequest`
- `UpdatePINRequest/Response`
- `CheckPINExistsRequest/Response`

### 2. Domain Model
**File**: `auth-service/internal/domain/user_pin.go`

```go
type UserPINEntity struct {
    IdDomain
    PIN          string     // Mã PIN đã hash (bcrypt)
    PINCheckTime int        // Số lần nhập sai
    PINDate      *time.Time // Thời gian tạo/update
    IsActive     bool       // Trạng thái hoạt động
    LockedUntil  *time.Time // Thời gian khóa
}
```

### 3. DTOs
**File**: `auth-service/internal/dto/auth.go`

DTOs đã thêm:
- `CreatePINRequest` - Tạo PIN (pin, phone)
- `CreatePINResponse` - Response tạo PIN
- `VerifyPINRequest` - Xác minh PIN
- `UpdatePINRequest` - Update PIN (oldPin, newPin)
- `UpdatePINResponse`
- `CheckPINExistsRequest/Response`

### 4. Repository
**Files**:
- `auth-service/internal/interface/repo/pin_repo.go` - Interface
- `auth-service/infra/postgre/pin_postgres.go` - Implementation

Methods:
- `GetByAuthID` - Lấy PIN theo authID
- `Create` - Tạo mới PIN
- `Update` - Cập nhật PIN
- `Delete` - Xóa PIN
- `CheckPINExists` - Kiểm tra PIN có tồn tại
- `ResetCheckCounter` - Reset số lần nhập sai
- `IncrementCheckTime` - Tăng số lần nhập sai
- `LockPIN` - Khóa PIN
- `UnlockPIN` - Mở khóa PIN

### 5. Usecase
**File**: `auth-service/internal/usecase/pin_usecase.go`

Logic nghiệp vụ:
- **CreatePIN**: Validate → Hash PIN (bcrypt) → Lưu DB
- **VerifyPIN**: Kiểm tra lock → Verify PIN → Reset counter/Lock → Tạo session → Generate token
- **UpdatePIN**: Verify PIN cũ → Hash PIN mới → Update
- **CheckPINExists**: Kiểm tra PIN có tồn tại

Bảo mật:
- Số lần thử tối đa: 3 lần
- Thời gian khóa: 15 phút
- PIN hash bằng bcrypt

### 6. Handler
**File**: `auth-service/infra/handler/auth_handler.go`

Updated:
- Added `PINUsecase *usecase.PINUsecase` to struct
- Updated constructor `NewAuthHandler`

Handler methods:
- `CreatePIN` - POST `/v2/auth/pin/create`
- `VerifyPIN` - POST `/v2/auth/pin/verify`
- `UpdatePIN` - PUT `/v2/auth/pin/update` (requires auth)
- `CheckPINExists` - GET `/v2/auth/pin/check` (requires auth)

### 7. Database Migration
**File**: `auth-service/infra/db/migrate_pin.sql`

Bảng `user_pin`:
```sql
CREATE TABLE user_pin (
    auth_id BIGINT PRIMARY KEY,
    pin VARCHAR(255) NOT NULL,
    pin_check_time INT DEFAULT 0,
    pin_date TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

## API Endpoints

### 1. Tạo mã PIN
```
POST /v2/auth/pin/create
Content-Type: application/json

Request:
{
  "pin": "123456",
  "phone": "+84987654321"
}

Response:
{
  "authId": 123,
  "message": "Tạo mã PIN thành công"
}
```

### 2. Xác minh/Đăng nhập bằng PIN
```
POST /v2/auth/pin/verify
Content-Type: application/json

Request:
{
  "phone": "+84987654321",
  "pin": "123456",
  "version": "1.0.0",
  "platform": "MOBILE",
  "os": "iOS",
  "deviceName": "iPhone 13 Pro"
}

Response:
{
  "authId": 123,
  "profileId": 456,
  "accessToken": "eyJhbGc...",
  "fullName": "Nguyen Van A",
  "role": "ROLE_USER",
  "referralCode": "ABC123"
}
```

### 3. Cập nhật mã PIN
```
PUT /v2/auth/pin/update
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
  "oldPin": "123456",
  "newPin": "654321"
}

Response:
{
  "message": "Cập nhật mã PIN thành công"
}
```

### 4. Kiểm tra PIN tồn tại
```
GET /v2/auth/pin/check
Authorization: Bearer {token}

Response:
{
  "exists": true,
  "authId": 123
}
```

## Security Features

1. **PIN Hashing**: Sử dụng bcrypt để hash PIN
2. **Rate Limiting**: Giới hạn 3 lần nhập sai
3. **Account Locking**: Khóa 15 phút sau khi nhập sai 3 lần
4. **Auto Reset**: Reset counter khi đăng nhập thành công
5. **Session Management**: Tạo session tracking cho mỗi lần đăng nhập

## Error Codes

- `400` - Bad Request (PIN không hợp lệ)
- `401` - Unauthorized (PIN sai)
- `403` - Forbidden (PIN bị vô hiệu hóa)
- `404` - Not Found (Tài khoản/PIN không tồn tại)
- `423` - Locked (PIN bị khóa)
- `500` - Internal Server Error

## Cần hoàn thiện

### 1. Generate Proto Code
```bash
cd shared/protobuf/schema/auth
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative auth_method.proto
```

### 2. Wire Dependency Injection
Thêm vào `wire/wire.go`:
```go
pinRepo := postgre.NewPINRepository(db)
pinUsecase := usecase.NewPINUsecase(
    pinRepo,
    statusRepo,
    authMethodRepo,
    sessionRepo,
    factory,
)
```

Update `handler.NewAuthHandler` call để include `pinUsecase`.

### 3. Run Migration
```bash
psql -U postgres -d auth_db -f auth-service/infra/db/migrate_pin.sql
```

### 4. Register in Gateway
Thêm PIN routes vào `gateway-service` nếu cần expose qua HTTP gateway.

### 5. Testing
- Unit tests cho usecase
- Integration tests cho repository
- E2E tests cho handlers

## Usage Example

```go
// Client code example
import authpb "pb/types/auth"

// 1. Tạo PIN
createReq := &authpb.CreatePINRequest{
    Pin:   "123456",
    Phone: "+84987654321",
}
createResp, err := authClient.CreatePIN(ctx, createReq)

// 2. Đăng nhập bằng PIN
verifyReq := &authpb.VerifyPINRequest{
    Phone:      "+84987654321",
    Pin:        "123456",
    Version:    "1.0.0",
    Platform:   "MOBILE",
    Os:         "iOS",
    DeviceName: "iPhone 13 Pro",
}
loginResp, err := authClient.VerifyPIN(ctx, verifyReq)

// 3. Cập nhật PIN
updateReq := &authpb.UpdatePINRequest{
    OldPin: "123456",
    NewPin: "654321",
}
updateResp, err := authClient.UpdatePIN(ctx, updateReq)

// 4. Kiểm tra PIN
checkReq := &authpb.CheckPINExistsRequest{}
checkResp, err := authClient.CheckPINExists(ctx, checkReq)
```

## Notes

- PIN phải là 6 số
- Mỗi tài khoản chỉ có 1 PIN
- PIN được hash bằng bcrypt (không lưu plain text)
- Session được tạo tự động khi verify PIN thành công
- Token generation hiện tại là placeholder, cần implement proper JWT generation

