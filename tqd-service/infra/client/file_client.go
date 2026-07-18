package client

import (
	"bytes"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"

	"github.com/spf13/viper"
)

const fileClientServiceName = "tqd-service"

// FileClient gọi sang file-service để upload/tải file. tqd-service chỉ lưu lại path trả về
// (không lưu binary), giống cách assistant-service dùng file-service.
type FileClient struct {
	baseURL        string
	serviceAuthKey string
	httpClient     *http.Client
}

// NewFileClient trả về provider.FileProvider trực tiếp nên không cần wire.Bind.
func NewFileClient() provider.FileProvider {
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

// setServiceAuthHeader gắn HMAC signature theo chuẩn X-Service-Auth của file-service.
func (c *FileClient) setServiceAuthHeader(req *http.Request) {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	data := timestampStr + ":" + fileClientServiceName
	signature := _utils.GenerateHmacSHA256(data, c.serviceAuthKey)

	req.Header.Set("X-Service-Auth", timestampStr+":"+signature)
	req.Header.Set("X-Service-Name", fileClientServiceName)
}

type uploadFileResponse struct {
	File struct {
		Path string `json:"path"`
	} `json:"file"`
}

// UploadFile gửi file lên file-service (POST /v1/file/upload), trả về path đã mã hoá
// (dùng để lưu vào QHPlanningDocument.Filepath và gọi lại GetFile sau này).
func (c *FileClient) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("copy file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/file/upload", body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.setServiceAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload file: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read upload response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file-service error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var parsed uploadFileResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("parse upload response: %w", err)
	}
	if parsed.File.Path == "" {
		return "", fmt.Errorf("file-service không trả về path hợp lệ")
	}
	return parsed.File.Path, nil
}

// GetFile tải nội dung file từ file-service theo path đã lưu (GET /v1/file/load/{p}/{s}).
// Signature (s) để trống vì đây là service-to-service request, đã xác thực qua X-Service-Auth.
func (c *FileClient) GetFile(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/file/load/%s/%s", c.baseURL, path, "")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setServiceAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file-service error: status=%d body=%s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *FileClient) DeleteFile(ctx context.Context, fileURL string) error {
	// Chưa cần dùng ở phase này (chưa có API xoá file trên file-service cho service nội bộ).
	return nil
}

func (c *FileClient) GetFileURL(ctx context.Context, filePath string) string {
	return fmt.Sprintf("%s/v1/file/load/%s/%s", c.baseURL, filePath, "")
}

func (c *FileClient) ValidateFile(ctx context.Context, file io.Reader, filename string, maxSize int64, allowedTypes []string) error {
	// Validate cơ bản do file-service tự xử lý phía server; chưa cần validate thêm ở tqd.
	return nil
}

func (c *FileClient) GetFileInfo(ctx context.Context, fileURL string) (*dto.FileInfoDTO, error) {
	return nil, nil
}

func (c *FileClient) GenerateThumbnail(ctx context.Context, fileURL string) (string, error) {
	return "", nil
}
