# ZNS Token Management - Example Usage

## Response từ Zalo API

Khi gọi refresh token API, bạn sẽ nhận được response như sau:

```json
{
    "access_token": "EXAMPLE_ACCESS_TOKEN",
    "refresh_token": "_JwjE4t_VdRyFg8pHl1zT_OwhprCx2nrcLoEE5IL",
    "expires_in": "90000"
}
```

## Cách 1: Tự động refresh (Scheduler - Recommended ✅)

Scheduler sẽ tự động gọi refresh token vào **00:00 mỗi ngày** và lưu vào Redis:

```go
// Không cần làm gì cả!
// Scheduler tự động chạy khi service start
// Tokens được tự động cập nhật và lưu vào Redis
```

**Logs khi scheduler chạy:**
```
⏰ [ZNS Scheduler] Next refresh token job will run at: 2025-10-22 00:00:00
🔄 [ZNS Scheduler] Starting refresh token job...
✅ [ZNS] Refresh token thành công. Access token expires in: 90000 seconds
💾 [ZNS] Saved new tokens to Redis successfully
```

---

## Cách 2: Update token manually

Nếu bạn có response từ API và muốn update ngay lập tức:

### Example 1: Trong Usecase/Handler

```go
package usecase

import (
    "context"
    "auth/internal/interface/provider"
)

type ZnsUsecase struct {
    znsProvider provider.IZnsProvider
}

func (u *ZnsUsecase) UpdateZnsTokens(ctx context.Context, accessToken, refreshToken string) error {
    // Update tokens với response từ Zalo API
    err := u.znsProvider.UpdateTokens(ctx, accessToken, refreshToken)
    if err != nil {
        return err
    }
    
    log.Println("✅ ZNS tokens updated successfully")
    return nil
}
```

### Example 2: Sử dụng response JSON

```go
package main

import (
    "context"
    "encoding/json"
)

type ZaloRefreshResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    string `json:"expires_in"`
}

func updateTokensFromResponse(znsProvider provider.IZnsProvider, responseJSON string) {
    ctx := context.Background()
    
    // Parse JSON response
    var resp ZaloRefreshResponse
    err := json.Unmarshal([]byte(responseJSON), &resp)
    if err != nil {
        log.Printf("Error parsing response: %v", err)
        return
    }
    
    // Update tokens
    err = znsProvider.UpdateTokens(ctx, resp.AccessToken, resp.RefreshToken)
    if err != nil {
        log.Printf("Error updating tokens: %v", err)
        return
    }
    
    log.Printf("✅ Tokens updated successfully")
    log.Printf("📝 Access token expires in: %s seconds", resp.ExpiresIn)
}
```

### Example 3: Direct update với token strings

```go
func main() {
    ctx := context.Background()
    
    // Response tokens từ Zalo API
    accessToken := "EXAMPLE_ACCESS_TOKEN"
    refreshToken := "_JwjE4t_VdRyFg8pHl1zT_OwhprCx2nrcLoEE5ILL2gn..."
    
    // Update vào service
    err := znsProvider.UpdateTokens(ctx, accessToken, refreshToken)
    if err != nil {
        log.Fatalf("Failed to update tokens: %v", err)
    }
    
    // Verify tokens đã được update
    currentAccessToken := znsProvider.GetAccessToken()
    log.Printf("Current access token: %s...", currentAccessToken[:50])
}
```

---

## Cách 3: Manual refresh (Gọi API refresh)

Nếu muốn gọi refresh API manually:

```go
func manualRefresh(znsProvider provider.IZnsProvider) {
    ctx := context.Background()
    
    // Gọi Zalo refresh API
    newAccessToken, newRefreshToken, err := znsProvider.RefreshAccessToken(ctx)
    if err != nil {
        log.Fatalf("Failed to refresh token: %v", err)
    }
    
    log.Printf("✅ Refreshed successfully!")
    log.Printf("New Access Token: %s...", newAccessToken[:50])
    log.Printf("New Refresh Token: %s...", newRefreshToken[:50])
    
    // Tokens đã được tự động:
    // 1. Cập nhật vào memory (znsProvider.accessToken)
    // 2. Lưu vào Redis (zns:access_token, zns:refresh_token)
}
```

---

## Flow hoàn chỉnh với Redis Persistence

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Service Start                                            │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. ZnsProvider.loadTokensFromRedis()                        │
│    - Redis có token? → Load từ Redis                        │
│    - Redis không có? → Load từ config → Save to Redis       │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. Service Running                                          │
│    - Tokens trong memory: ✅                                │
│    - Tokens trong Redis: ✅                                 │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. Scheduler chạy vào 00:00 (hoặc manual refresh)           │
│    - Call Zalo API: POST /v4/oa/access_token                │
│    - Response: new access_token + new refresh_token          │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. ZnsProvider.RefreshAccessToken()                         │
│    - Update tokens in memory                                │
│    - Save to Redis (zns:access_token, zns:refresh_token)    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. Service tiếp tục sử dụng tokens mới                      │
│    - SendOTPZNS() dùng accessToken mới                      │
│    - SendZNSMessage() dùng accessToken mới                  │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│ 7. Service restart (bất kỳ lúc nào)                         │
│    - Load tokens từ Redis                                   │
│    - ✅ Vẫn dùng được tokens mới nhất!                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Check token hiện tại

```go
// Lấy token hiện tại
accessToken := znsProvider.GetAccessToken()
refreshToken := znsProvider.GetRefreshToken()

log.Printf("Current Access Token: %s...", accessToken[:50])
log.Printf("Current Refresh Token: %s...", refreshToken[:50])
```

---

## Redis Keys

Tokens được lưu trong Redis với keys:

```
zns:access_token   → Access Token hiện tại
zns:refresh_token  → Refresh Token hiện tại
```

Check trong Redis CLI:
```bash
redis-cli
> GET zns:access_token
> GET zns:refresh_token
```

---

## Best Practices

1. ✅ **Sử dụng Scheduler (Auto refresh)** - Recommended
   - Tự động refresh vào 00:00 mỗi ngày
   - Không cần quan tâm đến token expiry
   - Tokens luôn được persist trong Redis

2. ✅ **Manual Update khi cần**
   - Dùng `UpdateTokens()` khi có response mới từ API
   - Dùng khi test hoặc khởi tạo tokens lần đầu

3. ✅ **Luôn check tokens từ Redis**
   - Service restart vẫn dùng tokens mới nhất
   - Không mất tokens sau khi deploy

4. ⚠️ **Không hardcode tokens trong code**
   - Dùng config file hoặc environment variables
   - Redis sẽ override config tokens

