package provider

import "context"

// AssistantProvider gọi sang assistant-service để sinh nội dung bằng AI (Gemini),
// dùng cho job phân loại tài liệu quy hoạch.
type AssistantProvider interface {
	// GenerateContent gửi prompt tới assistant-service (GenerateContentWithGemini), trả về content thô + tên model.
	GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (content string, model string, err error)

	// ClassifyDocument gửi prompt kèm nội dung file (multimodal) tới ClassifyDocumentWithGemini.
	// fileContent rỗng => assistant fallback text-only.
	ClassifyDocument(ctx context.Context, prompt string, fileContent []byte, mimeType string) (content string, model string, err error)

	// Close releases the provider-owned outbound resource after the worker has stopped.
	Close() error
}
