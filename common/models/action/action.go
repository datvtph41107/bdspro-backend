package c_action

import (
	c_config "common/config"

	"github.com/gin-gonic/gin"
)

func SaveModel(c *gin.Context, e any) error {
	return c_config.ConfigDB.WithContext(c).Save(e).Error
}
