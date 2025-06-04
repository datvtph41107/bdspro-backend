package user_usecase

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// func (uc *UserCase) CheckPermission(c context.Context, id uint64) error {
// 	profileId := _utils.GetProfileIdWithContext(c)
// 	assetCostType, err := uc.repo.GetByID(c, id)
// 	if err != nil {
// 		return err
// 	}

// 	if assetCostType.OwnerID != profileId || assetCostType.OwnerType != enums.EOwnerTypeMember {
// 		return &_routes.Except{
// 			Code:    403,
// 			Message: "Bạn không có quyền truy cập",
// 		}
// 	}
// 	return nil
// }
