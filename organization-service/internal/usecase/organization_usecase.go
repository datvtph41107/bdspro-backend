package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
}

type organizationUsecase struct {
	OrganizationRepository repository.OrganizationRepository
}

func NewOrganizationUsecase(organizationRepository repository.OrganizationRepository) OrganizationUsecase {
	return &organizationUsecase{OrganizationRepository: organizationRepository}
}

func (o *organizationUsecase) CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error) {
	exist, err := o.OrganizationRepository.ExistByTaxCode(ctx, organization.TaxCode)
	if err != nil {
		return nil, nil, err
	}
	if exist {
		return nil, nil, custom_error.InvalidTaxCode()
	}
	organization, member, err = o.OrganizationRepository.CreateOrganizationWithMember(ctx, organization, member)
	if err != nil {
		return nil, nil, err
	}
	return organization, member, nil
}
