package usecases

import (
	_dto "common/domain/dto"
	"context"
	"errors"
	"fmt"
	"time"

	"assistant/internal/dto"
	"assistant/internal/interface/provider"
)

var ErrProductSuggestUnavailable = errors.New("failed to get product suggest")

type AssistantUsecase struct {
	DeepseekClient provider.DeepseekProvider
	OpenAIClient   provider.OpenAIProvider
	GeminiClient   provider.GeminiProvider
	BdsproInternal provider.BdsproInternalProvider
	HubClient      provider.HubProvider
}

func NewAssistantUsecase(
	deepseekClient provider.DeepseekProvider,
	openaiClient provider.OpenAIProvider,
	geminiClient provider.GeminiProvider,
	bdsproInternal provider.BdsproInternalProvider,
	hubClient provider.HubProvider,
) *AssistantUsecase {
	return &AssistantUsecase{
		DeepseekClient: deepseekClient,
		OpenAIClient:   openaiClient,
		GeminiClient:   geminiClient,
		BdsproInternal: bdsproInternal,
		HubClient:      hubClient,
	}
}

// AnalyzeProductText - Phân tích văn bản sản phẩm BĐS
func (u *AssistantUsecase) AnalyzeProductText(ctx context.Context, content string, contextStr string) (*_dto.ProductV3DTO, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào content nếu có
	if contextStr != "" {
		content = fmt.Sprintf("%s\n\nContext bổ sung: %s", content, contextStr)
	}

	// Gọi Gemini để phân tích
	productV3, err := u.GeminiClient.AnalyzeProductText(ctx, content)
	if err != nil {
		return nil, "", 0, err
	}
	// productV3 := &_dto.ProductV3DTO{}

	processingTime := time.Since(startTime).Milliseconds()

	// Raw response để debug (optional)
	modelName := u.GeminiClient.GetModel() // Lấy từ config

	addressV3, err := u.HubClient.InferAddressFromText(ctx, content)
	if err != nil {
		return nil, "", 0, err
	}

	productV3_bdspro, err := u.BdsproInternal.GetSuggest(ctx, content)
	if err != nil {
		return nil, "", 0, ErrProductSuggestUnavailable
	}

	productV3.Address = addressV3
	productV3.PropertyType = productV3_bdspro.PropertyType
	productV3.DocType = productV3_bdspro.DocType
	productV3.Amenities = productV3_bdspro.Amenities

	return productV3, modelName, processingTime, nil
}

// Chat - Thực hiện conversation với AI
func (u *AssistantUsecase) Chat(ctx context.Context, message string, history []dto.DeepseekMessage, contextStr string, sessionID string) (string, string, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào message nếu có
	if contextStr != "" {
		message = fmt.Sprintf("%s\n\nContext: %s", message, contextStr)
	}

	// Gọi DeepSeek để chat
	reply, err := u.DeepseekClient.Chat(ctx, message, history)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to chat: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Generate hoặc use existing session ID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	modelName := u.DeepseekClient.GetModel() // Lấy từ config

	return reply, sessionID, modelName, processingTime, nil
}

// GenerateContent - Tạo nội dung theo template
func (u *AssistantUsecase) GenerateContent(ctx context.Context, prompt string, template string, variables map[string]string, maxTokens int, temperature float32) (string, string, int64, error) {
	startTime := time.Now()

	// Build prompt từ template và variables
	finalPrompt := prompt
	if template != "" {
		finalPrompt = u.buildPromptFromTemplate(template, variables, prompt)
	}

	// Gọi DeepSeek để generate
	content, err := u.DeepseekClient.GenerateContent(ctx, finalPrompt, maxTokens, temperature)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate content: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()
	modelName := u.DeepseekClient.GetModel() // Lấy từ config

	return content, modelName, processingTime, nil
}

// AnalyzeProductTextWithOpenAI - Phân tích văn bản sản phẩm BĐS bằng OpenAI
func (u *AssistantUsecase) AnalyzeProductTextWithOpenAI(ctx context.Context, content string, contextStr string) (*dto.ProductAnalysisResult, string, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào content nếu có
	if contextStr != "" {
		content = fmt.Sprintf("%s\n\nContext bổ sung: %s", content, contextStr)
	}

	// Gọi OpenAI để phân tích
	result, err := u.OpenAIClient.AnalyzeProductText(ctx, content)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("failed to analyze product text with openai: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Raw response để debug (optional)
	rawResponse := ""                      // Có thể marshal result về JSON nếu cần
	modelName := u.OpenAIClient.GetModel() // Lấy từ config

	return result, rawResponse, modelName, processingTime, nil
}

// ChatWithOpenAI - Thực hiện conversation với OpenAI
func (u *AssistantUsecase) ChatWithOpenAI(ctx context.Context, message string, history []dto.DeepseekMessage, contextStr string, sessionID string) (string, string, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào message nếu có
	if contextStr != "" {
		message = fmt.Sprintf("%s\n\nContext: %s", message, contextStr)
	}

	// Gọi OpenAI để chat
	reply, err := u.OpenAIClient.Chat(ctx, message, history)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to chat with openai: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Generate hoặc use existing session ID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	modelName := u.OpenAIClient.GetModel() // Lấy từ config

	return reply, sessionID, modelName, processingTime, nil
}

// GenerateContentWithOpenAI - Tạo nội dung theo template bằng OpenAI
func (u *AssistantUsecase) GenerateContentWithOpenAI(ctx context.Context, prompt string, template string, variables map[string]string, maxTokens int, temperature float32) (string, string, int64, error) {
	startTime := time.Now()

	// Build prompt từ template và variables
	finalPrompt := prompt
	if template != "" {
		finalPrompt = u.buildPromptFromTemplate(template, variables, prompt)
	}

	// Gọi OpenAI để generate
	content, err := u.OpenAIClient.GenerateContent(ctx, finalPrompt, maxTokens, temperature)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate content with openai: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()
	modelName := u.OpenAIClient.GetModel() // Lấy từ config

	return content, modelName, processingTime, nil
}

// buildPromptFromTemplate - Build prompt từ template và variables
func (u *AssistantUsecase) buildPromptFromTemplate(template string, variables map[string]string, basePrompt string) string {
	switch template {
	case "product_description":
		return u.buildProductDescriptionPrompt(variables, basePrompt)
	case "email":
		return u.buildEmailPrompt(variables, basePrompt)
	case "social_post":
		return u.buildSocialPostPrompt(variables, basePrompt)
	default:
		return basePrompt
	}
}

// buildProductDescriptionPrompt - Build prompt cho mô tả sản phẩm
func (u *AssistantUsecase) buildProductDescriptionPrompt(variables map[string]string, basePrompt string) string {
	prompt := "Viết mô tả chi tiết và hấp dẫn cho sản phẩm bất động sản sau:\n\n"

	if name, ok := variables["name"]; ok {
		prompt += fmt.Sprintf("Tên: %s\n", name)
	}
	if area, ok := variables["area"]; ok {
		prompt += fmt.Sprintf("Diện tích: %s m²\n", area)
	}
	if location, ok := variables["location"]; ok {
		prompt += fmt.Sprintf("Vị trí: %s\n", location)
	}
	if price, ok := variables["price"]; ok {
		prompt += fmt.Sprintf("Giá: %s\n", price)
	}

	if basePrompt != "" {
		prompt += fmt.Sprintf("\nYêu cầu thêm: %s", basePrompt)
	}

	return prompt
}

// buildEmailPrompt - Build prompt cho email
func (u *AssistantUsecase) buildEmailPrompt(variables map[string]string, basePrompt string) string {
	prompt := "Viết email chuyên nghiệp:\n\n"

	if to, ok := variables["to"]; ok {
		prompt += fmt.Sprintf("Người nhận: %s\n", to)
	}
	if subject, ok := variables["subject"]; ok {
		prompt += fmt.Sprintf("Chủ đề: %s\n", subject)
	}

	if basePrompt != "" {
		prompt += fmt.Sprintf("\nNội dung: %s", basePrompt)
	}

	return prompt
}

// buildSocialPostPrompt - Build prompt cho bài đăng mạng xã hội
func (u *AssistantUsecase) buildSocialPostPrompt(variables map[string]string, basePrompt string) string {
	prompt := "Viết bài đăng thu hút cho mạng xã hội:\n\n"

	if platform, ok := variables["platform"]; ok {
		prompt += fmt.Sprintf("Nền tảng: %s\n", platform)
	}
	if tone, ok := variables["tone"]; ok {
		prompt += fmt.Sprintf("Giọng điệu: %s\n", tone)
	}

	if basePrompt != "" {
		prompt += fmt.Sprintf("\nChủ đề: %s", basePrompt)
	}

	return prompt
}

// AnalyzeProductTextWithGemini - Phân tích văn bản sản phẩm BĐS bằng Gemini
func (u *AssistantUsecase) AnalyzeProductTextWithGemini(ctx context.Context, content string, contextStr string) (*_dto.ProductV3DTO, string, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào content nếu có
	if contextStr != "" {
		content = fmt.Sprintf("%s\n\nContext bổ sung: %s", content, contextStr)
	}

	// Gọi Gemini để phân tích
	result, err := u.GeminiClient.AnalyzeProductText(ctx, content)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("failed to analyze product text with gemini: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Raw response để debug (optional)
	rawResponse := ""                      // Có thể marshal result về JSON nếu cần
	modelName := u.GeminiClient.GetModel() // Lấy từ config

	return result, rawResponse, modelName, processingTime, nil
}

// ChatWithGemini - Thực hiện conversation với Gemini
func (u *AssistantUsecase) ChatWithGemini(ctx context.Context, message string, history []dto.DeepseekMessage, contextStr string, sessionID string) (string, string, string, int64, error) {
	startTime := time.Now()

	// Thêm context vào message nếu có
	if contextStr != "" {
		message = fmt.Sprintf("%s\n\nContext: %s", message, contextStr)
	}

	// Gọi Gemini để chat
	reply, err := u.GeminiClient.Chat(ctx, message, history)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("failed to chat with gemini: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Generate hoặc use existing session ID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	modelName := u.GeminiClient.GetModel() // Lấy từ config

	return reply, sessionID, modelName, processingTime, nil
}

// ClassifyDocumentWithGemini - Phân loại/đọc nội dung 1 tài liệu bằng Gemini (multimodal).
// Nếu fileContent rỗng sẽ tự fallback về text-only (chỉ dựa prompt).
func (u *AssistantUsecase) ClassifyDocumentWithGemini(ctx context.Context, prompt string, fileContent []byte, mimeType string, maxTokens int, temperature float32) (string, string, int64, error) {
	startTime := time.Now()

	content, err := u.GeminiClient.GenerateContentWithFile(ctx, prompt, fileContent, mimeType, maxTokens, temperature)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to classify document with gemini: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()
	modelName := u.GeminiClient.GetModel()

	return content, modelName, processingTime, nil
}

// GenerateContentWithGemini - Tạo nội dung theo template bằng Gemini
func (u *AssistantUsecase) GenerateContentWithGemini(ctx context.Context, prompt string, template string, variables map[string]string, maxTokens int, temperature float32) (string, string, int64, error) {
	startTime := time.Now()

	// Build prompt từ template và variables
	finalPrompt := prompt
	if template != "" {
		finalPrompt = u.buildPromptFromTemplate(template, variables, prompt)
	}

	// Gọi Gemini để generate
	content, err := u.GeminiClient.GenerateContent(ctx, finalPrompt, maxTokens, temperature)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate content with gemini: %w", err)
	}

	processingTime := time.Since(startTime).Milliseconds()
	modelName := u.GeminiClient.GetModel() // Lấy từ config

	return content, modelName, processingTime, nil
}
