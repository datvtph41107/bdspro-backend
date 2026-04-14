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

func (h *ParcelGrpcHandler) toGeoJSON(req *tqdpb.FindParcelsByPolygonRequest) (string, error) {
	if req.GeoJson != "" {
		if !isValidGeoJSON(req.GeoJson) {
			return "", errors.New("invalid GeoJSON polygon")
		}
		return req.GeoJson, nil
	}
	if len(req.Coordinates) >= 3 {
		return pointsToGeoJSON(req.Coordinates)
	}
	return "", errors.New("either geoJson or at least 3 coordinates are required")
}

func pointsToGeoJSON(points []*tqdpb.Point) (string, error) {
	if len(points) < 3 {
		return "", fmt.Errorf("need at least 3 points to form a polygon, got %d", len(points))
	}
	ring := make([][]float64, len(points)+1)
	for i, pt := range points {
		if pt == nil {
			return "", fmt.Errorf("point at index %d is nil", i)
		}
		ring[i] = []float64{pt.Lng, pt.Lat}
	}
	// Đóng ring
	ring[len(points)] = []float64{points[0].Lng, points[0].Lat}

	// GeoJSON Polygon
	polygon := map[string]interface{}{
		"type":        "Polygon",
		"coordinates": [][][]float64{ring},
	}
	jsonBytes, err := json.Marshal(polygon)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GeoJSON: %v", err)
	}
	return string(jsonBytes), nil
}

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
