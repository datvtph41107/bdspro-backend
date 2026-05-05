package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"log"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParcelGrpcHandler struct {
	tqdpb.UnimplementedParcelServiceServer
	usecase usecase.ParcelUsecase
	mapper  *mapper.ParcelMapper
}

func NewParcelGrpcHandler(
	usecase usecase.ParcelUsecase,
	mapper *mapper.ParcelMapper,
) *ParcelGrpcHandler {
	return &ParcelGrpcHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// =====================================================
// API Detail VIP - Get parcel detail information
// =====================================================

func (h *ParcelGrpcHandler) GetParcelInfo(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelInfoResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	data, err := h.usecase.GetParcelInfo(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelInfoDetail(data), nil
}

func (h *ParcelGrpcHandler) GetParcelLayers(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelLayersResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	data, err := h.usecase.GetParcelLayers(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelLayers(data), nil
}

// =====================================================
// Batch Operations
// =====================================================

// CreateBatchParcels - Tạo batch parcels
func (h *ParcelGrpcHandler) CreateBatchParcels(ctx context.Context, req *tqdpb.ParcelListDTO) (*tqdpb.ParcelListDTO, error) {
	// Get user ID from context
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	// Validate request
	if req == nil || len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "request data is required")
	}

	log.Printf("[CreateBatchParcels] UserID: %d, Batch size: %d", userID, len(req.Data))

	// Convert proto to request DTO
	createReq := h.mapper.ProtoToDomains(req)

	// Call usecase
	result, err := h.usecase.CreateBatch(ctx, &dto.ParcelList{
		Data: createReq,
	})
	if err != nil {
		log.Printf("[CreateBatchParcels] Error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create parcel: %v", err)
	}

	log.Printf("[CreateBatchParcels] Successfully created %d parcels", len(result.Data))

	return &tqdpb.ParcelListDTO{
		Data: h.mapper.ToProtoFromDomainList(result),
	}, nil
}

// =====================================================
// Location-based Queries
// =====================================================

// FindParcelByLocation - Tìm parcel theo tọa độ
func (h *ParcelGrpcHandler) FindParcelByLocation(ctx context.Context, req *tqdpb.FindParcelByLocationRequest) (*tqdpb.Parcel, error) {
	// Validate input
	if req.Latitude == 0 && req.Longitude == 0 {
		return nil, status.Error(codes.InvalidArgument, "latitude and longitude are required")
	}

	// Validate coordinate ranges
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}

	log.Printf("[FindParcelByLocation] Lat: %f, Lng: %f", req.Latitude, req.Longitude)

	// Call usecase
	result, err := h.usecase.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		log.Printf("[FindParcelByLocation] Error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find parcel by location: %v", err)
	}

	if result == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.DomainToProto(result), nil
}

// =====================================================
// Polygon-based Queries
// =====================================================

// FindParcelsByPolygon - Tìm parcels theo polygon
func (h *ParcelGrpcHandler) FindParcelsByPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	// Validate coordinates
	var points []repo.Point
	if len(req.Coordinates) == 0 {
		return nil, status.Error(codes.InvalidArgument, "coordinates are required")
	}

	// Build points list
	points = make([]repo.Point, len(req.Coordinates))
	for i, pt := range req.Coordinates {
		// Validate each point
		if pt.Lat < -90 || pt.Lat > 90 {
			return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
		}
		if pt.Lng < -180 || pt.Lng > 180 {
			return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
		}
		points[i] = repo.Point{Lat: pt.Lat, Lng: pt.Lng}
	}

	// Validate polygon has at least 3 points
	if len(points) < 3 {
		return nil, status.Error(codes.InvalidArgument, "polygon must have at least 3 points")
	}

	// Build pagable
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	log.Printf("[FindParcelsByPolygon] Points count: %d, Intersect: %v, MaxAreaKm2: %f, Page: %d, Size: %d",
		len(points), req.Intersect, req.MaxAreaKm2, pagable.GetPage(), pagable.GetSize())

	page := int(pagable.GetPage())
	pageSize := int(pagable.GetSize())

	parcels, total, err := h.usecase.FindParcelsByPolygon(
		ctx,
		points,
		req.Intersect,
		req.MaxAreaKm2,
		page,
		pageSize,
	)
	if err != nil {
		log.Printf("[FindParcelsByPolygon] Usecase error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[FindParcelsByPolygon] Found %d parcels, Total: %d", len(parcels), total)

	// Convert domain to proto
	data := make([]*tqdpb.Parcel, len(parcels))
	for i, p := range parcels {
		data[i] = h.mapper.DomainToProto(&p)
	}

	return &tqdpb.FindParcelsByPolygonResponse{
		Data:     data,
		Total:    int32(total),
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *ParcelGrpcHandler) SearchTxtClientParcels(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ListParcelResponse, error) {
	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "search text is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[SearchTxtClientParcels] UserID: %d, SearchText: %s, Page: %d, Size: %d",
		userID, req.Text, req.Page, req.Size)

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	// Gọi đúng method SearchParcelsByText (không phải Search)
	parcels, total, err := h.usecase.SearchParcelsByText(ctx, req.Text, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		log.Printf("[SearchTxtClientParcels] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert kết quả sang proto
	results := make([]*tqdpb.ParcelInfo, 0, len(parcels))
	for _, p := range parcels {
		results = append(results, h.mapper.ToProtoParcelInfo(p))
	}

	return &tqdpb.ListParcelResponse{
		Data:  results,
		Total: uint32(total),
	}, nil
}

// PreviewParcelsById - Xem trước thông tin parcel theo ID
func (h *ParcelGrpcHandler) PreviewParcelsById(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ParcelInfo, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id)

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		log.Printf("[PreviewParcelsById] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if info == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	// Convert sang proto
	return h.mapper.ToProtoParcelInfo(info), nil
}

func (h *ParcelGrpcHandler) PublicParcelsById(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ParcelInfo, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id)

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		log.Printf("[PreviewParcelsById] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if info == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	// Convert sang proto
	return h.mapper.ToProtoParcelInfo(info), nil
}
