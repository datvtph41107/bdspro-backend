package postgres

import (
	"context"
	"fmt"
	shared_enum "pb/enums"
	"time"

	"notification/internal/domain"
	"notification/internal/dto"
	usecase "notification/internal/usecase"

	"gorm.io/gorm"
)

type AdminHistoryPostgres struct {
	DB *gorm.DB
}

func NewAdminHistoryPostgres(db *gorm.DB) usecase.AdminHistoryStore {
	return &AdminHistoryPostgres{DB: db}
}

func (r *AdminHistoryPostgres) Search(c context.Context, adminID uint64, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.AdminHistoryEntity{}).
		Where("admin_id = ? and deleted_at is null", adminID)

	if dto.TargetType != nil {
		query = query.Where("target_type = ?", dto.TargetType)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.AdminRole != nil && *dto.AdminRole != "" {
		query = query.Where("admin_role = ?", dto.AdminRole)
	}

	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	if dto.IsInternal != nil {
		query = query.Where("is_internal = ?", dto.IsInternal)
	}

	query = query.Order("created_at DESC")

	var histories []domain.AdminHistoryEntity
	var total int64

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(dto.GetOffset()).Limit(dto.GetLimit()).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *AdminHistoryPostgres) SearchByOwner(c context.Context, ownerID uint64, ownerType shared_enum.EOwnerType, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.AdminHistoryEntity{}).
		Where("owner_id = ? and owner_type = ? and deleted_at is null", ownerID, ownerType)

	if dto.TargetType != nil {
		query = query.Where("target_type = ?", dto.TargetType)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.AdminRole != nil && *dto.AdminRole != "" {
		query = query.Where("admin_role = ?", dto.AdminRole)
	}

	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	query = query.Order("created_at DESC")

	var histories []domain.AdminHistoryEntity
	var total int64

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(dto.GetOffset()).Limit(dto.GetLimit()).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *AdminHistoryPostgres) CreateAdminHistory(c context.Context, entity *domain.AdminHistoryEntity) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *AdminHistoryPostgres) CreateAdminHistoryBatch(c context.Context, entities []*domain.AdminHistoryEntity) error {
	return r.DB.WithContext(c).CreateInBatches(entities, len(entities)).Error
}

func (r *AdminHistoryPostgres) GetAdminActions(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) ([]domain.AdminHistoryEntity, error) {
	var actions []domain.AdminHistoryEntity
	err := r.DB.WithContext(c).Model(&domain.AdminHistoryEntity{}).
		Where("admin_id = ? and created_at >= ? and created_at <= ? and deleted_at is null", adminID, fromDate, toDate).
		Order("created_at DESC").
		Find(&actions).Error
	return actions, err
}

func (r *AdminHistoryPostgres) GetRecentAdminActions(c context.Context, limit int) ([]domain.AdminHistoryEntity, error) {
	var actions []domain.AdminHistoryEntity
	err := r.DB.WithContext(c).Model(&domain.AdminHistoryEntity{}).
		Where("deleted_at is null").
		Order("created_at DESC").
		Limit(limit).
		Find(&actions).Error
	return actions, err
}

func (r *AdminHistoryPostgres) GetAdminActionStats(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) (map[string]int64, error) {
	var results []struct {
		ActionType int32 `json:"action_type"`
		Count      int64 `json:"count"`
	}

	err := r.DB.WithContext(c).Model(&domain.AdminHistoryEntity{}).
		Select("action_type, count(*) as count").
		Where("admin_id = ? and created_at >= ? and created_at <= ? and deleted_at is null", adminID, fromDate, toDate).
		Group("action_type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[fmt.Sprintf("action_%d", result.ActionType)] = result.Count
	}

	return stats, nil
}

func (r *AdminHistoryPostgres) SearchInternal(c context.Context, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Debug().Model(&domain.AdminHistoryEntity{}).
		Where("deleted_at is null")

	if dto.AdminID != nil {
		query = query.Where("admin_id = ?", dto.AdminID)
	}

	if dto.OwnerID != nil {
		query = query.Where("owner_id = ?", dto.OwnerID)
	}

	if dto.OwnerType != nil {
		query = query.Where("owner_type = ?", dto.OwnerType)
	}

	if dto.TargetType != nil {
		query = query.Where("target_type = ?", dto.TargetType)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.AdminRole != nil && *dto.AdminRole != "" {
		query = query.Where("admin_role = ?", dto.AdminRole)
	}

	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	if dto.IsInternal != nil {
		query = query.Where("is_internal = ?", dto.IsInternal)
	}

	query = query.Order("created_at DESC")

	var histories []domain.AdminHistoryEntity
	var total int64

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(dto.GetOffset()).Limit(dto.GetLimit()).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}
