package attribute_attribute

import (
	"bdspro/v1/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApartmentAttrRepo struct {
	DB *gorm.DB
}

func NewApartmentAttrRepo(db *gorm.DB) *ApartmentAttrRepo {
	return &ApartmentAttrRepo{
		DB: db,
	}
}

func (r *ApartmentAttrRepo) Create(c *gin.Context, e *models.ApartmentAttribute) error {
	return r.DB.WithContext(c).Create(&e).Error
}
