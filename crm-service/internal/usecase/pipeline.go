package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type PipelineUsecase struct {
	pipelineRepo      repo.PipelineRepo
	permissionUsecase PermissionUsecase
}

func NewPipelineUsecase(
	pipelineRepo repo.PipelineRepo,
	permissionUsecase PermissionUsecase,
) *PipelineUsecase {
	return &PipelineUsecase{
		pipelineRepo:      pipelineRepo,
		permissionUsecase: permissionUsecase,
	}
}

func (u *PipelineUsecase) Search(c context.Context, ownerId uint64, ownerType enums.EOwnerType, dto dto.PipelineSearchDTO) ([]domain.PipelineEntity, int64, error) {
	if err := u.permissionUsecase.UserInOwner(c, ownerId, ownerType); err != nil {
		return nil, 0, err
	}
	return u.pipelineRepo.Search(c, ownerId, ownerType, dto)
}

func (u *PipelineUsecase) Create(c context.Context, organizationId uint64, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, organizationId, enums.EOwnerTypeOrgnization); err != nil {
		return nil, err
	}
	return u.pipelineRepo.Create(c, entity)
}

func (u *PipelineUsecase) Update(c context.Context, id uint64, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}
	return u.pipelineRepo.Update(c, entity)
}

func (u *PipelineUsecase) Delete(c context.Context, id uint64) error {
	pipeline, err := u.pipelineRepo.GetByID(c, id)
	if err != nil {
		return err
	}
	if err := u.permissionUsecase.UserInOwner(c, pipeline.OwnerID, pipeline.OwnerType); err != nil {
		return err
	}
	return u.pipelineRepo.Delete(c, id)
}
