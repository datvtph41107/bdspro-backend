package handler

import (
	"assistant/config"
	"assistant/internal/dto"
	"assistant/internal/usecases"
	"context"
	"time"

	assistantpb "pb/types/assistant"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssistantHandler struct {
	assistantpb.UnimplementedAssistantServiceServer
	Usecase      *usecases.AssistantUsecase
	StartTime    time.Time
	providerMode string
}

func NewAssistantHandler(usecase *usecases.AssistantUsecase, runtime config.Runtime) *AssistantHandler {
	return &AssistantHandler{
		Usecase:      usecase,
		StartTime:    time.Now(),
		providerMode: runtime.ProviderMode,
	}
}

// AnalyzeProductText - Phân tích văn bản sản phẩm BĐS
// @Summary Phân tích văn bản sản phẩm BĐS sử dụng AI
// @Description Sử dụng AI để phân tích và trích xuất thông tin từ văn bản mô tả sản phẩm
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.AnalyzeRequest true "Nội dung cần phân tích"
// @Success 200 {object} assistantpb.SuggestResponse
// @Router /v2/assistant/bdspro/product [post]
func (h *AssistantHandler) AnalyzeProductText(ctx context.Context, req *assistantpb.AnalyzeRequest) (*assistantpb.SuggestResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	productV3, modelName, processingTime, err := h.Usecase.AnalyzeProductText(ctx, req.Content, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze: %v", err)
	}

	productInfo := &sharepb.ProductV3Proto{
		Name:        productV3.Name,
		Description: productV3.Description,
		Area:        productV3.Area,
	}

	if productV3.SourceType != nil {
		productInfo.SourceType = &sharepb.ItemV3Proto{
			Id:   productV3.SourceType.ID,
			Name: productV3.SourceType.Name,
		}
	}

	if productV3.HouseInfo != nil {
		productInfo.HouseInfo = &sharepb.HouseInfoV3Proto{
			NumBedroom:  productV3.HouseInfo.NumBedroom,
			NumBathroom: productV3.HouseInfo.NumBathroom,
			NumFloor:    productV3.HouseInfo.NumFloor,
			NumFront:    productV3.HouseInfo.NumFront,
			NumCarPark:  productV3.HouseInfo.NumCarPark,
			Orientation: productV3.HouseInfo.Orientation,
			Furniture:   productV3.HouseInfo.Furniture,
		}
	}

	if productV3.PriceData != nil {
		productInfo.PriceData = &sharepb.PriceV3Proto{
			Currency:           productV3.PriceData.Currency,
			SalePrice:          productV3.PriceData.SalePrice,
			SaleCommission:     productV3.PriceData.SaleCommission,
			SaleCommissionType: productV3.PriceData.SaleCommissionType,
			RentPrice:          productV3.PriceData.RentPrice,
		}
	}

	if productV3.TransactionType != nil {
		productInfo.TransactionType = &sharepb.ItemV3Proto{
			Id:   productV3.TransactionType.ID,
			Name: productV3.TransactionType.Name,
		}
	}

	if productV3.Address != nil {
		productInfo.Address = &sharepb.AddressV3Proto{
			Detail:       productV3.Address.Detail,
			ProvinceId:   productV3.Address.ProvinceID,
			DistrictId:   productV3.Address.DistrictID,
			WardId:       productV3.Address.WardID,
			ProvinceName: productV3.Address.ProvinceName,
			DistrictName: productV3.Address.DistrictName,
			WardName:     productV3.Address.WardName,
		}
	}
	if productV3.PropertyType != nil {
		productInfo.PropertyType = &sharepb.ItemV3Proto{
			Id:   productV3.PropertyType.ID,
			Name: productV3.PropertyType.Name,
		}
	}
	if productV3.DocType != nil {
		productInfo.DocType = &sharepb.ItemV3Proto{
			Id:   productV3.DocType.ID,
			Name: productV3.DocType.Name,
		}
	}
	if productV3.Amenities != nil {
		productInfo.Amenities = make([]*sharepb.ItemV3Proto, len(productV3.Amenities))
		for i, amenity := range productV3.Amenities {
			productInfo.Amenities[i] = &sharepb.ItemV3Proto{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}
	return &assistantpb.SuggestResponse{
		ModelName:      modelName,
		ProcessingTime: processingTime,
		Product:        productInfo,
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
	reply, sessionID, modelName, processingTime, err := h.Usecase.Chat(ctx, req.Message, history, req.Context, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to chat: %v", err)
	}

	return &assistantpb.ChatResponse{
		Reply:     reply,
		SessionId: sessionID,
		Metadata: &assistantpb.ChatMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from API response
			Model:            modelName,
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
	content, modelName, processingTime, err := h.Usecase.GenerateContent(
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
			TokensUsed:       0, // TODO: Get from API response
			Model:            modelName,
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
	providerStatus := h.providerMode

	checks := map[string]string{
		"deepseek_client": providerStatus,
		"openai_client":   providerStatus,
		"gemini_client":   providerStatus,
		"service":         "healthy",
	}

	return &assistantpb.HealthResponse{
		Status:        "healthy",
		UptimeSeconds: uptime,
		Version:       "1.0.0",
		Checks:        checks,
	}, nil
}

// AnalyzeProductTextWithDeepSeek - Phân tích văn bản sản phẩm BĐS sử dụng DeepSeek
// @Summary Phân tích văn bản sản phẩm BĐS sử dụng DeepSeek AI
// @Description Sử dụng DeepSeek AI để phân tích và trích xuất thông tin từ văn bản mô tả sản phẩm
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.AnalyzeRequest true "Nội dung cần phân tích"
// @Success 200 {object} assistantpb.AnalyzeResponse
// @Router /v2/assistant/deepseek/analyze/product [post]
func (h *AssistantHandler) AnalyzeProductTextWithDeepSeek(ctx context.Context, req *assistantpb.AnalyzeRequest) (*assistantpb.AnalyzeResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích với DeepSeek
	result, modelName, processingTime, err := h.Usecase.AnalyzeProductText(ctx, req.Content, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze with deepseek: %v", err)
	}

	// Map kết quả sang protobuf
	productInfo := &assistantpb.ProductInfo{
		Name: &result.Name,
		// Area:            float32(result.Area),
		// Description:     &result.Description,
		// PropertyType:    result.propertyType,
		// TransactionType: int32(result.TransactionType),
		// LegalDoc:        &result.LegalDoc,
		// Address: result.Address.Detail,
	}

	return &assistantpb.AnalyzeResponse{
		Product: productInfo,
		Metadata: &assistantpb.AnalysisMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			ConfidenceScore:  0.0,
			Model:            modelName,
		},
	}, nil
}

// ChatWithDeepSeek - Chat với DeepSeek AI assistant
// @Summary Chat với DeepSeek AI assistant
// @Description Thực hiện conversation với DeepSeek AI assistant
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.ChatRequest true "Message và history"
// @Success 200 {object} assistantpb.ChatResponse
// @Router /v2/assistant/deepseek/chat [post]
func (h *AssistantHandler) ChatWithDeepSeek(ctx context.Context, req *assistantpb.ChatRequest) (*assistantpb.ChatResponse, error) {
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

	// Gọi usecase để chat với DeepSeek
	reply, sessionID, modelName, processingTime, err := h.Usecase.Chat(ctx, req.Message, history, req.Context, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to chat with deepseek: %v", err)
	}

	return &assistantpb.ChatResponse{
		Reply:     reply,
		SessionId: sessionID,
		Metadata: &assistantpb.ChatMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			Model:            modelName,
		},
	}, nil
}

// GenerateContentWithDeepSeek - Generate nội dung với DeepSeek
// @Summary Generate nội dung với DeepSeek
// @Description Tạo nội dung theo template và prompt sử dụng DeepSeek
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.GenerateRequest true "Prompt và template"
// @Success 200 {object} assistantpb.GenerateResponse
// @Router /v2/assistant/deepseek/generate [post]
func (h *AssistantHandler) GenerateContentWithDeepSeek(ctx context.Context, req *assistantpb.GenerateRequest) (*assistantpb.GenerateResponse, error) {
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

	// Gọi usecase để generate với DeepSeek
	content, modelName, processingTime, err := h.Usecase.GenerateContent(
		ctx,
		req.Prompt,
		req.Template,
		req.Variables,
		maxTokens,
		temperature,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate with deepseek: %v", err)
	}

	return &assistantpb.GenerateResponse{
		Content: content,
		Metadata: &assistantpb.GenerateMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			Model:            modelName,
		},
	}, nil
}

// AnalyzeProductTextWithOpenAI - Phân tích văn bản sản phẩm BĐS sử dụng OpenAI
// @Summary Phân tích văn bản sản phẩm BĐS sử dụng OpenAI
// @Description Sử dụng OpenAI GPT để phân tích và trích xuất thông tin từ văn bản mô tả sản phẩm
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.AnalyzeRequest true "Nội dung cần phân tích"
// @Success 200 {object} assistantpb.AnalyzeResponse
// @Router /v2/assistant/openai/analyze/product [post]
func (h *AssistantHandler) AnalyzeProductTextWithOpenAI(ctx context.Context, req *assistantpb.AnalyzeRequest) (*assistantpb.AnalyzeResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích với OpenAI
	result, rawResponse, modelName, processingTime, err := h.Usecase.AnalyzeProductTextWithOpenAI(ctx, req.Content, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze with openai: %v", err)
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

	// Address
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

	// PriceData
	if result.PriceData != nil {
		if result.PriceData.SalePrice != nil {
			productInfo.SalePrice = result.PriceData.SalePrice
		}
		if result.PriceData.RentPrice != nil {
			productInfo.RentPrice = result.PriceData.RentPrice
		}
	}

	// HouseInfo
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
			TokensUsed:       0,
			ConfidenceScore:  0.0,
			Model:            modelName,
		},
	}, nil
}

// ChatWithOpenAI - Chat với OpenAI AI assistant
// @Summary Chat với OpenAI AI assistant
// @Description Thực hiện conversation với OpenAI GPT assistant
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.ChatRequest true "Message và history"
// @Success 200 {object} assistantpb.ChatResponse
// @Router /v2/assistant/openai/chat [post]
func (h *AssistantHandler) ChatWithOpenAI(ctx context.Context, req *assistantpb.ChatRequest) (*assistantpb.ChatResponse, error) {
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

	// Gọi usecase để chat với OpenAI
	reply, sessionID, modelName, processingTime, err := h.Usecase.ChatWithOpenAI(ctx, req.Message, history, req.Context, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to chat with openai: %v", err)
	}

	return &assistantpb.ChatResponse{
		Reply:     reply,
		SessionId: sessionID,
		Metadata: &assistantpb.ChatMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			Model:            modelName,
		},
	}, nil
}

// GenerateContentWithOpenAI - Generate nội dung với OpenAI
// @Summary Generate nội dung với OpenAI
// @Description Tạo nội dung theo template và prompt sử dụng OpenAI GPT
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.GenerateRequest true "Prompt và template"
// @Success 200 {object} assistantpb.GenerateResponse
// @Router /v2/assistant/openai/generate [post]
func (h *AssistantHandler) GenerateContentWithOpenAI(ctx context.Context, req *assistantpb.GenerateRequest) (*assistantpb.GenerateResponse, error) {
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

	// Gọi usecase để generate với OpenAI
	content, modelName, processingTime, err := h.Usecase.GenerateContentWithOpenAI(
		ctx,
		req.Prompt,
		req.Template,
		req.Variables,
		maxTokens,
		temperature,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate with openai: %v", err)
	}

	return &assistantpb.GenerateResponse{
		Content: content,
		Metadata: &assistantpb.GenerateMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			Model:            modelName,
		},
	}, nil
}

// AnalyzeProductTextWithGemini - Phân tích văn bản sản phẩm BĐS sử dụng Gemini
// @Summary Phân tích văn bản sản phẩm BĐS sử dụng Gemini AI
// @Description Sử dụng Google Gemini AI để phân tích và trích xuất thông tin từ văn bản mô tả sản phẩm
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.AnalyzeRequest true "Nội dung cần phân tích"
// @Success 200 {object} assistantpb.AnalyzeResponse
// @Router /v2/assistant/gemini/analyze/product [post]
func (h *AssistantHandler) AnalyzeProductTextWithGemini(ctx context.Context, req *assistantpb.AnalyzeRequest) (*assistantpb.AnalyzeResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích với Gemini
	_, rawResponse, modelName, processingTime, err := h.Usecase.AnalyzeProductTextWithGemini(ctx, req.Content, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze with gemini: %v", err)
	}

	// Map kết quả sang protobuf (tương tự như AnalyzeProductText)
	productInfo := &assistantpb.ProductInfo{}

	return &assistantpb.AnalyzeResponse{
		Product:     productInfo,
		RawResponse: rawResponse,
		Metadata: &assistantpb.AnalysisMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from API response
			ConfidenceScore:  0.0,
			Model:            modelName,
		},
	}, nil
}

// ChatWithGemini - Chat với Gemini AI assistant
// @Summary Chat với Gemini AI assistant
// @Description Thực hiện conversation với Google Gemini AI assistant
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.ChatRequest true "Message và history"
// @Success 200 {object} assistantpb.ChatResponse
// @Router /v2/assistant/gemini/chat [post]
func (h *AssistantHandler) ChatWithGemini(ctx context.Context, req *assistantpb.ChatRequest) (*assistantpb.ChatResponse, error) {
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

	// Gọi usecase để chat với Gemini
	reply, sessionID, modelName, processingTime, err := h.Usecase.ChatWithGemini(ctx, req.Message, history, req.Context, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to chat with gemini: %v", err)
	}

	return &assistantpb.ChatResponse{
		Reply:     reply,
		SessionId: sessionID,
		Metadata: &assistantpb.ChatMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from API response
			Model:            modelName,
		},
	}, nil
}

// GenerateContentWithGemini - Generate nội dung với Gemini
// @Summary Generate nội dung với Gemini
// @Description Tạo nội dung theo template và prompt sử dụng Google Gemini
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.GenerateRequest true "Prompt và template"
// @Success 200 {object} assistantpb.GenerateResponse
// @Router /v2/assistant/gemini/generate [post]
func (h *AssistantHandler) GenerateContentWithGemini(ctx context.Context, req *assistantpb.GenerateRequest) (*assistantpb.GenerateResponse, error) {
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

	// Gọi usecase để generate với Gemini
	content, modelName, processingTime, err := h.Usecase.GenerateContentWithGemini(
		ctx,
		req.Prompt,
		req.Template,
		req.Variables,
		maxTokens,
		temperature,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate with gemini: %v", err)
	}

	return &assistantpb.GenerateResponse{
		Content: content,
		Metadata: &assistantpb.GenerateMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0, // TODO: Get from API response
			Model:            modelName,
		},
	}, nil
}

// ClassifyDocumentWithGemini - Đọc & phân loại 1 tài liệu bằng Gemini (multimodal)
// @Summary Phân loại tài liệu bằng Gemini
// @Description Gửi prompt kèm nội dung file (PDF/ảnh/text) để Gemini đọc và trả về kết quả. file_content rỗng => text-only.
// @Tags Assistant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body assistantpb.ClassifyDocumentRequest true "Prompt và nội dung file"
// @Success 200 {object} assistantpb.GenerateResponse
// @Router /v2/assistant/gemini/classify-document [post]
func (h *AssistantHandler) ClassifyDocumentWithGemini(ctx context.Context, req *assistantpb.ClassifyDocumentRequest) (*assistantpb.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, status.Errorf(codes.InvalidArgument, "prompt is required")
	}

	maxTokens := int(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 2000
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.2
	}

	content, modelName, processingTime, err := h.Usecase.ClassifyDocumentWithGemini(
		ctx,
		req.Prompt,
		req.FileContent,
		req.MimeType,
		maxTokens,
		temperature,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to classify document with gemini: %v", err)
	}

	return &assistantpb.GenerateResponse{
		Content: content,
		Metadata: &assistantpb.GenerateMetadata{
			ProcessingTimeMs: processingTime,
			TokensUsed:       0,
			Model:            modelName,
		},
	}, nil
}
