package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type RatePostgres struct {
	db *gorm.DB
}

// @bind: crm/internal/repo.RateRepo
func NewRatePostgres(db *gorm.DB) repo.RateRepo {
	return &RatePostgres{db: db}
}

func (r *RatePostgres) Create(ctx context.Context, rate *domain.Rate) (*domain.Rate, error) {
	// Tạo rate trước để có ID
	if err := r.db.WithContext(ctx).Create(rate).Error; err != nil {
		return nil, err
	}

	// Lưu attachs nếu có
	if len(rate.Attachs) > 0 {
		for i := range rate.Attachs {
			rate.Attachs[i].RateID = rate.ID
		}
		if err := r.db.WithContext(ctx).Create(&rate.Attachs).Error; err != nil {
			return nil, err
		}
	}

	return rate, nil
}

func (r *RatePostgres) GetByID(ctx context.Context, id uint64) (*domain.Rate, error) {
	var rate domain.Rate
	if err := r.db.WithContext(ctx).
		Preload("Attachs").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatePostgres) DetailRate(ctx context.Context, id uint64) (*domain.Rate, error) {
	var rate domain.Rate
	if err := r.db.WithContext(ctx).
		Preload("Attachs").
		Where("id = ? AND deleted_at IS NULL and hidden = false", id).
		First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatePostgres) GetByProfileIdAndOwnerIdAndOwnerOf(ctx context.Context, profileId uint64, ownerId uint64, ownerOf uint8) (*domain.Rate, error) {
	var rate domain.Rate
	if err := r.db.WithContext(ctx).
		Where("created_by = ? AND owner_id = ? AND owner_of = ? AND deleted_at IS NULL", profileId, ownerId, ownerOf).
		First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatePostgres) GetRateList(ctx context.Context, id uint64, ownerOf uint32, req *dto.RateSearchDTO) ([]domain.Rate, int64, error) {
	var rates []domain.Rate
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Rate{}).
		Select("id, score, comment, parent_id, created_by, anonymous, created_at").
		Preload("Attachs").
		Where("owner_id = ? AND owner_of = ? AND deleted_at IS NULL", id, ownerOf)

	if req.ParentId != nil {
		query = query.Where("parent_id = ?", *req.ParentId)
	} else {
		query = query.Where("parent_id is null")
	}

	// Count total
	query.Count(&total)

	// Get results with pagination
	if err := query.
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&rates).Error; err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}

func (r *RatePostgres) Update(ctx context.Context, id uint64, rate *domain.Rate) error {
	// Update rate fields
	if err := r.db.WithContext(ctx).
		Model(&domain.Rate{}).
		Where("id = ? and deleted_at is null", id).
		Updates(map[string]interface{}{
			"score":   rate.Score,
			"comment": rate.Comment,
		}).Error; err != nil {
		return err
	}

	// Chỉ update attachs nếu rate.Attachs không phải là nil
	// nil = không update attachs, empty slice = xóa tất cả, có giá trị = thay thế
	if rate.Attachs != nil {
		// Xóa attachs cũ trước
		if err := r.db.WithContext(ctx).
			Where("rate_id = ?", id).
			Delete(&domain.FeedbackAttachEntity{}).Error; err != nil {
			return err
		}

		// Tạo attachs mới nếu có
		if len(rate.Attachs) > 0 {
			for i := range rate.Attachs {
				rate.Attachs[i].RateID = id
			}
			if err := r.db.WithContext(ctx).Create(&rate.Attachs).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RatePostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Rate{}).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *RatePostgres) UpdateHidden(ctx context.Context, id uint64, hidden bool) error {
	return r.db.WithContext(ctx).Model(&domain.Rate{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("hidden", hidden).Error
}

func (r *RatePostgres) CountInfo(ctx context.Context, ownerID uint64, ownerOf uint32) (domain.RateStats, error) {
	var stats domain.RateStats

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				COUNT(*) AS total_reviews,
				ROUND(AVG(score), 2) AS average_score,
				COUNT(CASE WHEN score = 1 THEN 1 END) AS star_1_count,
				COUNT(CASE WHEN score = 2 THEN 1 END) AS star_2_count,
				COUNT(CASE WHEN score = 3 THEN 1 END) AS star_3_count,
				COUNT(CASE WHEN score = 4 THEN 1 END) AS star_4_count,
				COUNT(CASE WHEN score = 5 THEN 1 END) AS star_5_count
			FROM feedback_rates
			WHERE owner_id = ? AND owner_of = ? AND deleted_at IS NULL
		`, ownerID, ownerOf).
		Scan(&stats).Error

	if err != nil {
		return domain.RateStats{}, err
	}
	return stats, nil
}

func (r *RatePostgres) GetHistoryRate(ctx context.Context, rateId uint64, page, size uint32) ([]domain.Rate, int64, error) {
	var rates []domain.Rate
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Rate{}).
		Select("id, score, comment, parent_id, created_by, anonymous, created_at").
		Preload("Attachs").
		Where("root_id = ? AND deleted_at IS NULL", rateId)

	// Count total
	query.Count(&total)

	// Get results with pagination
	if err := query.
		Order("created_at DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&rates).Error; err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}