package client

import (
	"assistant/config"
	"assistant/internal/dto"
	"bytes"
	_dto "common/domain/dto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ++@bind: assistant/internal/interface/provider.DeepseekProvider

// @bind: assistant/internal/interface/provider.DeepseekProvider
type DeepseekClient struct {
	apiKey       string
	baseURL      string
	model        string
	maxTokens    int
	providerMode string
	httpClient   *http.Client
}

func NewDeepseekClient(runtime config.Runtime) *DeepseekClient {
	cfg := runtime.Deepseek
	return &DeepseekClient{
		apiKey:       cfg.APIKey,
		baseURL:      cfg.BaseURL,
		model:        cfg.Model,
		maxTokens:    cfg.MaxTokens,
		providerMode: runtime.ProviderMode,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// "description": "Mô tả chi tiết sản phẩm",
// "note": "Ghi chú thêm",

// AnalyzeProductText - Phân tích văn bản sản phẩm BĐS
func (c *DeepseekClient) AnalyzeProductText(ctx context.Context, content string) (*_dto.ProductV3DTO, error) {
	systemPrompt := `Bạn là trợ lý AI chuyên phân tích thông tin bất động sản tiếng Việt. 
Nhiệm vụ của bạn là trích xuất thông tin sản phẩm bất động sản từ văn bản người dùng cung cấp.

Hãy phân tích và trả về JSON với cấu trúc sau:
{
  "name": "Tên/tiêu đề sản phẩm (bắt buộc)",
  "area": diện tích (số float, đơn vị m2, bắt buộc),
  "transactionType": 10 hoặc 20 (10=Bán, 20=Cho thuê),
  "googleMapLink": "Link Google Maps nếu có",
  "address": {
    "province": "Tỉnh/Thành phố",
    "district": "Quận/Huyện",
    "ward": "Phường/Xã",
    "detail": "Địa chỉ chi tiết"
  },
  "propertyType": "Loại BĐS (nhà riêng, chung cư, đất nền, biệt thự, ...)",
  "legalDoc": "Giấy tờ pháp lý (sổ đỏ/hồng, sổ chung, ...)",
  "priceData": {
    "salePrice": giá bán (số float, VND),
    "saleCommission": hoa hồng bán (số float, VND),
    "deposite": tiền đặt cọc (số float, VND),
    "rentPrice": giá thuê (số float, VND/tháng),
    "rentCommission": hoa hồng thuê (số float, VND),
    "rentPaymentCycle": chu kỳ thanh toán thuê (số: 1=tháng, 3=quý, 6=6 tháng, 12=năm),
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
  "amenities": ["tiện ích 1", "tiện ích 2", "..."]
}

Lưu ý:
- Các trường BẮT BUỘC: name, area
- Các trường khác chỉ trả về nếu có thông tin trong văn bản
- Trường nào không có thông tin thì KHÔNG trả về (omit)
- Chỉ trả về JSON thuần, KHÔNG có markdown, code block hay text khác
- Số liệu giá cả phải là số, không có ký tự đặc biệt (tỷ, triệu...)
- transactionType: 10 cho "Bán", 20 cho "Cho thuê"`

	response, err := c.sendRequest(ctx, systemPrompt, content, 0, 0)
	if err != nil {
		return nil, err
	}

	// Parse JSON từ response
	contentStr := response.Choices[0].Message.Content
	contentStr = cleanJSONFromMarkdown(contentStr)

	var result _dto.ProductV3DTO
	if err := json.Unmarshal([]byte(contentStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse product analysis: %w, content: %s", err, contentStr)
	}

	return &result, nil
}

// Chat - Thực hiện conversation với AI
func (c *DeepseekClient) Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error) {
	if err := validateProviderCall("deepseek", c.providerMode, c.apiKey); err != nil {
		return "", err
	}
	systemPrompt := `Bạn là trợ lý AI thông minh, nhiệt tình và hữu ích. 
Bạn có kiến thức sâu rộng về bất động sản tại Việt Nam.
Hãy trả lời câu hỏi một cách chính xác, ngắn gọn và dễ hiểu.`

	// Build messages với history
	messages := []dto.DeepseekMessage{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, dto.DeepseekMessage{
		Role:    "user",
		Content: message,
	})

	requestBody := dto.DeepseekRequest{
		Model:     c.model,
		Messages:  messages,
		Stream:    false,
		MaxTokens: c.maxTokens,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deepseek API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var deepseekResp dto.DeepseekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in deepseek response")
	}

	return deepseekResp.Choices[0].Message.Content, nil
}

// GenerateContent - Tạo nội dung theo prompt
func (c *DeepseekClient) GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error) {
	response, err := c.sendRequest(ctx, "", prompt, maxTokens, temperature)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

// sendRequest - Helper để gửi request đến DeepSeek API
func (c *DeepseekClient) sendRequest(ctx context.Context, systemPrompt, userMessage string, maxTokens int, temperature float32) (*dto.DeepseekResponse, error) {
	if err := validateProviderCall("deepseek", c.providerMode, c.apiKey); err != nil {
		return nil, err
	}
	messages := []dto.DeepseekMessage{}

	if systemPrompt != "" {
		messages = append(messages, dto.DeepseekMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	messages = append(messages, dto.DeepseekMessage{
		Role:    "user",
		Content: userMessage,
	})

	requestBody := dto.DeepseekRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}

	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}
	requestBody.MaxTokens = maxTokens
	if temperature > 0 {
		requestBody.Temperature = temperature
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
		return nil, fmt.Errorf("deepseek API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var deepseekResp dto.DeepseekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &deepseekResp, nil
}

// cleanJSONFromMarkdown - Loại bỏ markdown code block wrapper nếu có
func cleanJSONFromMarkdown(content string) string {
	content = strings.TrimSpace(content)

	// Remove ```json và ``` wrapper nếu có
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}

	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content)
}

// GetModel - Trả về model name đang sử dụng
func (c *DeepseekClient) GetModel() string {
	return c.model
}
