package usecase

import (
	base_enum "base/enum"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type PipelineUsecase struct {
	pipelineRepo       repo.PipelineRepo
	permissionUsecase  PermissionUsecase
	notificationClient provider.NotificationProvider
	customerRepo       repo.LeadRepo
	stageRepo          repo.StageRepo
	transaction        provider.ITransaction
}

func NewPipelineUsecase(
	pipelineRepo repo.PipelineRepo,
	permissionUsecase PermissionUsecase,
	notificationClient provider.NotificationProvider,
	customerRepo repo.LeadRepo,
	stageRepo repo.StageRepo,
	transaction provider.ITransaction,
) *PipelineUsecase {
	return &PipelineUsecase{
		pipelineRepo:       pipelineRepo,
		permissionUsecase:  permissionUsecase,
		notificationClient: notificationClient,
		customerRepo:       customerRepo,
		stageRepo:          stageRepo,
		transaction:        transaction,
	}
}

func (u *PipelineUsecase) Search(ctx context.Context, dto dto.PipelineSearchDTO, withStages bool) ([]domain.PipelineEntity, int64, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	if organizationId == 0 {
		return nil, 0, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập quy trình",
		}
	}

	// if err := u.permissionUsecase.UserInOwner(c, ownerId, ownerType); err != nil {
	// 	return nil, 0, err
	// }
	// ownerId
	return u.pipelineRepo.Search(ctx, organizationId, enums.EOwnerOfOrgnization, dto, withStages)
}

func (u *PipelineUsecase) Create(c context.Context, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	if organizationId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập quy trình",
		}
	}
	err := u.permissionUsecase.UserInOwner(c, organizationId, base_enum.EOwnerOfOrgnization)
	if err != nil {
		return nil, err
	}

	var pipeline *domain.PipelineEntity
	err = u.transaction.WithTransaction(c, func(c context.Context) error {
		entity.IsCustom = true
		entity.OwnerID = organizationId
		entity.OwnerType = base_enum.EOwnerOfOrgnization
		pipeline, err = u.pipelineRepo.Create(c, entity)
		if err != nil {
			return err
		}

		if entity.Stages != nil {
			for _, stage := range entity.Stages {
				stage.PipelineID = pipeline.ID
				_, err = u.stageRepo.Create(c, &stage)
				if err != nil {
					return err
				}
			}
		}

		pipeline.Stages = entity.Stages
		return nil
	})

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Tạo quy trình thành công",
	// 	Note:       []string{"Tạo quy trình", pipeline.PipelineName},
	// 	TargetId:   pipeline.ID,
	// 	TargetType: base_enum.TargetHistoryPipeline,
	// 	ActionType: base_enum.HistoryPipelineCreate,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })

	return pipeline, nil
}

func (u *PipelineUsecase) Update(c context.Context, id uint64, dto *dto.PipelineSaveDTO) (*domain.PipelineEntity, error) {
	entity, err := u.pipelineRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if !entity.IsCustom {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền cập nhật quy trình mặc định",
		}
	}

	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}

	organizationId := _utils.GetOrganizationIdFromContext(c)
	if organizationId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập quy trình",
		}
	}
	var pipeline *domain.PipelineEntity
	err = u.transaction.WithTransaction(c, func(c context.Context) error {
		entity.PipelineName = dto.PipelineName
		entity.Color = dto.Color
		entity.Active = dto.Active
		entity.OwnerID = organizationId
		entity.OwnerType = base_enum.EOwnerOfOrgnization
		pipeline, err = u.pipelineRepo.Update(c, entity)
		if err != nil {
			return err
		}

		stages := make([]domain.StageEntity, len(dto.Stages))
		for i, stage := range dto.Stages {
			stages[i] = *stage
			stages[i].PipelineID = pipeline.ID
		}
		err = u.stageRepo.BulkUpdate(c, dto.RemoveStageIds, stages)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Cập nhật quy trình thành công",
	// 	Note:       []string{"Cập nhật quy trình", pipeline.PipelineName},
	// 	TargetId:   pipeline.ID,
	// 	TargetType: base_enum.TargetHistoryPipeline,
	// 	ActionType: base_enum.HistoryPipelineUpdate,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })

	return pipeline, nil
}

func (u *PipelineUsecase) Delete(c context.Context, id uint64) error {
	pipeline, err := u.pipelineRepo.GetByID(c, id)
	if err != nil {
		return err
	}
	if err := u.permissionUsecase.UserInOwner(c, pipeline.OwnerID, pipeline.OwnerType); err != nil {
		return err
	}

	// check nếu pipeline có giai đoạn
	existed, err := u.customerRepo.ExistedWithPipelineID(c, id)
	if err != nil {
		return err
	}
	if existed {
		return &_routes.Except{
			Code:    400,
			Message: "Quy trình này đang có khách hàng",
		}
	}

	err = u.pipelineRepo.Delete(c, id)
	if err != nil {
		return err
	}

	// u.notificationClient.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	Title:      "Xóa quy trình thành công",
	// 	Note:       []string{"Xóa quy trình", pipeline.PipelineName},
	// 	TargetId:   pipeline.ID,
	// 	TargetType: base_enum.TargetHistoryPipeline,
	// 	ActionType: base_enum.HistoryPipelineDelete,
	// 	OwnerID:    &pipeline.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(pipeline.OwnerType),
	// })

	return nil
}

func (u *PipelineUsecase) GetByID(c context.Context, id uint64) (*domain.PipelineEntity, error) {
	pipeline, err := u.pipelineRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (u *PipelineUsecase) GetDefault(c context.Context) (*domain.PipelineEntity, *domain.StageEntity, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	if organizationId == 0 {
		return nil, nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập quy trình mặc định",
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)
	if profileId == 0 {
		return nil, nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập quy trình mặc định",
		}
	}

	pipeline, err := u.pipelineRepo.GetDefault(c, organizationId)
	if err != nil {
		return nil, nil, err
	}

	stages, err := u.stageRepo.GetByPipelineID(c, pipeline.ID)
	if err != nil {
		return nil, nil, err
	}

	var stage *domain.StageEntity
	if pipeline.StageDefaultID != nil {
		for _, s := range stages {
			if s.ID == *pipeline.StageDefaultID {
				stage = &s
				break
			}
		}
	} else {
		stage = &stages[0]
	}

	pipeline.Stages = stages

	return pipeline, stage, nil
}
