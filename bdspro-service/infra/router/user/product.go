package user_router

import (
	"bdspro/internal/dto"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	_utils "common/utils"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UserProductRouter struct {
	UC *shared_usecase.ProductUsecase
	// UProductProvider *providers.UProductProvider
}

func NewUserProductRouter(product *shared_usecase.ProductUsecase) *UserProductRouter {
	return &UserProductRouter{
		UC: product,
	}
}

func (route *UserProductRouter) RegisterRoutes(r *gin.RouterGroup, path string) {
	api := r.Group(path)
	api.GET("/market", route.SearchProductMarket)
	api.GET("/total", route.SearchSmartProduct)

	api.GET("/me", route.SearchProductWithOwner)
	api.GET("/link/:id", route.SearchProductWithOwner)
	api.GET("/public", route.SearchProductWithOwner)
	api.POST("/new", route.CreateNewProduct)
	api.GET("/suggest", route.SuggestProduct)
	api.GET("/:id", route.Detail)
	api.GET("/:id/posts", route.Posts)
	api.GET("/:id/members", route.Members)
	api.GET("/:id/history", route.History)
	api.PUT("/edit/:id", route.Edit)
	// api.PUT("/publish", route.PublishProduct)
	// api.PUT("/revert", route.RevertProduct)

	api.PUT("/archived", route.Archived)
	api.DELETE("/:id", route.Delete)
	api.PUT("/sale-status", route.SaleStatusUpdate)
	api.PUT("/rent-status", route.RentStatusUpdate)
	api.PUT("/sale-visibility", route.SaleVisibilityUpdate)
	api.PUT("/rent-visibility", route.RentVisibilityUpdate)
	// api.POST("/deposite", route.Deposite)
	api.DELETE("/deposite/:depositeId", route.CancelDeposite)
	api.POST("/export", route.ExportExcel)

	// api.DELETE("/:id", route.Delete)
}

// .
// // @Summary Tạo sản phẩm mới
// // @Description API để tạo một sản phẩm mới
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param product body dto.ProductSaveRequest true "Thông tin sản phẩm"
// // @Success 201 {object} dto.CreateResponse "Sản phẩm được tạo thành công"
// // @Router /v1/bdspro/v2/user/product/new [post]
func (route *UserProductRouter) CreateNewProduct(c *gin.Context) {
	body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UC.CreateProduct(c, body)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Tạo sản phẩm mới
// // @Description API để tạo một sản phẩm mới
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param product body dto.ProductSaveRequest true "Thông tin sản phẩm"
// // @Success 201 {object} dto.CreateResponse "Sản phẩm được tạo thành công"
// // @Router /v1/bdspro/v2/user/product/{id} [post]
func (route *UserProductRouter) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UC.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Lấy danh sách bài viết liên quan
// // @Description Returns a list of posts related to a product
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Router /v1/bdspro/v2/user/product/{id}/posts [get]
func (route *UserProductRouter) Posts(c *gin.Context) {
	var dto dto.PostLinkSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UC.Posts(c, id, dto)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Lấy danh sách thành viên của sản phẩm
// // @Description Returns a list of members of a product
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Router /v1/bdspro/v2/user/product/{id}/members [get]
func (route *UserProductRouter) Members(c *gin.Context) {
	var dto dto.SharingAccessSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UC.Members(c, id, dto)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Lấy danh sách lịch sử của sản phẩm
// // @Description Returns a list of history of a product
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Router /v1/bdspro/v2/user/product/{id}/history [get]
func (route *UserProductRouter) History(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	dto := dto.ProductHistorySearch{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, total, err := route.UC.History(c, id, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// .
// // SearchProductWithOwner godoc
// // @Summary Lấy danh sách sản phẩm của user hiện tại
// // @Description Returns a list of products based on search filters
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param name query string false "Product name"
// // @Param code query string false "Product code"
// // @Param area query number false "Product area"
// // @Param price query integer false "Product price"
// // @Param parentId query integer false "null=>all; 0=>parent; >0=>childrent"
// // @Param description query string false "Product description"
// // @Param note query string false "Additional notes"
// // @Param limit query int false "Limit results" default(10)
// // @Param skip query int false "Skip results" default(0)
// // @Param sort query string false "Sort order" default("price_asc")
// // @Router /v1/bdspro/v2/user/product/me [get]
func (route *UserProductRouter) SearchProductWithOwner(c *gin.Context) {
	body := dto.ProductSearchRequest{}
	err := _utils.ParseQuery2(c, &body)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := route.UC.GetProductOfCurrentUser(c, body)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// .
// // @Summary Lấy danh sách sản phẩm trang chủ
// // @Description Returns a list of products based on search filters
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Param name query string false "Product name"
// // @Param code query string false "Product code"
// // @Param area query number false "Product area"
// // @Param price query integer false "Product price"
// // @Param description query string false "Product description"
// // @Param note query string false "Additional notes"
// // @Param size query int false "Size results" default(10)
// // @Param page query int false "Page" default(0)
// // @Param sort query string false "Sort order" default("price_asc")
// // @Router /v1/bdspro/v2/user/product/market [get]
func (route *UserProductRouter) SearchProductMarket(c *gin.Context) {
	// _routes.RouteResult(c, nil, nil)
	body, _ := _utils.ParseQueryToStruct[dto.ProductSearchRequest](c.Request.URL.Query())
	result, total, err := route.UC.SearchProductMarket(c, body)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// .
// // @Summary Tìm kiếm tổng hợp
// // @Description Returns a list of products based on search filters
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param text query string true "Text"
// // @Param size query int false "Sizet results" default(10)
// // @Param page query int false "Page" default(0)
// // @Param sort query string false "Sort order" default("price_asc")
// // @Router /v1/bdspro/v2/user/product/total [get]
func (route *UserProductRouter) SearchSmartProduct(c *gin.Context) {
	query := dto.TextSearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := route.UC.SmartSearch(c, query.Text, &query)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// .
// // @Summary Gợi ý các trường thông tin
// // @Description Returns a list of products based on search filters
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param name query string false "Product name"
// // @Param code query string false "Product code"
// // @Param area query number false "Product area"
// // @Param price query integer false "Product price"
// // @Param description query string false "Product description"
// // @Param note query string false "Additional notes"
// // @Param limit query int false "Limit results" default(10)
// // @Param skip query int false "Skip results" default(0)
// // @Param sort query string false "Sort order" default("price_asc")
// // @Router /v1/bdspro/v2/user/product/suggest [get]
func (route *UserProductRouter) SuggestProduct(c *gin.Context) {
	_body, _ := _utils.ParseBody[dto.ProductSaveDetect](c)
	result, err := route.UC.SuggestFields(c, _body.Content)
	_routes.RouteResult(c, result, err)
}

// .
// // @Summary Cập nhật sản phẩm
// // @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Param productId path int true "Product ID"
// // @Security BearerAuth
// // @Param body body dto.ProductSaveRequest true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/{productId} [put]
func (route *UserProductRouter) Edit(c *gin.Context) {
	// id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }
	// result, err := route.UC.Update(c, id, body)
	// _routes.RouteResult(c, result, err)
}

// .
// // @Summary Cập nhật sản phẩm
// // @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Param productId path int true "Product ID"
// // @Security BearerAuth
// // @Param body body dto.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/archived [put]
func (route *UserProductRouter) Archived(c *gin.Context) {
	body := dto.ArchivedRequest{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UC.Archived(c, body)
	_routes.RouteResult(c, nil, err)
}

// .
// // @Summary Xóa sản phẩm
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Param productId path int true "Product ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/product/{productId} [delete]
func (route *UserProductRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err = route.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// .
// // @Summary Xóa sản phẩm
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/rent-status [put]
func (route *UserProductRouter) RentStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UC.RentStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// .
// // @Summary Xóa sản phẩm
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/sale-status [put]
func (route *UserProductRouter) SaleStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UC.SaleStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// .
// // @Summary Trạng thái sản phẩm
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/rent-visibility [put]
func (route *UserProductRouter) RentVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UC.RentVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// .
// // @Summary Trạng thái sản phẩm
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// // @Router /v1/bdspro/v2/user/product/sale-visibility [put]
func (route *UserProductRouter) SaleVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UC.SaleVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// // .
// // // @Summary Cọc
// // // @Description API
// // // @Tags User: Sản phẩm
// // // @Accept json
// // // @Produce json
// // // @Security BearerAuth
// // // @Param body body dto.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// // // @Router /v1/bdspro/v2/user/product/deposite [post]
// func (route *UserProductRouter) Deposite(c *gin.Context) {
// 	dto := &domain.Deposite{}

// 	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	deposite, err := route.UC.ProductDeposite(c, dto)
// 	_routes.RouteResult(c, deposite, err)
// }

// .
// // @Summary Hủy cọc
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param depositeId path int true "Deposite ID"
// // @Router /v1/bdspro/v2/user/product/deposite/{depositeId} [delete]
func (route *UserProductRouter) CancelDeposite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("depositeId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = route.UC.CancelDeposite(c, id)
	_routes.RouteResult(c, nil, err)
}

// // @Summary Chia sẻ co-broker
// // @Description API
// // @Tags User
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/product/access-share [post]
// func (route *UProductRouter) Share(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	err = route.UserProductService.CancelDeposite(c, id)
// 	_routes.RouteResult(c, nil, err)
// }

//	func (route *UProductRouter) ExportExcel(c *gin.Context) {
//		body := dto.ProductSearchRequest{}
//		if err := _utils.ParseQuery2(c, &body); err != nil {
//			_routes.RouteResult(c, nil, err)
//			return
//		}
//		products, total, err := route.UProductProvider.ExportExcel(c, body)
//		_routes.RouteResult(c, _routes.ResponseDTO{
//			Code:          0,
//			Data:          products,
//			TotalElements: &total,
//		}, err)
//	}
func nilIfNil[T any](ptr *T, getter func() any) any {
	if ptr == nil {
		return ""
	}
	return getter()
}

func dereference(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		return val.Elem().Interface()
	}
	return value
}

// .
// // @Summary Xuất danh sách sản phẩm ra Excel
// // @Description API này dùng để xuất dữ liệu sản phẩm ra file Excel
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce application/octet-stream
// // @Security BearerAuth
// // @Param filter query string false "Bộ lọc sản phẩm (nếu có)"
// // @Success 200 {file} file "Excel file"
// // @Failure 400 {object} _routes.ResponseDTO
// // @Failure 500 {object} _routes.ResponseDTO
// // @Router /v1/bdspro/v2/user/product/export [post]
func (route *UserProductRouter) ExportExcel(c *gin.Context) {
	var query dto.ProductSearchRequest
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	products, err := route.UC.SearchExport(c, &query)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	// Dữ liệu mẫu
	// users := []dto.ProductSaveRequest{
	// 	{Name: "Nguyễn Văn A", Area: 25},
	// 	{Name: "Trần Thị B", Area: 30},
	// }

	// Tạo file Excel
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Ghi tiêu đề
	// f.SetCellValue(sheet, "A1", "Họ tên")
	// f.SetCellValue(sheet, "B1", "Tuổi")

	// Ghi tiêu đề cột
	headers := []string{
		"ID", "Ngày tạo", "Ngày cập nhật", "Tên", "Mã", "Danh mục", "Diện tích", "Mô tả", "Loại giao dịch",
		"Trạng thái bán", "Hiển thị bán", "Trạng thái cho thuê", "Hiển thị thuê", "ID người sở hữu", "ID tổ chức sở hữu",
		"Link vị trí", "Địa chỉ", "Link Google Maps", "ID tài sản", "ID căn hộ", "Giá giao dịch", "Ngày xóa", "Đã lưu trữ",
		"Số phòng ngủ", "Số phòng tắm", "Số tầng", "Nội thất", "Hướng", "Mặt tiền", "Chỗ đậu xe", "Số toilet",
		"Tiện ích", "Ảnh", "Loại ảnh", "Ảnh chính",
		"Tỉnh/Thành", "Quận/Huyện", "Phường/Xã",
		"Loại giấy tờ", "Loại bất động sản",
		"Loại tiền", "Chủ sở hữu giá", "Giá bán", "Hoa hồng bán", "Loại hoa hồng bán",
		"Giá thuê", "Hoa hồng thuê", "Loại hoa hồng thuê", "Chu kỳ thanh toán thuê",
		"Giá nhập", "Chi phí vận hành", "Ghi chú nội bộ", "Lợi nhuận kỳ vọng", "Tài liệu riêng", "Tiền đặt cọc",
	}

	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}

	// Ghi dữ liệu sản phẩm
	for i, product := range products {
		row := i + 2
		mediaUrls := make([]string, len(product.MediaList))
		for i, media := range product.MediaList {
			mediaUrls[i] = media.MediaURL
		}
		joinedMediaUrls := strings.Join(mediaUrls, ",")

		amenities := make([]string, len(product.Amenities))
		for i, amenity := range product.Amenities {
			amenities[i] = amenity.Name
		}
		joinedAmenities := strings.Join(amenities, ",")

		values := []interface{}{
			product.ID,
			product.CreatedAt,
			product.UpdatedAt,
			product.Name,
			product.Code,
			product.CategoryID,
			product.Area,
			product.Description,
			product.TransactionType,
			nil,
			product.SaleVisibility,
			nil,
			product.RentVisibility,
			product.OwnerID,
			product.OwnerOf,
			product.PositionUrl,
			product.Address,
			product.GoogleMapLink,
			product.AssetId,
			product.ApartmentID,
			nil, //product.Price.TransactionPrice,
			product.DeletedAt,
			product.Archived,
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumBedroom) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumBathroom) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumFloor) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.Furniture) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.Orientation) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumFront) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumCarPark) }),
			nilIfNil(product.HouseInfo, func() any { return dereference(product.HouseInfo.NumToilet) }),
			joinedAmenities, //nilIfNil(product.Amenities, func() any { return product.Amenities }),
			joinedMediaUrls, //strings.Join(product.MediaList, ","),
			nil,             //strings.Join(product.MediaTypes, ","),
			nil,             //product.MediaList,
			nilIfNil(product.Province, func() any { return product.Province.Name }),
			// nilIfNil(product.District, func() any { return product.District.Name }),
			nilIfNil(product.Ward, func() any { return product.Ward.Name }),
			nilIfNil(product.DocType, func() any { return product.DocType.Name }),
			nilIfNil(product.PropertyType, func() any { return product.PropertyType.Name }),
			dereference(product.Price.Currency),
			nil,
			dereference(product.Price.SalePrice),
			product.Price.SaleCommission,
			nil, //product.Price.SaleCommissionType,
			dereference(product.Price.RentPrice),
			dereference(product.Price.RentCommission),
			nil, //product.Price.RentCommissionType,
			dereference(product.Price.RentPaymentCycle),
			nil, //product.Price.ImportPrice,
			nil, //product.Price.OperatingCost,
			nil, //product.Price.InternalNote,
			nil, //product.Price.TargetProfit,
			nil, //product.PrivateDocs,
			nil,
		}

		for j, val := range values {
			col, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), val)
		}
	}

	// Set active sheet
	sheetIndex, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(sheetIndex)

	// Header cho việc download
	filename := fmt.Sprintf("users_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// Ghi file vào response
	if err := f.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Lỗi khi tạo file Excel: %v", err)
	}
}

func toStr(i int) string {
	return fmt.Sprintf("%d", i)
}
