package org_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrgAccessRoute struct {
	uc *usecases.OrgMemAccessUsecase
}

func NewOrgAccessRoute(uc *usecases.OrgMemAccessUsecase) *OrgAccessRoute {
	return &OrgAccessRoute{uc: uc}
}

func (r *OrgAccessRoute) RegisterRoutes(router *gin.RouterGroup) {
	accessGroup := router.Group("/product-access")
	{
		accessGroup.GET("/:org_id", r.GetProductAccessByOrgID)
	}
}

// @Summary Get product access by org ID
// @Description Lấy danh sách sản phẩm được chia sẻ của tổ chức
// @Tags Org/Product-access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param org_id path int true "Org ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param productId query int false "Product ID"
// @Success 200 {object} []domain.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /v1/bdspro/v2/org/product-access/{org_id} [get]
func (r *OrgAccessRoute) GetProductAccessByOrgID(c *gin.Context) {
	orgID, err := strconv.ParseInt(c.Param("org_id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	var dto dto.ProductAccessSearch
	if err := c.ShouldBindQuery(&dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	productAccess, err := r.uc.GetProductAccessByOrgID(c, orgID, dto)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	_routes.RouteResult(c, productAccess, nil)
}
