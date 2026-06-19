package postgres

import (
	_db "common/db"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"
)

type QHPlanningDocumentRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningDocumentRepo @bind: internal/interface/repo.IQHPlanningDocumentRepo
func NewQHPlanningDocumentRepo(db *_db.TransactionRepo) repo.IQHPlanningDocumentRepo {
	return &QHPlanningDocumentRepo{db: db}
}

func (r *QHPlanningDocumentRepo) Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) error {
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningDocumentRepo) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) error {
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).Where("id = ?", id).Updates(entity).Error
}

func (r *QHPlanningDocumentRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Where("id = ?", id).Delete(&qh_domain.QHPlanningDocument{}).Error
}

func (r *QHPlanningDocumentRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	var entity qh_domain.QHPlanningDocument
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *QHPlanningDocumentRepo) GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	return r.GetByID(ctx, id)
}

func (r *QHPlanningDocumentRepo) GetAll(ctx context.Context) ([]qh_domain.QHPlanningDocument, error) {
	var entities []qh_domain.QHPlanningDocument
	err := r.db.GetDB(ctx).Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *QHPlanningDocumentRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error) {
	var documents []qh_domain.QHPlanningDocument
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).Where("deleted_at IS NULL")

	if req.PlanningProjectID != 0 {
		db = db.Where("planning_project_id = ?", req.PlanningProjectID)
	}
	if req.DocumentType != "" {
		db = db.Where("document_type = ?", req.DocumentType)
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
		Find(&documents).Error

	return documents, total, err
}
