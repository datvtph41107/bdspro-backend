package client

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/liyue201/goqr"
	"github.com/spf13/viper"
)

// @bind: assistant/internal/interface/provider.QRCodeProvider

type QRCodeClient struct {
	httpClient *http.Client
	tempDir    string
}

func NewQRCodeClient() *QRCodeClient {
	timeoutSeconds := viper.GetInt("ocr.timeout_seconds")
	if timeoutSeconds == 0 {
		timeoutSeconds = 60
	}

	// Tạo thư mục tạm
	tempDir := viper.GetString("ocr.temp_dir")
	if tempDir == "" {
		tempDir = os.TempDir() + "/assistant-qrcode"
	}
	os.MkdirAll(tempDir, 0755)

	return &QRCodeClient{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		tempDir: tempDir,
	}
}

// ReadQRCodeFromURL - Đọc QR code từ URL ảnh và trả về raw text
func (c *QRCodeClient) ReadQRCodeFromURL(ctx context.Context, imageURL string) (string, error) {
	// Download ảnh từ URL
	imageData, err := c.downloadImage(ctx, imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	// Decode ảnh từ bytes
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		// Nếu decode trực tiếp thất bại, lưu vào file tạm rồi decode
		imagePath := filepath.Join(c.tempDir, fmt.Sprintf("qr_%d.png", time.Now().UnixNano()))
		defer os.Remove(imagePath)

		if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
			return "", fmt.Errorf("failed to save image: %w", err)
		}

		file, err := os.Open(imagePath)
		if err != nil {
			return "", fmt.Errorf("failed to open image: %w", err)
		}
		defer file.Close()

		img, _, err = image.Decode(file)
		if err != nil {
			return "", fmt.Errorf("failed to decode image: %w", err)
		}
	}

	// Convert sang grayscale
	gray := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	// Decode QR code
	codes, err := goqr.Recognize(gray)
	if err != nil {
		return "", fmt.Errorf("failed to recognize QR code: %w", err)
	}

	if len(codes) == 0 {
		return "", fmt.Errorf("no QR code found in image")
	}

	// Trả về text của QR code đầu tiên
	return string(codes[0].Payload), nil
}

// downloadImage - Download ảnh từ URL
func (c *QRCodeClient) downloadImage(ctx context.Context, imageURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status=%d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageData, nil
}
