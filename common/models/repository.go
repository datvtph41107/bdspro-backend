package cmodels

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FileRepository thao tác với FileEntity trong DB
type BaseRepository struct {
	DB *gorm.DB
}

// Save lưu file vào DB
func (r *BaseRepository) Create(c *gin.Context, file *BaseEntity) error {
	return r.DB.WithContext(c).Create(file).Error
}
