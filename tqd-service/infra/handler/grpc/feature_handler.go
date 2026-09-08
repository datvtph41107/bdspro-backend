package handler_grpc

import (
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type FeatureGrpcHandler struct {
	tqdpb.UnimplementedFeatureServiceServer
	featureUsecase *usecase.FeatureUsecase
	SyncProvider   *_utils.SyncUtil
}

func NewFeatureGrpcHandler(usecase *usecase.FeatureUsecase, syncProvider *_utils.SyncUtil) *FeatureGrpcHandler {
	return &FeatureGrpcHandler{featureUsecase: usecase, SyncProvider: syncProvider}
}

func (h *FeatureGrpcHandler) GetFeaturesByPointRadius(ctx context.Context, req *tqdpb.GetFeaturesByPointRadiusRequest) (*tqdpb.ListFeaturesResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDFeatureList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListFeaturesResponse{}, nil
	}

	// Validate request
	if req.Radius <= 0 {
		return nil, status.Error(codes.InvalidArgument, "radius must be positive")
	}
	if req.Lat < -90 || req.Lat > 90 || req.Lng < -180 || req.Lng > 180 {
		return nil, status.Error(codes.InvalidArgument, "invalid coordinates")
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 1000
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	usecaseReq := &usecase.PointRadiusRequest{
		Lat:    req.Lat,
		Lng:    req.Lng,
		Radius: int(req.Radius),
		Limit:  limit,
		Offset: offset,
	}
	features, total, err := h.featureUsecase.GetFeaturesByPointRadius(ctx, usecaseReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query features: %v", err)
	}

	pbFeatures := make([]*tqdpb.Feature, 0, len(features))
	for _, f := range features {
		pbFeat, err := convertFeatureToProto(f)
		if err != nil {
			// Log error but continue
			continue
		}
		pbFeatures = append(pbFeatures, pbFeat)
	}

	return &tqdpb.ListFeaturesResponse{
		Features: pbFeatures,
		Total:    total,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

func (h *FeatureGrpcHandler) GetFeaturesByPolygon(ctx context.Context, req *tqdpb.GetFeaturesByPolygonRequest) (*tqdpb.ListFeaturesResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDFeatureList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListFeaturesResponse{}, nil
	}

	var geojson string
	var err error

	if req.PolygonGeojson != "" {
		geojson = req.PolygonGeojson
	} else if len(req.Coordinates) > 0 {
		geojson, err = pointsToPolygonGeoJSON(req.Coordinates)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid coordinates: %v", err)
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "either polygon_geojson or coordinates is required")
	}

	// Validate GeoJSON (optional)
	var geom map[string]interface{}
	if err := json.Unmarshal([]byte(geojson), &geom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid GeoJSON: %v", err)
	}
	geomType, ok := geom["type"].(string)
	if !ok || (geomType != "Polygon" && geomType != "MultiPolygon") {
		return nil, status.Error(codes.InvalidArgument, "geometry must be Polygon or MultiPolygon")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 1000
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	features, total, err := h.featureUsecase.GetFeaturesByPolygon(ctx, []byte(geojson), limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query features: %v", err)
	}

	pbFeatures := make([]*tqdpb.Feature, 0, len(features))
	for _, f := range features {
		pbFeat, err := convertFeatureToProto(f)
		if err != nil {
			continue
		}
		pbFeatures = append(pbFeatures, pbFeat)
	}

	return &tqdpb.ListFeaturesResponse{
		Features: pbFeatures,
		Total:    total,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

func convertFeatureToProto(f *domain.Feature) (*tqdpb.Feature, error) {
	var geomMap map[string]interface{}
	if err := json.Unmarshal([]byte(f.Geometry), &geomMap); err != nil {
		return nil, err
	}
	geomStruct, err := structpb.NewStruct(geomMap)
	if err != nil {
		return nil, err
	}

	propsBytes, err := json.Marshal(f.Properties)
	if err != nil {
		return nil, err
	}
	var propsMap map[string]interface{}
	if err := json.Unmarshal(propsBytes, &propsMap); err != nil {
		return nil, err
	}
	propsStruct, err := structpb.NewStruct(propsMap)
	if err != nil {
		return nil, err
	}

	sessionID := int32(0)
	if f.SessionID != nil {
		sessionID = int32(*f.SessionID)
	}
	extent := int32(0)
	if f.Extent != nil {
		extent = int32(*f.Extent)
	}

	return &tqdpb.Feature{
		Id:          uint64(f.ID),
		SessionId:   sessionID,
		Layer:       f.Layer,
		Geometry:    geomStruct,
		Properties:  propsStruct,
		Tile:        f.Tile,
		Extent:      extent,
		FeatureHash: f.FeatureHash,
		CreatedAt:   f.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   f.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func pointsToPolygonGeoJSON(points []*tqdpb.Point) (string, error) {
	if len(points) < 3 {
		return "", fmt.Errorf("at least 3 points required for polygon")
	}
	coords := make([][]float64, len(points))
	for i, p := range points {
		coords[i] = []float64{p.Lng, p.Lat}
	}
	// Ensure polygon is closed
	if coords[0][0] != coords[len(coords)-1][0] || coords[0][1] != coords[len(coords)-1][1] {
		coords = append(coords, coords[0])
	}
	geojson := map[string]interface{}{
		"type":        "Polygon",
		"coordinates": [][][]float64{coords},
	}
	bytes, err := json.Marshal(geojson)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
