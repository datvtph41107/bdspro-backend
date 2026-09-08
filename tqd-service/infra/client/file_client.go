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
	"net/url"
	"strconv"
	"strings"
	"time"

	"tqd/internal/interface/provider"
)

const fileClientServiceName = "tqd-service"

// FileClient gọi sang file-service để upload/tải file. tqd-service chỉ lưu durable File reference
// trả về và không sở hữu binary/storage semantics.
type FileClient struct {
	baseURL        string
	serviceAuthKey string
	httpClient     *http.Client
}

// FileConfig is process-owned configuration injected at composition.
type FileConfig struct {
	BaseURL        string
	ServiceAuthKey string
	Timeout        time.Duration
}

func NewFileClient(cfg FileConfig) provider.FileProvider {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8002"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &FileClient{
		baseURL:        baseURL,
		serviceAuthKey: strings.TrimSpace(cfg.ServiceAuthKey),
		httpClient:     &http.Client{Timeout: timeout},
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
// Segment s is a non-empty compatibility placeholder; service-to-service authorization is carried by X-Service-Auth.
func (c *FileClient) GetFile(ctx context.Context, path string) ([]byte, error) {
	requestURL := fmt.Sprintf("%s/v1/file/load/%s/local", c.baseURL, url.PathEscape(path))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
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
