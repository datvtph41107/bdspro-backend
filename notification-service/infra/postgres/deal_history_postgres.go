package postgres

import (
	"context"
	"fmt"
	"notification/internal/domain"
	usecase "notification/internal/usecase"

	"gorm.io/gorm"
)

type dealHistoryPostgres struct {
	db *gorm.DB
}

func NewDealHistoryPostgres(db *gorm.DB) usecase.DealHistoryStore {
	return &dealHistoryPostgres{db: db}
}

func (r *dealHistoryPostgres) Create(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error) {
	result := r.db.WithContext(ctx).Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create deal history: %w", result.Error)
	}

	return history, nil
}

func (r *dealHistoryPostgres) GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*domain.DealHistoryEntity, uint32, error) {
	var histories []*domain.DealHistoryEntity
	var total int64

	// Đếm tổng số records
	countResult := r.db.WithContext(ctx).Model(&domain.DealHistoryEntity{}).Where("deal_id = ?", dealID).Count(&total)
	if countResult.Error != nil {
		return nil, 0, fmt.Errorf("failed to count deal history: %w", countResult.Error)
	}

	// Lấy danh sách với pagination
	offset := (page - 1) * size
	result := r.db.WithContext(ctx).
		Where("deal_id = ?", dealID).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&histories)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to get deal history: %w", result.Error)
	}

	return histories, uint32(total), nil
}

func (r *dealHistoryPostgres) GetByID(ctx context.Context, id uint64) (*domain.DealHistoryEntity, error) {
	var history domain.DealHistoryEntity

	result := r.db.WithContext(ctx).First(&history, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deal history by id: %w", result.Error)
	}

	return &history, nil
}

func (r *dealHistoryPostgres) Update(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error) {
	result := r.db.WithContext(ctx).Save(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update deal history: %w", result.Error)
	}

	return history, nil
}

func (r *dealHistoryPostgres) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&domain.DealHistoryEntity{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete deal history: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("deal history with id %d not found", id)
	}

	return nil
}
