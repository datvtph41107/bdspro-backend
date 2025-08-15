package repo

import (
	_utils "common/utils"
	"context"
	"time"
	"user/enums"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"gorm.io/gorm"
)

// @bind: user/internal/repo.ReportRepo
type ReportRepo struct {
	DB *gorm.DB
}

func NewReportRepo(db *gorm.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

func (r *ReportRepo) Create(ctx context.Context, report *user_model.Report) error {
	return r.DB.WithContext(ctx).Create(report).Error
}

func (r *ReportRepo) GetByID(ctx context.Context, id uint64) (*user_model.Report, error) {
	var report user_model.Report
	err := r.DB.WithContext(ctx).
		Preload("Reason").
		Preload("User").
		Preload("Reporter").
		Preload("ProofDocs").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&report).Error
	return &report, err
}

func (r *ReportRepo) List(ctx context.Context, req dto.ReportListRequest) ([]user_model.Report, int64, error) {
	var reports []user_model.Report
	var total int64

	query := r.DB.WithContext(ctx).
		Model(&user_model.Report{}).
		Preload("Reason").
		Preload("User").
		Where("deleted_at IS NULL")

	// Filter by UserID
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
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

	// Paginate and get results
	err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&reports).Error

	return reports, total, err
}

func (r *ReportRepo) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.DB.WithContext(ctx).
		Model(&user_model.Report{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *ReportRepo) Delete(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&user_model.Report{}).Error
}

func (r *ReportRepo) UpdateStatus(ctx context.Context, id uint64, status uint32, response string) error {
	updates := map[string]interface{}{
		"report_status": status,
		"updated_at":    time.Now(),
	}

	if response != "" {
		updates["response"] = response
	}

	// Nếu báo cáo được duyệt hoặc từ chối, thêm thông tin admin xử lý
	if status == uint32(enums.StatusApproved) || status == uint32(enums.StatusRejected) {
		updates["resolved_at"] = time.Now()

		// Lấy admin ID từ context nếu có
		if adminID := _utils.GetProfileIdWithContext(ctx); adminID != 0 {
			updates["resolved_by"] = adminID
			updates["updated_by"] = adminID
		}
	}

	return r.Update(ctx, id, updates)
}

// CreateProof tạo bằng chứng báo cáo
func (r *ReportRepo) CreateProof(ctx context.Context, proof *user_model.ReportProof) error {
	return r.DB.WithContext(ctx).Create(proof).Error
}

// GetProofsByReportID lấy danh sách bằng chứng theo reportID
func (r *ReportRepo) GetProofsByReportID(ctx context.Context, reportID uint64) ([]user_model.ReportProof, error) {
	var proofs []user_model.ReportProof
	err := r.DB.WithContext(ctx).
		Where("report_id = ? AND deleted_at IS NULL", reportID).
		Find(&proofs).Error
	return proofs, err
}

// DeleteProofsByReportID xóa tất cả bằng chứng của một báo cáo
func (r *ReportRepo) DeleteProofsByReportID(ctx context.Context, reportID uint64) error {
	return r.DB.WithContext(ctx).
		Where("report_id = ?", reportID).
		Delete(&user_model.ReportProof{}).Error
}

// GetReasons lấy danh sách nguyên nhân báo cáo
func (r *ReportRepo) GetReasons(ctx context.Context) ([]user_model.ReportReason, error) {
	var reasons []user_model.ReportReason
	err := r.DB.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Find(&reasons).Error
	return reasons, err
}
