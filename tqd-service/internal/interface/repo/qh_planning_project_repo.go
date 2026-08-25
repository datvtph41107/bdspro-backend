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
	ListLayers(ctx context.Context, req *qh_dto.ListPlanningProjectLayersRequest) ([]qh_domain.QHLayer, int64, error)
	CountLayersByProject(ctx context.Context, projectID uint64) (int64, error)
	GetAll(ctx context.Context) ([]qh_domain.QHPlanningProject, error)
	GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error)
	// UpdateProcessStatus cập nhật riêng cột process_status (dùng cho job classify, tránh Updates() bỏ qua giá trị 0).
	UpdateProcessStatus(ctx context.Context, id uint64, status uint32) error

	GetPlanningProjectSeoSource(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectSeoSource, error)
	GetPlanningProjectSeoSourcesByIDs(ctx context.Context, planningProjectIDs []uint64) ([]qh_dto.PlanningProjectSeoSource, error)
	UpdatePlanningProjectSeoID(ctx context.Context, planningProjectID uint64, seoID uint64) error
	// Assign/Unassign provinces and wards
	AssignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error
	UnassignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error
	AssignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error
	UnassignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error
}
