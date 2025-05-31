package repository

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

// OrganizationProfileRepository defines the interface for organization profile data access
type OrganizationProfileRepository interface {
	// Create creates a new organization profile
	Create(ctx context.Context, profile *entity.OrganizationProfile) error

	// Update updates an existing organization profile
	Update(ctx context.Context, profile *entity.OrganizationProfile) error

	// Delete deletes an organization profile
	Delete(ctx context.Context, organizationID uint32) error

	// GetByOrganizationID retrieves an organization profile by organization ID
	GetByOrganizationID(ctx context.Context, organizationID uint32) (*entity.OrganizationProfile, error)

	// ListByIndustry retrieves organization profiles by industry with pagination
	ListByIndustry(ctx context.Context, industry string, offset, limit int) ([]*entity.OrganizationProfile, error)

	// UpdateDescription updates organization profile description
	UpdateDescription(ctx context.Context, organizationID uint32, description string) error

	// UpdateIndustry updates organization profile industry
	UpdateIndustry(ctx context.Context, organizationID uint32, industry string) error

	// UpdateNumberOfEmployees updates organization profile number of employees
	UpdateNumberOfEmployees(ctx context.Context, organizationID uint32, numberOfEmployees int32) error

	// UpdateFoundingDate updates organization profile founding date
	UpdateFoundingDate(ctx context.Context, organizationID uint32, foundingDate *time.Time) error
} 