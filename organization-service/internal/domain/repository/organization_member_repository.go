package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

// OrganizationMemberRepository defines the interface for organization member data access
type OrganizationMemberRepository interface {
	// Create creates a new organization member
	Create(ctx context.Context, member *entity.OrganizationMember) error

	// Update updates an existing organization member
	Update(ctx context.Context, member *entity.OrganizationMember) error

	// Delete deletes an organization member by ID
	Delete(ctx context.Context, id uint32) error

	// GetByID retrieves an organization member by ID
	GetByID(ctx context.Context, id uint32) (*entity.OrganizationMember, error)

	// GetByUserAndOrg retrieves a member by user ID and organization ID
	GetByUserAndOrg(ctx context.Context, userID, orgID uint32) (*entity.OrganizationMember, error)

	// ListByOrganization retrieves members of an organization with pagination
	ListByOrganization(ctx context.Context, orgID uint32, offset, limit int) ([]*entity.OrganizationMember, error)

	// ListByUser retrieves all organizations a user is member of
	ListByUser(ctx context.Context, userID uint32) ([]*entity.OrganizationMember, error)

	// UpdateRole updates member's role
	UpdateRole(ctx context.Context, id uint32, role string) error

	// UpdateStatus updates member's status
	UpdateStatus(ctx context.Context, id uint32, status string) error

	// RemoveMember marks a member as removed with current timestamp
	RemoveMember(ctx context.Context, id uint32) error

	// CountByOrganization returns total number of members in an organization
	CountByOrganization(ctx context.Context, orgID uint32) (int64, error)
} 