package postgre

import (
	"context"
	"crm/infra/impl"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type ReportPostgres struct {
	db *gorm.DB
}

func NewReportPostgres(db *gorm.DB) repo.ReportRepo {
	return &ReportPostgres{db: db}
}

func (r *ReportPostgres) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	if err := impl.GetDB(ctx, r.db).Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportPostgres) GetByID(ctx context.Context, id uint64) (*domain.Report, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).
		Preload("Reason").
		Preload("ProofDocs").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&report).Error
	return &report, err
}

func (r *ReportPostgres) GetByCreatedBy(ctx context.Context, createdBy uint64, req *dto.ReportListRequest) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Preload("Reason").
		Preload("ProofDocs").
		Where("created_by = ? AND deleted_at IS NULL", createdBy)

	// Count total
	query.Count(&total)

	// Get results with pagination
	if err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *ReportPostgres) GetList(ctx context.Context, req *dto.ReportListRequest) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Preload("Reason").
		Preload("ProofDocs").
		Where("deleted_at IS NULL")

	// Filter by UserID
	if req.OwnerId != nil {
		query = query.Where("owner_id = ? and owner_of = ?", *req.OwnerId, req.OwnerOf)
	}

	// Filter by Status
	if req.Status != nil {
		query = query.Where("report_status = ?", *req.Status)
	}

	// Filter by ReasonID
	if req.ReasonID != nil {
		query = query.Where("reason_id = ?", *req.ReasonID)
	}

	// Count total
	query.Count(&total)

	// Get results with pagination
	if err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *ReportPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *ReportPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *ReportPostgres) ExistByOwnerIdAndUserIdAndOwnerOf(ctx context.Context, ownerId uint64, userId uint64, ownerOf uint32) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Where("owner_id = ? AND created_by = ? AND owner_of = ? AND deleted_at IS NULL", ownerId, userId, ownerOf).
		Count(&count).Error
	return count > 0, err
}

func (r *ReportPostgres) GetIDByOwnerIdAndUserIdAndOwnerOf(ctx context.Context, ownerId uint64, userId uint64, ownerOf uint32) (*uint64, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).
		Model(&domain.Report{}).
		Select("id").
		Where("owner_id = ? AND created_by = ? AND owner_of = ? AND deleted_at IS NULL", ownerId, userId, ownerOf).
		First(&report).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &report.ID, nil
}

// ReportProofPostgres implementation
type ReportProofPostgres struct {
	db *gorm.DB
}

func NewReportProofPostgres(db *gorm.DB) *ReportProofPostgres {
	return &ReportProofPostgres{db: db}
}

func (r *ReportProofPostgres) Create(ctx context.Context, proof *domain.ReportProof) (*domain.ReportProof, error) {
	if err := r.db.WithContext(ctx).Create(proof).Error; err != nil {
		return nil, err
	}
	return proof, nil
}

func (r *ReportProofPostgres) GetByReportID(ctx context.Context, reportID uint64) ([]domain.ReportProof, error) {
	var proofs []domain.ReportProof
	err := r.db.WithContext(ctx).
		Where("report_id = ? AND deleted_at IS NULL", reportID).
		Order("created_at ASC").
		Find(&proofs).Error
	return proofs, err
}

func (r *ReportProofPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.ReportProof{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *ReportProofPostgres) DeleteByReportID(ctx context.Context, reportID uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.ReportProof{}).
		Where("report_id = ? AND deleted_at IS NULL", reportID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}