package routes

import (
	"auth/data"
	"auth/services"
	_jwt "common/jwt"
	_routes "common/routes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileRouter struct {
	Profile *services.ProfileService
}

func NewProfileRouter(profileService *services.ProfileService) *ProfileRouter {
	return &ProfileRouter{
		Profile: profileService,
	}
}

var publicRoutes = []string{"/user/profile/info/"}

func ProfileRoutes(r *gin.Engine, router *ProfileRouter) {
	profileGroup := r.Group("/user/profile")
	{
		profileGroup.Use(_jwt.JWTAuthMiddleware(publicRoutes))
		profileGroup.GET("/me", router.ProfileMe)
		profileGroup.PUT("/set-info", router.SetInfo)
		profileGroup.GET("/info/:profileId", router.PublicInfo)
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
