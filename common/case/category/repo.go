package _category

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ICategoryRepo[T ICategoryModel] interface {
	Create(c *gin.Context, entity *CategoryModel[T]) error
	Update(c *gin.Context, id uint64, entity *CategoryModel[T]) error
	Delete(c *gin.Context, id uint64) error
	Search(c *gin.Context, profileId uint64, offset, limit int) ([]CategoryModel[T], error)
	GetByID(c *gin.Context, id uint64) (*CategoryModel[T], error)
}

type CategoryRepo[M ICategoryModel] struct {
	DB *gorm.DB
}

// Tạo mới bản ghi
func (r *CategoryRepo[T]) Create(c *gin.Context, entity *CategoryModel[T]) error {
	return r.DB.WithContext(c).
		// Model(entity).
		Create(entity).
		Error
}

// Cập nhật bản ghi
func (r *CategoryRepo[T]) Update(c *gin.Context, id uint64, entity *CategoryModel[T]) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

// Xóa mềm bản ghi (cập nhật deleted_at thay vì xóa trực tiếp)
func (r *CategoryRepo[T]) Delete(c *gin.Context, id uint64) error {
	return r.DB.Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// Lấy tất cả bản ghi (bỏ qua những bản ghi đã bị xóa mềm)
func (r *CategoryRepo[T]) Search(c *gin.Context, profileId uint64, offset, limit int) ([]CategoryModel[T], error) {
	var entities []CategoryModel[T]
	err := r.DB.Where("deleted_at IS NULL and (owner_id = ? or is_custom = false)", profileId).
		Offset(offset).
		Limit(limit).
		Find(&entities).
		Error
	return entities, err
}

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *CategoryRepo[T]) GetByID(c *gin.Context, id uint64) (*CategoryModel[T], error) {
	var entity *CategoryModel[T]
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
