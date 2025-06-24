package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/dto"
	"crm/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeadRouter struct {
	UC *usecase.LeadUsecase
}

func NewLeadRouter(UC *usecase.LeadUsecase) *LeadRouter {
	return &LeadRouter{UC: UC}
}

func (r *LeadRouter) LeadRoutes(engine *gin.RouterGroup) {
	group := engine.Group("/lead")
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
func (r *LeadRouter) Search(c *gin.Context) {
	dto := dto.LeadSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	// ownerId, ownerType, err := utils.GetOwnerFromContext(c)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }

	result, total, err := r.UC.Search(c, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// .
// // @Summary Tạo khách hàng mới
// // @Description Tạo một khách hàng mới với thông tin được cung cấp
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param request body domain.LeadEntity true "Thông tin cập nhật"
// // @Failure 400 {string} string "Yêu cầu không hợp lệ"
// // @Failure 409 {string} string "Khách hàng đã tồn tại"
// // @Router /lead [post]
// // Thêm khách hàng
func (r *LeadRouter) Create(c *gin.Context) {
	entity := dto.LeadSaveDTO{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Create(c, &entity)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Cập nhật thông tin khách hàng
// // @Description Cập nhật thông tin khách hàng hiện tại
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param request body domain.LeadEntity true "Thông tin cập nhật"
// // @Param id path int true "ID của khách hàng"
// // @Security BearerAuth
// // @Router /lead/{id} [put]
func (r *LeadRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	entity := dto.LeadSaveDTO{}
	if err := _utils.ParseBodyWithValidator(c, &entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Update(c, id, &entity)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Xóa khách hàng
// // @Description Xóa một khách hàng khỏi hệ thống
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của khách hàng"
// // @Security BearerAuth
// // @Router /lead/{id} [delete]
func (r *LeadRouter) Delete(c *gin.Context) {
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

// .
// // @Summary Cập nhật ghi chú khách hàng
// // @Description Cập nhật ghi chú khách hàng hiện tại
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của khách hàng"
// // @Param request body dto.LeadDTO true "Thông tin cập nhật"
// // @Security BearerAuth
// // @Router /lead/{id}/note [put]
func (r *LeadRouter) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	dto := dto.LeadDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.UpdateNote(c, id, &dto)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Gán khách hàng
// // @Description Gán khách hàng hiện tại
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của khách hàng"
// // @Param request body dto.LeadDTO true "Thông tin cập nhật"
// // @Security BearerAuth
// // @Router /lead/{id}/assign [put]
func (r *LeadRouter) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}
	dto := dto.LeadDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Assign(c, id, &dto)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Chuyển giai đoạn khách hàng
// // @Description Chuyển giai đoạn khách hàng hiện tại
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của khách hàng"
// // @Param request body dto.LeadDTO true "Thông tin cập nhật"
// // @Security BearerAuth
// // @Router /lead/{id}/stage [put]
func (r *LeadRouter) SwitchStage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	dto := dto.LeadDTO{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.SwitchStage(c, id, dto.StageID)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Lấy thông tin khách hàng
// // @Description Lấy thông tin khách hàng hiện tại
// // @Tags Khách hàng
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của khách hàng"
// // @Security BearerAuth
// // @Router /lead/{id} [get]
func (r *LeadRouter) GetByID(c *gin.Context) {
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
