package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)
}

type organizationUsecase struct {
	OrganizationRepository repository.OrganizationRepository
}

func NewOrganizationUsecase(organizationRepository repository.OrganizationRepository) OrganizationUsecase {
	return &organizationUsecase{OrganizationRepository: organizationRepository}
}

func (o *organizationUsecase) CreateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	organization, err := o.OrganizationRepository.Create(ctx, organization)
	if err != nil {
		return nil, err
	}
	return organization, nil
}