package postgres

import (
	_db "common/db"
	"context"
	"strings"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type QHPlanningProjectRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningProjectRepo @bind: internal/interface/repo.IQHPlanningProjectRepo
func NewQHPlanningProjectRepo(db *_db.TransactionRepo) repo.IQHPlanningProjectRepo {
	return &QHPlanningProjectRepo{db: db}
}

func (r *QHPlanningProjectRepo) Create(ctx context.Context, entity *qh_domain.QHPlanningProject) error {
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningProjectRepo) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) error {
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningProject{}).Where("id = ?", id).Updates(entity).Error
}

func (r *QHPlanningProjectRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Where("id = ?", id).Delete(&qh_domain.QHPlanningProject{}).Error
}

func (r *QHPlanningProjectRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	var entity qh_domain.QHPlanningProject
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *QHPlanningProjectRepo) GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	return r.GetByID(ctx, id)
}

func (r *QHPlanningProjectRepo) GetAll(ctx context.Context) ([]qh_domain.QHPlanningProject, error) {
	var entities []qh_domain.QHPlanningProject
	err := r.db.GetDB(ctx).Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *QHPlanningProjectRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error) {
	var projects []qh_domain.QHPlanningProject
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningProject{}).Where("deleted_at IS NULL")

	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", searchTerm, searchTerm)
	}
	if req.PlanningType != "" {
		db = db.Where("planning_type = ?", req.PlanningType)
	}
	if req.PlanningLevel != "" {
		db = db.Where("planning_level = ?", req.PlanningLevel)
	}
	if req.ValidityStatus != "" {
		db = db.Where("validity_status = ?", req.ValidityStatus)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("created_at DESC").
		Find(&projects).Error

	return projects, total, err
}
