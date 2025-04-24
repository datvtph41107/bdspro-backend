package house_info

import (
	"bdspro/v1/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HouseInfoRepo struct {
	DB *gorm.DB
}

func NewHouseInfoRepo(db *gorm.DB) *HouseInfoRepo {
	return &HouseInfoRepo{
		DB: db,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *HouseInfoRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.HouseInfo{}).
		Where("created_by = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *HouseInfoRepo) Create(access *models.HouseInfo) error {
	return repo.DB.Create(&access).Error
}

func (r *HouseInfoRepo) UpdateHouseInfo(ctx *gin.Context, houseInfo *models.HouseInfo) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	// tx := r.DB.Begin()

	// houseInfo.ProductID = productID
	return r.DB.WithContext(ctx).Save(&houseInfo).Error
}
