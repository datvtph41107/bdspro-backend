package handler

import (
	"context"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type RegistryHandler struct {
	crmpb.UnimplementedRegistryServiceServer

	registryUsecase *usecase.RegistryUsecase
}

// @bind: crm/infra/handler.RegistryHandler
func NewRegistryHandler(registryUsecase *usecase.RegistryUsecase) *RegistryHandler {
	return &RegistryHandler{
		registryUsecase: registryUsecase,
	}
}

// List gets a list of registry entries
// @Summary Get registry list
// @Description Get all registry entries
// @Tags Registry
// @Accept json
// @Produce json
// @Success 200 {object} dto.RegistryListResponse
// @Router /registry [get]
func (h *RegistryHandler) List(ctx context.Context, req *sharepb.Empty) (*crmpb.GetRegistryListResponse, error) {
	// Get all registries
	registries, err := h.registryUsecase.GetAllRegistries(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	responses := make([]*crmpb.FeedbackRegistry, len(registries))
	for i, registry := range registries {
		responses[i] = &crmpb.FeedbackRegistry{
			RegistryKey:  registry.RegistryKey,
			RegistryName: registry.RegistryName,
			Port:         int32(registry.Port),
			Service:      registry.Service,
		}
	}

	return &crmpb.GetRegistryListResponse{
		Data: responses,
	}, nil
}
