package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_db "common/db"
	"context"
	"time"
)

type IntentPostgres struct {
	DB *_db.TransactionRepo
}

func NewIntentPostgres(db *_db.TransactionRepo) repo.IntentRepo {
	return &IntentPostgres{
		DB: db,
	}
}

// Create tạo intent mới
func (r *IntentPostgres) Create(ctx context.Context, intent *domain.Intent) error {
	db := r.DB.GetDB(ctx)
	return db.Create(intent).Error
}

// GetByID lấy intent theo ID
func (r *IntentPostgres) GetByID(ctx context.Context, id uint64) (*domain.Intent, error) {
	var intent domain.Intent
	db := r.DB.GetDB(ctx)

	err := db.Where("id = ? AND deleted_at IS NULL", id).First(&intent).Error
	if err != nil {
		return nil, err
	}

	return &intent, nil
}

// GetList lấy danh sách intents với phân trang và filter
func (r *IntentPostgres) GetList(ctx context.Context, req *dto.IntentSearchRequest) ([]domain.Intent, int64, error) {
	var intents []domain.Intent
	var total int64

	db := r.DB.GetDB(ctx)
	query := db.Model(&domain.Intent{}).Where("deleted_at IS NULL")

	// Apply filters
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}

	if req.Action != nil {
		query = query.Where("detected_intent = ?", *req.Action)
	}

	if req.ProfileID != nil {
		query = query.Where("profile_id = ?", *req.ProfileID)
	}

	if req.OrganizationID != nil {
		query = query.Where("organization_id = ?", *req.OrganizationID)
	}

	if req.SessionID != nil {
		query = query.Where("session_id = ?", *req.SessionID)
	}

	if req.MinConfidence != nil {
		query = query.Where("confidence >= ?", *req.MinConfidence)
	}

	if req.FromDate != nil {
		query = query.Where("created_at >= ?", *req.FromDate)
	}

	if req.ToDate != nil {
		query = query.Where("created_at <= ?", *req.ToDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := req.GetOffset()
	limit := req.GetLimit()
	query = query.Offset(offset).Limit(limit)

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&intents).Error; err != nil {
		return nil, 0, err
	}

	return intents, total, nil
}

// GetBySessionID lấy danh sách intents theo session
func (r *IntentPostgres) GetBySessionID(ctx context.Context, sessionID string) ([]domain.Intent, error) {
	var intents []domain.Intent
	db := r.DB.GetDB(ctx)

	err := db.Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Order("created_at ASC").
		Find(&intents).Error

	if err != nil {
		return nil, err
	}

	return intents, nil
}

// GetStatistics lấy thống kê intents
func (r *IntentPostgres) GetStatistics(ctx context.Context, req *dto.IntentStatisticsRequest) ([]domain.IntentStatistics, error) {
	var statistics []domain.IntentStatistics
	db := r.DB.GetDB(ctx)

	query := db.Model(&domain.Intent{}).
		Select(`
			category,
			detected_intent as action,
			COUNT(*) as total_count,
			AVG(confidence) as avg_confidence,
			MAX(created_at) as last_used
		`).
		Where("deleted_at IS NULL")

	// Apply filters
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}

	if req.OrganizationID != nil {
		query = query.Where("organization_id = ?", *req.OrganizationID)
	}

	if req.FromDate != nil {
		query = query.Where("created_at >= ?", *req.FromDate)
	}

	if req.ToDate != nil {
		query = query.Where("created_at <= ?", *req.ToDate)
	}

	// Group by
	groupBy := "category, detected_intent"
	if req.GroupBy != "" {
		switch req.GroupBy {
		case "category":
			groupBy = "category"
		case "action":
			groupBy = "detected_intent"
		case "day":
			groupBy = "DATE(created_at), category, detected_intent"
		case "week":
			groupBy = "DATE_TRUNC('week', created_at), category, detected_intent"
		case "month":
			groupBy = "DATE_TRUNC('month', created_at), category, detected_intent"
		}
	}

	query = query.Group(groupBy).Order("total_count DESC")

	if err := query.Scan(&statistics).Error; err != nil {
		return nil, err
	}

	return statistics, nil
}

// GetTopIntents lấy top intents được sử dụng nhiều nhất
func (r *IntentPostgres) GetTopIntents(ctx context.Context, category *enums.EIntentCategory, limit int) ([]domain.IntentStatistics, error) {
	var statistics []domain.IntentStatistics
	db := r.DB.GetDB(ctx)

	query := db.Model(&domain.Intent{}).
		Select(`
			category,
			detected_intent as action,
			COUNT(*) as total_count,
			AVG(confidence) as avg_confidence,
			MAX(created_at) as last_used
		`).
		Where("deleted_at IS NULL")

	if category != nil {
		query = query.Where("category = ?", *category)
	}

	query = query.Group("category, detected_intent").
		Order("total_count DESC").
		Limit(limit)

	if err := query.Scan(&statistics).Error; err != nil {
		return nil, err
	}

	return statistics, nil
}

// Update cập nhật intent
func (r *IntentPostgres) Update(ctx context.Context, id uint64, intent *domain.Intent) error {
	db := r.DB.GetDB(ctx)
	return db.Model(&domain.Intent{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(intent).Error
}

// Delete xóa intent (soft delete)
func (r *IntentPostgres) Delete(ctx context.Context, id uint64) error {
	db := r.DB.GetDB(ctx)
	now := time.Now()

	return db.Model(&domain.Intent{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

// GetRecentIntentsByProfile lấy intents gần đây của user
func (r *IntentPostgres) GetRecentIntentsByProfile(ctx context.Context, profileID uint64, limit int) ([]domain.Intent, error) {
	var intents []domain.Intent
	db := r.DB.GetDB(ctx)

	err := db.Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Order("created_at DESC").
		Limit(limit).
		Find(&intents).Error

	if err != nil {
		return nil, err
	}

	return intents, nil
}

// Ensure IntentPostgres implements IntentRepo
var _ repo.IntentRepo = (*IntentPostgres)(nil)
