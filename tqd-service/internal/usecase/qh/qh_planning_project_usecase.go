package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type IQHPlanningProjectUsecase interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error)
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error)
	// Approve phê duyệt đồ án (chỉ khi đã Classified) và phê duyệt luôn các tài liệu đã Classified thuộc đồ án.
	Approve(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
}

type qhPlanningProjectUsecase struct {
	repo         repo.IQHPlanningProjectRepo
	documentRepo repo.IQHPlanningDocumentRepo
}

// NewQHPlanningProjectUsecase @bind: internal/usecase/qh.IQHPlanningProjectUsecase
func NewQHPlanningProjectUsecase(repo repo.IQHPlanningProjectRepo, documentRepo repo.IQHPlanningDocumentRepo) IQHPlanningProjectUsecase {
	return &qhPlanningProjectUsecase{
		repo:         repo,
		documentRepo: documentRepo,
	}
}

func (u *qhPlanningProjectUsecase) Create(ctx context.Context, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error) {
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) (*qh_domain.QHPlanningProject, error) {
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningProjectUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningProjectUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningProjectUsecase) Approve(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	project, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy đồ án",
		}
	}
	if project.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Đồ án chưa được phân loại xong (chưa ở trạng thái Classified), không thể phê duyệt",
		}
	}

	if err := u.documentRepo.ApproveClassifiedByProject(ctx, id); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt các tài liệu của đồ án: " + err.Error(),
		}
	}

	if err := u.repo.UpdateProcessStatus(ctx, id, enums.PlanningProcessStatusApproved); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt đồ án",
		}
	}

	project.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusApproved)
	return project, nil
}
