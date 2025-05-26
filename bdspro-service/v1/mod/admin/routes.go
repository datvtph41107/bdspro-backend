package admin

import (
	a_amenity "bdspro/v1/mod/admin/amenity"
	a_apartment "bdspro/v1/mod/admin/apartment"
	a_build_room "bdspro/v1/mod/admin/build_room"
	a_developer "bdspro/v1/mod/admin/developer"
	"bdspro/v1/mod/admin/district"
	a_doc_type "bdspro/v1/mod/admin/doc_type"
	a_product "bdspro/v1/mod/admin/product"
	a_product_access "bdspro/v1/mod/admin/product_access"
	a_project "bdspro/v1/mod/admin/project"
	a_project_build "bdspro/v1/mod/admin/project_build"
	a_property_type "bdspro/v1/mod/admin/property_type"
	"bdspro/v1/mod/admin/province"
	"bdspro/v1/mod/admin/ward"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type AdminRoute struct {
	province      *province.ProvinceRouter
	district      *district.DistrictRouter
	ward          *ward.WardRouter
	product       *a_product.AdminProductRouter
	propertyType  *a_property_type.PropertyTypeRouter
	docType       *a_doc_type.DocTypeRouter
	amenity       *a_amenity.AmenityRouter
	project       *a_project.ProjectRouter
	projectBuild  *a_project_build.ProjectBuildRouter
	buildRoom     *a_build_room.BuildRoomRouter
	apartment     *a_apartment.ApartmentRouter
	developer     *a_developer.DeveloperRouter
	productAccess *a_product_access.ProductAccessRouter
}

func NewAdminRoute(province *province.ProvinceRouter,
	district *district.DistrictRouter,
	ward *ward.WardRouter,
	product *a_product.AdminProductRouter,
	propertyType *a_property_type.PropertyTypeRouter,
	docType *a_doc_type.DocTypeRouter,
	amenity *a_amenity.AmenityRouter,
	project *a_project.ProjectRouter,
	projectBuild *a_project_build.ProjectBuildRouter,
	buildRoom *a_build_room.BuildRoomRouter,
	apartment *a_apartment.ApartmentRouter,
	developer *a_developer.DeveloperRouter,
	productAccess *a_product_access.ProductAccessRouter,
) *AdminRoute {
	return &AdminRoute{
		province:      province,
		district:      district,
		ward:          ward,
		product:       product,
		propertyType:  propertyType,
		docType:       docType,
		amenity:       amenity,
		project:       project,
		projectBuild:  projectBuild,
		buildRoom:     buildRoom,
		apartment:     apartment,
		developer:     developer,
		productAccess: productAccess,
	}
}

func (r *AdminRoute) RegisterAdminRouter(e *gin.Engine) {
	adminRouter := e.Group("/bdspro/v1/admin")
	adminRouter.Use(_jwt.JWTAuthMiddleware(nil))
	r.product.RegisterRoutes(adminRouter)
	r.province.RegisterRoutes(adminRouter)
	r.district.RegisterRoutes(adminRouter)
	r.ward.RegisterRoutes(adminRouter)
	r.propertyType.RegisterRoutes(adminRouter)
	r.docType.RegisterRoutes(adminRouter)
	r.amenity.RegisterRoutes(adminRouter)
	r.project.RegisterRoutes(adminRouter)
	r.projectBuild.RegisterRoutes(adminRouter)
	r.buildRoom.RegisterRoutes(adminRouter)
	r.apartment.RegisterRoutes(adminRouter)
	r.developer.RegisterRoutes(adminRouter)
	r.productAccess.RegisterRoutes(adminRouter)
}
