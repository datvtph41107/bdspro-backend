package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	UpdateOrganization(ctx context.Context, currentUserId uint32, organization *entity.Organization) (*entity.Organization, error)
	GetOrganizations(ctx context.Context, userId uint32) ([]*entity.Organization, error)
}

type organizationUsecase struct {
	OrganizationRepository       repository.OrganizationRepository
	organizationAuthUsecase OrganizationAuthUsecase
}

func NewOrganizationUsecase(
	organizationRepository repository.OrganizationRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
) OrganizationUsecase {
	return &organizationUsecase{
		OrganizationRepository:       organizationRepository,
		organizationAuthUsecase: organizationAuthUsecase,
	}
}

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
	return org, nil
}


func (o *organizationUsecase) GetOrganizations(ctx context.Context, userId uint32) ([]*entity.Organization, error) {
	return o.OrganizationRepository.FindByUserId(ctx, userId)
}