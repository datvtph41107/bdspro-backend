package zns

import (
	"auth/config"
	"auth/env"
	"auth/internal/domain"
	"auth/internal/interface/repo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// @bind: auth/internal/interface/provider.IZnsProvider
type ZnsProvider struct {
	accessToken    string
	refreshToken   string
	templateID     int
	otpTemplateID  int
	httpClient     *http.Client
	redisClient    *redis.Client
	authConfigRepo repo.IAuthConfigRepo // Repository để lấy config từ DB
	mu             sync.RWMutex         // Mutex để đồng bộ khi đọc/ghi token
}

// NewZnsProvider tạo instance mới của ZnsProvider
func NewZnsProvider(authConfigRepo repo.IAuthConfigRepo) *ZnsProvider {
	// Khởi tạo Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Properties.Redis.Host, config.Properties.Redis.Port),
		Password: config.Properties.Redis.Password,
		DB:       0,
	})

	provider := &ZnsProvider{
		authConfigRepo: authConfigRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		redisClient: redisClient,
	}

	// Load tất cả config từ DB, nếu không có thì dùng config file
	provider.loadAllConfigFromDB()

	return provider
}

// loadAllConfigFromDB load tất cả config từ DB (tokens + template IDs)
// Nếu không có trong DB thì load từ config file
func (z *ZnsProvider) loadAllConfigFromDB() {
	ctx := context.Background()

	z.mu.Lock()
	defer z.mu.Unlock()

	// Load access token từ DB
	accessTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSToken)
	if err != nil {
		// Không có trong DB, dùng token từ config
		z.accessToken = config.Properties.ZNS.AccessToken
		fmt.Printf("⚠️ [ZNS] Error loading access token from DB: %v, using config\n", err)
	} else {
		z.accessToken = accessTokenConfig.ConfigValue
		fmt.Println("✅ [ZNS] Loaded access token from DB")
	}

	// Load refresh token từ DB
	refreshTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSRefresh)
	if err != nil {
		// Không có trong DB, dùng token từ config
		z.refreshToken = config.Properties.ZNS.RefreshToken
		fmt.Printf("⚠️ [ZNS] Error loading refresh token from DB: %v, using config\n", err)
	} else {
		z.refreshToken = refreshTokenConfig.ConfigValue
		fmt.Println("✅ [ZNS] Loaded refresh token from DB")
	}

	// Load template ID từ DB
	templateConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSTemplateID)
	if err != nil {
		// Không có trong DB, dùng từ config
		z.templateID = config.Properties.ZNS.TemplateID
		fmt.Printf("⚠️ [ZNS] Error loading template ID from DB: %v, using config\n", err)
	} else {
		// Parse string to int
		if templateID, err := strconv.Atoi(templateConfig.ConfigValue); err == nil {
			z.templateID = templateID
			fmt.Printf("✅ [ZNS] Loaded template ID from DB: %d\n", templateID)
		} else {
			z.templateID = config.Properties.ZNS.TemplateID
			fmt.Printf("⚠️ [ZNS] Error parsing template ID from DB: %v, using config\n", err)
		}
	}

	// Load OTP template ID từ DB
	otpTemplateConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSOTPTemplate)
	if err != nil {
		// Không có trong DB, dùng từ config
		z.otpTemplateID = config.Properties.ZNS.OTPTemplateID
		fmt.Printf("⚠️ [ZNS] Error loading OTP template ID from DB: %v, using config\n", err)
	} else {
		// Parse string to int
		if otpTemplateID, err := strconv.Atoi(otpTemplateConfig.ConfigValue); err == nil {
			z.otpTemplateID = otpTemplateID
			fmt.Printf("✅ [ZNS] Loaded OTP template ID from DB: %d\n", otpTemplateID)
		} else {
			z.otpTemplateID = config.Properties.ZNS.OTPTemplateID
			fmt.Printf("⚠️ [ZNS] Error parsing OTP template ID from DB: %v, using config\n", err)
		}
	}
}

// reloadConfigFromDB reload tất cả config từ DB để đảm bảo dùng config mới nhất
// Gọi method này trước khi gửi ZNS để sync config đã được cập nhật
func (z *ZnsProvider) reloadConfigFromDB(ctx context.Context) error {
	z.mu.Lock()
	defer z.mu.Unlock()

	// Load access token từ DB
	accessTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSToken)
	if err != nil {
		return fmt.Errorf("lỗi khi reload access token từ DB: %w", err)
	}
	if accessTokenConfig.ConfigValue != "" {
		z.accessToken = accessTokenConfig.ConfigValue
	}

	// Load refresh token từ DB
	refreshTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSRefresh)
	if err != nil {
		return fmt.Errorf("lỗi khi reload refresh token từ DB: %w", err)
	}
	if refreshTokenConfig.ConfigValue != "" {
		z.refreshToken = refreshTokenConfig.ConfigValue
	}

	// Load template ID từ DB
	templateConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSTemplateID)
	if err == nil && templateConfig.ConfigValue != "" {
		if templateID, err := strconv.Atoi(templateConfig.ConfigValue); err == nil {
			z.templateID = templateID
		}
	}

	// Load OTP template ID từ DB
	otpTemplateConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSOTPTemplate)
	if err == nil && otpTemplateConfig.ConfigValue != "" {
		if otpTemplateID, err := strconv.Atoi(otpTemplateConfig.ConfigValue); err == nil {
			z.otpTemplateID = otpTemplateID
		}
	}

	return nil
}

// saveTokensToDB lưu access token và refresh token mới vào DB
func (z *ZnsProvider) saveTokensToDB(ctx context.Context, accessToken, refreshToken string) error {
	z.mu.Lock()
	defer z.mu.Unlock()

	// Lưu access token vào DB
	err := z.authConfigRepo.UpdateValue(ctx, domain.ConfigKeyZNSToken, accessToken)
	if err != nil {
		return fmt.Errorf("lỗi khi lưu access token vào DB: %w", err)
	}

	// Lưu refresh token vào DB
	err = z.authConfigRepo.UpdateValue(ctx, domain.ConfigKeyZNSRefresh, refreshToken)
	if err != nil {
		return fmt.Errorf("lỗi khi lưu refresh token vào DB: %w", err)
	}

	// Cập nhật token trong memory luôn
	z.accessToken = accessToken
	z.refreshToken = refreshToken

	fmt.Println("💾 [ZNS] Saved new tokens to DB and updated in memory successfully")
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
	// Reload config từ DB để đảm bảo dùng config mới nhất (sau khi refresh)
	if err := z.reloadConfigFromDB(ctx); err != nil {
		fmt.Printf("⚠️ [ZNS] Warning: Failed to reload config from DB: %v, using cached config\n", err)
	}

	// return nil, nil
	// Chuẩn hóa số điện thoại sang format quốc tế
	phone = cleanPhoneNumber(phone)
	if env.ENV_RUNTIME == "" {
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

	// Set headers (dùng RLock để đọc token an toàn)
	z.mu.RLock()
	currentAccessToken := z.accessToken
	z.mu.RUnlock()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", currentAccessToken)

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
	ExpiresIn    string `json:"expires_in"`
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
	// ✅ Lấy refresh token mới nhất từ DB thay vì memory
	refreshTokenConfig, err := z.authConfigRepo.GetByKey(ctx, domain.ConfigKeyZNSRefresh)
	if err != nil {
		return "", "", fmt.Errorf("lỗi khi lấy refresh token từ DB: %w", err)
	}

	currentRefreshToken := refreshTokenConfig.ConfigValue
	if currentRefreshToken == "" {
		return "", "", fmt.Errorf("refresh token trong DB rỗng")
	}

	fmt.Printf("🔄 [ZNS] Using refresh token from DB: %s...\n", currentRefreshToken)

	// Chuẩn bị form data
	formData := url.Values{}
	formData.Set("refresh_token", currentRefreshToken)
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

	// Lưu tokens mới vào DB và cập nhật trong memory (thread-safe)
	if err := z.saveTokensToDB(ctx, znsResp.AccessToken, znsResp.RefreshToken); err != nil {
		return "", "", fmt.Errorf("lỗi khi lưu tokens mới: %w", err)
	}

	// Log thành công
	fmt.Printf("✅ [ZNS] Refresh token thành công. Access token expires in: %s seconds\n", znsResp.ExpiresIn)

	return znsResp.AccessToken, znsResp.RefreshToken, nil
}
