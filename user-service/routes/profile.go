package routes

import (
	_dto "common/dto"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"encoding/json"
	"net/http"
	"user/data"
	auth_service "user/services/auth"

	"github.com/gin-gonic/gin"
)

type ProfileRouter struct {
	Profile *auth_service.ProfileService
}

func NewProfileRouter(profileService *auth_service.ProfileService) *ProfileRouter {
	return &ProfileRouter{
		Profile: profileService,
	}
}

var publicRoutes = []string{"/user/profile/info/", "/user/profile/find-by-phone"}

func ProfileRoutes(r *gin.Engine, router *ProfileRouter) {
	profileGroup := r.Group("/user/profile")
	{
		profileGroup.Use(_jwt.JWTAuthMiddleware(publicRoutes))
		profileGroup.GET("/find-by-phone", router.FindByPhone)
		profileGroup.GET("/chat/find", router.FindOnChat)
		profileGroup.GET("/find/ids", router.FindByIds)

		profileGroup.GET("/me", router.ProfileMe)
		profileGroup.PUT("/set-info", router.SetInfo)
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
// @Success 200 {object} data.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/find-by-phone [get]
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
// @Success 200 {object} data.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/find-by-phone [get]
func (p *ProfileRouter) FindWithFriend(c *gin.Context) {
	var dto data.ProfileSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := p.Profile.SearchUserByFriend(c, dto)

	_routes.RouteResult(c, result, err)
}

func (p *ProfileRouter) FindWithGroup(c *gin.Context) {
	var dto _dto.Page
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
	var dto _dto.Page
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
// @Success 200 {object} data.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/chat/find [get]
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
// @Success 200 {object} data.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/find/ids [get]
func (p *ProfileRouter) FindByIds(c *gin.Context) {
	result, err := p.Profile.ProfileFindByIds(c)

	_routes.RouteResult(c, result, err)
}

func (p *ProfileRouter) Find(c *gin.Context) {
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
// @Success 200 {object} data.UserInfoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/me [get]
func (p *ProfileRouter) ProfileMe(c *gin.Context) {
	result, err := p.Profile.GetInfoMe(c)

	_routes.RouteResult(c, result, err)
}

// UpdateUserInfo cập nhật thông tin người dùng
// @Summary Cập nhật thông tin người dùng
// @Description API này cho phép người dùng cập nhật thông tin cá nhân
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body data.UserInfoRequest true "Thông tin cập nhật"
// @Success 200 {object} data.UserInfoResponse "Cập nhật thành công"
// @Router /profile/set-info [put]
func (p *ProfileRouter) SetInfo(c *gin.Context) {
	var dto data.UserInfoRequest
	r, _ := json.Marshal(&c.Request.Body)
	print(r)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := p.Profile.UpdateUserInfo(c, dto)
	_routes.RouteResult(c, result, err)
}

// UpdateUserInfo Thông tin người dùng
// @Summary Thông tin người dùng
// @Description API này cho phép người dùng cập nhật thông tin cá nhân
// @Tags Profile
// @Accept json
// @Produce json
// @Param profileId path string true "ID của Profile"
// @Success 200 {object} data.PublicUserInfo "Cập nhật thành công"
// @Router /profile/info/{profileId} [get]
func (p *ProfileRouter) PublicInfo(c *gin.Context) {
	result, err := p.Profile.PublicUserInfo(c)
	_routes.RouteResult(c, result, err)
}
