package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IQHPlanningProjectRepo interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningProject) error
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error)
	GetAll(ctx context.Context) ([]qh_domain.QHPlanningProject, error)
	GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
	// UpdateProcessStatus cập nhật riêng cột process_status (dùng cho job classify, tránh Updates() bỏ qua giá trị 0).
	UpdateProcessStatus(ctx context.Context, id uint64, status uint32) error
}
