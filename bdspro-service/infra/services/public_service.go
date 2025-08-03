package service

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	_dto "common/dto"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type BdsproPublicService struct {
	pb_bdspro.UnimplementedBdsproPublicServiceServer
	UC            *shared_usecase.ProductUsecase
	ListUC        *usecases.ListUsecase
	Mapper        *mapper.PublicMapper
	ProductMapper *mapper.ProductMapper
}

func NewBdsproPublicService(
	uc *shared_usecase.ProductUsecase,
	listUC *usecases.ListUsecase,
	mapper *mapper.PublicMapper,
	productMapper *mapper.ProductMapper) *BdsproPublicService {
	return &BdsproPublicService{
		UC:            uc,
		ListUC:        listUC,
		Mapper:        mapper,
		ProductMapper: productMapper,
	}
}

// @Summary Lấy danh sách sản phẩm thị trường
// @Description Lấy danh sách sản phẩm thị trường
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.ProductMarketRequest true "Request"
// @Router /v2/bdspro/v2/product/market [get]
func (s *BdsproPublicService) GetProductMarket(ctx context.Context, req *pb_bdspro.ProductMarketRequest) (*pb_bdspro.ProductMarketResponse, error) {
	query := s.ProductMapper.MapMarketProductToDTO(req)

	result, total, err := s.UC.SearchProductMarket(ctx, query)
	if err != nil {
		return nil, err
	}
	products := s.ProductMapper.MapProductMarketPbList(result)
	return &pb_bdspro.ProductMarketResponse{
		Data:  products,
		Total: uint32(total),
	}, err
}

// @Summary Lấy danh sách loại tài sản
// @Description Lấy danh sách loại tài sản
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/property-type [get]
func (s *BdsproPublicService) GetPropertyType(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := dto.PropertyTypeSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListPropertyType(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data:  items,
		Total: uint32(len(result)),
	}, nil
}

// @Summary Lấy danh sách khu vực
// @Description Lấy danh sách khu vực
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/region [get]
func (s *BdsproPublicService) GetRegion(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := s.Mapper.PbToRegionRequest(req)
	result, total, err := s.ListUC.ListRegion(ctx, query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách loại pháp lý
// @Description Lấy danh sách loại pháp lý
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/doc-type [get]
func (s *BdsproPublicService) GetDocType(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := dto.DocTypeSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListDocType(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data: items,
	}, nil
}

// @Summary Lấy danh sách dự án
// @Description Lấy danh sách dự án
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/project [get]
func (s *BdsproPublicService) GetProject(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := dto.ProjectSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, total, err := s.ListUC.ListProject(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách tiện ích
// @Description Lấy danh sách tiện ích
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/amenity [get]
func (s *BdsproPublicService) GetAmenity(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := dto.AmenitySearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListAmenity(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data:  items,
		Total: uint32(len(result)),
	}, nil
}

// @Summary Lấy tổng số lượng sản phẩm
// @Description Lấy tổng số lượng sản phẩm
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.ProductTotalRequest true "Request"
// @Router /v2/bdspro/v2/product/total [get]
func (s *BdsproPublicService) GetProductTotal(ctx context.Context, req *pb_bdspro.ProductTotalRequest) (*pb_bdspro.SearchResponse, error) {
	query := dto.TextSearchRequest{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Text: req.Text,
	}
	result, total, err := s.UC.SmartSearch(ctx, req.Text, &query)
	if err != nil {
		return nil, err
	}
	products := s.ProductMapper.MapProductPbList(result)

	return &pb_bdspro.SearchResponse{
		Data:          products,
		TotalElements: uint32(total),
	}, nil
}

// @Summary Lấy danh sách tòa nhà trong dự án
// @Description Lấy danh sách tòa nhà trong dự án
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query pb_bdspro.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/project-build [get]
func (s *BdsproPublicService) GetProjectBuild(ctx context.Context, req *pb_bdspro.SearchQueryRequest) (*pb_bdspro.ListItemResponse, error) {
	query := dto.TextSearchRequest{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Text: req.Text,
	}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, total, err := s.ListUC.ListProjectBuild(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb_bdspro.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &pb_bdspro.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}
