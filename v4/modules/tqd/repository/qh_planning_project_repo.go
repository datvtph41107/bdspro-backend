package repository

import (
	"context"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
)

type IQHPlanningProjectRepo interface {
	Create(ctx context.Context, entity *qh.QHPlanningProject) error
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningProject) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningProject, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh.QHPlanningProject, int64, error)
	ListLayers(ctx context.Context, req *qh_dto.ListPlanningProjectLayersRequest) ([]qh.QHLayer, int64, error)
	CountLayersByProject(ctx context.Context, projectID uint64) (int64, error)
	GetAll(ctx context.Context) ([]qh.QHPlanningProject, error)
	GetDetail(ctx context.Context, id uint64) (*qh.QHPlanningProject, error)
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
