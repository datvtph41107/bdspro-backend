package usecase

import (
	_errors "common/errors"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type StageUsecase struct {
	stageRepo          repo.StageRepo
	permissionUsecase  PermissionUsecase
	pipelineRepo       repo.PipelineRepo
	notificationClient provider.NotificationProvider
	customerRepo       repo.LeadRepo
}

func NewStageUsecase(
	stageRepo repo.StageRepo,
	permissionUsecase PermissionUsecase,
	pipelineRepo repo.PipelineRepo,
	notificationClient provider.NotificationProvider,
	customerRepo repo.LeadRepo,
) *StageUsecase {
	return &StageUsecase{
		stageRepo:          stageRepo,
		permissionUsecase:  permissionUsecase,
		pipelineRepo:       pipelineRepo,
		notificationClient: notificationClient,
		customerRepo:       customerRepo,
	}
}

func (u *StageUsecase) CheckPermission(c context.Context, pipelineId uint64, getRequest bool) (*domain.PipelineEntity, error) {
	pipeline, err := u.pipelineRepo.GetByID(c, pipelineId)
	if err != nil {
		return nil, err
	}
	if !pipeline.IsCustom && !getRequest {
		return nil, _errors.ReturnError(service.DefaultStageMutationDenied)
	}

	if err := u.permissionUsecase.UserInOwner(c, pipeline.OwnerID, pipeline.OwnerType); err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (u *StageUsecase) Search(c context.Context, pipelineId uint64, dto dto.StageSearchDTO) ([]domain.StageEntity, int64, error) {
	if _, err := u.CheckPermission(c, pipelineId, true); err != nil {
		return nil, 0, err
	}
	return u.stageRepo.Search(c, pipelineId, dto)
}

func (u *StageUsecase) Create(c context.Context, entity *domain.StageEntity) (*domain.StageEntity, error) {
	// pipeline, err := u.CheckPermission(c, entity.PipelineID, false)
	// if err != nil {
	// 	return nil, err
	// }
	stage, err := u.stageRepo.Create(c, entity)
	if err != nil {
		return nil, err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Tạo giai đoạn thành công",
	// 	Note:       []string{"Tạo giai đoạn", stage.StageName},
	// 	TargetId:   stage.ID,
	// 	TargetType: base_enum.TargetHistoryStage,
	// 	ActionType: base_enum.HistoryStageCreate,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })
	return stage, nil
}

func (u *StageUsecase) Update(c context.Context, id uint64, entity *domain.StageEntity) (*domain.StageEntity, error) {
	stage, err := u.stageRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	// pipeline, err := u.CheckPermission(c, stage.PipelineID, false)
	// if err != nil {
	// 	return nil, err
	// }

	stage.OrderNumber = entity.OrderNumber
	stage.StageName = entity.StageName
	stage.Step = entity.Step
	stage.Active = entity.Active
	stage.RuleID = entity.RuleID
	stage.PipelineID = entity.PipelineID
	stage.Color = entity.Color
	stage.ColorRGB = entity.ColorRGB

	stage, err = u.stageRepo.Update(c, stage)
	if err != nil {
		return nil, err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Cập nhật giai đoạn thành công",
	// 	Note:       []string{"Cập nhật giai đoạn", stage.StageName},
	// 	TargetId:   stage.ID,
	// 	TargetType: base_enum.TargetHistoryStage,
	// 	ActionType: base_enum.HistoryStageUpdate,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })

	return stage, nil
}

func (u *StageUsecase) Delete(c context.Context, id uint64) error {
	// stage, err := u.stageRepo.GetByID(c, id)
	// if err != nil {
	// 	return err
	// }

	// pipeline, err := u.CheckPermission(c, stage.PipelineID, false)
	// if err != nil {
	// 	return err
	// }

	existed, err := u.customerRepo.ExistedWithStageID(c, id)
	if err != nil {
		return err
	}
	if existed {
		return _errors.ReturnError(service.StageHasCustomers)
	}

	err = u.stageRepo.Delete(c, id)
	if err != nil {
		return err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Xóa giai đoạn thành công",
	// 	Note:       []string{"Xóa giai đoạn", stage.StageName},
	// 	TargetId:   stage.ID,
	// 	TargetType: base_enum.TargetHistoryStage,
	// 	ActionType: base_enum.HistoryStageDelete,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })

	return nil
}
