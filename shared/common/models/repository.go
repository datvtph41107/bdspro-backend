package _models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

// Save lưu file vào DB
func (r *BaseRepository) Create(c *gin.Context, e *BaseEntity) error {
	return r.DB.WithContext(c).Create(e).Error
}

func (r *BaseRepository) Save(c *gin.Context, e any) error {
	return r.DB.WithContext(c).Save(e).Error
}
