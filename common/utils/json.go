package _utils

import (
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

func ParseBody[T any](c *gin.Context) (*T, error) {
	var entity T

	// Đọc toàn bộ body từ request
	if err := c.ShouldBindQuery(&entity); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	return &entity, nil
}
