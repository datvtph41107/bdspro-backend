package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	shared_usecase "bdspro/internal/usecases/shared"
	_dto "common/domain/dto"
	"context"
	bdspropb "pb/types/bdspro"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	bdspropb.UnimplementedProductServiceServer
	UC               *shared_usecase.ProductUsecase
	ProductRepo      repo.ProductRepo
	ProductMapper    *mapper.ProductMapper
	ProductValidator *validator.ProductValidator
	UserClient       *client.UserClient
	DepositeMapper   *mapper.DepositeMapper
	// TransactionMapper mapper.TransactionTransformer
	PostMapper *mapper.PostMapper
}

func NewProductService(uc *shared_usecase.ProductUsecase,
	productRepo repo.ProductRepo,
	productMapper *mapper.ProductMapper,
	productValidator *validator.ProductValidator,
	userClient *client.UserClient,
	depositeMapper *mapper.DepositeMapper,
	// transactionMapper mapper.TransactionTransformer,
	postMapper *mapper.PostMapper,
) *ProductService {
	return &ProductService{
		UC:               uc,
		ProductRepo:      productRepo,
		ProductMapper:    productMapper,
		ProductValidator: productValidator,
		UserClient:       userClient,
		DepositeMapper:   depositeMapper,
		// TransactionMapper: transactionMapper,
		PostMapper: postMapper,
	}
}

// @Summary Lưu trữ sản phẩm
// @Description API này cho phép lưu trữ sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/archived [put]
func (s *ProductService) ArchiveProduct(c context.Context, req *bdspropb.ArchivedRequest) (*bdspropb.Response, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	err := s.UC.Archived(c, dto.ArchivedRequest{
		ProductId: &req.Id,
		Archived:  req.Archived,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.SaveChildRequest true "Thông tin sản phẩm con"
// @Router /v2/bdspro/v2/product/child/devide [post]
func (s *ProductService) ChildDevideProduct(ctx context.Context, req *bdspropb.DevideChildRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.MapDevideChildRequest(req)
	if err := s.ProductValidator.ValidateDevideChildRequest(dto); err != nil {
		return nil, err
	}

	_, err := s.UC.DevideProductChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// // @Summary Lấy lịch sử sản phẩm con
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Param body query dto.ProductHistorySearch true "Thông tin sản phẩm con"
// // @Router /v2/bdspro/v2/product/child/history/{id} [get]
// func (s *ProductService) ChildHistories(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.HistoryResponse, error) {
// 	id := req.Id

// 	_body := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&_body, req); err != nil {
// 		return nil, err
// 	}

// 	result, total, err := s.UC.ChildHistory(ctx, id, &_body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	histories := make([]*bdspropb.History, 0)
// 	copier.Copy(&histories, &result)
// 	return &bdspropb.GetHistoryResponse{
// 		Data:  histories,
// 		Total: uint64(total),
// 	}, nil
// }

// @Summary Gộp sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.MergeProductChildRequest true "Thông tin sản phẩm con"
// @Router /v2/bdspro/v2/product/child/merge [post]
func (s *ProductService) ChildMergeProduct(ctx context.Context, req *bdspropb.MergeChildRequest) (*bdspropb.Response, error) {
	_, err := s.UC.MergeProductChild(ctx, &dto.MergeProductChildRequest{
		ParentID: &req.ParentId,
		ChildIDs: &req.ChildIds,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Gộp toàn bộ sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/child/merge-all/{id} [post]
func (s *ProductService) ChildMergeAll(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	_, err := s.UC.MergeAllProductChild(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param body body bdspropb.SaveChildRequest true "Thông tin sản phẩm con"
// @Security BearerAuth
// @Router /v2/bdspro/v2/product/child/split [post]
func (s *ProductService) ChildSplit(ctx context.Context, req *bdspropb.SaveChildRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.SaveProductToEntity(req)
	if err := s.ProductValidator.ValidateSaveProductChild(dto); err != nil {
		return nil, err
	}
	_, err := s.UC.CreateChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/deposite [post]
func (s *ProductService) ProductDeposite(ctx context.Context, req *bdspropb.DepositeRequest) (*bdspropb.Transaction, error) {
	// dto := s.DepositeMapper.MapDepositePbToDomain(req)

	// if dto.DepositeAmount == nil || *dto.DepositeAmount <= 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "deposite amount is required")
	// }

	// if dto.ContactID == nil || *dto.ContactID == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "contactId is required")
	// }

	// if dto.TransactionType != enums.TransactionTypeSale && dto.TransactionType != enums.TransactionTypeRent {
	// 	return nil, status.Errorf(codes.InvalidArgument, "transactionType is invalid")
	// }

	// result, err := s.UC.ProductDeposite(ctx, dto)
	// if err != nil {
	// 	return nil, err
	// }
	// transaction := s.TransactionMapper.EntityToTransactionResponse(result)
	// return transaction, nil
	return nil, nil
}

// @Summary Hủy cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transactionId path int true "Transaction ID"
// @Router /v2/bdspro/v2/product/deposite/{transactionId} [delete]
func (s *ProductService) DeleteDeposite(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.UC.CancelDeposite(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Xuất danh sách sản phẩm ra Excel
// @Description API này dùng để xuất dữ liệu sản phẩm ra file Excel
// @Tags User: Sản phẩm
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param filter query string false "Bộ lọc sản phẩm (nếu có)"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} _routes.ResponseDTO
// @Failure 500 {object} _routes.ResponseDTO
// @Router /v2/bdspro/v2/product/export [post]
func (s *ProductService) ExportProduct(context.Context, *bdspropb.IdRequest) (*bdspropb.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportProduct not implemented")
}

// @Summary Lấy danh sách sản phẩm của user hiện tại
// @Description Returns a list of products based on search filters
// @Tags User: Sản phẩm
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
// @Router /v2/bdspro/v2/product/me [get]
func (s *ProductService) ProductMe(ctx context.Context, req *bdspropb.ProductSearch) (*bdspropb.SearchResponse, error) {
	body := &dto.ProductSearchRequest{}
	if err := copier.Copy(body, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.SearchProductByOwner(ctx, body)
	if err != nil {
		return nil, err
	}

	products := s.ProductMapper.MapProductPbList(result)

	// products := make([]*bdspropb.ProductResponse, 0)
	// copier.Copy(&products, &result)
	return &bdspropb.SearchResponse{
		Data:          products,
		TotalElements: total,
	}, nil
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dto.ProductSaveRequest true "Thông tin sản phẩm"
// @Success 201 {object} dto.CreateResponse "Sản phẩm được tạo thành công"
// @Router /v2/bdspro/v2/product/new [post]
func (s *ProductService) ProductNew(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	body := s.ProductMapper.ProductSavePbToDTO(req)
	if err := s.ProductValidator.ValidateCreateProduct(body); err != nil {
		return nil, err
	}

	result, err := s.UC.CreateProduct(ctx, body)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if result.Product != nil {
		id = result.Product.ID
	}
	return &bdspropb.Response{
		Id:      id,
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái thuê sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/rent-status [put]
func (s *ProductService) RentStatus(ctx context.Context, req *bdspropb.StatusRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID:  &req.Id,
		RentStatus: enums.EProductRentStatus(req.Status),
		ContactID:  &req.ContactId,
		Note:       req.Note,
		Amount:     req.Amount,
		Phone:      req.Phone,
	}

	err := s.UC.RentStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái thuê sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/rent-visibility [put]
func (s *ProductService) RentVisibility(ctx context.Context, req *bdspropb.VisibilityRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Visibility == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, visibility is required")
	}
	body := dto.VisibilityUpdate{
		ProductID:  &req.Id,
		Visibility: enums.EVisibility(req.Visibility),
	}

	err := s.UC.RentVisibilityUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái bán sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/sale-status [put]
func (s *ProductService) SaleStatus(ctx context.Context, req *bdspropb.StatusRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID:  &req.Id,
		SaleStatus: enums.EProductSaleStatus(req.Status),
		ContactID:  &req.ContactId,
		Note:       req.Note,
		Amount:     req.Amount,
		Phone:      req.Phone,
	}

	err := s.UC.SaleStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái bán sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/sale-visibility [put]
func (s *ProductService) SaleVisibility(ctx context.Context, req *bdspropb.VisibilityRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Visibility == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, visibility is required")
	}
	body := dto.VisibilityUpdate{
		ProductID:  &req.Id,
		Visibility: enums.EVisibility(req.Visibility),
	}

	err := s.UC.SaleVisibilityUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Gợi ý các trường thông tin
// @Description Returns a list of products based on search filters
// @Tags User: Sản phẩm
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
// @Router /v2/bdspro/v2/product/suggest [post]
func (s *ProductService) Suggest(ctx context.Context, req *bdspropb.SuggestRequest) (*bdspropb.ProductResponse, error) {
	products, err := s.UC.SuggestFields(ctx, req.Content)
	if err != nil {
		return nil, err
	}

	// Tạo ProductResponse từ TotalSearchParser
	productResponse := s.mapTotalSearchParserToResponse(products)

	return productResponse, nil
}

// @Summary Gợi ý các trường thông tin sử dụng DeepSeek AI
// @Description Phân tích văn bản mô tả bất động sản bằng DeepSeek AI và trả về các trường thông tin được gợi ý
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.SuggestRequest true "Nội dung mô tả sản phẩm"
// @Success 200 {object} bdspropb.ProductResponse
// @Router /v2/bdspro/v2/product/suggest/deepseek [post]
func (s *ProductService) SuggestDeepseek(ctx context.Context, req *bdspropb.SuggestRequest) (*bdspropb.ProductResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích bằng DeepSeek
	products, err := s.UC.SuggestFieldsWithDeepseek(ctx, req.Content)
	if err != nil {
		return nil, err
	}

	// Tạo ProductResponse từ TotalSearchParser
	productResponse := s.mapTotalSearchParserToResponse(products)

	return productResponse, nil
}

// mapProductToResponse maps domain.Product to bdspropb.ProductResponse
func (s *ProductService) mapTotalSearchParserToResponse(parser *dto.TotalSearchParser) *bdspropb.ProductResponse {
	response := &bdspropb.ProductResponse{
		Name:            parser.Name,
		TransactionType: int32(parser.TransactionType),
		HouseInfo: &bdspropb.HouseInfo{
			NumBedroom:  parser.Bedroom,
			NumBathroom: parser.Bathroom,
			NumFloor:    parser.Floor,
			NumFront:    parser.Frontage,
			NumCarPark:  parser.Park,
			NumToilet:   parser.Toilet,
			// Orientation: parser.Orientation,
			// Furniture:   parser.Furniture,
		},
	}

	// Map các trường từ parser
	if parser.Area != nil {
		response.Area = float64(*parser.Area)
	}

	if parser.PriceSuggest != nil {
		response.PriceSuggest = parser.PriceSuggest
	}

	if parser.Address != nil {
		response.Address = *parser.Address
	}

	// Map regions
	if parser.Province != nil {
		response.ProvinceId = &parser.Province.ID
	}

	if parser.District != nil {
		response.DistrictId = &parser.District.ID
	}

	if parser.Ward != nil {
		response.WardId = &parser.Ward.ID
	}

	// Map property types (lấy property type đầu tiên nếu có)
	if len(parser.PropertyTypes) > 0 {
		response.PropertyTypeId = &parser.PropertyTypes[0].ID
	}

	// Map doc types (lấy doc type đầu tiên nếu có)
	if len(parser.DocTypes) > 0 {
		response.DocTypeId = &parser.DocTypes[0].ID
	}

	// Map sale price
	if parser.SalePrice != nil {
		response.Price = &bdspropb.ProductPrice{
			SalePrice: parser.SalePrice,
		}
	}

	// Map rent price
	if parser.RentPrice != nil {
		if response.Price == nil {
			response.Price = &bdspropb.ProductPrice{}
		}
		response.Price.RentPrice = parser.RentPrice
	}

	// Map amenities
	if len(parser.Amenities) > 0 {
		response.Amenities = make([]*bdspropb.AmenityItem, len(parser.Amenities))
		for i, amenity := range parser.Amenities {
			response.Amenities[i] = &bdspropb.AmenityItem{
				Id:   amenity.ID,
				Name: amenity.Name,
			}
		}
	}

	if parser.RentPrice != nil {
		response.Price.RentPrice = parser.RentPrice
	}

	return response
}

// @Summary Lấy chi tiết sản phẩm
// @Description API để lấy chi tiết sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID của sản phẩm"
// @Success 200 {object} bdspropb.IdRequest "Chi tiết sản phẩm"
// @Router /v2/bdspro/v2/product/{id}/detail [get]
func (s *ProductService) GetDetail(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.ProductResponse, error) {
	product, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	productResponse := s.ProductMapper.MapDetailProductPb(ctx, product)
	if product.OwnerID != 0 && product.OwnerOf == enums.EOwnerOfMember {
		productResponse.CreatedUser = s.UserClient.GetProfileById(ctx, product.OwnerID)
	}
	return productResponse, nil
}

// // @Summary Lấy danh sách lịch sử của sản phẩm
// // @Description Returns a list of history of a product
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Router /v2/bdspro/v2/product/{id}/histories [get]
// func (s *ProductService) GetHistories(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.HistoryResponse, error) {
// 	query := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&query, req); err != nil {
// 		return nil, err
// 	}

// 	histories, total, err := s.UC.History(ctx, req.Id, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	historyResponse := make([]*bdspropb.History, 0)
// 	copier.Copy(&historyResponse, &histories)
// 	return &bdspropb.HistoryResponse{
// 		Data:  historyResponse,
// 		Total: uint64(total),
// 	}, nil
// }

// @Summary Lấy danh sách thành viên của sản phẩm
// @Description Returns a list of members of a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/members [get]
func (s *ProductService) GetMembers(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.MemberResponse, error) {
	query := dto.SharingAccessSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	members, err := s.UC.Members(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}
	memberResponse := make([]*bdspropb.ProductAccess, 0)
	copier.Copy(&memberResponse, &members)
	return &bdspropb.MemberResponse{
		Data: memberResponse,
	}, nil
}

// @Summary Lấy danh sách bài viết liên quan
// @Description Returns a list of posts related to a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/posts [get]
func (s *ProductService) GetPosts(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.PostResponse, error) {
	query := dto.PostLinkSearch{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		ProductID: &req.Id,
	}
	posts, err := s.UC.Posts(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}

	postResponse := s.PostMapper.PostToPbList(posts)
	// postResponse := make([]*bdspropb.PostItem, 0)
	// copier.Copy(&postResponse, &posts)
	return &bdspropb.PostResponse{
		Data: postResponse,
	}, nil
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Param body body dto.ProductSaveRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/{id} [put]
func (s *ProductService) UpdateProduct(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	body := s.ProductMapper.ProductSavePbToDTO(req)
	if err := s.ProductValidator.ValidateCreateProduct(body); err != nil {
		return nil, err
	}

	_, err := s.UC.Update(ctx, req.Id, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
	// id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }
	// result, err := route.UC.Update(c, id, body)
	// _routes.RouteResult(c, result, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/product/{id} [delete]
func (s *ProductService) DeleteProduct(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "ok",
	}, nil
}

// @Summary Lấy thông tin về diện tích sản phẩm
// @Description Returns information about the area of a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/area-info/{id} [get]
func (s *ProductService) AreaInfo(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.AreaInfoResponse, error) {
	areaInfo, err := s.UC.InfoArea(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.AreaInfoResponse{
		TotalArea:    areaInfo.Total,
		AvaiableArea: areaInfo.Avaiable,
		SplitArea:    areaInfo.Area,
	}, nil
}

func (s *ProductService) GetDetailByIds(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.ProductInternalResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := s.ProductMapper.MapProductPbList(products)
	return &bdspropb.ProductInternalResponse{
		Data: productResponse,
	}, nil
}

func (s *ProductService) GetProductAttachmentByIds(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.ProductAttachmentResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := mapper.MapProductAttachmentPbList(ctx, s.ProductMapper.CodeDataUsecase, products)
	return &bdspropb.ProductAttachmentResponse{
		Data: productResponse,
	}, nil
}

// @Summary Get group products
// @Description Get list of products in a group based on sharing access
// @Tags Product
// @Accept json
// @Produce json
// @Param groupId path int true "Group ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param text query string false "Search text"
// @Param parentId query int false "Parent ID"
// @Param propertyTypeId query int false "Property Type ID"
// @Param docTypeId query int false "Document Type ID"
// @Param projectId query int false "Project ID"
// @Param provinceId query int false "Province ID"
// @Param districtId query int false "District ID"
// @Param wardId query int false "Ward ID"
// @Param transactionType query int false "Transaction Type"
// @Param saleStatus query int false "Sale Status"
// @Param rentStatus query int false "Rent Status"
// @Param saleVisibility query int false "Sale Visibility"
// @Param rentVisibility query int false "Rent Visibility"
// @Param sourceType query int false "Source Type"
// @Param archived query bool false "Archived"
// @Success 200 {object} bdspropb.SearchResponse
// @Router /v2/bdspro/v2/product/group/{groupId} [get]
func (s *ProductService) GetGroupProducts(ctx context.Context, req *bdspropb.GroupProductSearch) (*bdspropb.SearchResponse, error) {
	// Convert protobuf request to DTO
	searchReq := &dto.GroupProductSearchRequest{
		GroupID:          req.GroupId,
		Text:             req.Text,
		ParentID:         req.ParentId,
		PropertyTypeID:   req.PropertyTypeId,
		DocTypeID:        req.DocTypeId,
		ProjectID:        req.ProjectId,
		ProvinceID:       req.ProvinceId,
		DistrictID:       req.DistrictId,
		WardID:           req.WardId,
		TransactionType:  req.TransactionType,
		SaleStatus:       req.SaleStatus,
		RentStatus:       req.RentStatus,
		SaleVisibility:   req.SaleVisibility,
		RentVisibility:   req.RentVisibility,
		SourceType:       req.SourceType,
		TransactionPrice: req.TransactionPrice,
		SalePrice:        req.SalePrice,
		RentPrice:        req.RentPrice,
		SaleCommission:   req.SaleCommission,
		RentCommission:   req.RentCommission,
		RentPaymentCycle: req.RentPaymentCycle,
		Archived:         req.Archived,
	}

	// Set pagination
	if req.Page > 0 {
		searchReq.Page = req.Page
	}
	if req.Size > 0 {
		searchReq.Size = req.Size
	}

	// Get products from usecase
	products, total, err := s.UC.GetGroupProducts(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert domain products to protobuf response
	productResponses := s.ProductMapper.MapProductPbList(products)

	return &bdspropb.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}

func (s *ProductService) CurrentOrganizationProducts(ctx context.Context, req *bdspropb.ProductSearch) (*bdspropb.SearchResponse, error) {
	dto := s.ProductMapper.PbSearchRequestToDTO(req)
	products, total, err := s.UC.GetProductOfCurrentOrganization(ctx, dto)
	if err != nil {
		return nil, err
	}

	productResponses := s.ProductMapper.MapProductPbList(products)
	return &bdspropb.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}

func (s *ProductService) CreateOrganizationProduct(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.ProductSavePbToDTO(req)
	result, err := s.UC.CreateOrganizationProduct(ctx, dto)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if result.Product != nil {
		id = result.Product.ID
	}
	return &bdspropb.Response{
		Id:      id,
		Message: "success",
	}, nil
}
