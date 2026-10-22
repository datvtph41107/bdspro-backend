package repository

import (
	"context"
	"errors"
	"time"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

// ErrAPIKeyNotFound is returned when an api_keys lookup fails.
var ErrAPIKeyNotFound = errors.New("api key: not found")

// APIKeyRepository manages api_keys persistence.
type APIKeyRepository interface {
	Create(ctx context.Context, key *entity.APIKey) error
	GetAll(ctx context.Context) ([]*entity.APIKey, error)
	GetByValue(ctx context.Context, value string) (*entity.APIKey, error)
	Delete(ctx context.Context, id uint64) error
}

type apiKeyPostgres struct{ db *gorm.DB }

func NewAPIKeyPostgres(db *gorm.DB) APIKeyRepository {
	return &apiKeyPostgres{db: db}
}

func (r *apiKeyPostgres) Create(ctx context.Context, key *entity.APIKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *apiKeyPostgres) GetAll(ctx context.Context) ([]*entity.APIKey, error) {
	var rows []*entity.APIKey
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *apiKeyPostgres) GetByValue(ctx context.Context, value string) (*entity.APIKey, error) {
	var row entity.APIKey
	err := r.db.WithContext(ctx).Where("api_key = ? AND deleted_at IS NULL", value).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAPIKeyNotFound
	}
	return &row, err
}

func (r *apiKeyPostgres) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&entity.APIKey{}).
		Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}
