package zns

import (
	"auth/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// @bind: auth/internal/interface/provider.IZnsProvider
type ZnsProvider struct {
	accessToken   string
	templateID    int
	otpTemplateID int
	httpClient    *http.Client
	redisClient   *redis.Client
}

// NewZnsProvider tạo instance mới của ZnsProvider
func NewZnsProvider() *ZnsProvider {
	// Khởi tạo Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Properties.Redis.Host, config.Properties.Redis.Port),
		Password: config.Properties.Redis.Password,
		DB:       0,
	})

	return &ZnsProvider{
		accessToken:   config.Properties.ZNS.AccessToken,
		templateID:    config.Properties.ZNS.TemplateID,
		otpTemplateID: config.Properties.ZNS.OTPTemplateID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		redisClient: redisClient,
	}
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
	// Chuẩn hóa số điện thoại sang format quốc tế
	phone = cleanPhoneNumber(phone)

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
