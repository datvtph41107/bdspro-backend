package routes

import (
	_dto "common/domain/dto"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"
	"user/internal/dto"
	"user/internal/usecases"

	"github.com/gin-gonic/gin"
)

type ProfileRouter struct {
	Profile *usecases.ProfileUsecase
}

func NewProfileRouter(profileUsecase *usecases.ProfileUsecase) *ProfileRouter {
	return &ProfileRouter{
		Profile: profileUsecase,
	}
}

func ProfileRoutes(r *gin.Engine, router *ProfileRouter) {
	profileGroup := r.Group("/user/profile")
	{
		profileGroup.GET("/find-by-phone", router.FindByPhone)
		profileGroup.GET("/chat/find", router.FindOnChat)
		profileGroup.GET("/find/ids", router.FindByIds)

		// profileGroup.GET("/me", router.ProfileMe)
		// profileGroup.PUT("/set-info", router.SetInfo)
		profileGroup.GET("/info/:profileId", router.PublicInfo)

		// tìm kiếm user
		profileGroup.GET("/find/global", router.FindWithFriend)
		profileGroup.GET("/find/friend", router.FindWithFriend)

		// mock data
		profileGroup.GET("/find/group", router.FindWithGroup)
		profileGroup.GET("/find/organization", router.FindWithOrg)

	}
}

// ProfileMe godoc
// @Summary Lấy thông tin cá nhân của user
// @Description API này trả về thông tin của user hiện tại
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/user/profile/find-by-phone [get]
func (p *ProfileRouter) FindByPhone(c *gin.Context) {
	result, err := p.Profile.ProfileGetByPhone(c)

	_routes.RouteResult(c, result, err)
}

// ProfileMe godoc
// @Summary Lấy thông tin cá nhân của user
// @Description API này trả về thông tin của user hiện tại
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/user/profile/find-by-phone [get]
func (p *ProfileRouter) FindWithFriend(c *gin.Context) {
	var dto dto.ProfileSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := p.Profile.SearchUserByFriend(c, dto)

	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

func (p *ProfileRouter) FindWithGroup(c *gin.Context) {
	var dto _dto.Pagable
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, _routes.ResponseDTO{
			Code:    0,
			Message: err.Error(),
		}, nil)
		return
	}
	if dto.GetPage() > 0 {
		total := int64(0)
		_routes.RouteResult(c, _routes.ResponseDTO{
			Code:          0,
			Data:          nil,
			TotalElements: &total,
		}, nil)
		return
	}
	data := []struct {
		GroupID uint64 `json:"groupId"`
		Name    string `json:"name"`
		Avatar  string `json:"avatar"`
	}{
		{
			GroupID: 1,
			Name:    "Nguồn hàng BĐS Sơn Tây",
		}, {
			GroupID: 2,
			Name:    "Nguồn hàng BĐS Ba vì",
		}, {
			GroupID: 3,
			Name:    "Hội mua bán tập thể cũ",
		}, {
			GroupID: 4,
			Name:    "Nhà đất hà nội",
		}, {
			GroupID: 5,
			Name:    "Cộng đồng môi giới",
		}, {
			GroupID: 6,
			Name:    "Đất rẻ dưới 1 tỷ",
		},
	}
	total := int64(3)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          data,
		TotalElements: &total,
	}, nil)
}

func (p *ProfileRouter) FindWithOrg(c *gin.Context) {
	var dto _dto.Pagable
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, _routes.ResponseDTO{
			Code:    0,
			Message: err.Error(),
		}, nil)
		return
	}
	if dto.GetPage() > 0 {
		total := int64(0)
		_routes.RouteResult(c, _routes.ResponseDTO{
			Code:          0,
			Data:          nil,
			TotalElements: &total,
		}, nil)
		return
	}
	data := []struct {
		OrganizationID uint64 `json:"organizationId"`
		Name           string `json:"name"`
		Avatar         string `json:"avatar"`
	}{
		{
			OrganizationID: 1,
			Name:           "Công ty Cổ phần Vinhomes – tập đoàn Vingroup",
		}, {
			OrganizationID: 2,
			Name:           "Công ty CP Tập đoàn Ecopark",
		}, {
			OrganizationID: 3,
			Name:           "Tập đoàn Sunshine Group",
		}, {
			OrganizationID: 4,
			Name:           "Tập đoàn đầu tư bất động sản Hà Đô",
		}, {
			OrganizationID: 5,
			Name:           "CTY CP BDSPro",
		}, {
			OrganizationID: 6,
			Name:           "Công ty BĐS CenGroup",
		},
	}
	total := int64(3)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          data,
		TotalElements: &total,
	}, nil)
}

// ProfileMe godoc
// @Summary Lấy thông tin cá nhân của user
// @Description API này trả về thông tin của user hiện tại
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query string true "ID của Profile"
// @Success 200 {object} dto.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/user/profile/chat/find [get]
func (p *ProfileRouter) FindOnChat(c *gin.Context) {
	result, err := p.Profile.ProfileChatFind(c)

	_routes.RouteResult(c, result, err)
}

// ProfileMe godoc
// @Summary Lấy thông tin cá nhân của user
// @Description API này trả về thông tin của user hiện tại
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ids query string true "ID của Profile"
// @Success 200 {object} dto.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/user/profile/find/ids [get]
func (p *ProfileRouter) FindByIds(c *gin.Context) {
	result, err := p.Profile.ProfileFindByIds(c)

	_routes.RouteResult(c, result, err)
}

func (p *ProfileRouter) Find(c *gin.Context) {
	result, err := p.Profile.ProfileGetByPhone(c)

	_routes.RouteResult(c, result, err)
}

// // ProfileMe godoc
// // @Summary Lấy thông tin cá nhân của user
// // @Description API này trả về thông tin của user hiện tại
// // @Tags Profile
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {object} dto.UserInfoResponse
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /v1/user/profile/me [get]
// func (p *ProfileRouter) ProfileMe(c *gin.Context) {
// 	result, err := p.Profile.GetInfoMe(c)

// 	_routes.RouteResult(c, result, err)
// }

// // UpdateUserInfo cập nhật thông tin hồ sơ
// // @Summary Cập nhật thông tin hồ sơ
// // @Description API này cho phép người dùng cập nhật thông tin hồ sơ
// // @Tags Profile
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param request body dto.UserInfoRequest true "Thông tin cập nhật"
// // @Success 200 {object} dto.UserInfoResponse "Cập nhật thành công"
// // @Router /v1/user/profile/set-info [put]
// func (p *ProfileRouter) SetInfo(c *gin.Context) {
// 	var dto dto.UserInfoRequest
// 	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	result, err := p.Profile.UpdateUserInfo(c, &dto)
// 	_routes.RouteResult(c, result, err)
// }

// UpdateUserInfo Thông tin hồ sơ
// @Summary Xem thông tin public của người dùng
// @Description API này gọi không cần Token. dùng trong màn preview profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path string true "ID của Profile"
// @Success 200 {object} dto.PublicUserInfo "Cập nhật thành công"
// @Router /v1/user/profile/info/{profileId} [get]
func (p *ProfileRouter) PublicInfo(c *gin.Context) {
	profileId, _ := strconv.ParseUint(c.Param("profileId"), 10, 64)
	result, err := p.Profile.PublicUserInfo(c, profileId, false)
	_routes.RouteResult(c, result, err)
}
