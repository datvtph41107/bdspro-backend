package postgres

import (
	_dto "common/dto"
	"context"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"gorm.io/gorm"
)

// @bind: user/internal/repo.RateRepo
type RatePostgres struct {
	db             *gorm.DB
	attachPostgres *AttachPostgres
}

func NewRatePostgres(db *gorm.DB,
	attachPostgres *AttachPostgres,
) *RatePostgres {
	return &RatePostgres{
		db:             db,
		attachPostgres: attachPostgres,
	}
}

func (r *RatePostgres) CreateRate(ctx context.Context, rate *user_model.Rate) (*user_model.Rate, error) {
	if err := r.db.WithContext(ctx).Create(rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *RatePostgres) GetRateList(ctx context.Context, targetId uint64, targetType uint8, dto *dto.RateSearchDTO) ([]user_model.RateItem, int64, error) {
	var rates []user_model.RateItem
	var total int64

	query := r.db.WithContext(ctx).Model(&user_model.Rate{}).
		Where("hidden = false AND root_id IS NULL AND deleted_at IS NULL AND target_id = ? AND target_type = ?", targetId, targetType)

	if dto.ParentId != nil {
		query = query.Where("parent_id = ?", dto.ParentId)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Attachs").
		Preload("CreatedUser").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&rates).
		Error; err != nil {
		return nil, 0, err
	}

	for _, rate := range rates {
		if rate.Anonymous {
			rate.CreatedUser = nil
		}
	}

	return rates, total, nil
}

func (r *RatePostgres) UpdateRate(ctx context.Context, rate *user_model.Rate, oldRate *user_model.Rate) (*user_model.Rate, error) {
	id := rate.ID

	attachs := oldRate.Attachs

	// tạo rate mới là clone của rate cũ, trừ attachs
	oldRate.ID = 0
	oldRate.RootID = &id
	oldRate.Attachs = nil
	if err := GetDB(ctx, r.db).
		Create(oldRate).
		Error; err != nil {
		return nil, err
	}
	for i := range attachs {
		attachs[i].RateID = &oldRate.ID
	}

	// cập nhật rate cũ
	newAttachs := rate.Attachs
	rate.ID = id
	rate.Attachs = nil
	if err := GetDB(ctx, r.db).
		Select("Anonymous", "Score", "Comment").
		Where("id = ? and deleted_at is null", id).
		Updates(rate).
		Error; err != nil {
		return nil, err
	}
	for i := range newAttachs {
		newAttachs[i].ID = 0
		newAttachs[i].RateID = &id
	}

	mergeAttachs := append(attachs, newAttachs...)
	if err := r.attachPostgres.CreateBatchAttach(ctx, &mergeAttachs); err != nil {
		return nil, err
	}

	return rate, nil
}

func (r *RatePostgres) DeleteRate(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&user_model.Rate{}, id).Error
}

func (r *RatePostgres) GetByID(ctx context.Context, id uint64) (*user_model.Rate, error) {
	var rate user_model.Rate
	if err := r.db.WithContext(ctx).
		Preload("Attachs").
		First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatePostgres) GetHistoryRate(ctx context.Context, rateId uint64, dto *_dto.Pagable) ([]user_model.RateItem, int64, error) {
	var rates []user_model.RateItem
	var total int64

	query := r.db.WithContext(ctx).Model(&user_model.Rate{}).
		Where("root_id = ? AND deleted_at IS NULL", rateId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Attachs").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&rates).
		Error; err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}

func (r *RatePostgres) GetByProfileIdAndTargetIdAndTargetType(ctx context.Context, profileId uint64, targetId uint64, targetType uint8) (*user_model.Rate, error) {
	var rate user_model.Rate
	if err := r.db.WithContext(ctx).
		Where("created_by = ? AND target_id = ? AND target_type = ? AND deleted_at IS NULL", profileId, targetId, targetType).
		First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatePostgres) UpdateHidden(ctx context.Context, id uint64, hidden bool) error {
	return r.db.WithContext(ctx).Model(&user_model.Rate{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("hidden", hidden).Error
}

func (r *RatePostgres) CountInfo(ctx context.Context, targetID uint64, targetType uint8) (user_model.RateStats, error) {
	var stats user_model.RateStats

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
			FROM rates
			WHERE target_id = ? AND target_type = ? AND deleted_at IS NULL
		`, targetID, targetType).
		Scan(&stats).Error

	if err != nil {
		return user_model.RateStats{}, err
	}
	return stats, nil
}
