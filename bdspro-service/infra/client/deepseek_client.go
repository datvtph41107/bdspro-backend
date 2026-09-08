package client

import (
	"bdspro/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// @bind: bdspro/internal/provider.DeepseekProvider
type DeepseekClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewDeepseekClient() *DeepseekClient {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	return &DeepseekClient{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SuggestProductInfo - Gọi DeepSeek để phân tích thông tin sản phẩm từ văn bản
func (c *DeepseekClient) SuggestProductInfo(ctx context.Context, content string) (*dto.DeepseekProductSuggestion, error) {
	// Tạo prompt cho DeepSeek
	systemPrompt := `Bạn là trợ lý AI chuyên phân tích thông tin bất động sản tiếng Việt. 
Nhiệm vụ của bạn là trích xuất thông tin sản phẩm bất động sản từ văn bản người dùng cung cấp.

Hãy phân tích và trả về JSON với cấu trúc sau:
{
  "name": "Tên/tiêu đề sản phẩm",
  "area": diện tích (số, đơn vị m2),
  "address": "Địa chỉ chi tiết",
  "province": "Tỉnh/Thành phố",
  "district": "Quận/Huyện",
  "ward": "Phường/Xã",
  "propertyType": "Loại BĐS (nhà riêng, chung cư, đất nền, ...)",
  "transactionType": 1 (1=Bán, 2=Cho thuê, 3=Cả hai),
  "salePrice": giá bán (số, VND),
  "rentPrice": giá thuê (số, VND/tháng),
  "bedroom": số phòng ngủ (số),
  "bathroom": số phòng tắm (số),
  "floor": số tầng (số),
  "frontage": mặt tiền (số, đơn vị m),
  "orientation": "Hướng nhà (Đông, Tây, Nam, Bắc, ...)",
  "legalDoc": "Giấy tờ pháp lý (sổ đỏ/hồng, sổ chung, ...)",
  "furniture": "Nội thất (đầy đủ, cơ bản, không có)",
  "description": "Mô tả tóm tắt"
}

Lưu ý:
- Chỉ trả về các trường có thông tin trong văn bản
- Trường nào không có thông tin thì không trả về
- Giá trị null cho các trường không xác định
- Chỉ trả về JSON thuần, không có markdown hay text khác`

	requestBody := dto.DeepseekRequest{
		Model: "deepseek-chat",
		Messages: []dto.DeepseekMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Tạo HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deepseek API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// Parse response
	var deepseekResp dto.DeepseekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(deepseekResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in deepseek response")
	}

	// Parse JSON từ message content
	content = deepseekResp.Choices[0].Message.Content

	// Clean markdown code block nếu có
	content = cleanJSONFromMarkdown(content)

	var suggestion dto.DeepseekProductSuggestion
	if err := json.Unmarshal([]byte(content), &suggestion); err != nil {
		return nil, fmt.Errorf("failed to parse suggestion JSON: %w, content: %s", err, content)
	}

	return &suggestion, nil
}

// cleanJSONFromMarkdown - Loại bỏ markdown code block wrapper nếu có
func cleanJSONFromMarkdown(content string) string {
	// Remove ```json và ``` wrapper nếu có
	if len(content) > 7 && content[:7] == "```json" {
		content = content[7:]
	} else if len(content) > 3 && content[:3] == "```" {
		content = content[3:]
	}

	if len(content) > 3 && content[len(content)-3:] == "```" {
		content = content[:len(content)-3]
	}

	// Trim whitespace
	return string(bytes.TrimSpace([]byte(content)))
}
