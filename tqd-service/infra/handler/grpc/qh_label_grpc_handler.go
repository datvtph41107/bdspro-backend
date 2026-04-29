package handler_grpc

import (
	"context"
	"fmt"

	tqdpb "pb/types/tqd"
	shared "pb/types/shared"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QHLabelGrpcHandler xử lý gRPC cho QHLabelService
type QHLabelGrpcHandler struct {
	tqdpb.UnimplementedQHLabelServiceServer
	labelUsecase usecase.QHLabelUsecase
}

func NewQHLabelGrpcHandler(labelUsecase usecase.QHLabelUsecase) *QHLabelGrpcHandler {
	return &QHLabelGrpcHandler{labelUsecase: labelUsecase}
}

func (h *QHLabelGrpcHandler) CreateQHLabel(ctx context.Context, req *tqdpb.CreateQHLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId không được để trống")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không được để trống")
	}

	label := &qh_domain.QHLabel{
		LayerID:     req.LayerId,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   int(req.SortOrder),
	}

	result, err := h.labelUsecase.Create(ctx, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoLabel(result), nil
}

func (h *QHLabelGrpcHandler) GetQHLabel(ctx context.Context, req *tqdpb.GetQHLabelRequest) (*tqdpb.QHLabelResponse, error) {
	result, err := h.labelUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoLabel(result), nil
}

func (h *QHLabelGrpcHandler) UpdateQHLabel(ctx context.Context, req *tqdpb.UpdateQHLabelRequest) (*tqdpb.QHLabelResponse, error) {
	label := &qh_domain.QHLabel{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   int(req.SortOrder),
	}

	result, err := h.labelUsecase.Update(ctx, req.Id, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoLabel(result), nil
}

func (h *QHLabelGrpcHandler) DeleteQHLabel(ctx context.Context, req *tqdpb.DeleteQHLabelRequest) (*shared.SubmitResponse, error) {
	if err := h.labelUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &shared.SubmitResponse{
		Id:      req.Id,
		Message: "Xóa nhãn thành công",
	}, nil
}

func (h *QHLabelGrpcHandler) ListQHLabelsByLayer(ctx context.Context, req *tqdpb.ListQHLabelsByLayerRequest) (*tqdpb.ListQHLabelsResponse, error) {
	labels, err := h.labelUsecase.ListByLayerID(ctx, req.LayerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.QHLabelResponse, len(labels))
	for i, l := range labels {
		data[i] = toProtoLabel(&l)
	}

	return &tqdpb.ListQHLabelsResponse{
		Data:  data,
		Total: int32(len(labels)),
	}, nil
}

// toProtoLabel convert domain → proto
func toProtoLabel(l *qh_domain.QHLabel) *tqdpb.QHLabelResponse {
	return &tqdpb.QHLabelResponse{
		Id:          l.ID,
		LayerId:     l.LayerID,
		Name:        l.Name,
		Color:       l.Color,
		Description: l.Description,
		SortOrder:   int32(l.SortOrder),
		CreatedAt:   fmt.Sprintf("%v", l.CreatedAt),
		UpdatedAt:   fmt.Sprintf("%v", l.UpdatedAt),
	}
}
