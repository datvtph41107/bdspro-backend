package iu_asset

import (
	"bdspro/v1/models"

	"github.com/gin-gonic/gin"
)

type IUserAssetService interface {
	CreateAsset(ctx *gin.Context, assetData *models.Asset) (*models.Asset, error)
	GetAssetByProductID(ctx *gin.Context, productId *uint64) (*models.Asset, error)
}
