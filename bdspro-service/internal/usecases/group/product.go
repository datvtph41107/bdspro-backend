package group_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	org_dto "bdspro/internal/dto/org"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	user_repo "bdspro/internal/repo/user"
	"bdspro/internal/usecases"
	iusecase "bdspro/internal/usecases/interface"
	shared_usecase "bdspro/internal/usecases/shared"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type GroupProductUsecase struct {
	*shared_usecase.ProductUsecase
	ProductAccessRepo repo.ProductAccessRepo
	ClientRepo        iusecase.IClientUC
}

func NewGroupProductUsecase(
	productRepo user_repo.UProductRepo,
	ProductHistoryRepo repo.ProductHistoryRepo,
	assetRepo user_repo.UAssetRepo,
	mediaRepo repo.ProductMediaRepo,
	productAccessRepo repo.ProductAccessRepo,
	houseInfoRepo repo.HouseInfoRepo,
	priceRepo repo.ProductPriceRepo,
	productPrivate repo.ProductPrivateRepo,
	region repo.RegionRepo,
	docTypeRepo repo.DocTypeRepo,
	propertyType repo.PropertyTypeRepo,
	amenity repo.AmenityRepo,
	memberPlan usecases.PlanUsecase,
	depositeRepo repo.DepositeRepo,
	clientRepo iusecase.IClientUC,
) *GroupProductUsecase {
	uc := &GroupProductUsecase{
		ProductUsecase: shared_usecase.MakeProductUsecase(
			productRepo,
			ProductHistoryRepo,
			assetRepo,
			mediaRepo,
			productAccessRepo,
			houseInfoRepo,
			priceRepo,
			productPrivate,
			region,
			docTypeRepo,
			propertyType,
			amenity,
			memberPlan,
			nil,
			depositeRepo,
			enums.EOwnerTypeGroup,
		),
		ClientRepo: clientRepo,
	}

	return uc
}

func (u *GroupProductUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	return nil, nil
}

func (u *GroupProductUsecase) GetOwnerType(c context.Context) enums.EOwnerType {
	return enums.EOwnerTypeGroup
}

func (u *GroupProductUsecase) HasPermission(c context.Context, productId *uint64, permission enums.EPermission) (bool, error) {
	return true, nil
}

func (u *GroupProductUsecase) GetProduct(c context.Context, groupId uint64, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	userInGroup, err := u.ClientRepo.UserInGroup(c, profileId, groupId)
	if err != nil {
		return nil, 0, err
	}
	if !userInGroup {
		return nil, 0, &_routes.Except{
			Code:    401,
			Message: "Bạn không trong nhóm này",
		}
	}
	return u.ProductUsecase.SearchByProductAccess(c, &groupId, enums.EOwnerTypeGroup, &dto)
}

func (u *GroupProductUsecase) ShareProduct(c context.Context, dto org_dto.ProductShareRequest) (*domain.ProductAccess, error) {
	_, err := u.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	userInGroup, err := u.ClientRepo.UserInGroup(c, profileId, dto.GroupID)
	if err != nil {
		return nil, err
	}
	if !userInGroup {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không trong nhóm này",
		}
	}

	productAccess := domain.ProductAccess{
		ProductID:  dto.ProductID,
		TargetId:   dto.GroupID,
		TargetType: enums.EOwnerTypeGroup,
	}

	err = _db.SaveWithContext(c, &productAccess)
	if err != nil {
		return nil, err
	}

	return &productAccess, nil
}
