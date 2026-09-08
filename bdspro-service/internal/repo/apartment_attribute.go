package repo

import (
	"bdspro/internal/domain"

	"github.com/gin-gonic/gin"
)

type ApartmentAttrRepo interface {
	Create(c *gin.Context, e *domain.ApartmentAttribute) error
}
