package handler

import (
	"context"
	"feedback/internal/dto"
	"feedback/internal/usecase"
	feedbackpb "pb/types/feedback"
	sharepb "pb/types/shared"
	"strings"
)

type InternalHandler struct {
	feedbackpb.UnimplementedInternalServiceServer
	RegistryUsecase *usecase.RegistryUsecase
	RateUsecase     *usecase.RateUsecase
}

func NewInternalHandler(registryUsecase *usecase.RegistryUsecase, rateUsecase *usecase.RateUsecase) *InternalHandler {
	return &InternalHandler{
		RegistryUsecase: registryUsecase,
		RateUsecase:     rateUsecase,
	}
}

// RegisterEntity registers a new entity in the registry
func (h *InternalHandler) RegisterEntity(ctx context.Context, req *feedbackpb.FeedbackRegistry) (*sharepb.SubmitResponse, error) {
	// Check if entity already exists
	existingEntity, err := h.RegistryUsecase.GetRegistryByName(ctx, req.RegistryName)
	if err == nil && existingEntity != nil {
		// Entity already exists, return success without creating new one
		return &sharepb.SubmitResponse{
			Message: "Entity already exists",
		}, nil
	}

	// Convert protobuf to DTO
	cleanName := strings.TrimSpace(strings.ToLower(req.RegistryName))
	createReq := &dto.RegistryCreateRequest{
		RegistryKey:  req.RegistryKey,
		RegistryName: cleanName,
		Port:         int(req.Port),
		Service:      req.Service,
	}

	// Create registry entry
	_, err = h.RegistryUsecase.CreateRegistry(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "Entity registered successfully",
	}, nil
}

func (h *InternalHandler) GetRateStats(ctx context.Context, req *sharepb.IdRequest) (*sharepb.RateStatsV3Proto, error) {
	rateStats, err := h.RateUsecase.GetRateStats(ctx, req.Id, "user")
	if err != nil {
		return nil, err
	}
	return &sharepb.RateStatsV3Proto{
		TotalReviews: uint64(rateStats.TotalReviews),
		AverageScore: float32(rateStats.AverageScore),
		Star1Count:   uint64(rateStats.Star1Count),
		Star2Count:   uint64(rateStats.Star2Count),
		Star3Count:   uint64(rateStats.Star3Count),
		Star4Count:   uint64(rateStats.Star4Count),
		Star5Count:   uint64(rateStats.Star5Count),
	}, nil
}
