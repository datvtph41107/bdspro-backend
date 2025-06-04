package user_router

// type UPropertyTypeRouter struct {
// 	UC *usecases.PropertyTypeUsecase
// }

// type PropertyTypeRouter struct {
// 	crud.BaseRouter[domain.PropertyType]
// }

// func NewPropertyTypeRouter(service *PropertyTypeService) *PropertyTypeRouter {
// 	return &PropertyTypeRouter{
// 		BaseRouter: crud.BaseRouter[domain.PropertyType]{
// 			Service: service.BaseService,
// 		},
// 	}
// }

// func (r *PropertyTypeRouter) RegisterRoutes(router *gin.RouterGroup) {
// 	PropertyTypeGroup := router.Group("/property-type")
// 	{
// 		PropertyTypeGroup.POST("", r.BaseRouter.Create)
// 		PropertyTypeGroup.PUT(":id", r.BaseRouter.Update)
// 		PropertyTypeGroup.DELETE(":id", r.BaseRouter.Delete)
// 		PropertyTypeGroup.GET(":id", r.BaseRouter.GetByID)
// 		PropertyTypeGroup.GET("", r.BaseRouter.GetAll)
// 	}
// }

// // @Summary Tạo Loại bất động sản mới
// // @Description Thêm một Loại bất động sản mới vào hệ thống
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param PropertyType body domain.PropertyType true "Thông tin Loại bất động sản"
// // @Success 201 {object} domain.PropertyType
// // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// // @Router /v1/bdspro/admin/property-type [post]
// func (r *PropertyTypeRouter) Create(c *gin.Context) {}

// // @Summary Cập nhật thông tin Loại bất động sản
// // @Description Cập nhật thông tin Loại bất động sản theo ID
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Loại bất động sản"
// // @Param PropertyType body domain.PropertyType true "Thông tin Loại bất động sản mới"
// // @Success 200 {object} domain.PropertyType
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/admin/property-type/{id} [put]
// func (r *PropertyTypeRouter) Update(c *gin.Context) {}

// // @Summary Xóa Loại bất động sản
// // @Description Xóa Loại bất động sản theo ID (soft delete)
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Loại bất động sản"
// // @Success 200 {object} map[string]bool "Trạng thái xóa"
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// // @Router /v1/bdspro/admin/property-type/{id} [delete]
// func (r *PropertyTypeRouter) Delete(c *gin.Context) {}

// // @Summary Lấy thông tin Loại bất động sản theo ID
// // @Description Trả về thông tin chi tiết của một Loại bất động sản
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Loại bất động sản"
// // @Success 200 {object} domain.PropertyType
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/admin/property-type/{id} [get]
// func (r *PropertyTypeRouter) GetByID(c *gin.Context) {}

// // @Summary Lấy danh sách Loại bất động sản
// // @Description Trả về danh sách tất cả các Loại bất động sản
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.PropertyType
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/admin/property-type [get]
// func (r *PropertyTypeRouter) GetAll(c *gin.Context) {}
