package client

import (
	"bytes"
	_dto "common/domain/dto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// GeminiClient calls Google Gemini generateContent for product text analysis.
type GeminiClient struct {
	apiKey     string
	baseURL    string
	model      string
	maxTokens  int
	httpClient *http.Client
}

type GeminiRequest struct {
	Contents         []GeminiContent         `json:"contents"`
	GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text,omitempty"`
}

type GeminiGenerationConfig struct {
	Temperature     *float32 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

func NewGeminiClient() *GeminiClient {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = viper.GetString("gemini.api_key")
	}

	baseURL := viper.GetString("gemini.base_url")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	model := viper.GetString("gemini.model")
	if model == "" {
		model = "gemini-2.5-flash-lite"
	}

	timeoutSeconds := viper.GetInt("gemini.timeout_seconds")
	if timeoutSeconds == 0 {
		timeoutSeconds = 60
	}
	maxTokens := viper.GetInt("gemini.max_tokens")
	if maxTokens <= 0 {
		maxTokens = 4000
	}

	return &GeminiClient{
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		maxTokens: maxTokens,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *GeminiClient) AnalyzeProductText(ctx context.Context, content string) (*_dto.ProductV3DTO, error) {
	systemPrompt := `Bạn là trợ lý AI chuyên phân tích thông tin bất động sản tiếng Việt. 
Nhiệm vụ của bạn là trích xuất thông tin sản phẩm bất động sản từ văn bản người dùng cung cấp.

Hãy phân tích và trả về JSON với cấu trúc sau:
{
  "name": "Tên/tiêu đề sản phẩm (bắt buộc)",
  "area": diện tích (số float, đơn vị m2, bắt buộc),
  "transactionType": {"name": "Bán/Cho thuê", "id": 10/20},
  "googleMapLink": "Link Google Maps nếu có",
  "address": {
    "provinceName": "Tỉnh/Thành phố",
    "districtName": "Quận/Huyện",
    "wardName": "Phường/Xã",
    "detail": "Địa chỉ chi tiết"
  },
  "sourceType": {"name": "Loại nguồn (Chủ đất, Môi giới, Sàn, Khác)", "id": 10/20/30/40},
  "propertyType": {"name":"Loại BĐS (nhà riêng, chung cư, đất nền, biệt thự, ...)"},
  "legalDoc": {"name":"Giấy tờ pháp lý (sổ đỏ/hồng, sổ chung, ...)"},
  "priceData": {
    "salePrice": giá bán (số float, VND),
    "saleCommission": hoa hồng bán (số float, VND),
    "deposite": tiền đặt cọc (số float, VND),
    "rentPrice": giá thuê (số float, VND/tháng),
    "rentCommission": hoa hồng thuê (số float, VND),
    "rentPaymentCycle": Chu kỳ thanh toán (1,2,3),
    "currency": "VND"
  },
  "houseInfo": {
    "numBedroom": số phòng ngủ (số int),
    "numBathroom": số phòng tắm/WC (số int),
    "numFloor": số tầng (số int),
    "numFront": số mặt tiền (số int),
    "numCarPark": số chỗ đậu xe (số int),
    "furniture": "Nội thất (full, partial, none)",
    "orientation": "Hướng nhà (Đông, Tây, Nam, Bắc, Đông Nam, ...)"
  },
  "productPrivate": {
    "importPrice": giá nhập (số float, VND),
    "operatingCost": chi phí vận hành (số float, VND),
    "internalNote": "Ghi chú nội bộ",
    "targetProfit": lợi nhuận mục tiêu (số float, VND)
  },
  "amenities": [{"name": "tiện ích 1", "id": 1}, {"name": "tiện ích 2", "id": 2}]
}

Lưu ý:
- Các trường BẮT BUỘC: name, area
- Các trường khác chỉ trả về nếu có thông tin trong văn bản
- Trường nào không có thông tin thì KHÔNG trả về (omit)
- Chỉ trả về JSON thuần, KHÔNG có markdown, code block hay text khác
- Số liệu giá cả phải là số, không có ký tự đặc biệt (tỷ, triệu...)
- transactionType: 10 cho "Bán", 20 cho "Cho thuê"`

	userMessage := fmt.Sprintf("%s\n\nVăn bản cần phân tích:\n%s", systemPrompt, content)

	response, err := c.sendRequest(ctx, userMessage, 0, 0)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in gemini response")
	}

	contentStr := response.Candidates[0].Content.Parts[0].Text
	contentStr = cleanJSONFromMarkdownGemini(contentStr)

	var result _dto.ProductV3DTO
	if err := json.Unmarshal([]byte(contentStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse product analysis: %w, content: %s", err, contentStr)
	}

	return &result, nil
}

func (c *GeminiClient) GetModel() string {
	return c.model
}

func (c *GeminiClient) sendRequest(ctx context.Context, userMessage string, maxTokens int, temperature float32) (*GeminiResponse, error) {
	if err := validateProviderCall("gemini", c.apiKey); err != nil {
		return nil, err
	}
	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}
	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: userMessage}}},
		},
	}
	config := &GeminiGenerationConfig{MaxOutputTokens: &maxTokens}
	if temperature > 0 {
		config.Temperature = &temperature
	}
	requestBody.GenerationConfig = config

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &geminiResp, nil
}

func cleanJSONFromMarkdownGemini(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}
