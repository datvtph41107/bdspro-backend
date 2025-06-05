package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
)

type StageUsecase struct {
	stageRepo         repo.StageRepo
	permissionUsecase PermissionUsecase
	pipelineRepo      repo.PipelineRepo
}

func NewStageUsecase(
	stageRepo repo.StageRepo,
	permissionUsecase PermissionUsecase,
	pipelineRepo repo.PipelineRepo,
) *StageUsecase {
	return &StageUsecase{
		stageRepo:         stageRepo,
		permissionUsecase: permissionUsecase,
		pipelineRepo:      pipelineRepo,
	}
}

func (u *StageUsecase) CheckPermission(c context.Context, pipelineId uint64) (*domain.PipelineEntity, error) {
	pipeline, err := u.pipelineRepo.GetByID(c, pipelineId)
	if err != nil {
		return nil, err
	}

	if err := u.permissionUsecase.UserInOwner(c, pipeline.OwnerID, pipeline.OwnerType); err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (u *StageUsecase) Search(c context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error) {
	if _, err := u.CheckPermission(c, pipelineId); err != nil {
		return nil, 0, err
	}
	return u.stageRepo.Search(c, pipelineId, dto)
}

func (u *StageUsecase) Create(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	if _, err := u.CheckPermission(c, entity.PipelineID); err != nil {
		return nil, err
	}
	return u.stageRepo.Create(c, entity)
}

func (u *StageUsecase) Update(c context.Context, id uint64, entity *domain.StageEntity) (*domain.StageEntity, error) {
	stage, err := u.stageRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	if _, err := u.CheckPermission(c, stage.PipelineID); err != nil {
		return nil, err
	}

	stage.OrderNumber = entity.OrderNumber
	stage.StageName = entity.StageName
	stage.Step = entity.Step
	stage.Active = entity.Active
	stage.RuleID = entity.RuleID
	stage.PipelineID = entity.PipelineID

	return u.stageRepo.Update(c, stage)
}

func (u *StageUsecase) Delete(c context.Context, id uint64) error {
	if _, err := u.CheckPermission(c, id); err != nil {
		return err
	}
	return u.stageRepo.Delete(c, id)
}
