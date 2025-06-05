package router

import (
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactRouter struct {
	service *usecase.ContactUsecase
}

func NewContactRouter(service *usecase.ContactUsecase) *ContactRouter {
	return &ContactRouter{
		service: service,
	}
}

// Khai báo API
func ContactRoutes(r *gin.Engine, router *ContactRouter) {
	api := r.Group("/crm/contacts")

	api.Use(_jwt.JWTAuthMiddleware(nil))
	api.GET("/user", router.Search)
	// api.GET("/org/:id", router.Get)
	api.POST("", router.Create)
	api.PUT("/:id", router.Update)
	api.DELETE("/:id", router.Delete)
}

// @Summary Lấy danh sách liên hệ
// @Description Lấy danh sách
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /contact [get]
// Lấy danh sách liên hệ
func (r *ContactRouter) Search(c *gin.Context) {
	dto := dto.ContactSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := r.service.Search(c, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Tạo liên hệ mới
// @Description Tạo một liên hệ mới với thông tin được cung cấp
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.ContactEntity true "Thông tin cập nhật"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 409 {string} string "Liên hệ đã tồn tại"
// @Router /contact [post]
func (r *ContactRouter) Create(c *gin.Context) {
	entity := domain.ContactEntity{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.service.Create(c, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật thông tin liên hệ
// @Description Cập nhật thông tin liên hệ hiện tại
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param request body domain.GroupEntity true "Thông tin cập nhật"
// @Param id path int true "ID của lời mời kết bạn"
// @Security BearerAuth
// @Router /contact/{id} [put]
func (r *ContactRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	entity := domain.ContactEntity{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.service.Update(c, id, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa liên hệ
// @Description Xóa một liên hệ khỏi hệ thống
// @Tags Liên hệ
// @Accept json
// @Produce json
// @Param id path int true "ID của liên hệ"
// @Security BearerAuth
// @Router /contact/{id} [delete]
func (r *ContactRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	err = r.service.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
