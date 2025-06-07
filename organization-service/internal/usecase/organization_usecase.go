package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

// OrganizationUsecase defines the business logic for organization management
type OrganizationUsecase interface {
	// CreateOrganization creates a new organization and adds the creator as an admin member
	// Returns the created organization, member info, and any error that occurred
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)

	// UpdateOrganization updates an existing organization's details
	// Only users with admin role can update the organization
	// Returns the updated organization and any error that occurred
	UpdateOrganization(ctx context.Context, currentUserId uint32, organization *entity.Organization) (*entity.Organization, error)

	// GetOrganizations retrieves all organizations where the user is either a creator or a member
	// Results are sorted by creation time in descending order
	// Returns a list of organizations, total count and any error that occurred
	GetOrganizations(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error)
}

// organizationUsecase implements OrganizationUsecase interface
type organizationUsecase struct {
	OrganizationRepository      repository.OrganizationRepository
	OrganizationMemberRepository repository.OrganizationMemberRepository
	organizationAuthUsecase     OrganizationAuthUsecase
	logWorker                   *OrganizationLogWorker
}

// NewOrganizationUsecase creates a new instance of OrganizationUsecase
func NewOrganizationUsecase(
	organizationRepository repository.OrganizationRepository,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
) OrganizationUsecase {
	return &organizationUsecase{
		OrganizationRepository:      organizationRepository,
		OrganizationMemberRepository: organizationMemberRepository,
		organizationAuthUsecase:     organizationAuthUsecase,
		logWorker:                   logWorker,
	}
}

// CreateOrganization creates a new organization and adds the creator as an admin member
// It performs the following steps:
// 1. Validates the organization data
// 2. Checks for duplicate tax code
// 3. Creates the organization
// 4. Creates the creator as an admin member
// 5. Logs the creation event
func (o *organizationUsecase) CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error) {
	org, err := o.OrganizationRepository.FindByTaxCode(ctx, organization.TaxCode)
	if err != nil {
		return nil, nil, err
	}
	if org != nil {
		return nil, nil, custom_error.InvalidTaxCode()
	}

	organization, member, err = o.OrganizationRepository.CreateOrganizationWithMember(ctx, organization, member)
	if err != nil {
		return nil, nil, err
	}

	// Log organization creation
	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organization.ID,
		ActorId:       organization.CreatedBy,
		LogType:       "CREATE",
		LogData:       "Created new organization",
	})

	return organization, member, nil
}

// UpdateOrganization updates an existing organization's details
// It performs the following steps:
// 1. Verifies the user has admin role
// 2. Validates the organization exists
// 3. Checks for duplicate tax code if tax code is being changed
// 4. Updates the organization
// 5. Logs the update event
func (o *organizationUsecase) UpdateOrganization(ctx context.Context, currentUserId uint32, organization *entity.Organization) (*entity.Organization, error) {
	hasRole, err := o.organizationAuthUsecase.HasRole(ctx, currentUserId, organization.ID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	org, err := o.OrganizationRepository.FindById(ctx, organization.ID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, custom_error.InvalidOrganizationID()
	}

	if organization.TaxCode != org.TaxCode {
		org, err = o.OrganizationRepository.FindByTaxCode(ctx, organization.TaxCode)
		if err != nil {
			return nil, err
		}
		if org != nil && org.ID != organization.ID {
			return nil, custom_error.InvalidTaxCode()
		}
	}

	org, err = o.OrganizationRepository.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}

	// Log organization update
	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organization.ID,
		ActorId:       currentUserId,
		LogType:       "UPDATE",
		LogData:       "Updated organization details",
	})

	return org, nil
}

// GetOrganizations retrieves all organizations where the user is either a creator or a member
// Results are sorted by creation time in descending order
// Returns a list of organizations, total count and any error that occurred
func (o *organizationUsecase) GetOrganizations(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error) {
	return o.OrganizationRepository.FindByUserIdWithPagination(ctx, userId, offset, limit)
}
