package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
)

type IQHPlanningEventRepo interface {
	Create(ctx context.Context, entity *qh_domain.QHPlanningEvent) error
	Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningEvent) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningEvent, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error)
	CountByProject(ctx context.Context, projectID uint64) (int64, error)
	GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningRelationItem, int64, error)
	CountRelationsByProject(ctx context.Context, projectID uint64) (int64, error)
	CreateRelation(ctx context.Context, entity *qh_domain.QHPlanningRelation) error
	UpdateRelation(ctx context.Context, id uint64, entity *qh_domain.QHPlanningRelation) error
	DeleteRelation(ctx context.Context, id uint64) error
	GetRelationByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningRelationItem, error)
}
