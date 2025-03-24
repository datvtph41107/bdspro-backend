package admin

import (
	"bdspro/mod/admin/district"
	a_doc_type "bdspro/mod/admin/doc_type"
	a_property_type "bdspro/mod/admin/property_type"
	"bdspro/mod/admin/province"
	"bdspro/mod/admin/ward"
	u_product "bdspro/mod/user/product"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type AdminRoute struct {
	province     *province.ProvinceRouter
	district     *district.DistrictRouter
	ward         *ward.WardRouter
	product      *u_product.UserProductRoute
	propertyType *a_property_type.PropertyTypeRouter
	docType      *a_doc_type.DocTypeRouter
}

func NewAdminRoute(province *province.ProvinceRouter,
	district *district.DistrictRouter,
	ward *ward.WardRouter,
	product *u_product.UserProductRoute,
	propertyType *a_property_type.PropertyTypeRouter,
	docType *a_doc_type.DocTypeRouter,
) *AdminRoute {
	return &AdminRoute{
		province:     province,
		district:     district,
		ward:         ward,
		product:      product,
		propertyType: propertyType,
		docType:      docType,
	}
}

func (r *AdminRoute) RegisterAdminRouter(e *gin.Engine) {
	adminRouter := e.Group("/bdspro/admin")
	adminRouter.Use(_jwt.JWTAuthMiddleware(nil))
	r.product.RegisterRouteUserProduct(adminRouter)
	r.province.RegisterRoutes(adminRouter)
	r.district.RegisterRoutes(adminRouter)
	r.ward.RegisterRoutes(adminRouter)
	r.propertyType.RegisterRoutes(adminRouter)
	r.docType.RegisterRoutes(adminRouter)
}
