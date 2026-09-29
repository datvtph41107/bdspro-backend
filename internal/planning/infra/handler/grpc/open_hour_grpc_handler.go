package handler_grpc

import (
	"context"
	"fmt"
	"log/slog"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OpenHourGrpcHandler struct {
	tqdpb.UnimplementedOpenHourServiceServer
	usecase usecase.OpenHourUsecase
	mapper  *mapper.OpenHourMapper
}

func NewOpenHourGrpcHandler(
	usecase usecase.OpenHourUsecase,
	mapper *mapper.OpenHourMapper,
) *OpenHourGrpcHandler {
	return &OpenHourGrpcHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

func (h *OpenHourGrpcHandler) CreateOpenHour(ctx context.Context, req *tqdpb.OpenHour) (*tqdpb.OpenHour, error) {
	// Validate required fields
	if req.PoiId == nil || *req.PoiId == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	// Convert proto to request DTO
	createReq, err := h.mapper.FromProto(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Call usecase
	result, err := h.usecase.Create(ctx, createReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error creating open hour: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to create open hour: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

func (h *OpenHourGrpcHandler) GetOpenHour(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.OpenHour, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	result, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting open hour: %v", err))
		return nil, status.Errorf(codes.NotFound, "open hour not found: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

func (h *OpenHourGrpcHandler) UpdateOpenHour(ctx context.Context, req *tqdpb.OpenHour) (*tqdpb.OpenHour, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	updateReq := &dto.UpdateOpenHourRequest{}

	if req.DayOfWeek != "" {
		day, err := dto.StringToDayOfWeek(req.DayOfWeek)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid day of week: %v", err)
		}
		updateReq.DayOfWeek = &day
	}
	if req.StartTime != "" {
		updateReq.OpenTime = &req.StartTime
	}
	if req.EndTime != "" {
		updateReq.CloseTime = &req.EndTime
	}
	if req.Notes != "" {
		updateReq.Note = &req.Notes
	}
	updateReq.IsOpen = &req.IsOpen
	updateReq.IsRecurring = &req.IsRecurring

	if req.StartDate != "" {
		updateReq.StartDate = &req.StartDate
	}
	if req.EndDate != "" {
		updateReq.EndDate = &req.EndDate
	}
	if req.Status != 0 {
		status := int(req.Status)
		updateReq.Status = &status
	}
	if req.Type != 0 {
		typ := int(req.Type)
		updateReq.Type = &typ
	}

	result, err := h.usecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error updating open hour: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to update open hour: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

func (h *OpenHourGrpcHandler) DeleteOpenHour(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "open hour id is required")
	}

	if err := h.usecase.Delete(ctx, req.Id); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error deleting open hour: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete open hour: %v", err)
	}

	return &sharepb.SubmitResponse{
		Message: "Open hour deleted successfully",
	}, nil
}

func (h *OpenHourGrpcHandler) ListOpenHours(ctx context.Context, req *tqdpb.ListOpenHoursRequest) (*tqdpb.ListOpenHoursResponse, error) {
	filter := &dto.OpenHourFilter{
		Pagable: _dto.Pagable{
			Page: normalizePage(req.Page),
			Size: normalizeSize(req.Size),
		},
		POIID: req.PoiId,
	}

	if req.DayOfWeek != "" {
		day, err := dto.StringToDayOfWeek(req.DayOfWeek)
		if err == nil {
			filter.DayOfWeek = &day
		}
	}

	if req.IsOpen {
		filter.IsOpen = &req.IsOpen
	}

	result, err := h.usecase.List(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error listing open hours: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to list open hours: %v", err)
	}

	hours := h.mapper.ToProtoListFromResponses(result.Data)

	return &tqdpb.ListOpenHoursResponse{
		Data:  hours,
		Total: result.Total,
	}, nil
}

func (h *OpenHourGrpcHandler) GetOpenHourTimesByDay(ctx context.Context, req *tqdpb.GetOpenHourTimesByDayRequest) (*tqdpb.GetOpenHourTimesByDayResponse, error) {
	if req.DayOfWeek == "" {
		return nil, status.Error(codes.InvalidArgument, "day of week is required")
	}

	getReq := &dto.GetTimesByDayRequest{
		DayOfWeek: req.DayOfWeek,
	}

	results, err := h.usecase.GetTimesByDay(ctx, getReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting open hour times: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get open hour times: %v", err)
	}

	var response []*tqdpb.OpenHourTimesByDay
	for _, dayGroup := range results {
		times := make([]*tqdpb.OpenHourTime, len(dayGroup.Times))
		for i, t := range dayGroup.Times {
			times[i] = &tqdpb.OpenHourTime{
				Id:        t.ID,
				PoiId:     t.POIID,
				StartTime: t.OpenTime,
				EndTime:   t.CloseTime,
				IsOpen:    t.IsOpen,
				Notes:     t.Note,
			}
		}
		response = append(response, &tqdpb.OpenHourTimesByDay{
			DayOfWeek: dayGroup.DayName,
			Times:     times,
		})
	}

	return &tqdpb.GetOpenHourTimesByDayResponse{
		Data: response,
	}, nil
}

// Helper functions
func normalizePage(page uint32) uint32 {
	if page < 1 {
		return 1
	}
	return page
}

func normalizeSize(size uint32) uint32 {
	if size < 1 {
		return 10
	}
	if size > 100 {
		return 100
	}
	return size
}
