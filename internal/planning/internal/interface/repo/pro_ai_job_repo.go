package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IProAIJobRepo interface {
	Create(ctx context.Context, entity *qh_domain.ProAIJob) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.ProAIJob, error)
	GetList(ctx context.Context, req *qh_dto.ListProAIJobsRequest) ([]qh_domain.ProAIJob, int64, error)
	// UpdateProcessStatusByProjectID sync process_status của job gắn planning_project_id.
	UpdateProcessStatusByProjectID(ctx context.Context, planningProjectID uint64, status uint32) error
}
