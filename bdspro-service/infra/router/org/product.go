package org_router

import (
	"bdspro/internal/dto"
	org_usecase "bdspro/internal/usecases/org"
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

type OrgProductRouter struct {
	OrgProductUsecase *org_usecase.OrgProductUC
}

func NewOrgProductRouter(uc *org_usecase.OrgProductUC) *OrgProductRouter {
	return &OrgProductRouter{OrgProductUsecase: uc}
}

func (r *OrgProductRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	api := router.Group(path)
	{
		api.GET("/total", r.SearchSmartProduct)
		api.GET("/me", r.SearchProductWithOwner)
		api.GET("/link/:id", r.SearchProductWithOwner)
		api.GET("/public", r.SearchProductWithOwner)
		api.POST("/new", r.CreateNewProduct) //
		api.GET("/suggest", r.SuggestFields) //
		api.GET("/:id", r.Detail)            //
		api.GET("/:id/posts", r.Posts)
		api.GET("/:id/members", r.Members)
		api.GET("/:id/history", r.History)
		api.PUT("/edit/:id", r.Edit)
		// api.PUT("/publish", r.PublishProduct)
		// api.PUT("/revert", r.RevertProduct)

		api.PUT("/archived", r.Archived)
		api.DELETE("/:id", r.Delete)
		api.PUT("/sale-status", r.SaleStatusUpdate)
		api.PUT("/rent-status", r.RentStatusUpdate)
		api.PUT("/sale-visibility", r.SaleVisibilityUpdate)
		api.PUT("/rent-visibility", r.RentVisibilityUpdate)
		api.POST("/deposite", r.Deposite)
		api.DELETE("/deposite/:depositeId", r.CancelDeposite)
		api.POST("/export", r.ExportExcel)
	}
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dto.NewProductRequest true "Thông tin sản phẩm"
// @Router /v1/bdspro/v2/org/product/new [post]
func (route *OrgProductRouter) CreateNewProduct(c *gin.Context) {
	body := dto.ProductSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := route.OrgProductUsecase.CreateProduct(c, &body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của sản phẩm"
// @Success 200 {object} domain.Product "Sản phẩm được tạo thành công"
// @Router /v1/bdspro/v2/org/product/{id} [get]
func (route *OrgProductRouter) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.OrgProductUsecase.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách bài viết liên quan
// @Description Returns a list of posts related to a product
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của sản phẩm"
// @Param post body dto.PostLinkSearch true "Thông tin bài viết"
// @Success 200 {object} domain.Post "Bài viết được tạo thành công"
// @Router /v1/bdspro/v2/org/product/{id}/posts [get]
func (route *OrgProductRouter) Posts(c *gin.Context) {
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

	result, err := route.OrgProductUsecase.Posts(c, id, dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách thành viên của sản phẩm
// @Description Returns a list of members of a product
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của sản phẩm"
// @Param members body dto.ProductAccessSearch true "Thông tin thành viên"
// @Router /v1/bdspro/v2/org/product/{id}/members [get]
func (route *OrgProductRouter) Members(c *gin.Context) {
	var dto dto.ProductAccessSearch
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

	result, err := route.OrgProductUsecase.Members(c, id, dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách lịch sử của sản phẩm
// @Description Returns a list of history of a product
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của sản phẩm"
// @Param history body dto.ProductHistorySearch true "Thông tin lịch sử"
// @Success 200 {object} domain.ProductHistory "Danh sách lịch sử"
// @Router /v1/bdspro/v2/org/product/{id}/history [get]
func (route *OrgProductRouter) History(c *gin.Context) {
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

	result, total, err := route.OrgProductUsecase.History(c, id, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// SearchProductWithOwner godoc
// @Summary Lấy danh sách sản phẩm của user hiện tại
// @Description Returns a list of products based on search filters
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param parentId query integer false "null=>all; 0=>parent; >0=>childrent"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v1/bdspro/v2/org/product/me [get]
func (route *OrgProductRouter) SearchProductWithOwner(c *gin.Context) {
	body := dto.ProductSearchRequest{}
	err := _utils.ParseQuery2(c, &body)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := route.OrgProductUsecase.SearchProductByOwner(c, &body)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Tìm kiếm tổng hợp
// @Description Returns a list of products based on search filters
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query string true "Text"
// @Param size query int false "Size results" default(10)
// @Param page query int false "Page" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Success 200 {object} _routes.ResponseDTO "Danh sách sản phẩm"
// @Router /v1/bdspro/v2/org/product/total [get]
func (route *OrgProductRouter) SearchSmartProduct(c *gin.Context) {
	body := dto.TextSearchRequest{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := route.OrgProductUsecase.SmartSearch(c, body.Text, &body)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Gợi ý các trường thông tin
// @Description Returns a list of products based on search filters
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v1/bdspro/v2/org/product/suggest [get]
func (route *OrgProductRouter) SuggestFields(c *gin.Context) {
	body := dto.ProductSaveDetect{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := route.OrgProductUsecase.SuggestFields(c, &body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Param body body dto.ProductSaveRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/{productId} [put]
func (route *OrgProductRouter) Edit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.OrgProductUsecase.Update(c, id, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Param body body dto.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/archived [put]
func (route *OrgProductRouter) Archived(c *gin.Context) {
	body := dto.ArchivedRequest{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.OrgProductUsecase.Archived(c, body)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Router /v1/bdspro/v2/org/product/{productId} [delete]
func (route *OrgProductRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err = route.OrgProductUsecase.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/rent-status [put]
func (route *OrgProductRouter) RentStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.OrgProductUsecase.RentStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/sale-status [put]
func (route *OrgProductRouter) SaleStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.OrgProductUsecase.SaleStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Trạng thái sản phẩm
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/rent-visibility [put]
func (route *OrgProductRouter) RentVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.OrgProductUsecase.RentVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Trạng thái sản phẩm
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/sale-visibility [put]
func (route *OrgProductRouter) SaleVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.OrgProductUsecase.SaleVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Cọc
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/v2/org/product/deposite [post]
func (route *OrgProductRouter) Deposite(c *gin.Context) {
	dto := dto.DepositeRequest{}

	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	deposite, err := route.OrgProductUsecase.ProductDeposite(c, dto)
	_routes.RouteResult(c, deposite, err)
}

// @Summary Hủy cọc
// @Description API
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param depositeId path int true "Deposite ID"
// @Router /v1/bdspro/v2/org/product/deposite/{depositeId} [delete]
func (route *OrgProductRouter) CancelDeposite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("depositeId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = route.OrgProductUsecase.CancelDeposite(c, id)
	_routes.RouteResult(c, nil, err)
}

// // @Summary Chia sẻ co-broker
// // @Description API
// // @Tags Tổ chức: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/product/access-share [post]
// func (route *OrgProductRouter) Share(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	err = route.UserProductService.CancelDeposite(c, id)
// 	_routes.RouteResult(c, nil, err)
// }

//	func (route *OrgProductRouter) ExportExcel(c *gin.Context) {
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

// @Summary Xuất danh sách sản phẩm ra Excel
// @Description API này dùng để xuất dữ liệu sản phẩm ra file Excel
// @Tags Tổ chức: Sản phẩm
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param filter query string false "Bộ lọc sản phẩm (nếu có)"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} _routes.ResponseDTO
// @Failure 500 {object} _routes.ResponseDTO
// @Router /v1/bdspro/v2/org/product/export [post]
func (route *OrgProductRouter) ExportExcel(c *gin.Context) {
	var query dto.ProductSearchRequest
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	products, err := route.OrgProductUsecase.SearchExport(c, &query)
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
			product.SaleStatus,
			product.SaleVisibility,
			product.RentStatus,
			product.RentVisibility,
			product.OwnerUserID,
			product.OwnerOrgID,
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
			nilIfNil(product.District, func() any { return product.District.Name }),
			nilIfNil(product.Ward, func() any { return product.Ward.Name }),
			nilIfNil(product.DocType, func() any { return product.DocType.Name }),
			nilIfNil(product.PropertyType, func() any { return product.PropertyType.Name }),
			dereference(product.Price.Currency),
			dereference(product.Price.PriceOwner),
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
			nilIfNil(product.Deposite, func() any { return product.Deposite.Amount }),
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
