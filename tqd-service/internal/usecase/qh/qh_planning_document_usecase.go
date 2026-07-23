package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type IQHPlanningDocumentUsecase interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error)
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error)
	// Approve phê duyệt 1 tài liệu đã được AI phân loại (chỉ khi đang ở trạng thái Classified).
	Approve(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
	// Retry đưa tài liệu Failed về Pending để job classify chạy lại.
	Retry(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
}

type qhPlanningDocumentUsecase struct {
	repo        repo.IQHPlanningDocumentRepo
	projectRepo repo.IQHPlanningProjectRepo
	aiJobRepo   repo.IProAIJobRepo
}

// NewQHPlanningDocumentUsecase @bind: internal/usecase/qh.IQHPlanningDocumentUsecase
func NewQHPlanningDocumentUsecase(
	repo repo.IQHPlanningDocumentRepo,
	projectRepo repo.IQHPlanningProjectRepo,
	aiJobRepo repo.IProAIJobRepo,
) IQHPlanningDocumentUsecase {
	return &qhPlanningDocumentUsecase{
		repo:        repo,
		projectRepo: projectRepo,
		aiJobRepo:   aiJobRepo,
	}
}

func (u *qhPlanningDocumentUsecase) Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error) {
	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entity, nil
}

func (u *qhPlanningDocumentUsecase) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error) {
	if err := u.repo.Update(ctx, id, entity); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể cập nhật dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningDocumentUsecase) Delete(ctx context.Context, id uint64) (bool, error) {
	if err := u.repo.Delete(ctx, id); err != nil {
		return false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}
	return true, nil
}

func (u *qhPlanningDocumentUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}
	return entity, nil
}

func (u *qhPlanningDocumentUsecase) GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: err.Error(),
		}
	}
	return entities, total, nil
}

func (u *qhPlanningDocumentUsecase) Approve(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	document, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy tài liệu",
		}
	}
	if document.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Tài liệu chưa được AI phân loại xong (chưa ở trạng thái Classified), không thể phê duyệt",
		}
	}

	document.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusApproved)
	if err := u.repo.Update(ctx, id, document); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể phê duyệt tài liệu",
		}
	}
	return document, nil
}

func (u *qhPlanningDocumentUsecase) Retry(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	document, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy tài liệu",
		}
	}
	if document.ProcessStatus != enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed) {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Chỉ có thể retry tài liệu đang ở trạng thái Lỗi (Failed)",
		}
	}

	if err := u.repo.ResetClassifyForRetry(ctx, id); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể đưa tài liệu về chờ xử lý: " + err.Error(),
		}
	}

	// Đưa đồ án về Pending nếu đang Failed để worker/admin theo dõi lại.
	if project, err := u.projectRepo.GetByID(ctx, document.PlanningProjectID); err == nil && project != nil {
		if project.ProcessStatus == enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed) {
			_ = u.projectRepo.UpdateProcessStatus(ctx, project.ID, enums.PlanningProcessStatusPending)
			_ = u.aiJobRepo.UpdateProcessStatusByProjectID(ctx, project.ID, enums.PlanningProcessStatusPending)
		}
	}

	return u.repo.GetByID(ctx, id)
}
