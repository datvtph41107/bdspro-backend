package u_list

import (
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type ListRoute struct {
	listService *ListService
}

func NewListRoute(listService *ListService) *ListRoute {
	return &ListRoute{
		listService: listService,
	}
}

func (route *ListRoute) RegisterRouteListRoute(r *gin.RouterGroup) {
	api := r.Group("/list")
	api.GET("/property-type", route.ListProductProperty)
	api.GET("/region", route.ListRegion)
	api.GET("/doc-type", route.ListDocType)
	api.GET("/project", route.ListProject)
	api.GET("/amenity", route.ListAmenity)
}

func (route *ListRoute) ListProductProperty(c *gin.Context) {
	result, err := route.listService.ListPropertyType(c)
	_routes.RouteResult(c, result, err)
}

func (route *ListRoute) ListRegion(c *gin.Context) {
	result, err := route.listService.ListRegion(c)
	_routes.RouteResult(c, result, err)
}

func (route *ListRoute) ListDocType(c *gin.Context) {
	result, err := route.listService.ListDocType(c)
	_routes.RouteResult(c, result, err)
}
func (route *ListRoute) ListProject(c *gin.Context) {
	result, err := route.listService.ListProject(c)
	_routes.RouteResult(c, result, err)
}
func (route *ListRoute) ListAmenity(c *gin.Context) {
	result, err := route.listService.ListAmenity(c)
	_routes.RouteResult(c, result, err)
}
