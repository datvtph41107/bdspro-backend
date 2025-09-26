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
}

func NewInternalHandler(registryUsecase *usecase.RegistryUsecase) *InternalHandler {
	return &InternalHandler{
		RegistryUsecase: registryUsecase,
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
