package u_asset_cost_type

import (
	"bdspro/v1/models"
	_category "common/case/category"
	_jwt "common/jwt"
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type AssetCostTypeService struct {
	_category.CategoryService[models.AssetCostType, SearchDTO]
	AssetCostTypeRepo *AssetCostTypeRepo
}

func NewAssetCostTypeService(Repo *AssetCostTypeRepo) *AssetCostTypeService {
	return &AssetCostTypeService{
		CategoryService: _category.CategoryService[models.AssetCostType, SearchDTO]{
			CategoryRepo: Repo,
		},
		AssetCostTypeRepo: Repo,
	}
}

func (s *AssetCostTypeService) Search2(c *gin.Context, dto SearchDTO) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	// skip := dto.Page * dto.Size
	// limit := dto.Size
	// if limit <= 0 {
	// 	limit = 20
	// }
	entities, total, err := s.AssetCostTypeRepo.Search2(c,
		profileId,
		dto,
	)
	if err != nil {
		return _routes.ResponseDTO{}, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return _routes.ResponseDTO{
		Code:          0,
		Data:          entities,
		TotalElements: &total,
	}, nil
}
