package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProductAccessRouter struct {
	// crud2.BaseRouter[domain.ProductAccess, *usecases.ProductAccessUsecase]
	Usecase *shared_usecase.SharingAccessUsecase
}

func NewUProductAccessRouter(Usecase *shared_usecase.SharingAccessUsecase) *UserProductAccessRouter {
	return &UserProductAccessRouter{
		// BaseRouter: crud2.BaseRouter[domain.ProductAccess, *ProductAccessService]{
		// 	Service: service,
		// },
		Usecase: Usecase,
	}
}

func (r *UserProductAccessRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	group := router.Group(path)
	{
		group.POST("/:productId/bulk", r.BulkSave)
		// group.POST("/:productId/batch", r.CreateBatch)
		group.PUT("/commission", r.UpdateCommission)
		// group.DELETE(":id", r.Delete)
		group.GET("", r.Search)
		group.GET("/private-fields", r.GetFields)
		// group.PUT(":id", r.BaseRouter.Update)
		// group.GET(":id", r.BaseRouter.GetByID)
	}
}

// .
// // @Summary Tạo chia sẻ mới sản phẩm cho người dùng
// // @Description Thêm một sản phẩm mới vào hệ thống
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param productId path int true "Thông tin sản phẩm"
// // @Success 201 {object} domain.ProductAccess
// // @Router /v1/bdspro/v2/user/product/access/{productId}/bulk [post]
func (r *UserProductAccessRouter) BulkSave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := dto.SharingAccessBulk{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Usecase.BulkSave(c, id, enums.EOwnerOfMember, &body)
	_routes.RouteResult(c, result, err)
}

// // .
// // // @Summary Tạo chia sẻ mới sản phẩm cho người dùng
// // // @Description Thêm một sản phẩm mới vào hệ thống
// // // @Tags User: Sản phẩm
// // // @Accept json
// // // @Produce json
// // // @Security BearerAuth
// // // @Param productId path int true "Thông tin sản phẩm"
// // // @Success 201 {object} domain.ProductAccess
// // // @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// // // @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// // // @Router /v1/bdspro/v2/user/product/access/{productId}/batch [post]
// func (r *UserProductAccessRouter) CreateBatch(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	body := []dto.SharingAccessRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UC.CreateBatch(c, id, enums.EOwnerOfMember, body)
// 	_routes.RouteResult(c, result, err)
// }

// .
// // @Summary Cập nhật hoa hồng
// // @Description Cập nhật hoa hồng của sản phẩm
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.ProductAccessUpdate true "Thông tin hoa hồng"
// // @Success 200 {object} domain.ProductAccess
// // @Router /v1/bdspro/v2/user/product/access/commission [put]
func (r *UserProductAccessRouter) UpdateCommission(c *gin.Context) {

	body := dto.CommissionUpdate{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.Usecase.UpdateCommission(c, body)
	_routes.RouteResult(c, nil, err)
}

// // @Summary Cập nhật thông tin sản phẩm
// // @Description Cập nhật thông tin sản phẩm theo ID
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của sản phẩm"
// // @Param ProductAccess body domain.ProductAccess true "Thông tin sản phẩm mới"
// // @Success 200 {object} domain.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/v2/user/product-access/{id} [put]
// func (r *UserProductAccessRoute) Update(c *gin.Context) {}

// // .
// // // @Summary Xóa sản phẩm
// // // @Description Xóa sản phẩm theo ID (soft delete)
// // // @Tags User: Sản phẩm
// // // @Accept json
// // // @Produce json
// // // @Security BearerAuth
// // // @Param id path int true "ID của sản phẩm"
// // // @Success 200 {object} map[string]bool "Trạng thái xóa"
// // // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // // @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// // // @Router /v1/bdspro/v2/user/product/access/{id} [delete]
// func (r *UserProductAccessRouter) Delete(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	err = r.UC.ProductDelete(c, dto.SharingAccessDeleted{
// 		DomainID: id,
// 	})
// 	_routes.RouteResult(c, nil, err)
// }

// // @Summary Lấy thông tin Dự án theo ID
// // @Description Trả về thông tin chi tiết của một Dự án
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của sản phẩm"
// // @Success 200 {object} domain.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/v2/user/product-access/{id} [get]
// func (r *UserProductAccessRoute) GetByID(c *gin.Context) {}

// .
// // @Summary Lấy danh sách sản phẩm
// // @Description Trả về danh sách tất cả các sản phẩm
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.ProductAccess
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/v2/user/product/access [get]
func (r *UserProductAccessRouter) Search(c *gin.Context) {
	var dto dto.SharingAccessSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	results, err := r.Usecase.SearchForProduct(c, enums.EOwnerOfMember, &dto)
	_routes.RouteResult(c, results, err)
}

// .
// // @Summary Lấy danh sách sản phẩm
// // @Description Trả về danh sách tất cả các sản phẩm
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {array} domain.ProductAccess
// // @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// // @Router /v1/bdspro/v2/user/product/access/private-fields [get]
func (r *UserProductAccessRouter) GetFields(c *gin.Context) {
	// fields := []struct {
	// 	Value string `json:"value"`
	// 	Label string `json:"label"`
	// }{
	// 	{Value: "importPrice", Label: "Giá nhập"},
	// 	{Value: "operatingCost", Label: "Chi phí vận hành"},
	// }

	_routes.RouteResult(c, enums.PrivateFields, nil)
}
