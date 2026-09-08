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

	// Convert proto to CreateContactLabelRequestDTO
	createReq := &dto.CreateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		SortOrder:   int(req.SortOrder),
		IsSystem:    req.IsSystem,
		Type:        int(req.Type),
	}

	// Call usecase
	result, err := h.contactLabelUsecase.Create(ctx, createReq)
	if err != nil {
		log.Printf("error creating contact label: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create contact label: %v", err)
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// GetContactLabel handles gRPC GetContactLabel request
func (h *ContactLabelGrpcHandler) GetContactLabel(ctx context.Context, req *tqdpb.GetContactLabelRequest) (*tqdpb.ContactLabel, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	result, err := h.contactLabelUsecase.GetByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting contact label: %v", err)
		return nil, status.Errorf(codes.NotFound, "contact label not found: %v", err)
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// UpdateContactLabel handles gRPC UpdateContactLabel request
func (h *ContactLabelGrpcHandler) UpdateContactLabel(ctx context.Context, req *tqdpb.ContactLabel) (*tqdpb.ContactLabel, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	// Convert proto to UpdateContactLabelRequestDTO
	updateReq := &dto.UpdateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		SortOrder:   int(req.SortOrder),
		IsSystem:    req.IsSystem,
		Type:        int(req.Type),
	}

	result, err := h.contactLabelUsecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		log.Printf("error updating contact label: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update contact label: %v", err)
	}

	return h.contactLabelMapper.ToProto(result), nil
}

// DeleteContactLabel handles gRPC DeleteContactLabel request
func (h *ContactLabelGrpcHandler) DeleteContactLabel(ctx context.Context, req *tqdpb.DeleteContactLabelRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "contact label id is required")
	}

	err := h.contactLabelUsecase.Delete(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting contact label: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete contact label: %v", err)
	}

	return &sharepb.SubmitResponse{
		Message: "Contact label deleted successfully",
	}, nil
}

// ListContactLabels handles gRPC ListContactLabels request
func (h *ContactLabelGrpcHandler) ListContactLabels(ctx context.Context, req *tqdpb.ListContactLabelsRequest) (*tqdpb.ListContactLabelsResponse, error) {
	page := uint32(req.Page)
	if page < 1 {
		page = 1
	}
	size := uint32(req.Size)
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	listReq := &dto.ListContactLabelsRequestDTO{
		Pagable: _dto.Pagable{
			Page: page,
			Size: size,
		},
		Search: req.Search,
	}

	// Handle optional filters
	if req.IsActive {
		listReq.IsActive = &req.IsActive
	}
	if req.IsSystem {
		listReq.IsSystem = &req.IsSystem
	}

	result, err := h.contactLabelUsecase.List(ctx, listReq)
	if err != nil {
		log.Printf("error listing contact labels: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list contact labels: %v", err)
	}

	// Convert DTO list to proto list
	contactLabels := make([]*tqdpb.ContactLabel, len(result.Data))
	for i, item := range result.Data {
		contactLabels[i] = h.contactLabelMapper.ToProto(&item)
	}

	return &tqdpb.ListContactLabelsResponse{
		Data:  contactLabels,
		Total: result.Total,
	}, nil
}
