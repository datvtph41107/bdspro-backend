package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
	"crm/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerRouter struct {
	UC *usecase.CustomerUsecase
}

func NewCustomerRouter(UC *usecase.CustomerUsecase) *CustomerRouter {
	return &CustomerRouter{UC: UC}
}

func (r *CustomerRouter) CustomerRoutes(engine *gin.RouterGroup) {
	group := engine.Group("/customer")
	group.GET("/list/:ownerId/:ownerType", r.Search)
	group.POST("", r.Create)
	group.PUT("/:id", r.Update)
	group.DELETE("/:id", r.Delete)
	group.GET("/:id", r.GetByID)
	group.PUT("/:id/note", r.UpdateNote)
	group.PUT("/:id/assign", r.Assign)
	group.PUT("/:id/stage", r.SwitchStage)
}

// Lấy danh sách liên hệ
func (r *CustomerRouter) Search(c *gin.Context) {
	dto := dto.CustomerSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	ownerId, ownerType, err := utils.GetOwnerFromContext(c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.Search(c, ownerId, ownerType, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Tạo khách hàng mới
// @Description Tạo một khách hàng mới với thông tin được cung cấp
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CustomerEntity true "Thông tin cập nhật"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 409 {string} string "Khách hàng đã tồn tại"
// @Router /customer [post]
// Thêm khách hàng
func (r *CustomerRouter) Create(c *gin.Context) {
	entity := domain.CustomerEntity{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Create(c, &entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật thông tin khách hàng
// @Description Cập nhật thông tin khách hàng hiện tại
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param request body domain.CustomerEntity true "Thông tin cập nhật"
// @Param id path int true "ID của khách hàng"
// @Security BearerAuth
// @Router /customer/{id} [put]
func (r *CustomerRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	entity := domain.CustomerEntity{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Update(c, id, &entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa khách hàng
// @Description Xóa một khách hàng khỏi hệ thống
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID của khách hàng"
// @Security BearerAuth
// @Router /customer/{id} [delete]
func (r *CustomerRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// @Summary Cập nhật ghi chú khách hàng
// @Description Cập nhật ghi chú khách hàng hiện tại
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID của khách hàng"
// @Param request body dto.CustomerDTO true "Thông tin cập nhật"
// @Security BearerAuth
// @Router /customer/{id}/note [put]
func (r *CustomerRouter) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	dto := dto.CustomerDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.UpdateNote(c, id, &dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Gán khách hàng
// @Description Gán khách hàng hiện tại
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID của khách hàng"
// @Param request body dto.CustomerDTO true "Thông tin cập nhật"
// @Security BearerAuth
// @Router /customer/{id}/assign [put]
func (r *CustomerRouter) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	dto := dto.CustomerDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Assign(c, id, &dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Chuyển giai đoạn khách hàng
// @Description Chuyển giai đoạn khách hàng hiện tại
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID của khách hàng"
// @Param request body dto.CustomerDTO true "Thông tin cập nhật"
// @Security BearerAuth
// @Router /customer/{id}/stage [put]
func (r *CustomerRouter) SwitchStage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	dto := dto.CustomerDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.SwitchStage(c, id, dto.StageID)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy thông tin khách hàng
// @Description Lấy thông tin khách hàng hiện tại
// @Tags Khách hàng
// @Accept json
// @Produce json
// @Param id path int true "ID của khách hàng"
// @Security BearerAuth
// @Router /customer/{id} [get]
func (r *CustomerRouter) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	result, err := r.UC.GetByID(c, id)
	_routes.RouteResult(c, result, err)
}
