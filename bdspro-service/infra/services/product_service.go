package service

import (
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb_bdspro.UnimplementedProductServiceServer
	UC               *shared_usecase.ProductUsecase
	ProductRepo      repo.ProductRepo
	ProductMapper    *mapper.ProductMapper
	ProductValidator *validator.ProductValidator
}

func NewProductService(uc *shared_usecase.ProductUsecase,
	productRepo repo.ProductRepo,
	productMapper *mapper.ProductMapper,
	productValidator *validator.ProductValidator,
) *ProductService {
	return &ProductService{
		UC:               uc,
		ProductRepo:      productRepo,
		ProductMapper:    productMapper,
		ProductValidator: productValidator,
	}
}

// @Summary Lưu trữ sản phẩm
// @Description API này cho phép lưu trữ sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body pb_bdspro.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/archived [put]
func (s *ProductService) ArchiveProduct(c context.Context, req *pb_bdspro.ArchivedRequest) (*pb_bdspro.Response, error) {
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
	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body pb_bdspro.SaveChildRequest true "Thông tin sản phẩm con"
// @Router /v2/bdspro/v2/product/child/devide [post]
func (s *ProductService) ChildDevideProduct(ctx context.Context, req *pb_bdspro.DevideChildRequest) (*pb_bdspro.Response, error) {
	dto := s.ProductMapper.MapDevideChildRequest(req)
	if err := s.ProductValidator.ValidateDevideChildRequest(dto); err != nil {
		return nil, err
	}

	_, err := s.UC.DevideProductChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
// func (s *ProductService) ChildHistories(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.HistoryResponse, error) {
// 	id := req.Id

// 	_body := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&_body, req); err != nil {
// 		return nil, err
// 	}

// 	result, total, err := s.UC.ChildHistory(ctx, id, &_body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	histories := make([]*pb_bdspro.History, 0)
// 	copier.Copy(&histories, &result)
// 	return &pb_bdspro.GetHistoryResponse{
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
func (s *ProductService) ChildMergeProduct(ctx context.Context, req *pb_bdspro.MergeChildRequest) (*pb_bdspro.Response, error) {
	_, err := s.UC.MergeProductChild(ctx, &dto.MergeProductChildRequest{
		ParentID: &req.ParentId,
		ChildIDs: &req.ChildIds,
	})
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) ChildMergeAll(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	_, err := s.UC.MergeAllProductChild(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param body body pb_bdspro.SaveChildRequest true "Thông tin sản phẩm con"
// @Security BearerAuth
// @Router /v2/bdspro/v2/product/child/split [post]
func (s *ProductService) ChildSplit(ctx context.Context, req *pb_bdspro.SaveChildRequest) (*pb_bdspro.Response, error) {
	dto := s.ProductMapper.SaveProductToEntity(req)
	if err := s.ProductValidator.ValidateSaveProductChild(dto); err != nil {
		return nil, err
	}
	_, err := s.UC.CreateChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}

// @Summary Cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/deposite [post]
func (s *ProductService) Deposite(context.Context, *pb_bdspro.DepositRequest) (*pb_bdspro.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposite not implemented")
}

// @Summary Hủy cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param depositeId path int true "Deposite ID"
// @Router /v2/bdspro/v2/product/deposite/{depositeId} [delete]
func (s *ProductService) DeleteDeposite(context.Context, *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDeposite not implemented")
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
func (s *ProductService) ExportProduct(context.Context, *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
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
func (s *ProductService) ProductMe(ctx context.Context, req *pb_bdspro.ProductSearch) (*pb_bdspro.SearchResponse, error) {
	body := &dto.ProductSearchRequest{}
	if err := copier.Copy(body, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.SearchProductByOwner(ctx, body)
	if err != nil {
		return nil, err
	}

	products := make([]*pb_bdspro.ProductResponse, 0)
	copier.Copy(&products, &result)
	return &pb_bdspro.SearchResponse{
		Data:          products,
		TotalElements: uint32(total),
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
func (s *ProductService) ProductNew(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.Response, error) {
	body := &dto.ProductSaveRequest{}
	if err := copier.Copy(body, req); err != nil {
		return nil, err
	}

	_, err := s.UC.CreateProduct(ctx, body)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) RentStatus(ctx context.Context, req *pb_bdspro.StatusRequest) (*pb_bdspro.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID: &req.Id,
		Status:    enums.EProductStatus(req.Status),
	}

	err := s.UC.RentStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) RentVisibility(ctx context.Context, req *pb_bdspro.VisibilityRequest) (*pb_bdspro.Response, error) {
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
	return &pb_bdspro.Response{
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
func (s *ProductService) SaleStatus(ctx context.Context, req *pb_bdspro.StatusRequest) (*pb_bdspro.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID: &req.Id,
		Status:    enums.EProductStatus(req.Status),
	}

	err := s.UC.SaleStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) SaleVisibility(ctx context.Context, req *pb_bdspro.VisibilityRequest) (*pb_bdspro.Response, error) {
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
	return &pb_bdspro.Response{
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
func (s *ProductService) Suggest(ctx context.Context, req *pb_bdspro.SuggestRequest) (*pb_bdspro.ProductResponse, error) {
	products, err := s.UC.SuggestFields(ctx, req.Content)
	if err != nil {
		return nil, err
	}
	productResponse := &pb_bdspro.ProductResponse{}
	copier.Copy(productResponse, products)
	return productResponse, nil
}

// @Summary Lấy chi tiết sản phẩm
// @Description API để lấy chi tiết sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID của sản phẩm"
// @Success 200 {object} pb_bdspro.IdRequest "Chi tiết sản phẩm"
// @Router /v2/bdspro/v2/product/{id}/detail [get]
func (s *ProductService) GetDetail(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.ProductResponse, error) {
	product, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	productResponse := &pb_bdspro.ProductResponse{}
	copier.Copy(productResponse, product)
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
// func (s *ProductService) GetHistories(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.HistoryResponse, error) {
// 	query := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&query, req); err != nil {
// 		return nil, err
// 	}

// 	histories, total, err := s.UC.History(ctx, req.Id, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	historyResponse := make([]*pb_bdspro.History, 0)
// 	copier.Copy(&historyResponse, &histories)
// 	return &pb_bdspro.HistoryResponse{
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
func (s *ProductService) GetMembers(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.MemberResponse, error) {
	query := dto.SharingAccessSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	members, err := s.UC.Members(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}
	memberResponse := make([]*pb_bdspro.ProductAccess, 0)
	copier.Copy(&memberResponse, &members)
	return &pb_bdspro.MemberResponse{
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
func (s *ProductService) GetPosts(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.PostResponse, error) {
	query := dto.PostLinkSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	posts, err := s.UC.Posts(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}
	postResponse := make([]*pb_bdspro.PostItem, 0)
	copier.Copy(&postResponse, &posts)
	return &pb_bdspro.PostResponse{
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
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.Response, error) {
	body := &dto.ProductSaveRequest{}
	if err := copier.Copy(body, req); err != nil {
		return nil, err
	}

	_, err := s.UC.Update(ctx, req.Id, body)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.Response{
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
func (s *ProductService) AreaInfo(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AreaInfoResponse, error) {
	areaInfo, err := s.UC.InfoArea(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &pb_bdspro.AreaInfoResponse{
		TotalArea:    areaInfo.Total,
		AvaiableArea: areaInfo.Avaiable,
		SplitArea:    areaInfo.Area,
	}, nil
}

func (s *ProductService) GetDetailByIds(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.ProductInternalResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := s.ProductMapper.MapProductPbList(products)
	return &pb_bdspro.ProductInternalResponse{
		Data: productResponse,
	}, nil
}

func (s *ProductService) GetProductAttachmentByIds(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.ProductAttachmentResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := mapper.MapProductAttachmentPbList(products)
	return &pb_bdspro.ProductAttachmentResponse{
		Data: productResponse,
	}, nil
}
