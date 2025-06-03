package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UProductAccessRoute struct {
	// crud2.BaseRouter[models.ProductAccess, *usecases.ProductAccessUsecase]
	UC *usecases.ProductAccessUsecase
}

func NewUProductAccessRoute(UC *usecases.ProductAccessUsecase) *UProductAccessRoute {
	return &UProductAccessRoute{
		// BaseRouter: crud2.BaseRouter[models.ProductAccess, *ProductAccessService]{
		// 	Service: service,
		// },
		UC: UC,
	}
}

func (r *UProductAccessRoute) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/product-access")
	{
		group.POST("/:productId/bulk", r.BulkSave)
		group.POST("/:productId/batch", r.CreateBatch)
		group.PUT("/commission", r.UpdateCommission)
		// group.PUT(":id", r.BaseRouter.Update)
		group.DELETE(":id", r.Delete)
		// group.GET(":id", r.BaseRouter.GetByID)
		group.GET("", r.Search)
		group.GET("/private-fields", r.GetFields)
	}
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Thông tin Dự án"
// @Success 201 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/v2/user/product-access/{productId}/bulk [post]
func (r *UProductAccessRoute) BulkSave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := dto.ProductAccessBulkRequest{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.BulkSave(c, id, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Thông tin Dự án"
// @Success 201 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/v2/user/product-access/{productId}/batch [post]
func (r *UProductAccessRoute) CreateBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := []dto.ProductAccessRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.CreateBatch(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *UProductAccessRoute) UpdateCommission(c *gin.Context) {

	body := dto.ProductAccessUpdate{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.UC.UpdateCommission(c, body)
	_routes.RouteResult(c, nil, err)
}

// // @Summary Cập nhật thông tin Dự án
// // @Description Cập nhật thông tin Dự án theo ID
// // @Tags User/Product-share
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Param ProductAccess body models.ProductAccess true "Thông tin Dự án mới"
// // @Success 200 {object} models.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/v2/user/product-access/{id} [put]
// func (r *UserProductAccessRoute) Update(c *gin.Context) {}

// @Summary Xóa Dự án
// @Description Xóa Dự án theo ID (soft delete)
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /v1/bdspro/v2/user/product-access/{id} [delete]
func (r *UProductAccessRoute) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// // @Summary Lấy thông tin Dự án theo ID
// // @Description Trả về thông tin chi tiết của một Dự án
// // @Tags User/Product-share
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Success 200 {object} models.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/v2/user/product-access/{id} [get]
// func (r *UserProductAccessRoute) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/v2/user/product-access [get]
func (r *UProductAccessRoute) Search(c *gin.Context) {
	var dto dto.ProductAccessSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	results, err := r.UC.Search(c, dto)
	_routes.RouteResult(c, results, err)
}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/v2/user/product-access/private-fields [get]
func (r *UProductAccessRoute) GetFields(c *gin.Context) {
	// fields := []struct {
	// 	Value string `json:"value"`
	// 	Label string `json:"label"`
	// }{
	// 	{Value: "importPrice", Label: "Giá nhập"},
	// 	{Value: "operatingCost", Label: "Chi phí vận hành"},
	// }

	_routes.RouteResult(c, enums.PrivateFields, nil)
}
