package zns

import (
	"auth/config"
	"auth/env"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// @bind: auth/internal/interface/provider.IZnsProvider
type ZnsProvider struct {
	accessToken   string
	refreshToken  string
	templateID    int
	otpTemplateID int
	httpClient    *http.Client
	redisClient   *redis.Client
}

// Redis keys cho ZNS tokens
const (
	ZNSAccessTokenKey  = "zns:access_token"
	ZNSRefreshTokenKey = "zns:refresh_token"
)

// NewZnsProvider tạo instance mới của ZnsProvider
func NewZnsProvider() *ZnsProvider {
	// Khởi tạo Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Properties.Redis.Host, config.Properties.Redis.Port),
		Password: config.Properties.Redis.Password,
		DB:       0,
	})

	provider := &ZnsProvider{
		templateID:    config.Properties.ZNS.TemplateID,
		otpTemplateID: config.Properties.ZNS.OTPTemplateID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		redisClient: redisClient,
	}

	// Load tokens từ Redis, nếu không có thì dùng config
	provider.loadTokensFromRedis()

	return provider
}

// loadTokensFromRedis load access token và refresh token từ Redis
// Nếu không có trong Redis thì load từ config và lưu vào Redis
func (z *ZnsProvider) loadTokensFromRedis() {
	ctx := context.Background()

	// Load access token từ Redis
	accessToken, err := z.redisClient.Get(ctx, ZNSAccessTokenKey).Result()
	if err == redis.Nil {
		// Không có trong Redis, dùng token từ config
		accessToken = config.Properties.ZNS.AccessToken
		// Lưu vào Redis
		z.redisClient.Set(ctx, ZNSAccessTokenKey, accessToken, 0)
		fmt.Println("📝 [ZNS] Loaded access token from config and saved to Redis")
	} else if err != nil {
		// Lỗi khi đọc Redis, dùng token từ config
		accessToken = config.Properties.ZNS.AccessToken
		fmt.Printf("⚠️ [ZNS] Error loading access token from Redis: %v, using config\n", err)
	} else {
		fmt.Println("✅ [ZNS] Loaded access token from Redis")
	}
	z.accessToken = accessToken

	// Load refresh token từ Redis
	refreshToken, err := z.redisClient.Get(ctx, ZNSRefreshTokenKey).Result()
	if err == redis.Nil {
		// Không có trong Redis, dùng token từ config
		refreshToken = config.Properties.ZNS.RefreshToken
		// Lưu vào Redis
		z.redisClient.Set(ctx, ZNSRefreshTokenKey, refreshToken, 0)
		fmt.Println("📝 [ZNS] Loaded refresh token from config and saved to Redis")
	} else if err != nil {
		// Lỗi khi đọc Redis, dùng token từ config
		refreshToken = config.Properties.ZNS.RefreshToken
		fmt.Printf("⚠️ [ZNS] Error loading refresh token from Redis: %v, using config\n", err)
	} else {
		fmt.Println("✅ [ZNS] Loaded refresh token from Redis")
	}
	z.refreshToken = refreshToken
}

// saveTokensToRedis lưu access token và refresh token mới vào Redis
func (z *ZnsProvider) saveTokensToRedis(ctx context.Context, accessToken, refreshToken string) error {
	// Lưu access token
	err := z.redisClient.Set(ctx, ZNSAccessTokenKey, accessToken, 0).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi lưu access token vào Redis: %w", err)
	}

	// Lưu refresh token
	err = z.redisClient.Set(ctx, ZNSRefreshTokenKey, refreshToken, 0).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi lưu refresh token vào Redis: %w", err)
	}

	fmt.Println("💾 [ZNS] Saved new tokens to Redis successfully")
	return nil
}

// ZNSRequest định nghĩa cấu trúc payload gửi đến ZNS API
type ZNSRequest struct {
	Phone        string                 `json:"phone"`
	TemplateID   string                 `json:"template_id"`
	TemplateData map[string]interface{} `json:"template_data"`
	TrackingID   string                 `json:"tracking_id"`
}

// ZNSResponse định nghĩa cấu trúc response từ ZNS API
type ZNSResponse struct {
	Error   int                    `json:"error"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// checkDailyLimit kiểm tra và tăng counter ZNS daily limit
// Giới hạn: 100 lần gọi/ngày, reset vào 00:00 mỗi ngày
func (z *ZnsProvider) checkDailyLimit(ctx context.Context) error {
	// Tạo key theo ngày hiện tại (format: zns:daily:limit:2025-10-16)
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("zns:daily:limit:%s", today)

	// Lấy giá trị counter hiện tại
	count, err := z.redisClient.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("lỗi khi lấy ZNS counter từ Redis: %w", err)
	}

	// Kiểm tra đã vượt quá limit chưa
	if count >= 100 {
		return fmt.Errorf("đã vượt quá giới hạn gửi ZNS hàng ngày (100 lần/ngày)")
	}

	// Tăng counter
	newCount, err := z.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("lỗi khi tăng ZNS counter: %w", err)
	}

	// Nếu là lần đầu tiên trong ngày, set TTL = đến hết ngày (00:00 ngày mai)
	if newCount == 1 {
		now := time.Now()
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		ttl := tomorrow.Sub(now)
		z.redisClient.Expire(ctx, key, ttl)
	}

	return nil
}

// cleanPhoneNumber chuẩn hóa số điện thoại sang format quốc tế (84xxx)
// Ví dụ:
//   - 0967771234 => 84967771234
//   - +84967771234 => 84967771234
//   - 84967771234 => 84967771234
//   - 0962 771 234 => 84962771234
func cleanPhoneNumber(phone string) string {
	// Loại bỏ khoảng trắng và các ký tự đặc biệt
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, ".", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Loại bỏ dấu + nếu có
	phone = strings.TrimPrefix(phone, "+")

	// Nếu số bắt đầu bằng 0, thay thế bằng 84
	if strings.HasPrefix(phone, "0") {
		phone = "84" + phone[1:]
	}

	// Nếu đã có 84 ở đầu, giữ nguyên
	// Nếu không có 84, thêm vào (trường hợp số không có mã vùng)
	if !strings.HasPrefix(phone, "84") {
		phone = "84" + phone
	}

	return phone
}

// SendZNSMessage gửi tin nhắn ZNS với template tùy chỉnh
func (z *ZnsProvider) SendZNSMessage(
	ctx context.Context,
	phone string,
	templateData map[string]interface{},
	templateID *int,
	trackingID *string,
) (map[string]interface{}, error) {
	// return nil, nil
	// Chuẩn hóa số điện thoại sang format quốc tế
	phone = cleanPhoneNumber(phone)
	if env.ENV_RUNTIME == "local" {
		return nil, nil
	}

	// return map[string]interface{}{
	// 	"success":     true,
	// 	"message":     "Gửi ZNS thành công",
	// 	"data":        nil,
	// 	"phone":       phone,
	// 	"tracking_id": "",
	// }, nil

	// TODO: Tạm thời bỏ rate limit, sẽ bật lại sau
	// Kiểm tra rate limit (100 lần/ngày)
	// if err := z.checkDailyLimit(ctx); err != nil {
	// 	return map[string]interface{}{
	// 		"success": false,
	// 		"message": err.Error(),
	// 		"phone":   phone,
	// 	}, err
	// }

	// Sử dụng template ID mặc định nếu không được cung cấp
	tplID := z.templateID
	if templateID != nil {
		tplID = *templateID
	}

	// Tạo tracking ID mặc định nếu không được cung cấp
	trkID := fmt.Sprintf("zns_%s_%d", phone, time.Now().Unix())
	if trackingID != nil {
		trkID = *trackingID
	}

	// Chuẩn bị payload
	payload := ZNSRequest{
		Phone:        phone,
		TemplateID:   fmt.Sprintf("%d", tplID),
		TemplateData: templateData,
		TrackingID:   trkID,
	}

	// Marshal payload thành JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi marshal payload ZNS: %w", err)
	}

	// Tạo HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://business.openapi.zalo.me/message/template",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo request ZNS: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", z.accessToken)

	// Gửi request
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gửi request ZNS: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var znsResp ZNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&znsResp); err != nil {
		return nil, fmt.Errorf("lỗi khi parse response ZNS: %w", err)
	}

	// Kiểm tra response status
	if resp.StatusCode != http.StatusOK {
		return map[string]interface{}{
			"success":     false,
			"message":     fmt.Sprintf("Gửi ZNS thất bại: %d", resp.StatusCode),
			"phone":       phone,
			"tracking_id": trkID,
		}, fmt.Errorf("ZNS API trả về status code %d: %s", resp.StatusCode, znsResp.Message)
	}

	// Kiểm tra error code từ ZNS
	if znsResp.Error != 0 {
		return map[string]interface{}{
			"success":     false,
			"message":     fmt.Sprintf("ZNS error: %s", znsResp.Message),
			"phone":       phone,
			"tracking_id": trkID,
		}, fmt.Errorf("ZNS API trả về error code %d: %s", znsResp.Error, znsResp.Message)
	}

	// Thành công
	return map[string]interface{}{
		"success":     true,
		"message":     "Gửi ZNS thành công",
		"data":        znsResp.Data,
		"phone":       phone,
		"tracking_id": trkID,
	}, nil
}

// SendOTPZNS gửi OTP qua ZNS
// TODO: Implement logic cập nhật OTP tracking (last_otp_at, otp_times) vào database
func (z *ZnsProvider) SendOTPZNS(
	ctx context.Context,
	phone string,
	otp string,
) (map[string]interface{}, error) {
	// Chuẩn bị template data cho OTP
	templateData := map[string]interface{}{
		"otp": otp,
	}

	// Sử dụng OTP template ID
	otpTemplateID := z.otpTemplateID

	// Gửi ZNS message
	result, err := z.SendZNSMessage(ctx, phone, templateData, &otpTemplateID, nil)

	// Log kết quả
	if err == nil {
		if success, ok := result["success"].(bool); ok && success {
			fmt.Printf("📱 [ZNS] Đã gửi OTP thành công cho phone %s\n", phone)
			// TODO: Cập nhật OTP tracking vào database (UserOTPEntity hoặc AuthMethod)
		}
	}

	return result, err
}

// ZNSRefreshTokenResponse định nghĩa cấu trúc response từ API refresh token
type ZNSRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Error        int    `json:"error"`
	Message      string `json:"message"`
}

// RefreshAccessToken làm mới access token cho ZNS sử dụng refresh token
// API: https://oauth.zaloapp.com/v4/oa/access_token
// Method: POST
// Headers: secret_key, Content-Type: application/x-www-form-urlencoded
// Body: refresh_token, app_id, grant_type='refresh_token'
// Response: access_token mới và refresh_token mới
func (z *ZnsProvider) RefreshAccessToken(ctx context.Context) (string, string, error) {
	// Chuẩn bị form data
	formData := url.Values{}
	formData.Set("refresh_token", z.refreshToken)
	formData.Set("app_id", config.Properties.OAuth.Zalo.AppID)
	formData.Set("grant_type", "refresh_token")

	// Tạo HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://oauth.zaloapp.com/v4/oa/access_token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return "", "", fmt.Errorf("lỗi khi tạo request refresh token: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("secret_key", config.Properties.OAuth.Zalo.SecretKey)

	// Gửi request
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("lỗi khi gửi request refresh token: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("lỗi khi đọc response refresh token: %w", err)
	}

	// Parse response
	var znsResp ZNSRefreshTokenResponse
	if err := json.Unmarshal(body, &znsResp); err != nil {
		return "", "", fmt.Errorf("lỗi khi parse response refresh token: %w", err)
	}

	// Kiểm tra response status
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("refresh token thất bại với status code %d: %s", resp.StatusCode, znsResp.Message)
	}

	// Kiểm tra error code từ Zalo
	if znsResp.Error != 0 {
		return "", "", fmt.Errorf("zalo API trả về error code %d: %s", znsResp.Error, znsResp.Message)
	}

	// Kiểm tra access token và refresh token có hợp lệ không
	if znsResp.AccessToken == "" || znsResp.RefreshToken == "" {
		return "", "", fmt.Errorf("refresh token trả về access_token hoặc refresh_token rỗng")
	}

	// Cập nhật access token và refresh token mới vào instance
	z.accessToken = znsResp.AccessToken
	z.refreshToken = znsResp.RefreshToken

	// Lưu tokens mới vào Redis để persist
	if err := z.saveTokensToRedis(ctx, znsResp.AccessToken, znsResp.RefreshToken); err != nil {
		fmt.Printf("⚠️ [ZNS] Warning: Failed to save tokens to Redis: %v\n", err)
		// Không return error vì token đã được cập nhật trong memory
	}

	// Log thành công
	fmt.Printf("✅ [ZNS] Refresh token thành công. Access token expires in: %d seconds\n", znsResp.ExpiresIn)

	return znsResp.AccessToken, znsResp.RefreshToken, nil
}

// UpdateTokens cập nhật access token và refresh token mới
// Dùng để update token manually từ response API
func (z *ZnsProvider) UpdateTokens(ctx context.Context, accessToken, refreshToken string) error {
	if accessToken == "" || refreshToken == "" {
		return fmt.Errorf("access token và refresh token không được để trống")
	}

	// Cập nhật tokens trong memory
	z.accessToken = accessToken
	z.refreshToken = refreshToken

	// Lưu vào Redis
	if err := z.saveTokensToRedis(ctx, accessToken, refreshToken); err != nil {
		return fmt.Errorf("lỗi khi lưu tokens vào Redis: %w", err)
	}

	fmt.Println("✅ [ZNS] Tokens updated successfully")
	return nil
}

// GetAccessToken lấy access token hiện tại
func (z *ZnsProvider) GetAccessToken() string {
	return z.accessToken
}

// GetRefreshToken lấy refresh token hiện tại
func (z *ZnsProvider) GetRefreshToken() string {
	return z.refreshToken
}
