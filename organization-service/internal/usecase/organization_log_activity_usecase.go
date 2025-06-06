package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationLogActivityUsecase interface {
	CreateOrganizationLogActivity(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error)
	GetOrganizationLogActivityByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error)
	GetOrganizationLogActivitiesByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error)
	GetOrganizationLogActivitiesByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error)
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

func (u *organizationLogActivityUsecase) GetOrganizationLogActivitiesByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error) {
	return u.organizationLogActivityRepository.GetByOrganizationID(ctx, organizationID)
}

func (u *organizationLogActivityUsecase) GetOrganizationLogActivitiesByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error) {
	return u.organizationLogActivityRepository.GetByActorID(ctx, actorID)
} 