package org_router

type OrgAccessRoute struct {
	// uc *usecases.OrgMemAccessUsecase
}

// func NewOrgAccessRoute(
// // uc *usecases.OrgMemAccessUsecase,
// ) *OrgAccessRoute {
// 	return &OrgAccessRoute{
// 		// uc: uc
// 	}
// }

// func (r *OrgAccessRoute) RegisterRoutes(router *gin.RouterGroup) {
// 	accessGroup := router.Group("/product-access")
// 	{
// 		accessGroup.GET("/:org_id", r.GetProductAccessByOrgID)
// 	}
// }

// // @Summary Get product access by org ID
// // @Description Lấy danh sách sản phẩm được chia sẻ của tổ chức
// // @Tags Org/Product-access
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param org_id path int true "Org ID"
// // @Param page query int false "Page"
// // @Param limit query int false "Limit"
// // @Param productId query int false "Product ID"
// // @Success 200 {object} []domain.ProductAccess
// // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // @Failure 500 {object} map[string]string "Lỗi server"
// // @Router /v1/bdspro/v2/org/product-access/{org_id} [get]
// func (r *OrgAccessRoute) GetProductAccessByOrgID(c *gin.Context) {
// 	orgID, err := strconv.ParseInt(c.Param("org_id"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	var dto dto.ProductAccessSearch
// 	if err := c.ShouldBindQuery(&dto); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	productAccess, err := r.uc.GetProductAccessByOrgID(c, orgID, dto)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	_routes.RouteResult(c, productAccess, nil)
// }

// // func (r *OrgAccessRoute) RegisterRoutes(router *gin.RouterGroup) {
// // 	group := router.Group("/product-access")
// // 	{
// // 		group.POST("", r.BaseRouter.Create)
// // 		group.PUT(":id", r.BaseRouter.Update)
// // 		group.DELETE(":id", r.BaseRouter.Delete)
// // 		group.GET(":id", r.BaseRouter.GetByID)
// // 		group.GET("", r.BaseRouter.GetAll)
// // 		group.GET("/private-fields", r.GetFields)
// // 	}
// // }

// // @Summary Tạo Dự án mới
// // @Description Thêm một Dự án mới vào hệ thống
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param ProductAccess body domain.ProductAccess true "Thông tin Dự án"
// // @Success 201 {object} domain.ProductAccess
// // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// // @Router /v1/bdspro/admin/product-access [post]
// func (r *OrgAccessRoute) Create(c *gin.Context) {}

// // @Summary Cập nhật thông tin Dự án
// // @Description Cập nhật thông tin Dự án theo ID
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Param ProductAccess body domain.ProductAccess true "Thông tin Dự án mới"
// // @Success 200 {object} domain.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/admin/product-access/{id} [put]
// func (r *OrgAccessRoute) Update(c *gin.Context) {}

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
// // @Router /v1/bdspro/admin/product-access/{id} [delete]
// func (r *OrgAccessRoute) Delete(c *gin.Context) {}

// // @Summary Lấy thông tin Dự án theo ID
// // @Description Trả về thông tin chi tiết của một Dự án
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Success 200 {object} domain.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/admin/product-access/{id} [get]
// func (r *OrgAccessRoute) GetByID(c *gin.Context) {}

// // @Summary Lấy danh sách Dự án
// // @Description Trả về danh sách tất cả các Dự án
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.ProductAccess
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/admin/product-access [get]
// func (r *OrgAccessRoute) GetAll(c *gin.Context) {}

// // @Summary Lấy danh sách Dự án
// // @Description Trả về danh sách tất cả các Dự án
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.ProductAccess
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/admin/product-access/private-fields [get]
// func (r *OrgAccessRoute) GetFields(c *gin.Context) {
// 	// fields := []struct {
// 	// 	Value string `json:"value"`
// 	// 	Label string `json:"label"`
// 	// }{
// 	// 	{Value: "importPrice", Label: "Giá nhập"},
// 	// 	{Value: "operatingCost", Label: "Chi phí vận hành"},
// 	// }

// 	_routes.RouteResult(c, enums.PrivateFields, nil)
// }
