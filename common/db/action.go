package _db

import (
	"github.com/gin-gonic/gin"
)

func SaveWithAudit(c *gin.Context, e any) error {
	return DB.WithContext(c).Save(e).Error
}
