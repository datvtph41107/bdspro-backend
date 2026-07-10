package postgres

import (
	_db "common/db"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type QHPlanningEventRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningEventRepo @bind: internal/interface/repo.IQHPlanningEventRepo
func NewQHPlanningEventRepo(db *_db.TransactionRepo) repo.IQHPlanningEventRepo {
	return &QHPlanningEventRepo{db: db}
}

func (r *QHPlanningEventRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error) {
	var events []qh_domain.QHPlanningEvent
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningEvent{}).Where("deleted_at IS NULL")
	if req.PlanningProjectID != 0 {
		db = db.Where("planning_project_id = ?", req.PlanningProjectID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("event_date DESC NULLS LAST, created_at DESC").
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
