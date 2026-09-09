package client

import (
	"assistant/config"
	"assistant/internal/dto"
	"bytes"
	_dto "common/domain/dto"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// @bind: assistant/internal/interface/provider.GeminiProvider
type GeminiClient struct {
	apiKey       string
	baseURL      string
	model        string
	maxTokens    int
	providerMode string
	httpClient   *http.Client
}

// GeminiRequest - Request structure for Gemini API
type GeminiRequest struct {
	Contents         []GeminiContent         `json:"contents"`
	GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"`
}

// GeminiContent - Content trong request
type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart - Part trong content
type GeminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *GeminiInlineData `json:"inlineData,omitempty"`
}

// GeminiInlineData - Image data inline
type GeminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // Base64 encoded image
}

// GeminiGenerationConfig - Generation config
type GeminiGenerationConfig struct {
	Temperature     *float32 `json:"temperature,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
	TopP            *float32 `json:"topP,omitempty"`
	TopK            *int     `json:"topK,omitempty"`
}

// GeminiResponse - Response từ Gemini API
type GeminiResponse struct {
	Candidates     []GeminiCandidate     `json:"candidates"`
	UsageMetadata  GeminiUsageMetadata   `json:"usageMetadata,omitempty"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
}

// GeminiCandidate - Candidate trong response
type GeminiCandidate struct {
	Content       GeminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason,omitempty"`
	Index         int                  `json:"index"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

// GeminiSafetyRating - Safety rating
type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// GeminiUsageMetadata - Usage metadata
type GeminiUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// GeminiPromptFeedback - Prompt feedback
type GeminiPromptFeedback struct {
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

func NewGeminiClient(runtime config.Runtime) *GeminiClient {
	cfg := runtime.Gemini
	return &GeminiClient{
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

// AnalyzeProductText - Phân tích văn bản sản phẩm BĐS
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

	// Parse JSON từ response
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

// Chat - Thực hiện conversation với AI
func (c *GeminiClient) Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error) {
	if err := validateProviderCall("gemini", c.providerMode, c.apiKey); err != nil {
		return "", err
	}
	systemPrompt := `Bạn là trợ lý AI thông minh, nhiệt tình và hữu ích. 
Bạn có kiến thức sâu rộng về bất động sản tại Việt Nam.
Hãy trả lời câu hỏi một cách chính xác, ngắn gọn và dễ hiểu.`

	// Build contents từ history
	contents := []GeminiContent{}

	// Add system prompt as first user message (Gemini không có system role)
	if len(history) == 0 {
		contents = append(contents, GeminiContent{
			Role:  "user",
			Parts: []GeminiPart{{Text: systemPrompt}},
		})
		contents = append(contents, GeminiContent{
			Role:  "model",
			Parts: []GeminiPart{{Text: "Được ạ, tôi sẵn sàng hỗ trợ bạn về bất động sản."}},
		})
	}

	// Add history
	for _, msg := range history {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		} else if role == "system" {
			// Skip system messages hoặc convert thành user message
			continue
		}
		contents = append(contents, GeminiContent{
			Role:  role,
			Parts: []GeminiPart{{Text: msg.Content}},
		})
	}

	// Add current message
	contents = append(contents, GeminiContent{
		Role:  "user",
		Parts: []GeminiPart{{Text: message}},
	})

	requestBody := GeminiRequest{
		Contents:         contents,
		GenerationConfig: &GeminiGenerationConfig{MaxOutputTokens: &c.maxTokens},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
		return "", fmt.Errorf("gemini API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in gemini response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// GenerateContent - Tạo nội dung theo prompt
func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error) {
	response, err := c.sendRequest(ctx, prompt, maxTokens, temperature)
	if err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// GenerateContentWithFile - Gửi prompt kèm nội dung 1 file (multimodal) cho Gemini.
// fileContent được base64-encode và gắn vào inlineData; mimeType ví dụ "application/pdf", "image/png".
func (c *GeminiClient) GenerateContentWithFile(ctx context.Context, prompt string, fileContent []byte, mimeType string, maxTokens int, temperature float32) (string, error) {
	if len(fileContent) == 0 {
		return c.GenerateContent(ctx, prompt, maxTokens, temperature)
	}

	parts := []GeminiPart{
		{Text: prompt},
		{InlineData: &GeminiInlineData{
			MimeType: mimeType,
			Data:     base64.StdEncoding.EncodeToString(fileContent),
		}},
	}

	response, err := c.sendRequestWithParts(ctx, parts, maxTokens, temperature)
	if err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// sendRequest - Helper để gửi request text-only đến Gemini API
func (c *GeminiClient) sendRequest(ctx context.Context, userMessage string, maxTokens int, temperature float32) (*GeminiResponse, error) {
	return c.sendRequestWithParts(ctx, []GeminiPart{{Text: userMessage}}, maxTokens, temperature)
}

// sendRequestWithParts - Helper gửi request với nhiều part (text + inlineData) đến Gemini API
func (c *GeminiClient) sendRequestWithParts(ctx context.Context, parts []GeminiPart, maxTokens int, temperature float32) (*GeminiResponse, error) {
	if err := validateProviderCall("gemini", c.providerMode, c.apiKey); err != nil {
		return nil, err
	}
	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}
	contents := []GeminiContent{
		{
			Parts: parts,
		},
	}

	requestBody := GeminiRequest{
		Contents: contents,
	}

	// max_tokens luôn có default từ process config; temperature là optional.
	{
		config := &GeminiGenerationConfig{MaxOutputTokens: &maxTokens}
		if temperature > 0 {
			config.Temperature = &temperature
		}
		requestBody.GenerationConfig = config
	}

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

// cleanJSONFromMarkdown - Loại bỏ markdown code block wrapper nếu có (helper function)
func cleanJSONFromMarkdownGemini(content string) string {
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
func (c *GeminiClient) GetModel() string {
	return c.model
}
