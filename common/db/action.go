package _db

import (
	"github.com/gin-gonic/gin"
)

func SaveModel(c *gin.Context, e any) error {
	return DB.WithContext(c).Save(e).Error
}
