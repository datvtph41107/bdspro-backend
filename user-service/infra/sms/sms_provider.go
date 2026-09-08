package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"user/config"
	"user/internal/interface/providers"
)

type SMSProvider struct {
	httpClient *http.Client
	baseURL    string
}

const (
	defaultTimeout = 30 * time.Second
)

// NewSMSProvider khởi tạo SMSProvider với HTTP client và endpoint cấu hình
func NewSMSProvider() providers.ISmsProvider {
	baseURL := config.Properties.SMS.Endpoint
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://your-domain.com/api/v1/sms/send"
	}

	return &SMSProvider{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}
}

type smsRequest struct {
	Phone string `json:"phone"`
	Brand string `json:"brand"`
	Msg   string `json:"msg"`
}

// SendSMS gửi tin nhắn SMS tới nhà cung cấp thông qua API REST
func (s *SMSProvider) SendSMS(
	ctx context.Context,
	phone string,
	brand string,
	message string,
) (map[string]interface{}, error) {
	normalizedPhone := normalizePhoneNumber(phone)

	if !config.Properties.OutboundMessaging.Enabled {
		return map[string]interface{}{
			"success": true,
			"message": "Skip sending SMS because outbound messaging is disabled",
			"phone":   normalizedPhone,
			"brand":   brand,
		}, nil
	}

	payload := smsRequest{
		Phone: normalizedPhone,
		Brand: brand,
		Msg:   message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("sms marshal payload error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("sms create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", config.Properties.SMS.AuthToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sms send request error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sms read response error: %w", err)
	}

	var data map[string]interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &data); err != nil {
			data = map[string]interface{}{
				"raw_response": string(respBody),
			}
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if data == nil {
			data = make(map[string]interface{})
		}

		data["success"] = false
		data["status_code"] = resp.StatusCode
		data["message"] = fmt.Sprintf("sms provider returned non-success status: %d", resp.StatusCode)

		return data, fmt.Errorf("sms provider error: status %d", resp.StatusCode)
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	data["success"] = true
	data["phone"] = normalizedPhone
	data["brand"] = brand

	return data, nil
}

func normalizePhoneNumber(phone string) string {
	if phone == "" {
		return phone
	}

	phone = strings.TrimSpace(phone)
	replacer := strings.NewReplacer(" ", "", "-", "", ".", "", "(", "", ")", "")
	phone = replacer.Replace(phone)
	phone = strings.TrimPrefix(phone, "+")

	if strings.HasPrefix(phone, "0") {
		phone = "84" + phone[1:]
	}

	if !strings.HasPrefix(phone, "84") {
		phone = "84" + phone
	}

	return phone
}
