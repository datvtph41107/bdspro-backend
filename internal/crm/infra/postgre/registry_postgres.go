package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
	"strings"

	"gorm.io/gorm"
)

type RegistryPostgres struct {
	db *gorm.DB
}

// @bind: crm/internal/repo.RegistryRepo
func NewRegistryPostgres(db *gorm.DB) repo.RegistryRepo {
	return &RegistryPostgres{db: db}
}

// Create creates a new registry entry
func (r *RegistryPostgres) Create(ctx context.Context, registry *domain.FeedbackRegistry) error {
	err := r.db.WithContext(ctx).Create(registry).Error
	if err != nil {
		return err
	}
	return nil
}

// List gets all registry entries
func (r *RegistryPostgres) List(ctx context.Context) ([]*domain.FeedbackRegistry, error) {
	var registries []*domain.FeedbackRegistry
	err := r.db.WithContext(ctx).Find(&registries).Error
	if err != nil {
		return nil, err
	}
	return registries, nil
}

// GetByRegistryName gets a registry entry by registry name
func (r *RegistryPostgres) GetByRegistryName(ctx context.Context, registryName string) (*domain.FeedbackRegistry, error) {
	var registry domain.FeedbackRegistry

	cleanName := strings.TrimSpace(strings.ToLower(registryName))
	err := r.db.WithContext(ctx).
		Where("registry_name = ?", cleanName).
		First(&registry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found, return nil without error
		}
		return nil, err
	}
	return &registry, nil
}