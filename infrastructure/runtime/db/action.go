package _db

import (
	"context"

	"github.com/gin-gonic/gin"
)

func CreateWithAudit(c *gin.Context, e any) error {
	return DB.WithContext(c).Create(e).Error
}

func SaveWithAudit(c *gin.Context, e any) error {
	return DB.WithContext(c).Save(e).Error
}

func SaveWithContext(c context.Context, e any) error {
	return DB.WithContext(c).Save(e).Error
}
