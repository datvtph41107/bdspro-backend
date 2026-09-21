package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type IntentUsecase struct {
	IntentRepo repo.IntentRepo
}

func NewIntentUsecase(intentRepo repo.IntentRepo) *IntentUsecase {
	return &IntentUsecase{
		IntentRepo: intentRepo,
	}
}

// DetectIntent - Phát hiện intent từ user input
func (u *IntentUsecase) DetectIntent(ctx context.Context, req *dto.IntentDetectionRequest) (*dto.IntentDetectionResponse, error) {
	userInput := strings.TrimSpace(strings.ToLower(req.UserInput))
	if userInput == "" {
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("User input không được rỗng"))
	}

	// Detect intent bằng keyword matching (simple version)
	// Trong thực tế có thể dùng AI/ML model
	detectedIntent, confidence := u.detectIntentByKeywords(userInput)

	// Extract entities từ input
	entities := u.extractEntities(userInput, detectedIntent)

	// Lấy metadata của intent
	metadata, exists := enums.GetIntentByAction(detectedIntent)
	if !exists {
		metadata = enums.IntentMetadata{
			Category:    enums.IntentCategoryGeneral,
			Action:      enums.IntentActionUnknown,
			Description: "Intent không xác định",
		}
	}

	// Lưu vào database
	profileID := _utils.GetProfileIdWithContext(ctx)
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	entitiesJSON, _ := json.Marshal(entities)
	contextJSON, _ := json.Marshal(req.Context)

	intent := &domain.Intent{
		UserInput:      req.UserInput,
		DetectedIntent: detectedIntent,
		Category:       metadata.Category,
		Confidence:     confidence,
		Entities:       string(entitiesJSON),
		Context:        string(contextJSON),
		ProfileID:      &profileID,
		OrganizationID: &organizationID,
		SessionID:      req.SessionID,
		ProcessedBy:    stringPtr("ai"),
	}

	now := time.Now()
	intent.ProcessedAt = &now

	if err := u.IntentRepo.Create(ctx, intent); err != nil {
		// Log error nhưng không fail request
		fmt.Printf("Warning: Failed to save intent: %v\n", err)
	}

	// Tạo suggestions dựa trên intent
	suggestions := u.generateSuggestions(detectedIntent, entities)

	// Lấy required params
	requiredParams := u.getRequiredParams(detectedIntent)

	return &dto.IntentDetectionResponse{
		DetectedIntent: detectedIntent,
		Category:       metadata.Category,
		Confidence:     confidence,
		Entities:       entities,
		Suggestions:    suggestions,
		RequiredParams: requiredParams,
		Metadata:       metadata,
	}, nil
}

// ExecuteIntent - Thực thi intent
func (u *IntentUsecase) ExecuteIntent(ctx context.Context, req *dto.IntentExecutionRequest) (*dto.IntentExecutionResponse, error) {
	// Validate intent
	if !enums.IsValidIntent(req.Intent) {
		return &dto.IntentExecutionResponse{
			Success: false,
			Message: "Intent không hợp lệ",
			Error:   stringPtr("Invalid intent"),
		}, nil
	}

	// TODO: Implement actual execution logic based on intent
	// Đây là nơi route đến các usecase khác tương ứng
	switch req.Intent {
	case enums.IntentActionProductSearch:
		return u.executeProductSearch(ctx, req.Parameters)
	case enums.IntentActionProductDetail:
		return u.executeProductDetail(ctx, req.Parameters)
	case enums.IntentActionProductCreate:
		return u.executeProductCreate(ctx, req.Parameters)
	case enums.IntentActionHelp:
		return u.executeHelp(ctx)
	default:
		return &dto.IntentExecutionResponse{
			Success: false,
			Message: "Intent chưa được implement",
			Error:   stringPtr("Not implemented"),
		}, nil
	}
}

// GetIntentHistory - Lấy lịch sử intents
func (u *IntentUsecase) GetIntentHistory(ctx context.Context, req *dto.IntentSearchRequest) ([]domain.Intent, int64, error) {
	// Set default organization ID nếu không có
	if req.OrganizationID == nil {
		orgID := _utils.GetOrganizationIdFromContext(ctx)
		req.OrganizationID = &orgID
	}

	intents, total, err := u.IntentRepo.GetList(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return intents, total, nil
}

// GetStatistics - Lấy thống kê intents
func (u *IntentUsecase) GetStatistics(ctx context.Context, req *dto.IntentStatisticsRequest) (*dto.IntentStatisticsResponse, error) {
	statistics, err := u.IntentRepo.GetStatistics(ctx, req)
	if err != nil {
		return nil, err
	}

	// Tính total
	var total int64
	for _, stat := range statistics {
		total += stat.TotalCount
	}

	return &dto.IntentStatisticsResponse{
		Statistics: statistics,
		Total:      total,
		Period:     fmt.Sprintf("%s to %s", *req.FromDate, *req.ToDate),
	}, nil
}

// GetSuggestions - Lấy gợi ý intents
func (u *IntentUsecase) GetSuggestions(ctx context.Context, req *dto.IntentSuggestionRequest) (*dto.IntentSuggestionResponse, error) {
	partialInput := strings.TrimSpace(strings.ToLower(req.PartialInput))
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}

	var suggestions []dto.IntentSuggestion

	// Search through intent registry
	for action, metadata := range enums.IntentRegistry {
		// Filter by category if specified
		if req.Category != nil && metadata.Category != *req.Category {
			continue
		}

		// Calculate score based on keyword matching
		score := u.calculateMatchScore(partialInput, metadata)
		if score > 0 {
			example := ""
			if len(metadata.Examples) > 0 {
				example = metadata.Examples[0]
			}

			suggestions = append(suggestions, dto.IntentSuggestion{
				Action:      action,
				Category:    metadata.Category,
				Description: metadata.Description,
				Example:     example,
				Score:       score,
			})
		}
	}

	// Sort by score (descending)
	// TODO: Implement proper sorting

	// Limit results
	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}

	return &dto.IntentSuggestionResponse{
		Suggestions: suggestions,
	}, nil
}

// Private helper methods

func (u *IntentUsecase) detectIntentByKeywords(input string) (enums.EIntentAction, float32) {
	input = strings.ToLower(input)
	bestMatch := enums.IntentActionUnknown
	bestScore := float32(0.0)

	for action, metadata := range enums.IntentRegistry {
		score := float32(0.0)
		matchCount := 0

		// Check keywords
		for _, keyword := range metadata.Keywords {
			if strings.Contains(input, strings.ToLower(keyword)) {
				matchCount++
			}
		}

		if matchCount > 0 {
			score = float32(matchCount) / float32(len(metadata.Keywords))
		}

		if score > bestScore {
			bestScore = score
			bestMatch = action
		}
	}

	// Normalize confidence to 0-1 range
	confidence := bestScore
	if confidence > 1.0 {
		confidence = 1.0
	}

	return bestMatch, confidence
}

func (u *IntentUsecase) extractEntities(input string, intent enums.EIntentAction) []domain.ExtractedEntity {
	entities := []domain.ExtractedEntity{}

	// Simple entity extraction (có thể improve bằng NLP)
	input = strings.ToLower(input)

	// Extract product ID pattern (SP001, SP002, etc)
	if strings.Contains(input, "sp") {
		// TODO: Extract product ID
	}

	// Extract price
	if strings.Contains(input, "tỷ") || strings.Contains(input, "triệu") {
		// TODO: Extract price
	}

	// Extract area
	if strings.Contains(input, "m2") || strings.Contains(input, "m²") {
		// TODO: Extract area
	}

	return entities
}

func (u *IntentUsecase) generateSuggestions(intent enums.EIntentAction, entities []domain.ExtractedEntity) []string {
	suggestions := []string{}

	switch intent {
	case enums.IntentActionProductSearch:
		suggestions = append(suggestions, "Tìm kiếm sản phẩm theo khu vực", "Lọc theo giá", "Lọc theo diện tích")
	case enums.IntentActionProductDetail:
		suggestions = append(suggestions, "Xem thông tin chi tiết", "Xem lịch sử giao dịch", "Chia sẻ sản phẩm")
	case enums.IntentActionProductCreate:
		suggestions = append(suggestions, "Nhập thông tin cơ bản", "Thêm hình ảnh", "Đăng tin")
	}

	return suggestions
}

func (u *IntentUsecase) getRequiredParams(intent enums.EIntentAction) []string {
	switch intent {
	case enums.IntentActionProductDetail:
		return []string{"productId"}
	case enums.IntentActionProductCreate:
		return []string{"name", "area", "price", "address"}
	case enums.IntentActionProductSearch:
		return []string{}
	default:
		return []string{}
	}
}

func (u *IntentUsecase) calculateMatchScore(input string, metadata enums.IntentMetadata) float32 {
	if input == "" {
		return 0
	}

	matchCount := 0
	for _, keyword := range metadata.Keywords {
		if strings.Contains(input, strings.ToLower(keyword)) {
			matchCount++
		}
	}

	if matchCount == 0 {
		return 0
	}

	return float32(matchCount) / float32(len(metadata.Keywords))
}

// Execute methods (stubs - implement actual logic)

func (u *IntentUsecase) executeProductSearch(ctx context.Context, params map[string]interface{}) (*dto.IntentExecutionResponse, error) {
	return &dto.IntentExecutionResponse{
		Success: true,
		Message: "Tìm kiếm sản phẩm thành công",
		Data:    params,
	}, nil
}

func (u *IntentUsecase) executeProductDetail(ctx context.Context, params map[string]interface{}) (*dto.IntentExecutionResponse, error) {
	return &dto.IntentExecutionResponse{
		Success: true,
		Message: "Lấy chi tiết sản phẩm thành công",
		Data:    params,
	}, nil
}

func (u *IntentUsecase) executeProductCreate(ctx context.Context, params map[string]interface{}) (*dto.IntentExecutionResponse, error) {
	return &dto.IntentExecutionResponse{
		Success: true,
		Message: "Tạo sản phẩm thành công",
		Data:    params,
	}, nil
}

func (u *IntentUsecase) executeHelp(ctx context.Context) (*dto.IntentExecutionResponse, error) {
	helpMessage := `Tôi có thể giúp bạn:
- Tìm kiếm sản phẩm: "tìm nhà ở Hà Nội"
- Xem chi tiết: "cho tôi xem sản phẩm SP001"
- Tạo sản phẩm: "tạo sản phẩm mới"
- Đăng tin: "đăng bán nhà"
`
	return &dto.IntentExecutionResponse{
		Success: true,
		Message: helpMessage,
	}, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
