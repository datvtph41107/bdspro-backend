package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type IQHPlanningDocumentUsecase interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error)
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) (*qh_domain.QHPlanningDocument, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error)
}

type qhPlanningDocumentUsecase struct {
	repo repo.IQHPlanningDocumentRepo
}

// NewQHPlanningDocumentUsecase @bind: internal/usecase/qh.IQHPlanningDocumentUsecase
func NewQHPlanningDocumentUsecase(repo repo.IQHPlanningDocumentRepo) IQHPlanningDocumentUsecase {
	return &qhPlanningDocumentUsecase{
		repo: repo,
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
