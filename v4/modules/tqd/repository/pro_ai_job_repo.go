package repository

import (
	"context"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
)

type IProAIJobRepo interface {
	Create(ctx context.Context, entity *qh.ProAIJob) error
	GetByID(ctx context.Context, id uint64) (*qh.ProAIJob, error)
	GetList(ctx context.Context, req *qh_dto.ListProAIJobsRequest) ([]qh.ProAIJob, int64, error)
	// UpdateProcessStatusByProjectID sync process_status của job gắn planning_project_id.
	UpdateProcessStatusByProjectID(ctx context.Context, planningProjectID uint64, status uint32) error
}
