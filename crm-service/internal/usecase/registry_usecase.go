package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"
)

type RegistryUsecase struct {
	registryData map[string]uint32
	registryRepo repo.RegistryRepo
}

func NewRegistryUsecase(registryRepo repo.RegistryRepo) *RegistryUsecase {
	ctx := context.Background()
	registryData := make(map[string]uint32)
	registries, err := registryRepo.List(ctx)
	if err == nil {
		// Load data into memory
		for _, registry := range registries {
			registryData[registry.RegistryName] = registry.RegistryKey
		}
	}

	return &RegistryUsecase{
		registryData: registryData,
		registryRepo: registryRepo,
	}
}

// CreateRegistry creates a new registry entry
func (u *RegistryUsecase) CreateRegistry(ctx context.Context, req *dto.RegistryCreateRequest) (*dto.RegistryResponse, error) {
	// Check if registry already exists
	existingRegistry, err := u.registryRepo.GetByRegistryName(ctx, req.RegistryName)
	if err != nil {
		return nil, err
	}
	if existingRegistry != nil {
		fmt.Printf("registry with name '%s' already exists\n", req.RegistryName)
		return nil, fmt.Errorf("registry with name '%s' already exists", req.RegistryName)
	}

	registry := &domain.FeedbackRegistry{
		RegistryKey:  req.RegistryKey,
		RegistryName: req.RegistryName,
		Port:         req.Port,
		Service:      req.Service,
	}

	err = u.registryRepo.Create(ctx, registry)
	if err != nil {
		return nil, err
	}

	// Update in-memory data
	u.registryData[req.RegistryName] = req.RegistryKey

	return &dto.RegistryResponse{
		RegistryKey:  registry.RegistryKey,
		RegistryName: registry.RegistryName,
		Port:         registry.Port,
		Service:      registry.Service,
	}, nil
}

// GetAllRegistries gets all registry entries
func (u *RegistryUsecase) GetAllRegistries(ctx context.Context) ([]*dto.RegistryResponse, error) {
	registries, err := u.registryRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.RegistryResponse, len(registries))
	for i, registry := range registries {
		responses[i] = &dto.RegistryResponse{
			RegistryKey:  registry.RegistryKey,
			RegistryName: registry.RegistryName,
			Port:         registry.Port,
			Service:      registry.Service,
		}
	}

	return responses, nil
}

// GetRegistryByName gets a registry entry by registry name
func (u *RegistryUsecase) GetRegistryByName(ctx context.Context, registryName string) (*dto.RegistryResponse, error) {
	registry, err := u.registryRepo.GetByRegistryName(ctx, registryName)
	if err != nil {
		return nil, err
	}
	if registry == nil {
		return nil, nil // Not found
	}

	return &dto.RegistryResponse{
		RegistryKey:  registry.RegistryKey,
		RegistryName: registry.RegistryName,
		Port:         registry.Port,
		Service:      registry.Service,
	}, nil
}

// GetOwnerOfByRegistry gets owner of by registry name from memory
func (u *RegistryUsecase) GetOwnerOfByRegistry(registryName string) uint32 {
	key, exists := u.registryData[registryName]
	if !exists {
		return 0
	}
	return key
}