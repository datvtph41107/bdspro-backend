package u_asset_share

import (
	"bdspro/v1/models"
	_owner "common/case/owner"
	_jwt "common/jwt"
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type AssetShareService struct {
	_owner.OwnerService[models.AssetShare, SearchDTO]
	Repo *AssetShareRepo
}

func NewAssetShareService(Repo *AssetShareRepo) *AssetShareService {
	return &AssetShareService{
		OwnerService: _owner.OwnerService[models.AssetShare, SearchDTO]{
			OwnerRepo: Repo,
		},
		Repo: Repo,
	}
}

func (s AssetShareService) Create(c *gin.Context, md *_owner.OwnerModel[models.AssetShare]) (*_owner.OwnerModel[models.AssetShare], error) {
	profileId := _jwt.GetProfileId(c)
	ok, err := s.Repo.ExistByTarget(c, profileId, md.M.TargetID, md.M.AssetID)
	if ok {
		return nil, &_routes.Except{
			Message: "đã chia sẻ",
		}
	}
	if err != nil {
		return nil, err
	}
	md.OwnerID = profileId
	if err := s.Repo.Create(c, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

// func (s AssetShareService) Update(c *gin.Context, id uint64, md *_owner.OwnerModel[models.AssetShare]) (*_owner.OwnerModel[models.AssetShare], error) {
// 	if _, err := s.RequiredOwner(c, id); err != nil {
// 		return nil, err
// 	}

// 	if err := s.OwnerRepo.Update(c, id, md); err != nil {
// 		return nil, &_routes.Except{
// 			Code:    500,
// 			Message: err.Error(),
// 		}
// 	}

// 	return md, nil
// }
