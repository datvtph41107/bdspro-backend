package dto

// DeepseekRequest - DTO để gọi API DeepSeek
type DeepseekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepseekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

// DeepseekMessage - Message trong request gửi đến DeepSeek
type DeepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepseekResponse - Response từ DeepSeek API
type DeepseekResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []DeepseekChoice `json:"choices"`
	Usage   DeepseekUsage    `json:"usage"`
}

// DeepseekChoice - Choice trong response
type DeepseekChoice struct {
	Index        int             `json:"index"`
	Message      DeepseekMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

// DeepseekUsage - Thông tin về usage
type DeepseekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// DeepseekProductSuggestion - Kết quả parse từ DeepSeek response
type DeepseekProductSuggestion struct {
	Name            *string  `json:"name"`
	Area            *float32 `json:"area"`
	Address         *string  `json:"address"`
	Province        *string  `json:"province"`
	District        *string  `json:"district"`
	Ward            *string  `json:"ward"`
	PropertyType    *string  `json:"propertyType"`
	TransactionType *int32   `json:"transactionType"`
	SalePrice       *float64 `json:"salePrice"`
	RentPrice       *float64 `json:"rentPrice"`
	Bedroom         *int32   `json:"bedroom"`
	Bathroom        *int32   `json:"bathroom"`
	Floor           *int32   `json:"floor"`
	Frontage        *float32 `json:"frontage"`
	Orientation     *string  `json:"orientation"`
	LegalDoc        *string  `json:"legalDoc"`
	Furniture       *string  `json:"furniture"`
	Description     *string  `json:"description"`
}
