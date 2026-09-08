package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/usecases"
	admin_usecases "bdspro/internal/usecases/admin"
	shared_usecase "bdspro/internal/usecases/shared"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"sort"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type BdsproPublicService struct {
	bdspropb.UnimplementedBdsproPublicServiceServer
	UC                *shared_usecase.ProductUsecase
	ListUC            *usecases.ListUsecase
	AreaRegionUsecase *admin_usecases.AreaRegionUsecase
	Mapper            *mapper.PublicMapper
	ProductMapper     *mapper.ProductMapper
	AreaRegionMapper  *mapper.AreaRegionMapper
	StartUpTime       int64
}

func NewBdsproPublicService(
	uc *shared_usecase.ProductUsecase,
	listUC *usecases.ListUsecase,
	areaRegionUsecase *admin_usecases.AreaRegionUsecase,
	mapper *mapper.PublicMapper,
	productMapper *mapper.ProductMapper,
	areaRegionMapper *mapper.AreaRegionMapper,
) *BdsproPublicService {
	return &BdsproPublicService{
		UC:                uc,
		ListUC:            listUC,
		AreaRegionUsecase: areaRegionUsecase,
		Mapper:            mapper,
		ProductMapper:     productMapper,
		AreaRegionMapper:  areaRegionMapper,
		StartUpTime:       time.Now().UnixMilli(),
	}
}

// @Summary Cấu hình app quy hoạch (public)
// @Description Trả JSON object: backgroundMaps, qhLayers — mỗi phần tử có id, name, image
// @Tags Quy hoạch
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v2/qh/config-app [get]
func (s *BdsproPublicService) GetQhConfigApp(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	payload := map[string]interface{}{
		"backgroundMaps": []interface{}{
			map[string]interface{}{"id": "basemap-light", "name": "Bản đồ nền sáng", "image": ""},
			map[string]interface{}{"id": "basemap-dark", "name": "Bản đồ nền tối", "image": ""},
			map[string]interface{}{"id": "satellite", "name": "Ảnh vệ tinh", "image": ""},
		},
		"qhLayers": []interface{}{
			map[string]interface{}{"id": "planning-use", "name": "Quy hoạch sử dụng đất", "image": ""},
			map[string]interface{}{"id": "planning-transport", "name": "Quy hoạch giao thông", "image": ""},
		},
	}
	st, err := structpb.NewStruct(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "qh config-app: %v", err)
	}
	return st, nil
}

// @Summary Lấy danh sách sản phẩm thị trường
// @Description Lấy danh sách sản phẩm thị trường
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.ProductMarketRequest true "Request"
// @Router /v2/bdspro/v2/product/market [get]
func (s *BdsproPublicService) GetProductMarket(ctx context.Context, req *bdspropb.ProductMarketRequest) (*bdspropb.ProductMarketResponse, error) {
	query := s.ProductMapper.MapMarketProductToDTO(req)

	result, total, err := s.UC.SearchProductMarket(ctx, query)
	if err != nil {
		return nil, err
	}
	products := s.ProductMapper.MapProductMarketPbList(ctx, result)
	return &bdspropb.ProductMarketResponse{
		Data:  products,
		Total: total,
	}, err
}

// @Summary Lấy danh sách loại tài sản
// @Description Lấy danh sách loại tài sản
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/property-type [get]
func (s *BdsproPublicService) GetPropertyType(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.PropertyTypeListResponse, error) {
	query := dto.PropertyTypeSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListPropertyType(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.PropertyType, 0)
	for _, item := range result {
		items = append(items, &bdspropb.PropertyType{
			Id:               item.ID,
			Name:             item.Name,
			CreatedAt:        _utils.FormatTimeToString(item.CreatedAt),
			UpdatedAt:        _utils.FormatTimeToString(item.UpdatedAt),
			Active:           item.Active,
			ClassifyProperty: uint32(item.ClassifyProperty),
			FieldStrs:        item.FieldStrs,
		})
	}
	return &bdspropb.PropertyTypeListResponse{
		Data:          items,
		TotalElements: int64(len(result)),
	}, nil
}

// @Summary Lấy danh sách area-region (khu vực - bảng regions)
// @Description Lấy toàn bộ danh sách area-region (public, không cần auth)
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Router /v2/bdspro/v2/list/area-region [get]
func (s *BdsproPublicService) GetAreaRegions(ctx context.Context, req *sharepb.Empty) (*bdspropb.ListAreaRegionsResponse, error) {
	regions, err := s.AreaRegionUsecase.List(ctx)
	if err != nil {
		return nil, err
	}
	return &bdspropb.ListAreaRegionsResponse{
		Data: s.AreaRegionMapper.ToProtoList(regions),
	}, nil
}

// @Summary Lấy danh sách khu vực
// @Description Lấy danh sách khu vực
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/region [get]
func (s *BdsproPublicService) GetRegion(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := s.Mapper.PbToRegionRequest(req)
	result, total, err := s.ListUC.ListRegion(ctx, query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &bdspropb.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách khu vực
// @Description Lấy danh sách khu vực
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/area-region [get]
func (s *BdsproPublicService) GetAreaRegion(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
		Text: req.Text,
	}
	result, total, err := s.ListUC.ListAreaRegion(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.ItemResponse, len(result))
	for i, region := range result {
		items[i] = &bdspropb.ItemResponse{
			Id:   region.ID,
			Name: region.Name,
		}
	}
	return &bdspropb.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách loại pháp lý
// @Description Lấy danh sách loại pháp lý
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/doc-type [get]
func (s *BdsproPublicService) GetDocType(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := dto.DocTypeSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListDocType(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &bdspropb.ListItemResponse{
		Data: items,
	}, nil
}

// @Summary Lấy danh sách dự án
// @Description Lấy danh sách dự án
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/project [get]
func (s *BdsproPublicService) GetProject(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := dto.ProjectSearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, total, err := s.ListUC.ListProject(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &bdspropb.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách tiện ích
// @Description Lấy danh sách tiện ích
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/amenity [get]
func (s *BdsproPublicService) GetAmenity(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := dto.AmenitySearchDTO{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	result, err := s.ListUC.ListAmenity(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &bdspropb.ListItemResponse{
		Data:  items,
		Total: uint32(len(result)),
	}, nil
}

// @Summary Lấy tổng số lượng sản phẩm
// @Description Lấy tổng số lượng sản phẩm
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.ProductTotalRequest true "Request"
// @Router /v2/bdspro/v2/product/total [get]
func (s *BdsproPublicService) GetProductTotal(ctx context.Context, req *bdspropb.ProductTotalRequest) (*bdspropb.SearchResponse, error) {
	query := dto.TextSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Text: req.Text,
	}
	result, total, err := s.UC.SmartSearch(ctx, req.Text, &query)
	if err != nil {
		return nil, err
	}
	products := s.ProductMapper.MapProductPbList(result)

	return &bdspropb.SearchResponse{
		Data:          products,
		TotalElements: total,
	}, nil
}

// @Summary Lấy danh sách tòa nhà trong dự án
// @Description Lấy danh sách tòa nhà trong dự án
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param req query bdspropb.SearchQueryRequest true "Request"
// @Router /v2/bdspro/v2/list/project-build [get]
func (s *BdsproPublicService) GetProjectBuild(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	query := dto.TextSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
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
	items := make([]*bdspropb.ItemResponse, 0)
	copier.Copy(&items, &result)
	return &bdspropb.ListItemResponse{
		Data:  items,
		Total: uint32(total),
	}, nil
}

// @Summary Lấy danh sách sản phẩm của user (Public API)
// @Description Lấy danh sách sản phẩm của một user cụ thể - API public không cần authentication
// @Tags Sản phẩm
// @Accept json
// @Produce json
// @Param profileId path int true "Profile ID"
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
// @Router /v2/bdspro/public/product/user/{profileId} [get]
func (s *BdsproPublicService) GetUserProductsPublic(ctx context.Context, req *bdspropb.PublicUserProductSearch) (*bdspropb.SearchResponse, error) {
	// Convert protobuf request to DTO
	searchReq := &dto.UserProductSearchRequest{
		UserID:           req.ProfileId,
		Text:             req.Text,
		ParentID:         req.ParentId,
		PropertyTypeID:   req.PropertyTypeId,
		DocTypeID:        req.DocTypeId,
		ProjectID:        req.ProjectId,
		ProvinceID:       req.ProvinceId,
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
	products, total, err := s.UC.GetUserProducts(ctx, searchReq)
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

var enumProviders = map[string][]*bdspropb.EnumItem{
	"product_sale_status":     enumMapToList(enums.EProductSaleStatusNames),
	"product_rent_status":     enumMapToList(enums.EProductRentStatusNames),
	"asset_history":           enumMapToList(enums.AssetHistoryNames),
	"target_type":             enumMapToList(enums.TargetTypeNames),
	"visibility":              enumMapToList(enums.EVisibilityNames),
	"producthistory":          enumMapToList(enums.EProductHistoryNames),
	"poststatus":              enumMapToList(enums.EPostStatusNames),
	"doctype":                 enumMapToList(enums.EDocTypeNames),
	"assetstatus":             enumMapToList(enums.EAssetStatusNames),
	"house_orient":            enumMapToList(enums.EHouseOrientNames),
	"house_certificate":       enumMapToList(enums.EHouseCertificateNames),
	"property_status":         enumMapToList(enums.EPropertyStatusNames),
	"building_type":           enumMapToList(enums.BuildingTypeNames),
	"transaction_status":      enumMapToList(enums.TransactionStatusNames),
	"transaction_type":        enumMapToList(enums.TransactionTypeNames),
	"owner_type":              enumMapToList(enums.OwnerTypeNames),
	"post_transaction_status": enumMapToList(enums.EPostTransactionStatusNames),
	"distribution_status":     enumMapToList(enums.DistributionStatusNames),
	"distribution_reason":     enumMapToList(enums.DistributionReasonNameMap),
	"deal_type":               enumMapToList(enums.DealTypeNames),
	"source_status":           enumMapToList(enums.SourceStatusNames),
	"source_type":             enumMapToList(enums.SourceTypeMap),
	"mining_mode":             enumMapToList(enums.MiningModeNames),
	"mining_scope":            enumMapToList(enums.MiningScopeNames),
	"classify_property":       enumMapToList(enums.EClassifyPropertyNames),
	"purpose_used":            enumMapToList(enums.EPurposeUsedNames),
}

// @Summary Lấy enum dạng name-value
// @Description Lấy danh sách enum với value và name tương ứng
// @Tags Enum
// @Accept json
// @Produce json
// @Param text path string true "Tên loại enum"
// @Success 200 {object} bdspropb.GetEnumResponse
// @Failure 404 {object} map[string]string
// @Router /v2/bdspro/public/enums/{text} [get]
func (s *BdsproPublicService) GetEnumByName(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.GetEnumResponse, error) {
	enumType := req.Text
	fmt.Println("enumType: ", enumType)

	// Khai báo enumProviders nội bộ

	if values, ok := enumProviders[enumType]; ok {
		return &bdspropb.GetEnumResponse{
			Data: values,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "enum not found: %s", enumType)
}

func (s *BdsproPublicService) GetEnums(ctx context.Context, req *sharepb.SyncRequest) (*bdspropb.GetEnumsResponse, error) {
	if req.Timestamp > 0 && req.Timestamp < s.StartUpTime {
		return &bdspropb.GetEnumsResponse{
			Timestamp: s.StartUpTime,
			Lastest:   false,
		}, nil
	}
	result := []*bdspropb.GetEnumResponse{}
	for k, v := range enumProviders {
		result = append(result, &bdspropb.GetEnumResponse{
			Key:  k,
			Data: v,
		})
	}
	return &bdspropb.GetEnumsResponse{
		Enums:     result,
		Timestamp: s.StartUpTime,
		Lastest:   true,
	}, nil
}

func enumMapToList[T ~uint32 | ~uint | ~int](enumMap map[T]string) []*bdspropb.EnumItem {
	keys := make([]T, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}

	// Sắp xếp keys theo giá trị tăng dần
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := make([]*bdspropb.EnumItem, 0, len(keys))
	for _, k := range keys {
		result = append(result, &bdspropb.EnumItem{
			Value: uint32(k),
			Name:  enumMap[k],
		})
	}
	return result
}
