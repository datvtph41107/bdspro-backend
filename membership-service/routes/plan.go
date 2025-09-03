package routes

import (
	"common/case/crud"
	_jwt "common/jwt"
	"net/http"

	"membership/models"
	"membership/services"

	"github.com/gin-gonic/gin"
)

type PlanRouter struct {
	crud.BaseRouter[models.Plan, crud.ICrudRepo[models.Plan]]
}

func NewPlanRouter(service *services.PlanService) *PlanRouter {
	return &PlanRouter{
		BaseRouter: crud.BaseRouter[models.Plan, crud.ICrudRepo[models.Plan]]{
			Service: service.BaseUsecase,
		},
	}
}

func (r *PlanRouter) RegisterRoutes(router *gin.Engine) {
	PlanGroup := router.Group("/membership/plan")
	PlanGroup.Use(_jwt.JWTAuthMiddleware([]string{"/membership/plan/list"}))
	{
		PlanGroup.POST("", r.BaseRouter.Create)
		PlanGroup.PUT(":id", r.BaseRouter.Update)
		PlanGroup.DELETE(":id", r.BaseRouter.Delete)
		PlanGroup.GET(":id", r.BaseRouter.GetByID)
		PlanGroup.GET("", r.BaseRouter.GetAll)
		PlanGroup.GET("/list", r.List)
	}
}

// @BasePath /membership

// @Summary Tạo Gói dịch vụ mới
// @Description Thêm một Gói dịch vụ mới vào hệ thống
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan body models.Plan true "Thông tin Gói dịch vụ"
// @Success 201 {object} models.Plan
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /plan [post]
func (r *PlanRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Gói dịch vụ
// @Description Cập nhật thông tin Gói dịch vụ theo ID
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Gói dịch vụ"
// @Param plan body models.Plan true "Thông tin Gói dịch vụ mới"
// @Success 200 {object} models.Plan
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /plan/{id} [put]
func (r *PlanRouter) Detail(c *gin.Context) {}

// @Summary Xóa Gói dịch vụ
// @Description Xóa Gói dịch vụ theo ID (soft delete)
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Gói dịch vụ"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /plan/{id} [delete]
func (r *PlanRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Gói dịch vụ theo ID
// @Description Trả về thông tin chi tiết của một Gói dịch vụ
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Gói dịch vụ"
// @Success 200 {object} models.Plan
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /plan/{id} [get]
func (r *PlanRouter) Get(c *gin.Context) {}

// @Summary Lấy danh sách Gói dịch vụ
// @Description Trả về danh sách tất cả các Gói dịch vụ
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Plan
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /plan [get]
func (r *PlanRouter) Search(c *gin.Context) {}

// @Summary Lấy danh sách Gói dịch vụ
// @Description Trả về danh sách tất cả các Gói dịch vụ
// @Tags Gói dịch vụ
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Plan
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /plan/list [get]
func (r *PlanRouter) List(c *gin.Context) {
	result, _ := r.Service.GetAll(c)
	c.JSON(http.StatusOK, result)
}
