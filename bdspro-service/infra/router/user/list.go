package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type ListRouter struct {
	UC *usecases.ListUsecase
}

func NewListRouter(UC *usecases.ListUsecase) *ListRouter {
	return &ListRouter{
		UC: UC,
	}
}

func (route *ListRouter) RegisterRouter(r *gin.RouterGroup) {
	api := r.Group("/list")
	api.GET("/property-type", route.ListProductProperty)
	api.GET("/region", route.ListRegion)
	api.GET("/doc-type", route.ListDocType)
	api.GET("/project", route.ListProject)
	api.GET("/amenity", route.ListAmenity)
}

// @Summary Lấy danh sách loại tài sản
// @Description Lấy danh sách loại tài sản
// @Tags List
// @Param text query string false "Tên loại tài sản"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Accept json
// @Produce json
// @Success 200 {object} _routes.ResponseDTO{data=[]domain.PropertyType,totalElements=int64}
// @Failure 400 {object} _routes.ResponseDTO{errors=map[string]string}
// @Failure 500 {object} _routes.ResponseDTO{errors=map[string]string}
// @Router /v1/bdspro/v2/user/list/property-type [get]
func (route *ListRouter) ListProductProperty(c *gin.Context) {
	dto := dto.PropertyTypeSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UC.ListPropertyType(c, &dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách khu vực
// @Description Lấy danh sách khu vực
// @Tags List
// @Param text query string false "Tên khu vực"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Accept json
// @Produce json
// @Success 200 {object} _routes.ResponseDTO{data=[]domain.Region,totalElements=int64}
// @Failure 400 {object} _routes.ResponseDTO{errors=map[string]string}
// @Failure 500 {object} _routes.ResponseDTO{errors=map[string]string}
// @Router /v1/bdspro/v2/user/list/region [get]
func (route *ListRouter) ListRegion(c *gin.Context) {
	dto := dto.RegionRequest{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := route.UC.ListRegion(c, &dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Lấy danh sách loại pháp lý
// @Description Lấy danh sách loại pháp lý
// @Tags List
// @Param text query string false "Tên loại pháp lý"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Accept json
// @Produce json
// @Success 200 {object} _routes.ResponseDTO{data=[]domain.DocTypeItem,totalElements=int64}
// @Failure 400 {object} _routes.ResponseDTO{errors=map[string]string}
// @Failure 500 {object} _routes.ResponseDTO{errors=map[string]string}
// @Router /v1/bdspro/v2/user/list/doc-type [get]
func (route *ListRouter) ListDocType(c *gin.Context) {
	dto := dto.DocTypeSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UC.ListDocType(c, &dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code: 0,
		Data: result,
		// TotalElements: &len(result),
	}, err)
}

// @Summary Lấy danh sách dự án
// @Description Lấy danh sách dự án
// @Tags List
// @Param text query string false "Tên dự án"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Accept json
// @Produce json
// @Success 200 {object} _routes.ResponseDTO{data=[]domain.Project,totalElements=int64}
// @Failure 400 {object} _routes.ResponseDTO{errors=map[string]string}
// @Failure 500 {object} _routes.ResponseDTO{errors=map[string]string}
// @Router /v1/bdspro/v2/user/list/project [get]
func (route *ListRouter) ListProject(c *gin.Context) {
	dto := dto.ProjectSearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := route.UC.ListProject(c, &dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Lấy danh sách tiện ích
// @Description Lấy danh sách tiện ích
// @Tags List
// @Param text query string false "Tên tiện ích"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Accept json
// @Produce json
// @Success 200 {object} _routes.ResponseDTO{data=[]domain.Amenity,totalElements=int64}
// @Failure 400 {object} _routes.ResponseDTO{errors=map[string]string}
// @Failure 500 {object} _routes.ResponseDTO{errors=map[string]string}
// @Router /v1/bdspro/v2/user/list/amenity [get]
func (route *ListRouter) ListAmenity(c *gin.Context) {
	dto := dto.AmenitySearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UC.ListAmenity(c, &dto)
	_routes.RouteResult(c, result, err)
}
