package crud3

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type BaseRepo[T any] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *BaseRepo[T]) Create(c context.Context, entity *T) error {
	return r.DB.WithContext(c).Create(entity).Error
}

// Cập nhật bản ghi
func (r *BaseRepo[T]) Update(c context.Context, id uint64, entity *T) error {
	return r.DB.WithContext(c).Model(entity).Where("id = ? and deleted_at is null", id).Updates(entity).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *BaseRepo[T]) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).Model(new(T)).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetAll() ([]T, error) {
	var entities []T
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *BaseRepo[T]) GetByID(id uint64) (*T, error) {
	var entity T
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
