package usecase

import (
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type StageUsecase struct {
	stageRepo         repo.StageRepo
	permissionUsecase PermissionUsecase
	pipelineRepo      repo.PipelineRepo
	historyRepo       repo.HistoryRepo
	customerRepo      repo.CustomerRepo
}

func NewStageUsecase(
	stageRepo repo.StageRepo,
	permissionUsecase PermissionUsecase,
	pipelineRepo repo.PipelineRepo,
	historyRepo repo.HistoryRepo,
	customerRepo repo.CustomerRepo,
) *StageUsecase {
	return &StageUsecase{
		stageRepo:         stageRepo,
		permissionUsecase: permissionUsecase,
		pipelineRepo:      pipelineRepo,
		historyRepo:       historyRepo,
		customerRepo:      customerRepo,
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
	pipeline, err := u.CheckPermission(c, entity.PipelineID)
	if err != nil {
		return nil, err
	}
	stage, err := u.stageRepo.Create(c, entity)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Tạo giai đoạn " + stage.StageName,
		TargetId:   stage.ID,
		TargetType: enums.TargetHistoryStage,
		ActionType: enums.HistoryStageCreate,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})
	return stage, nil
}

func (u *StageUsecase) Update(c context.Context, id uint64, entity *domain.StageEntity) (*domain.StageEntity, error) {
	stage, err := u.stageRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	pipeline, err := u.CheckPermission(c, stage.PipelineID)
	if err != nil {
		return nil, err
	}

	stage.OrderNumber = entity.OrderNumber
	stage.StageName = entity.StageName
	stage.Step = entity.Step
	stage.Active = entity.Active
	stage.RuleID = entity.RuleID
	stage.PipelineID = entity.PipelineID

	stage, err = u.stageRepo.Update(c, stage)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Cập nhật giai đoạn " + stage.StageName,
		TargetId:   stage.ID,
		TargetType: enums.TargetHistoryStage,
		ActionType: enums.HistoryStageUpdate,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})

	return stage, nil
}

func (u *StageUsecase) Delete(c context.Context, id uint64) error {
	stage, err := u.stageRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	pipeline, err := u.CheckPermission(c, stage.PipelineID)
	if err != nil {
		return err
	}

	existed, err := u.customerRepo.ExistedWithStageID(c, id)
	if err != nil {
		return err
	}
	if existed {
		return &_routes.Except{
			Code:    400,
			Message: "Giai đoạn này đang có khách hàng",
		}
	}

	err = u.stageRepo.Delete(c, id)
	if err != nil {
		return err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Xóa giai đoạn " + stage.StageName,
		TargetId:   stage.ID,
		TargetType: enums.TargetHistoryStage,
		ActionType: enums.HistoryStageDelete,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})

	return nil
}
