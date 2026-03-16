package handler_grpc

import (
	"context"
	"log"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OpenHourGrpcHandler implements gRPC handlers for OpenHour service
type OpenHourGrpcHandler struct {
	tqdpb.UnimplementedOpenHourServiceServer
	openHourUsecase *usecase.OpenHourUsecase
	openHourMapper  *mapper.OpenHourMapper
}

// NewOpenHourGrpcHandler creates a new OpenHourGrpcHandler
func NewOpenHourGrpcHandler(
	openHourUsecase *usecase.OpenHourUsecase,
	openHourMapper *mapper.OpenHourMapper,
) *OpenHourGrpcHandler {
	return &OpenHourGrpcHandler{
		openHourUsecase: openHourUsecase,
		openHourMapper:  openHourMapper,
	}
}

// CreateOpenHour handles gRPC CreateOpenHour request
func (h *OpenHourGrpcHandler) CreateOpenHour(ctx context.Context, req *tqdpb.OpenHour) (*tqdpb.OpenHour, error) {
	if req.PoiId == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	createReq := &dto.CreateOpenHourRequestDTO{
		PoiID:       req.PoiId,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsOpen:      req.IsOpen,
		IsRecurring: req.IsRecurring,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Notes:       req.Notes,
	}

	result, err := h.openHourUsecase.CreateOpenHour(ctx, createReq)
	if err != nil {
		log.Printf("error creating open hour: %v", err)
		return nil, status.Error(codes.Internal, "failed to create open hour")
	}

	return h.openHourMapper.ToProto(result), nil
}

// GetOpenHour handles gRPC GetOpenHour request
func (h *OpenHourGrpcHandler) GetOpenHour(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.OpenHour, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	result, err := h.openHourUsecase.GetOpenHourByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting open hour: %v", err)
		return nil, status.Error(codes.NotFound, "open hour not found")
	}

	return h.openHourMapper.ToProto(result), nil
}

// UpdateOpenHour handles gRPC UpdateOpenHour request
func (h *OpenHourGrpcHandler) UpdateOpenHour(ctx context.Context, req *tqdpb.OpenHour) (*tqdpb.OpenHour, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	updateReq := &dto.OpenHourDTO{
		ID:          req.Id,
		PoiID:       req.PoiId,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsOpen:      req.IsOpen,
		IsRecurring: req.IsRecurring,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Notes:       req.Notes,
	}

	result, err := h.openHourUsecase.UpdateOpenHour(ctx, updateReq)
	if err != nil {
		log.Printf("error updating open hour: %v", err)
		return nil, status.Error(codes.Internal, "failed to update open hour")
	}

	return h.openHourMapper.ToProto(result), nil
}

// DeleteOpenHour handles gRPC DeleteOpenHour request
func (h *OpenHourGrpcHandler) DeleteOpenHour(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	err := h.openHourUsecase.DeleteOpenHour(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting open hour: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete open hour")
	}

	return &sharepb.SubmitResponse{
		Code:    200,
		Message: "Open hour deleted successfully",
		Data:    "",
	}, nil
}

// ListOpenHours handles gRPC ListOpenHours request
func (h *OpenHourGrpcHandler) ListOpenHours(ctx context.Context, req *tqdpb.ListOpenHoursRequest) (*tqdpb.ListOpenHoursResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.OpenHourFilterDTO{
		Page:      int(req.Page),
		Size:      int(req.Size),
		PoiID:     req.PoiId,
		DayOfWeek: req.DayOfWeek,
		IsOpen:    req.IsOpen,
	}

	results, total, err := h.openHourUsecase.ListOpenHours(ctx, filters)
	if err != nil {
		log.Printf("error listing open hours: %v", err)
		return nil, status.Error(codes.Internal, "failed to list open hours")
	}

	hours := make([]*tqdpb.OpenHour, len(results))
	for i, result := range results {
		hours[i] = h.openHourMapper.ToProto(result)
	}

	return &tqdpb.ListOpenHoursResponse{
		Data:  hours,
		Total: int64(total),
	}, nil
}

// GetOpenHourTimesByDay handles gRPC GetOpenHourTimesByDay request
func (h *OpenHourGrpcHandler) GetOpenHourTimesByDay(ctx context.Context, req *tqdpb.GetOpenHourTimesByDayRequest) (*tqdpb.GetOpenHourTimesByDayResponse, error) {
	if req.DayOfWeek == "" {
		return nil, status.Error(codes.InvalidArgument, "day of week is required")
	}

	results, err := h.openHourUsecase.GetOpenHourTimesByDay(ctx, req.DayOfWeek)
	if err != nil {
		log.Printf("error getting open hour times by day: %v", err)
		return nil, status.Error(codes.Internal, "failed to get open hour times by day")
	}

	// Group results by dayOfWeek
	dayMap := make(map[string][]*tqdpb.OpenHourTime)
	for _, result := range results {
		day := result.DayOfWeek
		if day == "" {
			continue
		}
		dayMap[day] = append(dayMap[day], &tqdpb.OpenHourTime{
			Id:        result.ID,
			PoiId:     result.PoiID,
			StartTime: result.StartTime,
			EndTime:   result.EndTime,
			IsOpen:    result.IsOpen,
			Notes:     result.Notes,
		})
	}

	var response []*tqdpb.OpenHourTimesByDay
	for day, times := range dayMap {
		response = append(response, &tqdpb.OpenHourTimesByDay{
			DayOfWeek: day,
			Times:     times,
		})
	}

	return &tqdpb.GetOpenHourTimesByDayResponse{
		Data: response,
	}, nil
}
