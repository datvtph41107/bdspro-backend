package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)

	UpdateOrganization(ctx context.Context, currentUserId uint32, organization *entity.Organization) (*entity.Organization, error)

	GetOrganizations(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error)

	GetOrganizationByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error)

	CheckOrganizationIsExist(ctx context.Context, id uint32) (bool, error)
}

type organizationUsecase struct {
	OrganizationRepository       repository.OrganizationRepository
	OrganizationMemberRepository repository.OrganizationMemberRepository
	organizationAuthUsecase      OrganizationAuthUsecase
	logWorker                    *OrganizationLogWorker
}

func NewOrganizationUsecase(
	organizationRepository repository.OrganizationRepository,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
) OrganizationUsecase {
	return &organizationUsecase{
		OrganizationRepository:       organizationRepository,
		OrganizationMemberRepository: organizationMemberRepository,
		organizationAuthUsecase:      organizationAuthUsecase,
		logWorker:                    logWorker,
	}
}

func (o *organizationUsecase) CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	
	organizationByUserId, err := o.OrganizationRepository.FindByUserId(ctx, currentUserId)
	if err != nil {
		return nil, nil, err
	}
	if len(organizationByUserId) >= 1 {
		return nil, nil, custom_error.UserOnlyOneOrganization()
	}
	
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

	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organization.ID,
		ActorId:        organization.CreatedBy,
		LogType:        "CREATE",
		LogData:        "Created new organization",
	})

	return organization, member, nil
}

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

	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organization.ID,
		ActorId:        currentUserId,
		LogType:        "UPDATE",
		LogData:        "Updated organization details",
	})

	return org, nil
}

func (o *organizationUsecase) GetOrganizations(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error) {
	return o.OrganizationRepository.FindByUserIdWithPagination(ctx, userId, offset, limit)
}

func (o *organizationUsecase) GetOrganizationByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error) {
	organizations, err := o.OrganizationRepository.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (o *organizationUsecase) CheckOrganizationIsExist(ctx context.Context, id uint32) (bool, error) {
	org, err := o.OrganizationRepository.FindById(ctx, id)
	if err != nil {
		return false, err
	}
	return org != nil, nil
}
