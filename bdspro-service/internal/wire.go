package internal

import (
	repo_impl "bdspro/internal/repo/impl"
	"bdspro/internal/service"
	usecase_impl "bdspro/internal/usecase/impl"

	"gorm.io/gorm"
)

func InitializeAssetService(db *gorm.DB) service.AssetService {
	assetRepo := repo_impl.NewAssetRepo(db)
	assetUseCase := usecase_impl.NewAssetUseCase(assetRepo)
	return service.NewAssetService(assetUseCase)
}
