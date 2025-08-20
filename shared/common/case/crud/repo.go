package crud

import (
	_dto "common/dto"
	"context"
	"time"

	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	GetAll(c context.Context) ([]T, error)
	GetByID(c context.Context, id uint64) (*T, error)
	GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error)
	GetDetail(c context.Context, id uint64) (*T, error)
}

type BaseRepo[T any] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *BaseRepo[T]) Create(c context.Context, entity *T) error {
	return r.DB.WithContext(c).
		Create(entity).
		Error
}

// Cập nhật bản ghi
func (r *BaseRepo[T]) Update(c context.Context, id uint64, entity *T) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).
		Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *BaseRepo[T]) Delete(c context.Context, id uint64) error {
	return r.DB.Model(new(T)).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", time.Now()).
		Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetAll(c context.Context) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(c).
		Where("deleted_at is null").
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetByID(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.WithContext(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepo[T]) GetList(c context.Context, pagable _dto.IPagable) ([]T, int64, error) {
	var entities []T
	var total int64
	err := r.DB.WithContext(c).
		Where("deleted_at is null").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.DB.WithContext(c).
		Model(new(T)).
		Where("deleted_at is null").
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	return entities, total, err
}

func (r *BaseRepo[T]) GetDetail(c context.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.WithContext(c).
		Where("id = ? and deleted_at is null", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
