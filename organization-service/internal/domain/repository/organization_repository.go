package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

// OrganizationRepository defines the interface for organization data access
type OrganizationRepository interface {
	// Create creates a new organization
	Create(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)

	// Update updates an existing organization
	Update(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)

	// Delete deletes an organization by ID
	Delete(ctx context.Context, id uint32) (*entity.Organization, error)

	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id uint32) (*entity.Organization, error)

	// GetByTaxCode retrieves an organization by tax code
	GetByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error)

	// GetByEmail retrieves an organization by email
	GetByEmail(ctx context.Context, email string) (*entity.Organization, error)

	// List retrieves organizations with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.Organization, error)

	// Count returns total number of organizations
	Count(ctx context.Context) (int64, error)

	// UpdateStatus updates organization status
	UpdateStatus(ctx context.Context, id uint32, status string) error
} 