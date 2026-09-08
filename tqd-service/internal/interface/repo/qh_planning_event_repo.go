package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IQHPlanningEventRepo interface {
	GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error)
	GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningRelationItem, int64, error)
}
