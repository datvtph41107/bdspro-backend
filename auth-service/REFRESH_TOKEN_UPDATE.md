# Cập Nhật Refresh Token & Auth Config - Summary

## Tổng Quan Thay Đổi

Hệ thống đã được cập nhật để:
1. **RefreshToken API** trả về cả `access_token` và `refresh_token` mới
2. **Lưu trữ động**: ZNS tokens được lấy từ DB thay vì file config cứng
3. **Tự động cập nhật**: Sau khi gọi API ZNS refresh token, tokens mới tự động lưu vào DB

---

## 📋 Các File Đã Tạo/Sửa

### 1. **Entity & Repository**

#### `auth-service/internal/domain/auth_config.go`
```go
type AuthConfig struct {
    BaseEntity
    ConfigKey   string `gorm:"column:config_key;type:varchar(10);uniqueIndex;not null"`
    ConfigValue string `gorm:"column:config_value;type:text"`
    IsActive    bool   `gorm:"column:is_active;default:true"`
}

// Constants
const (
    ConfigKeyZNSToken       = "ZNS_TOKEN"
    ConfigKeyZNSRefresh     = "ZNS_REFRESH"
    ConfigKeyZNSTemplateID  = "ZNS_TEMPLATE_ID"
    ConfigKeyZNSOTPTemplate = "ZNS_OTP_TEMPLATE_ID"
)
```

#### `auth-service/internal/interface/repo/auth_config_repo.go`
Interface định nghĩa các method:
- `GetByKey(ctx, key)` - Lấy config theo key
- `UpdateValue(ctx, key, value)` - Cập nhật giá trị config
- `Create/Update/Delete` - CRUD cơ bản

#### `auth-service/infra/postgre/auth_config_postgres.go`
Implementation PostgreSQL cho `IAuthConfigRepo`

---

### 2. **Migration SQL**

#### `auth-service/migrate/create_auth_config_table.sql`
```sql
CREATE TABLE IF NOT EXISTS auth_config (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE,
    config_value TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ...
);

INSERT INTO auth_config (config_key, config_value, is_active) VALUES
('ZNS_TOKEN', 'RRoFMY3v...', true),
('ZNS_REFRESH', '3AIiBOb4...', true),
('ZNS_TEMPLATE_ID', '469663', true),
('ZNS_OTP_TEMPLATE_ID', '497058', true);
```

**Chạy migration**:
```bash
psql -U your_user -d auth_db -f auth-service/migrate/create_auth_config_table.sql
```

---

### 3. **RefreshToken API**

#### `auth-service/internal/dto/auth_extended.go`
```go
type RefreshResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"` // ✅ Thêm mới
}
```

#### `auth-service/internal/usecase/auth_user_usecase.go`
**Thay đổi logic RefreshToken**:
```go
func (s *AuthUsecase) RefreshToken(c context.Context, refreshToken string) (*dto.RefreshResponse, error) {
    // ... validate refresh token cũ

    // Tạo access token mới
    newAccessToken, _ := _jwt.GenerateToken(...)

    // ✅ Tạo refresh token MỚI
    newRefreshToken, _ := _jwt.GenerateToken(...)

    // ✅ Lưu timestamp vào cache để invalidate refresh token cũ
    currentTime := time.Now().Format(time.RFC3339)
    s.CacheProvider.SaveToken(c, props.AuthID, currentTime, ...)

    // ✅ Lưu refresh token mới vào DB để ZNS provider có thể dùng
    err = s.AuthConfigRepo.UpdateValue(c, domain.ConfigKeyZNSRefresh, newRefreshToken)
    if err != nil {
        log.Printf("⚠️ [RefreshToken] Không thể lưu refresh token vào DB: %v", err)
    } else {
        log.Printf("✅ [RefreshToken] Đã lưu refresh token mới vào DB")
    }

    return &dto.RefreshResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken, // ✅ Trả về refresh token mới
    }, nil
}
```

#### `auth-service/infra/handler/auth_handler.go`
```go
func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
    result, err := h.AuthUsecase.RefreshToken(ctx, req.RefreshToken)
    
    return &authpb.RefreshTokenResponse{
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken, // ✅ Trả về refresh token mới
    }, nil
}
```

---

### 4. **ZNS Provider**

#### `auth-service/infra/zns/zns_provider.go`

**Thay đổi chính**:

1. **Inject AuthConfigRepo**:
```go
type ZnsProvider struct {
    ...
    authConfigRepo IAuthConfigRepo // ✅ Thêm dependency
    ...
}

func NewZnsProvider(authConfigRepo IAuthConfigRepo) *ZnsProvider {
    // ✅ Inject repo
}
```

2. **Load tokens từ DB thay vì Redis**:
```go
func (z *ZnsProvider) loadTokensFromDB() {
    // Load ZNS_TOKEN từ DB
    accessTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSToken)
    z.accessToken = accessTokenConfig.ConfigValue

    // Load ZNS_REFRESH từ DB
    refreshTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSRefresh)
    z.refreshToken = refreshTokenConfig.ConfigValue
}
```

3. **Lấy refresh token từ DB khi refresh**:
```go
func (z *ZnsProvider) RefreshAccessToken(ctx context.Context) (string, string, error) {
    // ✅ Lấy refresh token mới nhất từ DB thay vì memory
    refreshTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSRefresh)
    if err != nil {
        return "", "", fmt.Errorf("lỗi khi lấy refresh token từ DB: %w", err)
    }
    
    currentRefreshToken := refreshTokenConfig.ConfigValue
    fmt.Printf("🔄 [ZNS] Using refresh token from DB: %s...\n", currentRefreshToken[:20])
    
    // ... gọi API Zalo refresh token
    
    // ✅ Lưu tokens mới vào DB
    if err := z.saveTokensToDB(ctx, znsResp.AccessToken, znsResp.RefreshToken); err != nil {
        return "", "", err
    }
    
    return znsResp.AccessToken, znsResp.RefreshToken, nil
}
```

4. **Reload tokens trước khi gửi ZNS**:
```go
func (z *ZnsProvider) SendZNSMessage(...) (map[string]interface{}, error) {
    // ✅ Reload tokens từ DB để đảm bảo dùng token mới nhất
    if err := z.reloadTokensFromDB(ctx); err != nil {
        fmt.Printf("⚠️ Warning: Failed to reload tokens from DB\n")
    }
    
    // ... gửi ZNS
}
```

---

## 🔄 Workflow Mới

### Workflow RefreshToken
```
Client gửi refresh_token cũ
    ↓
[Auth Service]
    ├─ Validate refresh_token cũ
    ├─ Tạo access_token mới
    ├─ Tạo refresh_token mới  ← ✅ MỚI
    ├─ Lưu timestamp vào cache (invalidate token cũ) ← ✅ MỚI
    ├─ Lưu refresh_token mới vào DB ← ✅ MỚI
    └─ Trả về {access_token, refresh_token} ← ✅ MỚI

Client nhận cả 2 tokens mới
Client lưu refresh_token mới để dùng cho lần refresh sau ← ✅ MỚI
ZNS Provider tự động sync refresh_token mới từ DB ← ✅ MỚI
```

### Workflow ZNS Send
```
User gửi OTP
    ↓
[Auth Service]
    ├─ Reload ALL ZNS config từ DB ← ✅ MỚI (tokens + template IDs)
    ├─ Gửi ZNS với access_token và template_id từ DB
    │
    ├─ Nếu thành công → Done
    │
    └─ Nếu thất bại (token hết hạn):
        ├─ Lấy refresh_token từ DB ← ✅ MỚI
        ├─ Gọi RefreshAccessToken()
        ├─ Lưu tokens mới vào DB ← ✅ MỚI
        └─ Retry gửi ZNS với token mới
```

---

## 🎯 Lợi Ích

### 1. **Bảo Mật Cao Hơn**
- Refresh token được rotate sau mỗi lần sử dụng
- Token cũ bị invalidate ngay lập tức
- Ngăn chặn tấn công replay với refresh token bị đánh cắp

### 2. **Quản Lý Tập Trung**
- ZNS config (tokens + template IDs) lưu trong DB, dễ quản lý
- Có thể update config qua admin panel (nếu có)
- Không cần restart service khi đổi config

### 3. **Tự Động Sync**
- ZNS provider tự động lấy config mới nhất từ DB
- Sau khi refresh ZNS token, tự động lưu lại vào DB
- Sau khi refresh user token, tự động lưu refresh_token vào DB ← ✅ MỚI
- Các instance service khác tự động sync config mới

---

## 📝 Cách Sử Dụng

### Client Integration

**Trước đây**:
```javascript
const response = await fetch('/v2/auth/token/refresh', {
    method: 'POST',
    body: JSON.stringify({ refreshToken })
});
const { accessToken } = await response.json();
// Chỉ lưu access_token, dùng refresh_token cũ mãi
```

**Bây giờ**:
```javascript
const response = await fetch('/v2/auth/token/refresh', {
    method: 'POST',
    body: JSON.stringify({ refreshToken })
});
const { accessToken, refreshToken: newRefreshToken } = await response.json();

// ✅ Lưu CẢ 2 tokens mới
localStorage.setItem('access_token', accessToken);
localStorage.setItem('refresh_token', newRefreshToken); // ← Quan trọng!
```

### Admin Update ZNS Config

Nếu cần update ZNS config thủ công:

```sql
-- Update access token
UPDATE auth_config 
SET config_value = 'new_access_token_here', updated_at = NOW()
WHERE config_key = 'ZNS_TOKEN';

-- Update refresh token
UPDATE auth_config 
SET config_value = 'new_refresh_token_here', updated_at = NOW()
WHERE config_key = 'ZNS_REFRESH';

-- Update template ID
UPDATE auth_config 
SET config_value = '469664', updated_at = NOW()
WHERE config_key = 'ZNS_TEMPLATE_ID';

-- Update OTP template ID
UPDATE auth_config 
SET config_value = '497059', updated_at = NOW()
WHERE config_key = 'ZNS_OTP_TEMPLATE_ID';
```

Hoặc qua API (nếu có):
```bash
curl -X PUT /admin/config/ZNS_TOKEN \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"value": "new_token"}'

curl -X PUT /admin/config/ZNS_TEMPLATE_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"value": "469664"}'
```

---

## ⚠️ Breaking Changes

### Client phải cập nhật code:
- **BẮT BUỘC**: Lưu `refreshToken` mới sau mỗi lần gọi refresh
- Nếu không, refresh token cũ sẽ bị invalidate và user phải login lại

### Wire Dependencies:
- `ZnsProvider` bây giờ cần inject `IAuthConfigRepo`
- Cần update `wire.go` hoặc `wire_gen.go`:
```go
wire.Build(
    postgre.NewAuthConfigRepository,  // ← Thêm
    zns.NewZnsProvider,               // Đã có repo dependency
    ...
)
```

---

## 🧪 Testing

### 1. Test Migration
```bash
# Run migration
psql -U postgres -d auth_db -f auth-service/migrate/create_auth_config_table.sql

# Verify
SELECT * FROM auth_config WHERE config_key LIKE 'ZNS%';
```

### 2. Test RefreshToken API
```bash
# Lấy refresh token từ login
curl -X POST /v2/auth/otp/verify -d '{"phone":"0967771234","otp":"123456",...}'
# Response: {"accessToken":"...", "refreshToken":"abc123"}

# Gọi refresh
curl -X POST /v2/auth/token/refresh -d '{"refreshToken":"abc123"}'
# Response: {"accessToken":"new_access", "refreshToken":"xyz789"}
# Logs: "✅ [RefreshToken] Đã lưu refresh token mới vào DB"

# Thử dùng refresh token CŨ lại (phải FAIL)
curl -X POST /v2/auth/token/refresh -d '{"refreshToken":"abc123"}'
# Response: 401 Unauthorized ✅

# Kiểm tra DB đã có refresh token mới
SELECT config_value FROM auth_config WHERE config_key = 'ZNS_REFRESH';
# Phải trả về: xyz789 (refresh token mới)
```

### 3. Test ZNS
```bash
# Gửi OTP → ZNS provider tự động load ALL config từ DB
curl -X POST /v2/auth/otp/request -d '{"phone":"0967771234"}'
# Kiểm tra log: 
# "✅ [ZNS] Loaded access token from DB"
# "✅ [ZNS] Loaded template ID from DB: 469663"
# "✅ [ZNS] Loaded OTP template ID from DB: 497058"
```

---

## 🔧 Wire Generation

Sau khi thay đổi, cần regenerate wire:

```bash
cd auth-service
wire ./wire
```

Đảm bảo `wire.go` có:
```go
func InitializeAuthService() *AuthService {
    wire.Build(
        // ... existing providers
        postgre.NewAuthConfigRepository,  // ← Thêm
        zns.NewZnsProvider,               // Inject authConfigRepo
        // ...
    )
    return &AuthService{}
}
```

---

## 📊 Database Schema

```sql
auth_config
├─ id (BIGSERIAL PRIMARY KEY)
├─ config_key (VARCHAR(100) UNIQUE) → "ZNS_TOKEN", "ZNS_REFRESH", "ZNS_TEMPLATE_ID", "ZNS_OTP_TEMPLATE_ID"
├─ config_value (TEXT) → Token value hoặc Template ID
├─ is_active (BOOLEAN)
├─ created_at (TIMESTAMP)
└─ updated_at (TIMESTAMP)
```

**Sample Data**:
```sql
SELECT config_key, config_value FROM auth_config WHERE config_key LIKE 'ZNS%';

ZNS_TOKEN           | RRoFMY3vhLDil81HA_QkKcNwDHpNk4rIm-H0O_6E04MW-7DJkQzyJQRN0LgvZ5Dqj-K4Gtzy6ZecxgTO
ZNS_REFRESH         | 3AIiBOb4LqKDfACcWajKH0UWhYY8Lo4mIDUi1iXKE2r61IfJWO54XJiTJ7kCvcQQRKuBURVPmrkFKrrE
ZNS_TEMPLATE_ID     | 463
ZNS_OTP_TEMPLATE_ID | 498
```

---

## ✅ Checklist Deploy

- [ ] Run migration: `create_auth_config_table.sql`
- [ ] Verify data: Check `ZNS_TOKEN`, `ZNS_REFRESH`, `ZNS_TEMPLATE_ID`, `ZNS_OTP_TEMPLATE_ID` tồn tại trong DB
- [ ] Update wire dependencies (nếu dùng wire)
- [ ] Rebuild service: `go build`
- [ ] Cập nhật client code để lưu `refreshToken` mới
- [ ] Test RefreshToken API
- [ ] Test ZNS sending với config từ DB
- [ ] Monitor logs để ensure config được load từ DB

---

**Hoàn thành!** 🎉

