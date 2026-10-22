package dto

// AnalyzeProductRequest is POST /v4/assistant/bdspro/product (parity assistantpb.AnalyzeRequest).
type AnalyzeProductRequest struct {
	Content string `json:"content"`
	Context string `json:"context"`
}
