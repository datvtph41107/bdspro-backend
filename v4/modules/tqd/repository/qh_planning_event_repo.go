package repository

import (
	"context"
	qh_dto "mdn4/modules/tqd/dto/qh"
	qh "mdn4/modules/tqd/entity/qh"
)

type IQHPlanningEventRepo interface {
	Create(ctx context.Context, entity *qh.QHPlanningEvent) error
	Update(ctx context.Context, id uint64, entity *qh.QHPlanningEvent) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHPlanningEvent, error)
	GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningEvent, int64, error)
	CountByProject(ctx context.Context, projectID uint64) (int64, error)
	GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh.QHPlanningRelationItem, int64, error)
	CountRelationsByProject(ctx context.Context, projectID uint64) (int64, error)
	CreateRelation(ctx context.Context, entity *qh.QHPlanningRelation) error
	UpdateRelation(ctx context.Context, id uint64, entity *qh.QHPlanningRelation) error
	DeleteRelation(ctx context.Context, id uint64) error
	GetRelationByID(ctx context.Context, id uint64) (*qh.QHPlanningRelationItem, error)
}
