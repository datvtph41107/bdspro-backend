package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
)

// IntentDetectionRequest - Request để detect intent
type IntentDetectionRequest struct {
	UserInput string            `json:"userInput" binding:"required"`
	Context   map[string]string `json:"context"`
	SessionID *string           `json:"sessionId"`
}

// IntentDetectionResponse - Response sau khi detect intent
type IntentDetectionResponse struct {
	DetectedIntent enums.EIntentAction      `json:"detectedIntent"`
	Category       enums.EIntentCategory    `json:"category"`
	Confidence     float32                  `json:"confidence"`
	Entities       []domain.ExtractedEntity `json:"entities"`
	Suggestions    []string                 `json:"suggestions"`
	RequiredParams []string                 `json:"requiredParams"`
	Metadata       enums.IntentMetadata     `json:"metadata"`
}

// IntentExecutionRequest - Request để execute intent
type IntentExecutionRequest struct {
	Intent     enums.EIntentAction    `json:"intent" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
	Context    map[string]string      `json:"context"`
}

// IntentExecutionResponse - Response sau khi execute intent
type IntentExecutionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *string     `json:"error"`
}

// IntentSearchRequest - Request để search intent history
type IntentSearchRequest struct {
	_dto.Pagable
	Category       *enums.EIntentCategory `json:"category"`
	Action         *enums.EIntentAction   `json:"action"`
	ProfileID      *uint64                `json:"profileId"`
	OrganizationID *uint64                `json:"organizationId"`
	SessionID      *string                `json:"sessionId"`
	FromDate       *string                `json:"fromDate"`
	ToDate         *string                `json:"toDate"`
	MinConfidence  *float32               `json:"minConfidence"`
}

// IntentStatisticsRequest - Request để lấy thống kê intent
type IntentStatisticsRequest struct {
	Category       *enums.EIntentCategory `json:"category"`
	OrganizationID *uint64                `json:"organizationId"`
	FromDate       *string                `json:"fromDate"`
	ToDate         *string                `json:"toDate"`
	GroupBy        string                 `json:"groupBy"` // "category", "action", "day", "week", "month"
}

// IntentStatisticsResponse - Response thống kê intent
type IntentStatisticsResponse struct {
	Statistics []domain.IntentStatistics `json:"statistics"`
	Total      int64                     `json:"total"`
	Period     string                    `json:"period"`
}

// IntentTrainingRequest - Request để train/improve intent detection
type IntentTrainingRequest struct {
	UserInput      string               `json:"userInput" binding:"required"`
	CorrectIntent  enums.EIntentAction  `json:"correctIntent" binding:"required"`
	DetectedIntent *enums.EIntentAction `json:"detectedIntent"`
	Feedback       *string              `json:"feedback"`
}

// IntentSuggestionRequest - Request để lấy gợi ý intent
type IntentSuggestionRequest struct {
	PartialInput string                 `json:"partialInput"`
	Category     *enums.EIntentCategory `json:"category"`
	Limit        int                    `json:"limit"`
}

// IntentSuggestionResponse - Response gợi ý intent
type IntentSuggestionResponse struct {
	Suggestions []IntentSuggestion `json:"suggestions"`
}

// IntentSuggestion - Một gợi ý intent
type IntentSuggestion struct {
	Action      enums.EIntentAction   `json:"action"`
	Category    enums.EIntentCategory `json:"category"`
	Description string                `json:"description"`
	Example     string                `json:"example"`
	Score       float32               `json:"score"`
}
