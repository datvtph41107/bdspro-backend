package usecases

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
)

type AssetShareUC struct {
	common.CrudUsecase[domain.AssetShare, *dto.AssetShareSearchDTO]
}

func NewAssetShareUC(repo repo.AssetShareRepo) *AssetShareUC {
	return &AssetShareUC{
		CrudUsecase: common.CrudUsecase[domain.AssetShare, *dto.AssetShareSearchDTO]{
			Repo: repo,
		},
	}
}

// func (s AssetShareUC) Create(c *gin.Context, md *_owner.OwnerModel[domain.AssetShare]) (*_owner.OwnerModel[domain.AssetShare], error) {
// 	profileId := _utils.GetProfileIdWithContext(c)
// 	ok, err := s.Repo.ExistByTarget(c, profileId, md.M.TargetID, md.M.AssetID)
// 	if ok {
// 		return nil, &_utils.Except{
// 			Message: "đã chia sẻ",
// 		}
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	md.OwnerID = profileId
// 	if err := s.Repo.Create(c, md); err != nil {
// 		return nil, &_utils.Except{
// 			Code:    500,
// 			Message: err.Error(),
// 		}
// 	}

// 	return md, nil
// }

// func (s AssetShareService) Update(c *gin.Context, id uint64, md *_owner.OwnerModel[domain.AssetShare]) (*_owner.OwnerModel[domain.AssetShare], error) {
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
