package mapper

import (
	"crm/internal/domain"
	"crm/internal/dto"
)

// ConvertRegistryCreateRequestToDomain converts DTO to domain model
func ConvertRegistryCreateRequestToDomain(req *dto.RegistryCreateRequest) *domain.FeedbackRegistry {
	if req == nil {
		return nil
	}

	return &domain.FeedbackRegistry{
		RegistryKey:  req.RegistryKey,
		RegistryName: req.RegistryName,
		Port:         req.Port,
		Service:      req.Service,
	}
}

// ConvertRegistryToResponse converts domain model to DTO
func ConvertRegistryToResponse(registry *domain.FeedbackRegistry) *dto.RegistryResponse {
	if registry == nil {
		return nil
	}

	return &dto.RegistryResponse{
		RegistryKey:  registry.RegistryKey,
		RegistryName: registry.RegistryName,
		Port:         registry.Port,
		Service:      registry.Service,
	}
}