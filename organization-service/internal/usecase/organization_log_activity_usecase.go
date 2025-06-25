package usecase

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
)

type OrganizationLogActivityUsecase interface {
	CreateOrganizationLogActivity(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error)
	GetOrganizationLogActivityByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error)
	GetOrganizationLogActivitiesByOrganizationID(ctx context.Context, organizationID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error)
	GetOrganizationLogActivitiesByActorID(ctx context.Context, actorID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error)
}

type organizationLogActivityUsecase struct {
	organizationLogActivityRepository repository.OrganizationLogActivityRepository
}

func NewOrganizationLogActivityUsecase(organizationLogActivityRepository repository.OrganizationLogActivityRepository) OrganizationLogActivityUsecase {
	return &organizationLogActivityUsecase{
		organizationLogActivityRepository: organizationLogActivityRepository,
	}
}

func (u *organizationLogActivityUsecase) CreateOrganizationLogActivity(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error) {
	return u.organizationLogActivityRepository.Create(ctx, log)
}

func (u *organizationLogActivityUsecase) GetOrganizationLogActivityByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error) {
	return u.organizationLogActivityRepository.GetByID(ctx, id)
}

func (u *organizationLogActivityUsecase) GetOrganizationLogActivitiesByOrganizationID(ctx context.Context, organizationID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error) {
	return u.organizationLogActivityRepository.GetByOrganizationIDWithPagination(ctx, organizationID, page, size)
}

func (u *organizationLogActivityUsecase) GetOrganizationLogActivitiesByActorID(ctx context.Context, actorID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error) {
	return u.organizationLogActivityRepository.GetByActorIDWithPagination(ctx, actorID, page, size)
}
