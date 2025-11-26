package client

import (
	_utils "common/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// @bind: assistant/internal/interface/provider.FileProvider

type FileClient struct {
	baseURL        string
	serviceAuthKey string
	httpClient     *http.Client
}

func NewFileClient() *FileClient {
	baseURL := viper.GetString("file_service.base_url")
	if baseURL == "" {
		baseURL = "http://localhost:8002"
	}

	serviceAuthKey := viper.GetString("file_service.service_auth_key")
	if serviceAuthKey == "" {
		serviceAuthKey = os.Getenv("FILE_SERVICE_AUTH_KEY")
	}

	timeoutSeconds := viper.GetInt("file_service.timeout_seconds")
	if timeoutSeconds == 0 {
		timeoutSeconds = 30
	}

	return &FileClient{
		baseURL:        baseURL,
		serviceAuthKey: serviceAuthKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// GetFile lấy file từ file-service với quyền truy cập service internal
// p: encoded path (ví dụ: "p" + encoded path cho public, hoặc encoded path cho private)
// s: signature (có thể để trống nếu là service internal request)
func (c *FileClient) GetFile(ctx context.Context, p string, s string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/file/load/%s/%s", c.baseURL, p, s)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Tạo HMAC signature với timestamp để chống replay attack
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	serviceName := "assistant-service"

	// Signature = HMAC-SHA256(timestamp:service_name, secret_key)
	data := timestampStr + ":" + serviceName
	signature := _utils.GenerateHmacSHA256(data, c.serviceAuthKey)

	// Format: timestamp:signature
	authHeader := timestampStr + ":" + signature

	// Thêm header X-Service-Auth để xác thực service internal
	req.Header.Set("X-Service-Auth", authHeader)
	req.Header.Set("X-Service-Name", serviceName)

	fmt.Printf("\nExtractNationalCard: authHeader: %s", authHeader)
	fmt.Printf("\nExtractNationalCard: serviceName: %s", serviceName)
	fmt.Printf("\nExtractNationalCard: timestamp: %s", timestampStr)
	fmt.Printf("\nExtractNationalCard: signature: %s", signature)
	fmt.Printf("\nExtractNationalCard: data: %s", data)
	fmt.Printf("\nExtractNationalCard: url: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file service error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// GetFileURL trả về URL đầy đủ để truy cập file (dùng cho các trường hợp cần URL thay vì download)
func (c *FileClient) GetFileURL(p string, s string) string {
	return fmt.Sprintf("%s/v1/file/load/%s/%s", c.baseURL, p, s)
}
