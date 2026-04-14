package handler_grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
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

func (h *ParcelGrpcHandler) CreateBatchParcels(ctx context.Context, req *tqdpb.ParcelListDTO) (*tqdpb.ParcelListDTO, error) {

	// Convert proto to request DTO
	createReq := h.mapper.ProtoToDomains(req)

	// Call usecase
	result, err := h.usecase.CreateBatch(ctx, &dto.ParcelList{
		Data: createReq,
	})
	if err != nil {
		log.Printf("error creating parcel: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create parcel: %v", err)
	}

	return &tqdpb.ParcelListDTO{
		Data: h.mapper.ToProtoFromDomainList(result),
	}, nil
}

func (h *ParcelGrpcHandler) FindParcelByLocation(ctx context.Context, req *tqdpb.FindParcelByLocationRequest) (*tqdpb.Parcel, error) {

	// Call usecase
	result, err := h.usecase.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		log.Printf("error finding parcel by location: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find parcel by location: %v", err)
	}

	return h.mapper.DomainToProto(result), nil
}

func (h *ParcelGrpcHandler) FindParcelsByPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	geoJSON, err := h.toGeoJSON(req)
	if err != nil {
		log.Printf("Invalid polygon input: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("Polygon request: area_limit=%.2f km², page=%d, pageSize=%d, autoShrink=%v",
		req.MaxAreaKm2, req.Page, req.PageSize, req.AllowAutoShrink)

	dtoReq := &dto.PolygonQueryRequest{
		GeoJSON:         geoJSON,
		MaxAreaKm2:      req.MaxAreaKm2,
		Page:            int(req.Page),
		PageSize:        int(req.PageSize),
		AllowAutoShrink: req.AllowAutoShrink,
	}

	resp, err := h.usecase.FindByPolygon(ctx, dtoReq)
	if err != nil {
		log.Printf("FindByPolygon error: %v", err)
		return nil, h.grpcError(err)
	}

	parcelsProto := make([]*tqdpb.Parcel, len(resp.Parcels))
	for i, p := range resp.Parcels {
		parcelsProto[i] = h.mapper.DomainToProto(&p)
	}

	return &tqdpb.FindParcelsByPolygonResponse{
		Parcels:           parcelsProto,
		Total:             int32(resp.Total),
		Page:              int32(resp.Page),
		PageSize:          int32(resp.PageSize),
		AdjustedGeoJson:   resp.AdjustedGeoJSON,
		AdjustmentReason:  resp.AdjustmentReason,
		OriginalAreaKm2:   resp.OriginalAreaKm2,
		AdjustedAreaKm2:   resp.AdjustedAreaKm2,
		MaxAllowedAreaKm2: resp.MaxAllowedAreaKm2,
		Warning:           resp.Warning,
	}, nil
}

// toGeoJSON chuyển request về GeoJSON string (chỉ convert khi dùng linearRings)
func (h *ParcelGrpcHandler) toGeoJSON(req *tqdpb.FindParcelsByPolygonRequest) (string, error) {
	if req.GeoJson != "" {
		if !isValidGeoJSON(req.GeoJson) {
			return "", errors.New("invalid GeoJSON polygon")
		}
		return req.GeoJson, nil
	}
	if len(req.Coordinates) > 0 {
		return linearRingsToGeoJSON(req.Coordinates)
	}
	return "", errors.New("either geoJson or linearRings must be provided")
}

// linearRingsToGeoJSON chuyển danh sách linear ring thành GeoJSON Polygon
func linearRingsToGeoJSON(rings []*tqdpb.CoordinatePolygon) (string, error) {
	if len(rings) == 0 {
		return "", errors.New("linearRings cannot be empty")
	}
	// Kiểm tra outer ring
	if rings[0] == nil || len(rings[0].Points) < 3 {
		return "", fmt.Errorf("outer ring must have at least 3 points")
	}

	// Tạo coordinates theo định dạng GeoJSON
	coords := make([][][]float64, len(rings))
	for i, ring := range rings {
		if ring == nil {
			return "", fmt.Errorf("ring %d is nil", i)
		}
		if len(ring.Points) < 3 {
			return "", fmt.Errorf("ring %d has less than 3 points", i)
		}
		// Tạo mảng điểm cho ring (cần khép kín)
		points := make([][]float64, len(ring.Points)+1)
		for j, pt := range ring.Points {
			if pt == nil {
				return "", fmt.Errorf("ring %d has nil point at index %d", i, j)
			}
			points[j] = []float64{pt.Lng, pt.Lat}
		}
		// Đóng ring bằng điểm đầu tiên
		points[len(ring.Points)] = []float64{ring.Points[0].Lng, ring.Points[0].Lat}
		coords[i] = points
	}

	geoJSON := map[string]interface{}{
		"type":        "Polygon",
		"coordinates": coords,
	}
	jsonBytes, err := json.Marshal(geoJSON)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GeoJSON: %v", err)
	}
	return string(jsonBytes), nil
}

// isValidGeoJSON kiểm tra nhanh GeoJSON polygon
func isValidGeoJSON(geoJSON string) bool {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(geoJSON), &data); err != nil {
		return false
	}
	typ, _ := data["type"].(string)
	if typ != "Polygon" {
		return false
	}
	_, ok := data["coordinates"].([]interface{})
	return ok
}

// grpcError chuyển lỗi business sang mã gRPC
func (h *ParcelGrpcHandler) grpcError(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "exceeds limit"),
		strings.Contains(msg, "invalid polygon"),
		strings.Contains(msg, "self-intersects"),
		strings.Contains(msg, "must have at least"):
		return status.Error(codes.InvalidArgument, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}
