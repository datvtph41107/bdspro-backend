package admin_router

// type DeveloperRouter struct {
// 	crud2.BaseRouter[domain.Developer, *DeveloperService]
// }

// func NewDeveloperRouter(service *DeveloperService) *DeveloperRouter {
// 	return &DeveloperRouter{
// 		BaseRouter: crud2.BaseRouter[domain.Developer, *DeveloperService]{
// 			Service: service,
// 		},
// 	}
// }

// func (r *DeveloperRouter) RegisterRoutes(router *gin.RouterGroup) {
// 	group := router.Group("/developer")
// 	{
// 		group.POST("", r.BaseRouter.Create)
// 		group.PUT(":id", r.BaseRouter.Update)
// 		group.DELETE(":id", r.BaseRouter.Delete)
// 		group.GET(":id", r.BaseRouter.GetByID)
// 		group.GET("", r.BaseRouter.GetAll)
// 	}
// }

// // @Summary Tạo Dự án mới
// // @Description Thêm một Dự án mới vào hệ thống
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param Developer body domain.Developer true "Thông tin Dự án"
// // @Success 201 {object} domain.Developer
// // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// // @Router /v1/bdspro/admin/developer [post]
// func (r *DeveloperRouter) Create(c *gin.Context) {}

// // @Summary Cập nhật thông tin Dự án
// // @Description Cập nhật thông tin Dự án theo ID
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Param Developer body domain.Developer true "Thông tin Dự án mới"
// // @Success 200 {object} domain.Developer
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/admin/developer/{id} [put]
// func (r *DeveloperRouter) Update(c *gin.Context) {}

// // @Summary Xóa Dự án
// // @Description Xóa Dự án theo ID (soft delete)
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Success 200 {object} map[string]bool "Trạng thái xóa"
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// // @Router /v1/bdspro/admin/developer/{id} [delete]
// func (r *DeveloperRouter) Delete(c *gin.Context) {}

// // @Summary Lấy thông tin Dự án theo ID
// // @Description Trả về thông tin chi tiết của một Dự án
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Success 200 {object} domain.Developer
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/admin/developer/{id} [get]
// func (r *DeveloperRouter) GetByID(c *gin.Context) {}

// // @Summary Lấy danh sách Dự án
// // @Description Trả về danh sách tất cả các Dự án
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.Developer
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/admin/developer [get]
// func (r *DeveloperRouter) GetAll(c *gin.Context) {}
