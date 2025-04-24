package u_asset

// import (
// 	"bdspro/v1/models"
// 	"common/case/crud"

// 	"github.com/gin-gonic/gin"
// )

// type AssetRouter struct {
// 	crud.BaseRouter[models.Asset]
// }

// func NewAssetRouter(service *AssetService) *AssetRouter {
// 	return &AssetRouter{
// 		BaseRouter: crud.BaseRouter[models.Asset]{
// 			Service: service.BaseService,
// 		},
// 	}
// }

// func (r *AssetRouter) RegisterRoutes(router *gin.Engine) {
// 	assetGroup := router.Group("/bdspro/v1/assets")
// 	{
// 		assetGroup.POST("", r.BaseRouter.Create)
// 		assetGroup.PUT(":id", r.BaseRouter.Update)
// 		assetGroup.DELETE(":id", r.BaseRouter.Delete)
// 		assetGroup.GET(":id", r.BaseRouter.GetByID)
// 		assetGroup.GET("", r.BaseRouter.GetAll)
// 	}
// }

// // @BasePath /api

// // @Summary Tạo tỉnh mới
// // @Description Tạo một bản ghi tỉnh mới
// // @Tags assets
// // @Accept json
// // @Produce json
// // @Param Ward body models.Ward true "Thông tin tỉnh"
// // @Success 201 {object} models.Ward
// // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// // @Router /assets [post]
// func (r *WardRouter) create(c *gin.Context) {

// }

// // @Summary Cập nhật thông tin tỉnh
// // @Description Cập nhật thông tin tỉnh theo ID
// // @Tags assets
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của tỉnh"
// // @Param Ward body models.Ward true "Thông tin tỉnh mới"
// // @Success 200 {object} models.Ward
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /assets/{id} [put]
// func (r *WardRouter) update(c *gin.Context) {

// }

// // @Summary Xóa tỉnh
// // @Description Xóa tỉnh theo ID (soft delete)
// // @Tags assets
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của tỉnh"
// // @Success 200 {object} map[string]bool "Trạng thái xóa"
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// // @Router /assets/{id} [delete]
// func (r *WardRouter) delete(c *gin.Context) {

// }

// // @Summary Lấy thông tin tỉnh theo ID
// // @Description Lấy thông tin chi tiết của một tỉnh
// // @Tags assets
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của tỉnh"
// // @Success 200 {object} models.Ward
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /assets/{id} [get]
// func (r *WardRouter) getByID(c *gin.Context) {

// }

// // @Summary Lấy danh sách các tỉnh
// // @Description Trả về danh sách tất cả các tỉnh
// // @Tags assets
// // @Accept json
// // @Produce json
// // @Success 200 {array} models.Ward
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /assets [get]
// func (r *WardRouter) getAll(c *gin.Context) {

// }
