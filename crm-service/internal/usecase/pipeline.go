package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type PipelineUsecase struct {
	pipelineRepo      repo.PipelineRepo
	permissionUsecase PermissionUsecase
	historyRepo       repo.HistoryRepo
	customerRepo      repo.LeadRepo
	stageRepo         repo.StageRepo
}

func NewPipelineUsecase(
	pipelineRepo repo.PipelineRepo,
	permissionUsecase PermissionUsecase,
	historyRepo repo.HistoryRepo,
	customerRepo repo.LeadRepo,
	stageRepo repo.StageRepo,
) *PipelineUsecase {
	return &PipelineUsecase{
		pipelineRepo:      pipelineRepo,
		permissionUsecase: permissionUsecase,
		historyRepo:       historyRepo,
		customerRepo:      customerRepo,
		stageRepo:         stageRepo,
	}
}

func (u *PipelineUsecase) Search(ctx context.Context, dto dto.PipelineSearchDTO) ([]domain.PipelineEntity, int64, error) {
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
	return u.pipelineRepo.Search(ctx, organizationId, enums.EOwnerTypeOrgnization, dto)
}

func (u *PipelineUsecase) Create(c context.Context, organizationId uint64, entity *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	if err := u.permissionUsecase.UserInOwner(c, organizationId, enums.EOwnerTypeOrgnization); err != nil {
		return nil, err
	}

	entity.IsCustom = true
	pipeline, err := u.pipelineRepo.Create(c, entity)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Tạo quy trình " + pipeline.PipelineName,
		TargetId:   pipeline.ID,
		TargetType: enums.TargetHistoryPipeline,
		ActionType: enums.HistoryPipelineCreate,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})

	return pipeline, nil
}

func (u *PipelineUsecase) Update(c context.Context, id uint64, dto *domain.PipelineEntity) (*domain.PipelineEntity, error) {
	entity, err := u.pipelineRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	if err := u.permissionUsecase.UserInOwner(c, entity.OwnerID, entity.OwnerType); err != nil {
		return nil, err
	}

	entity.PipelineName = dto.PipelineName
	entity.Color = dto.Color
	entity.Active = dto.Active
	pipeline, err := u.pipelineRepo.Update(c, entity)
	if err != nil {
		return nil, err
	}

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Cập nhật quy trình " + pipeline.PipelineName,
		TargetId:   pipeline.ID,
		TargetType: enums.TargetHistoryPipeline,
		ActionType: enums.HistoryPipelineUpdate,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})

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

	u.historyRepo.CreateHistory(c, &domain.HistoryEntity{
		Note:       "Xóa quy trình " + pipeline.PipelineName,
		TargetId:   pipeline.ID,
		TargetType: enums.TargetHistoryPipeline,
		ActionType: enums.HistoryPipelineDelete,
		OwnerID:    &pipeline.OwnerID,
		OwnerType:  pipeline.OwnerType,
	})

	return nil
}

func (u *PipelineUsecase) GetByID(c context.Context, id uint64) (*domain.PipelineEntity, error) {
	return u.pipelineRepo.GetByID(c, id)
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

	var stage *domain.StageEntity
	if pipeline.StageDefaultID != 0 {
		stage, _ = u.stageRepo.GetByID(c, pipeline.StageDefaultID)
	} else {
		stage, _ = u.stageRepo.GetOneStageDefault(c, pipeline.ID)
	}

	return pipeline, stage, nil
}
