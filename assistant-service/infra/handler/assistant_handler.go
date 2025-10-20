package handler

import (
	"assistant/internal/dto"
	"assistant/internal/usecases"
	"context"
	"time"

	assistantpb "pb/types/assistant"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssistantHandler struct {
	assistantpb.UnimplementedAssistantServiceServer
	Usecase   *usecases.AssistantUsecase
	StartTime time.Time
}

func NewAssistantHandler(usecase *usecases.AssistantUsecase) *AssistantHandler {
	return &AssistantHandler{
		Usecase:   usecase,
		StartTime: time.Now(),
	}
}

// AnalyzeProductText - Phân tích văn bản sản phẩm BĐS
// @Summary Phân tích văn bản sản phẩm BĐS sử dụng AI
// @Description Sử dụng DeepSeek AI để phân tích và trích xuất thông tin từ văn bản mô tả sản phẩm
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.AnalyzeRequest true "Nội dung cần phân tích"
// @Success 200 {object} assistantpb.AnalyzeResponse
// @Router /v2/assistant/analyze/product [post]
func (h *AssistantHandler) AnalyzeProductText(ctx context.Context, req *assistantpb.AnalyzeRequest) (*assistantpb.AnalyzeResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích
	result, rawResponse, processingTime, err := h.Usecase.AnalyzeProductText(ctx, req.Content, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze: %v", err)
	}

	// Map kết quả sang protobuf
	productInfo := &assistantpb.ProductInfo{}

	// Basic fields
	if result.Name != nil {
		productInfo.Name = result.Name
	}
	if result.Area != nil {
		areaFloat32 := float32(*result.Area)
		productInfo.Area = &areaFloat32
	}
	if result.Description != nil {
		productInfo.Description = result.Description
	}
	if result.PropertyType != nil {
		productInfo.PropertyType = result.PropertyType
	}
	if result.TransactionType != nil {
		productInfo.TransactionType = result.TransactionType
	}
	if result.LegalDoc != nil {
		productInfo.LegalDoc = result.LegalDoc
	}

	// Address - map từ nested structure sang flat
	if result.Address != nil {
		if result.Address.Detail != nil {
			productInfo.Address = result.Address.Detail
		}
		if result.Address.Province != nil {
			productInfo.Province = result.Address.Province
		}
		if result.Address.District != nil {
			productInfo.District = result.Address.District
		}
		if result.Address.Ward != nil {
			productInfo.Ward = result.Address.Ward
		}
	}

	// PriceData - map từ nested structure
	if result.PriceData != nil {
		if result.PriceData.SalePrice != nil {
			productInfo.SalePrice = result.PriceData.SalePrice
		}
		if result.PriceData.RentPrice != nil {
			productInfo.RentPrice = result.PriceData.RentPrice
		}
	}

	// HouseInfo - map từ nested structure
	if result.HouseInfo != nil {
		if result.HouseInfo.NumBedroom != nil {
			productInfo.Bedroom = result.HouseInfo.NumBedroom
		}
		if result.HouseInfo.NumBathroom != nil {
			productInfo.Bathroom = result.HouseInfo.NumBathroom
		}
		if result.HouseInfo.NumFloor != nil {
			productInfo.Floor = result.HouseInfo.NumFloor
		}
		if result.HouseInfo.NumFront != nil {
			numFrontFloat32 := float32(*result.HouseInfo.NumFront)
			productInfo.Frontage = &numFrontFloat32
		}
		if result.HouseInfo.Orientation != nil {
			productInfo.Orientation = result.HouseInfo.Orientation
		}
		if result.HouseInfo.Furniture != nil {
			productInfo.Furniture = result.HouseInfo.Furniture
		}
	}

	// Amenities
	if len(result.Amenities) > 0 {
		productInfo.Amenities = result.Amenities
	}

	return &assistantpb.AnalyzeResponse{
		Product:     productInfo,
		RawResponse: rawResponse,
		Metadata: &assistantpb.AnalysisMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from DeepSeek response
			ConfidenceScore:  0.0,
			Model:            "deepseek-chat",
		},
	}, nil
}

// Chat - Chat với AI assistant
// @Summary Chat với AI assistant
// @Description Thực hiện conversation với AI assistant
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.ChatRequest true "Message và history"
// @Success 200 {object} assistantpb.ChatResponse
// @Router /v2/assistant/chat [post]
func (h *AssistantHandler) Chat(ctx context.Context, req *assistantpb.ChatRequest) (*assistantpb.ChatResponse, error) {
	if req.Message == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message is required")
	}

	// Convert history từ protobuf sang DTO
	history := make([]dto.DeepseekMessage, len(req.History))
	for i, msg := range req.History {
		history[i] = dto.DeepseekMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Gọi usecase để chat
	reply, sessionID, processingTime, err := h.Usecase.Chat(ctx, req.Message, history, req.Context, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to chat: %v", err)
	}

	return &assistantpb.ChatResponse{
		Reply:     reply,
		SessionId: sessionID,
		Metadata: &assistantpb.ChatMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from DeepSeek response
			Model:            "deepseek-chat",
		},
	}, nil
}

// GenerateContent - Generate nội dung
// @Summary Generate nội dung
// @Description Tạo nội dung theo template và prompt
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.GenerateRequest true "Prompt và template"
// @Success 200 {object} assistantpb.GenerateResponse
// @Router /v2/assistant/generate [post]
func (h *AssistantHandler) GenerateContent(ctx context.Context, req *assistantpb.GenerateRequest) (*assistantpb.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, status.Errorf(codes.InvalidArgument, "prompt is required")
	}

	maxTokens := int(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 2000
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	// Gọi usecase để generate
	content, processingTime, err := h.Usecase.GenerateContent(
		ctx,
		req.Prompt,
		req.Template,
		req.Variables,
		maxTokens,
		temperature,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate: %v", err)
	}

	return &assistantpb.GenerateResponse{
		Content: content,
		Metadata: &assistantpb.GenerateMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from DeepSeek response
			Model:            "deepseek-chat",
		},
	}, nil
}

// Health - Health check
// @Summary Health check
// @Description Kiểm tra trạng thái service
// @Tags Assistant
// @Accept json
// @Produce json
// @Success 200 {object} assistantpb.HealthResponse
// @Router /v2/assistant/health [get]
func (h *AssistantHandler) Health(ctx context.Context, req *assistantpb.HealthRequest) (*assistantpb.HealthResponse, error) {
	uptime := int64(time.Since(h.StartTime).Seconds())

	checks := map[string]string{
		"deepseek_client": "healthy",
		"service":         "healthy",
	}

	return &assistantpb.HealthResponse{
		Status:        "healthy",
		UptimeSeconds: uptime,
		Version:       "1.0.0",
		Checks:        checks,
	}, nil
}
