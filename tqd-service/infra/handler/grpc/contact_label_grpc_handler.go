package handler_grpc

import (
	"context"
	"log"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ContactLabelGrpcHandler implements gRPC handlers for ContactLabel service
type ContactLabelGrpcHandler struct {
	tqdpb.UnimplementedContactLabelServiceServer
	contactLabelUsecase usecase.ContactLabelUsecase
	contactLabelMapper  *mapper.ContactLabelMapper
}

// NewContactLabelGrpcHandler creates a new ContactLabelGrpcHandler
func NewContactLabelGrpcHandler(
	contactLabelUsecase usecase.ContactLabelUsecase,
	contactLabelMapper *mapper.ContactLabelMapper,
) *ContactLabelGrpcHandler {
	return &ContactLabelGrpcHandler{
		contactLabelUsecase: contactLabelUsecase,
		contactLabelMapper:  contactLabelMapper,
	}
}

// CreateContactLabel handles gRPC CreateContactLabel request
func (h *ContactLabelGrpcHandler) CreateContactLabel(ctx context.Context, req *tqdpb.ContactLabel) (*tqdpb.ContactLabel, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "contact label name is required")
	}

	label := &dto.CreateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		IsSystem:    req.IsSystem,
		Type:        req.Type,
	}

	result, err := h.contactLabelUsecase.Create(ctx, label)
	if err != nil {
		log.Printf("error creating contact label: %v", err)
		return nil, status.Error(codes.Internal, "failed to create contact label")
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// GetContactLabel handles gRPC GetContactLabel request
func (h *ContactLabelGrpcHandler) GetContactLabel(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.ContactLabel, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	result, err := h.contactLabelUsecase.GetByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting contact label: %v", err)
		return nil, status.Error(codes.NotFound, "contact label not found")
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// UpdateContactLabel handles gRPC UpdateContactLabel request
func (h *ContactLabelGrpcHandler) UpdateContactLabel(ctx context.Context, req *tqdpb.ContactLabel) (*tqdpb.ContactLabel, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	label := &dto.UpdateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		IsSystem:    req.IsSystem,
		Type:        req.Type,
	}

	result, err := h.contactLabelUsecase.Update(ctx, req.Id, label)
	if err != nil {
		log.Printf("error updating contact label: %v", err)
		return nil, status.Error(codes.Internal, "failed to update contact label")
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// DeleteContactLabel handles gRPC DeleteContactLabel request
func (h *ContactLabelGrpcHandler) DeleteContactLabel(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	err := h.contactLabelUsecase.Delete(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting contact label: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete contact label")
	}

	return &sharepb.SubmitResponse{
		Message: "Contact label deleted successfully",
	}, nil
}

// ListContactLabels handles gRPC ListContactLabels request
func (h *ContactLabelGrpcHandler) ListContactLabels(ctx context.Context, req *tqdpb.ListContactLabelsRequest) (*tqdpb.ListContactLabelsResponse, error) {
	listReq := &dto.ListContactLabelsRequestDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Search:   req.Search,
		IsActive: &req.IsActive,
		IsSystem: &req.IsSystem,
	}

	result, err := h.contactLabelUsecase.List(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contact labels: %v", err)
	}

	contactLabels := make([]*tqdpb.ContactLabel, len(result.Data))
	for i, item := range result.Data {
		contactLabels[i] = &tqdpb.ContactLabel{
			Id:           uint64(item.ID),
			Name:         item.Name,
			Code:         item.Code,
			Description:  item.Description,
			Color:        item.Color,
			Icon:         item.Icon,
			IsActive:     item.IsActive,
			SortOrder:    int32(item.SortOrder),
			ContactCount: int32(item.ContactCount),
			IsSystem:     item.IsSystem,
			Status:       int32(item.Status),
			Type:         int32(item.Type),
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	// Return response
	return &tqdpb.ListContactLabelsResponse{
		Data:  contactLabels,
		Total: result.Total,
	}, nil
}
