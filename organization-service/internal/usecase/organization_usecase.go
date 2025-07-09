package usecase

import (
	_utils "common/utils"
	"context"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/pkg/utils"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)
	GetOrganizations(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error)
	GetOrganizationByUserId(ctx context.Context) ([]*entity.Organization, error)
	CheckOrganizationIsExist(ctx context.Context) (bool, error)
	GetOwnedOrganization(ctx context.Context) (*entity.Organization, error)
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
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	organizationByUserId, err := o.OrganizationRepository.FindByUserId(ctx, uint32(currentUserId))
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

func (o *organizationUsecase) UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	org, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
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

	organization.ID = uint32(organizationId)
	org, err = o.OrganizationRepository.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}

	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "UPDATE",
		LogData:        "Updated organization details",
	})

	return org, nil
}

func (o *organizationUsecase) GetOrganizations(ctx context.Context, page, size int) ([]*entity.Organization, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	return o.OrganizationRepository.FindByUserIdWithPagination(ctx, uint32(currentUserId), page, size)
}

func (o *organizationUsecase) GetOrganizationByUserId(ctx context.Context) ([]*entity.Organization, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	organizations, err := o.OrganizationRepository.FindByUserId(ctx, uint32(currentUserId))
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (o *organizationUsecase) CheckOrganizationIsExist(ctx context.Context) (bool, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	org, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return false, err
	}
	return org != nil, nil
}

func (o *organizationUsecase) GetOwnedOrganization(ctx context.Context) (*entity.Organization, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organization, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.InvalidOrganizationID()
	}
	return organization, nil
}
