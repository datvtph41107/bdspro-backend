package a_apartment

import (
	"bdspro/v1/models"
	"common/case/crud2"
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type ApartmentRouter struct {
	crud2.BaseRouter[models.Apartment, *ApartmentService]
}

func NewApartmentRouter(service *ApartmentService) *ApartmentRouter {
	return &ApartmentRouter{
		BaseRouter: crud2.BaseRouter[models.Apartment, *ApartmentService]{
			Service: service,
		},
	}
}

func (r *ApartmentRouter) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/apartment")
	{
		group.POST("", r.UpdateApartments)
		group.PUT("/status", r.UpdateStatusApartments)
		group.POST("/import", r.ImportExcel)
		// group.DELETE(":id", r.BaseRouter.Delete)
		group.GET("", r.GetAll)
	}
}

// @Summary Tạo Tòa nhà trong dự án mới
// @Description Thêm một Tòa nhà trong dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Apartments body data.ApartmentRequest true "Thông tin Tòa nhà trong dự án"
// @Success 201 {object} data.ApartmentRequest
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/apartment [post]
func (r *ApartmentRouter) UpdateApartments(c *gin.Context) {
	result, err := r.Service.UpdateApartments(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo Tòa nhà trong dự án mới
// @Description Thêm một Tòa nhà trong dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Apartments body data.ApartmentRequest true "Thông tin Tòa nhà trong dự án"
// @Success 201 {object} data.ApartmentRequest
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/apartment/status [put]
func (r *ApartmentRouter) UpdateStatusApartments(c *gin.Context) {
	result, err := r.Service.UpdateStatusApartments(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách Tòa nhà trong dự án
// @Description Trả về danh sách tất cả các Tòa nhà trong dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Apartment
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/apartment/import [post]
func (r *ApartmentRouter) ImportExcel(c *gin.Context) {
	result, err := r.Service.ImportApartments(c)
	_routes.RouteResult(c, result, err)
}
