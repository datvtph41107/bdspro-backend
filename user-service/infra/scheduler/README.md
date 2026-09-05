# ZNS Scheduler

## Mô tả
Scheduler tự động làm mới (refresh) ZNS Access Token vào lúc **00:00** mỗi ngày và lưu vào Redis để persist.

## Cấu trúc

### `ZnsScheduler`
- **Start()**: Khởi động scheduler
- **Stop()**: Dừng scheduler
- **scheduleDaily()**: Schedule job chạy vào 00:00 mỗi ngày
- **refreshTokenJob()**: Job thực hiện refresh token

### `ZnsProvider` 
- **RefreshAccessToken()**: Gọi API refresh token và lưu vào Redis
- **UpdateTokens()**: Cập nhật token manually
- **GetAccessToken()**: Lấy access token hiện tại
- **GetRefreshToken()**: Lấy refresh token hiện tại

## Cơ chế hoạt động

### 1. **Khởi động Service**
Khi service khởi động, `ZnsProvider` sẽ:
- Load tokens từ Redis (nếu có)
- Nếu không có trong Redis, load từ config file và lưu vào Redis
- Log: `✅ [ZNS] Loaded access token from Redis`

### 2. **Auto Refresh vào 00:00**
Scheduler tự động:
- Tính toán thời gian còn lại đến 00:00 ngày tiếp theo
- Sử dụng `time.Timer` để chờ đến 00:00
- Gọi API refresh token của Zalo
- **Cập nhật tokens vào memory**
- **Lưu tokens mới vào Redis** (persist)
- Lặp lại cho ngày tiếp theo

### 3. **Persistence với Redis**
Tokens được lưu trong Redis với keys:
- `zns:access_token`: Access Token hiện tại
- `zns:refresh_token`: Refresh Token hiện tại

**Lợi ích**: Token luôn được persist, ngay cả khi service restart thì vẫn dùng token mới nhất từ Redis

## Logs

Scheduler log các sự kiện quan trọng:

```
🔧 [ZNS Scheduler] Starting ZNS Token Refresh Scheduler...
⏰ [ZNS Scheduler] Next refresh token job will run at: 2025-10-22 00:00:00 (in 8h 15m 30s)
🔄 [ZNS Scheduler] Starting refresh token job...
✅ [ZNS Scheduler] Successfully refreshed ZNS token at 2025-10-22 00:00:00
📝 [ZNS Scheduler] New Access Token: v_Dj0AnM8IRcys8BfraY9-ptRmMjHIrKl...
📝 [ZNS Scheduler] New Refresh Token: fiUw390e4sAdnxK7coSGOwt3k0cHErqMqTYD6g5N...
```

## Graceful Shutdown

Scheduler hỗ trợ graceful shutdown:

```go
// Trong cmd/grpc.go hoặc nơi cần stop scheduler
app.ZnsScheduler.Stop()
```

## Cấu hình

Token được load từ config file (`local.yml`, `develop.yml`):

```yaml
zns:
  access_token: "v_Dj0AnM8IRcys8BfraY9-ptRmMjHIrKl..."
  refresh_token: "fiUw390e4sAdnxK7coSGOwt3k0cHErqM..."
  template_id: 497058
  otp_template_id: 497058
```

## Dependency Injection

Scheduler được inject tự động qua Wire:

```go
// wire/wire.go
auth_infra_scheduler.NewZnsScheduler,

// initial/runtime.go
func NewApp(..., znsScheduler *scheduler.ZnsScheduler) *InitialApp {
    app.ZnsScheduler = znsScheduler
    app.StartZnsScheduler()
    return app
}
```

## Test thủ công

Để test scheduler ngay lập tức (không cần chờ đến 00:00):

1. Uncomment dòng trong `Start()`:
   ```go
   // Chạy ngay lần đầu tiên khi khởi động (optional)
   s.refreshTokenJob()
   ```

2. Restart service và xem logs

## Cách sử dụng

### 1. Lấy token hiện tại
```go
accessToken := znsProvider.GetAccessToken()
refreshToken := znsProvider.GetRefreshToken()
log.Printf("Current Access Token: %s", accessToken)
```

### 2. Update token manually
```go
ctx := context.Background()
err := znsProvider.UpdateTokens(ctx,
    "EXAMPLE_ACCESS_TOKEN", // new access token
    "_JwjE4t_VdRyFg8pHl1zT_OwhprCx2nr...", // new refresh token
)
if err != nil {
    log.Printf("Error updating tokens: %v", err)
} else {
    log.Println("✅ Tokens updated successfully")
}
```

### 3. Refresh token manually
```go
ctx := context.Background()
newAccessToken, newRefreshToken, err := znsProvider.RefreshAccessToken(ctx)
if err != nil {
    log.Printf("Error refreshing token: %v", err)
} else {
    log.Printf("✅ New Access Token: %s", newAccessToken)
    log.Printf("✅ New Refresh Token: %s", newRefreshToken)
}
```

## Flow hoàn chỉnh

```
Service Start
    ↓
Load tokens from Redis
    ↓ (if not exist)
Load tokens from config → Save to Redis
    ↓
Service Running (tokens in memory + Redis)
    ↓
00:00 Daily
    ↓
Call Zalo Refresh API
    ↓
Update tokens in memory
    ↓
Save tokens to Redis ← PERSIST!
    ↓
Continue using new tokens
    ↓
(Service restart)
    ↓
Load tokens from Redis ← Use latest tokens!
```

## TODO

- [x] Lưu token mới vào Redis để persist giữa các lần restart
- [ ] Thêm retry logic nếu refresh token thất bại
- [ ] Thêm alert/notification khi refresh token thất bại
- [ ] Monitor số lần refresh thành công/thất bại

