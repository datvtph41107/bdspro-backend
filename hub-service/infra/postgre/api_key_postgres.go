package postgres

import (
	"context"
	"errors"
	"time"

	_db "common/db"
	_provider "common/provider"
	"hub/internal/domain"
	_repo "hub/internal/repo"

	"gorm.io/gorm"
)

type ApiKeyPostgres struct {
	_provider.CrudRepo[domain.ApiKeyEntity]
}

func NewApiKeyRepo(db *_db.TransactionRepo) _repo.IApiKeyRepo {
	repo := &ApiKeyPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *ApiKeyPostgres) BeforeSave(ctx context.Context, id *uint64, entity *domain.ApiKeyEntity) error {
	if entity.Name == "" {
		return errors.New("tên key không được để trống")
	}
	if entity.AppName == "" {
		return errors.New("tên ứng dụng không được để trống")
	}
	if entity.ApiKey == "" {
		return errors.New("api key không được để trống")
	}

	query := r.GetDB(ctx).
		Model(&domain.ApiKeyEntity{}).
		Where("name = ? AND deleted_at IS NULL", entity.Name)

	if id != nil {
		query = query.Where("id <> ?", *id)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("tên key đã tồn tại")
	}

	valueQuery := r.GetDB(ctx).
		Model(&domain.ApiKeyEntity{}).
		Where("api_key = ? AND deleted_at IS NULL", entity.ApiKey)

	if id != nil {
		valueQuery = valueQuery.Where("id <> ?", *id)
	}

	if err := valueQuery.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("api key đã tồn tại")
	}

	return nil
}

func (r *ApiKeyPostgres) GetByID(ctx context.Context, id uint64) (*domain.ApiKeyEntity, error) {
	var entity domain.ApiKeyEntity
	err := r.GetDB(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _repo.ErrApiKeyNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *ApiKeyPostgres) GetByName(ctx context.Context, name string) (*domain.ApiKeyEntity, error) {
	var entity domain.ApiKeyEntity
	err := r.GetDB(ctx).
		Where("name = ? AND deleted_at IS NULL", name).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _repo.ErrApiKeyNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *ApiKeyPostgres) GetByValue(ctx context.Context, value string) (*domain.ApiKeyEntity, error) {
	var entity domain.ApiKeyEntity
	err := r.GetDB(ctx).
		Where("api_key = ? AND deleted_at IS NULL", value).
		First(&entity).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _repo.ErrApiKeyNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *ApiKeyPostgres) Delete(ctx context.Context, id uint64) error {
	result := r.GetDB(ctx).
		Model(&domain.ApiKeyEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return _repo.ErrApiKeyNotFound
	}
	return nil
}
