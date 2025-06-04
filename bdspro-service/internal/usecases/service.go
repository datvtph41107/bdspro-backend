package usecases

// type AssetCostService struct {
// 	_owner.OwnerService[domain.AssetCost, SearchDTO]
// 	AssetCostRepo *AssetCostRepo
// }

// func NewAssetCostService(Repo *AssetCostRepo) *AssetCostService {
// 	return &AssetCostService{
// 		OwnerService: _owner.OwnerService[domain.AssetCost, SearchDTO]{
// 			OwnerRepo: Repo,
// 		},
// 		AssetCostRepo: Repo,
// 	}
// }

// func (s *AssetCostService) Search2(c *gin.Context, dto SearchDTO) (_routes.ResponseDTO, error) {
// 	profileId := _jwt.GetProfileId(c)
// 	entities, total, err := s.AssetCostRepo.Search2(c,
// 		profileId,
// 		dto,
// 	)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, &_routes.Except{
// 			Code:    500,
// 			Message: err.Error(),
// 		}
// 	}

// 	return _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          entities,
// 		TotalElements: &total,
// 	}, nil
// }
