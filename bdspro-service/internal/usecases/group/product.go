package group_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	org_dto "bdspro/internal/dto/org"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	shared_usecase "bdspro/internal/usecases/shared"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type GroupProductUsecase struct {
	*shared_usecase.ProductUsecase
	ProductAccessRepo repo.ProductAccessRepo
	ClientRepo        repo.IClientRepo
}

func NewGroupProductUsecase(
	productUsecase *shared_usecase.ProductUsecase,
	clientRepo repo.IClientRepo,
) *GroupProductUsecase {
	uc := &GroupProductUsecase{}
	clone := *productUsecase
	newPtr := &clone

	newPtr.ProductUsecase = uc
	uc.ProductUsecase = newPtr
	uc.ClientRepo = clientRepo

	return uc
}

func (u *GroupProductUsecase) GetOwnerID(c context.Context) (*uint64, error) {
	return nil, nil
}

func (u *GroupProductUsecase) GetOwnerType(c context.Context) enums.ETargetType {
	return enums.ETargetTypeGroup
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
	return u.ProductUsecase.SearchByProductAccess(c, &groupId, enums.ETargetTypeGroup, &dto)
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
		TargetType: enums.ETargetTypeGroup,
	}

	err = _db.SaveWithContext(c, &productAccess)
	if err != nil {
		return nil, err
	}

	return &productAccess, nil
}
