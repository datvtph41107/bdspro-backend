package handler_grpc

import (
	"context"
	"errors"

	tqdpb "pb/types/tqd"
	domain "tqd/internal/domain/mappoint"
	usecase "tqd/internal/usecase/mappoint"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MapPointGrpcHandler struct {
	tqdpb.UnimplementedMapPointServiceServer
	service *usecase.Service
}

func NewMapPointGrpcHandler(service *usecase.Service) *MapPointGrpcHandler {
	return &MapPointGrpcHandler{service: service}
}

func mapPointError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, usecase.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "map point operation failed")
	}
}

func pointMessage(point domain.Point) *tqdpb.MapPoint {
	return &tqdpb.MapPoint{Id: point.ID, Lat: point.Latitude, Lng: point.Longitude, Name: point.Name}
}

func pointList(points []domain.Point) *tqdpb.MapPointList {
	result := &tqdpb.MapPointList{Data: make([]*tqdpb.MapPoint, 0, len(points))}
	for _, point := range points {
		result.Data = append(result.Data, pointMessage(point))
	}
	return result
}

func (h *MapPointGrpcHandler) CreateMapPoint(ctx context.Context, request *tqdpb.CreateMapPointRequest) (*tqdpb.MapPoint, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	point, err := h.service.Create(ctx, domain.Point{Latitude: request.Lat, Longitude: request.Lng, Name: request.Name})
	if err != nil {
		return nil, mapPointError(err)
	}
	return pointMessage(point), nil
}

func (h *MapPointGrpcHandler) GetMapPoint(ctx context.Context, request *tqdpb.GetMapPointRequest) (*tqdpb.MapPoint, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	point, err := h.service.Get(ctx, request.Id)
	if err != nil {
		return nil, mapPointError(err)
	}
	return pointMessage(point), nil
}

func (h *MapPointGrpcHandler) UpdateMapPoint(ctx context.Context, request *tqdpb.UpdateMapPointRequest) (*tqdpb.MapPoint, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	point, err := h.service.Update(ctx, domain.Point{ID: request.Id, Latitude: request.Lat, Longitude: request.Lng, Name: request.Name})
	if err != nil {
		return nil, mapPointError(err)
	}
	return pointMessage(point), nil
}

func (h *MapPointGrpcHandler) DeleteMapPoint(ctx context.Context, request *tqdpb.DeleteMapPointRequest) (*tqdpb.DeleteMapPointResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := h.service.Delete(ctx, request.Id); err != nil {
		return nil, mapPointError(err)
	}
	return &tqdpb.DeleteMapPointResponse{Message: "Map point deleted successfully"}, nil
}

func (h *MapPointGrpcHandler) FindNearbyMapPoints(ctx context.Context, request *tqdpb.FindNearbyMapPointsRequest) (*tqdpb.MapPointList, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	points, err := h.service.FindNearby(ctx, domain.Coordinate{Latitude: request.Lat, Longitude: request.Lng}, request.Radius)
	if err != nil {
		return nil, mapPointError(err)
	}
	return pointList(points), nil
}

func (h *MapPointGrpcHandler) FindMapPointsByPolygon(ctx context.Context, request *tqdpb.FindMapPointsByPolygonRequest) (*tqdpb.MapPointList, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	coordinates := make([]domain.Coordinate, 0, len(request.Coordinates))
	for _, coordinate := range request.Coordinates {
		if coordinate != nil {
			coordinates = append(coordinates, domain.Coordinate{Latitude: coordinate.Lat, Longitude: coordinate.Lng})
		}
	}
	points, err := h.service.FindWithinPolygon(ctx, coordinates)
	if err != nil {
		return nil, mapPointError(err)
	}
	return pointList(points), nil
}
