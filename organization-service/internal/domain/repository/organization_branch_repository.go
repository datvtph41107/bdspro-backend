package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

// OrganizationBranchRepository defines the interface for organization branch data access
type OrganizationBranchRepository interface {
	// Create creates a new organization branch
	Create(ctx context.Context, branch *entity.OrganizationBranch) error

	// Update updates an existing organization branch
	Update(ctx context.Context, branch *entity.OrganizationBranch) error

	// Delete deletes an organization branch
	Delete(ctx context.Context, id uint32) error

	// GetByID retrieves an organization branch by ID
	GetByID(ctx context.Context, id uint32) (*entity.OrganizationBranch, error)

	// ListByOrganization retrieves branches of an organization with pagination
	ListByOrganization(ctx context.Context, organizationID uint32, offset, limit int) ([]*entity.OrganizationBranch, error)

	// CountByOrganization returns total number of branches in an organization
	CountByOrganization(ctx context.Context, organizationID uint32) (int64, error)

	// SearchByName searches branches by name with pagination
	SearchByName(ctx context.Context, organizationID uint32, name string, offset, limit int) ([]*entity.OrganizationBranch, error)

	// GetByManager retrieves branches by manager name
	GetByManager(ctx context.Context, organizationID uint32, manager string) ([]*entity.OrganizationBranch, error)

	// UpdateManager updates branch manager
	UpdateManager(ctx context.Context, id uint32, manager string) error

	// UpdateContactInfo updates branch contact information (phone and email)
	UpdateContactInfo(ctx context.Context, id uint32, phone, email string) error
} 