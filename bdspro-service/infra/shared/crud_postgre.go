package shared

import (
	"bdspro/internal/common"
	"bdspro/internal/enums"
	"context"
	"time"

	"gorm.io/gorm"
)

type PostgreCrud[T any, D common.DTO] struct {
	DB   *gorm.DB
	Repo common.IOwnerRepo[T, D]
}

func (r *PostgreCrud[T, D]) QuerySearch(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) (*gorm.DB, error) {
	return r.DB.Model(new(T)).Where("deleted_at IS NULL"), nil
}

func (r *PostgreCrud[T, D]) Search(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) ([]T, int64, error) {
	var entities []T
	var total int64

	query, err := r.Repo.QuerySearch(c, profileId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&entities).Error; err != nil {
		return nil, 0, nil
	}
	query, _ = r.Repo.QuerySearch(c, profileId, ownerType, dto)
	err = query.Count(&total).Error
	return entities, total, err
}

func (r *PostgreCrud[T, D]) GetByID(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreCrud[T, D]) Detail(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreCrud[T, D]) Create(c context.Context, entity *T) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *PostgreCrud[T, D]) Update(c context.Context, id uint64, entity *T) error {
	return r.DB.WithContext(c).Model(new(T)).
		Where("id = ?", id).
		Updates(entity).Error
}

func (r *PostgreCrud[T, D]) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *PostgreCrud[T, D]) GetAll(c context.Context, dto D) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(c).Model(new(T)).
		Where("deleted_at IS NULL").
		Scan(&entities).Error
	return entities, err
}
