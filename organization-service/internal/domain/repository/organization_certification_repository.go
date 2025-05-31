package repository

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

// OrganizationCertificationRepository defines the interface for organization certification data access
type OrganizationCertificationRepository interface {
	// Create creates a new organization certification
	Create(ctx context.Context, certification *entity.OrganizationCertification) error

	// Update updates an existing organization certification
	Update(ctx context.Context, certification *entity.OrganizationCertification) error

	// Delete deletes an organization certification
	Delete(ctx context.Context, id uint32) error

	// GetByID retrieves an organization certification by ID
	GetByID(ctx context.Context, id uint32) (*entity.OrganizationCertification, error)

	// ListByOrganization retrieves certifications of an organization with pagination
	ListByOrganization(ctx context.Context, organizationID uint32, offset, limit int) ([]*entity.OrganizationCertification, error)

	// ListByType retrieves certifications by type with pagination
	ListByType(ctx context.Context, organizationID uint32, certificationType string) ([]*entity.OrganizationCertification, error)

	// ListExpired retrieves expired certifications of an organization
	ListExpired(ctx context.Context, organizationID uint32) ([]*entity.OrganizationCertification, error)

	// ListExpiringWithin retrieves certifications that will expire within the specified duration
	ListExpiringWithin(ctx context.Context, organizationID uint32, duration time.Duration) ([]*entity.OrganizationCertification, error)

	// UpdateStatus updates certification status
	UpdateStatus(ctx context.Context, id uint32, status string) error

	// UpdateExpiryDate updates certification expiry date
	UpdateExpiryDate(ctx context.Context, id uint32, expiryDate *time.Time) error

	// CountByOrganization returns total number of certifications for an organization
	CountByOrganization(ctx context.Context, organizationID uint32) (int64, error)

	// CountByStatus returns number of certifications by status for an organization
	CountByStatus(ctx context.Context, organizationID uint32, status string) (int64, error)
} 