package _owner

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IOwnerRepo[T IOwnerModel] interface {
	Create(c *gin.Context, entity *OwnerModel[T]) error
	Update(c *gin.Context, id uint64, entity *OwnerModel[T]) error
	Delete(c *gin.Context, id uint64) error
	Search(c *gin.Context, profileId uint64, offset, limit int) ([]OwnerModel[T], error)
	GetByID(c *gin.Context, id uint64) (*OwnerModel[T], error)
}

type OwnerRepo[M IOwnerModel] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *OwnerRepo[T]) Create(c *gin.Context, entity *OwnerModel[T]) error {
	return r.DB.WithContext(c).
		// Model(entity).
		Create(entity).
		Error
}

// Cập nhật bản ghi
func (r *OwnerRepo[T]) Update(c *gin.Context, id uint64, entity *OwnerModel[T]) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *OwnerRepo[T]) Delete(c *gin.Context, id uint64) error {
	return r.DB.Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *OwnerRepo[T]) Search(c *gin.Context, profileId uint64, offset, limit int) ([]OwnerModel[T], error) {
	var entities []OwnerModel[T]
	err := r.DB.Where("deleted_at IS NULL and owner_id = ?", profileId).
		Offset(offset).
		Limit(limit).
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *OwnerRepo[T]) GetByID(c *gin.Context, id uint64) (*OwnerModel[T], error) {
	var entity *OwnerModel[T]
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
