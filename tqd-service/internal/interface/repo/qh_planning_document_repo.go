package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IQHPlanningDocumentRepo interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) error
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error)
	GetAll(ctx context.Context) ([]qh_domain.QHPlanningDocument, error)
	GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error)
}
