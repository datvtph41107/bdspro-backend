package repo

import (
	"context"
	"time"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RateRepo struct {
	DB *gorm.DB
}

func NewRateRepo(db *gorm.DB) *RateRepo {
	return &RateRepo{DB: db}
}

func (r *RateRepo) Create(c *gin.Context, rate *user_model.Rate) error {
	return r.DB.WithContext(c).Create(rate).Error
}

func (r *RateRepo) List(c *gin.Context, targetID uint64, targetType uint8, dto dto.RateSearch) ([]user_model.RateItem, error) {
	var rates []user_model.RateItem

	query := r.DB.WithContext(c).
		Model(&user_model.RateItem{}).
		Select("id, score, comment, parent_id, created_by, anonymous").
		Preload("CreatedUser").
		Where("target_id = ? AND target_type = ? AND deleted_at IS NULL", targetID, targetType)

	// Lọc theo điểm nếu có truyền vào
	if dto.Score != nil {
		query = query.Where("score = ?", dto.Score)
	}

	if dto.ParentID != nil {
		query = query.Where("parent_id = ?", dto.ParentID)
	} else {
		query = query.Where("parent_id is null")
	}

	if err := query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&rates).Error; err != nil {
		return nil, err
	}

	return rates, nil
}

func (r *RateRepo) Update(c *gin.Context, id uint64, rate *user_model.Rate) error {
	return r.DB.WithContext(c).
		Model(&user_model.Rate{}).
		Where("id = ? and deleted_at is null", id).
		Updates(map[string]interface{}{
			"score":   rate.Score,
			"comment": rate.Comment,
		}).Error
}

func (r *RateRepo) Delete(c *gin.Context, profileID, id uint64) error {
	return r.DB.WithContext(c).
		Model(&user_model.Rate{}).
		Where("id = ? and created_by = ? and deleted_at is null", id, profileID).
		Update("deleted_at", time.Now()).Error
}

func (r *RateRepo) CountInfo(c context.Context, targetID uint64, targetType uint8) (user_model.RateStats, error) {
	var stats user_model.RateStats

	err := r.DB.WithContext(c).
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
