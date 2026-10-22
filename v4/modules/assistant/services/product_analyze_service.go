package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"

	assistantpb "pb/types/assistant"

	"mdn4/modules/assistant/client"
	"mdn4/modules/assistant/dto"
	"mdn4/modules/assistant/provider"
	hubservices "mdn4/modules/hub/services"
)

// ProductAnalyzeService analyzes free-text property listings (migrated from assistant-service).
type ProductAnalyzeService interface {
	AnalyzeProductText(ctx context.Context, req dto.AnalyzeProductRequest) (*assistantpb.SuggestResponse, error)
}

// @bind: mdn4/modules/assistant/services.ProductAnalyzeService
type LocalProductAnalyzeService struct {
	gemini   *client.GeminiClient
	suggest  provider.ProductSuggestProvider
	location hubservices.LocationV2Service
}

func NewLocalProductAnalyzeService(
	gemini *client.GeminiClient,
	suggest provider.ProductSuggestProvider,
	location hubservices.LocationV2Service,
) *LocalProductAnalyzeService {
	return &LocalProductAnalyzeService{
		gemini:   gemini,
		suggest:  suggest,
		location: location,
	}
}

func (s *LocalProductAnalyzeService) AnalyzeProductText(ctx context.Context, req dto.AnalyzeProductRequest) (*assistantpb.SuggestResponse, error) {
	content := req.Content
	if content == "" {
		return nil, _errors.BadRequestException("content is required")
	}
	if req.Context != "" {
		content = fmt.Sprintf("%s\n\nContext bổ sung: %s", content, req.Context)
	}

	start := time.Now()

	productV3, err := s.gemini.AnalyzeProductText(ctx, content)
	if err != nil {
		return nil, err
	}

	addressV3, err := s.location.InferAddressFromText(ctx, content)
	if err != nil {
		return nil, err
	}

	suggest, err := s.suggest.GetSuggest(ctx, content)
	if err != nil {
		return nil, _errors.ReturnError(500, "failed to get product suggest")
	}

	productV3.Address = addressV3
	productV3.PropertyType = suggest.PropertyType
	productV3.DocType = suggest.DocType
	productV3.Amenities = suggest.Amenities

	modelName := s.gemini.GetModel()
	processingTime := time.Since(start).Milliseconds()

	return mapSuggestResponse(productV3, modelName, processingTime), nil
}
