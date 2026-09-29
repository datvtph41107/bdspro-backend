// internal/handler/identifier_handler.go
package handler

import (
	"context"
	"errors"

	_utils "common/utils"

	"bdspro/infra/mapper"
	"bdspro/internal/domain"
	"bdspro/internal/usecases"

	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyIdentifierHandler struct {
	bdspropb.UnimplementedIdentifierServiceServer
	identifierUsecase *usecases.PropertyIdentifierUsecase
	mapper            *mapper.PropertyIdentifierMapper
}

// NewIdentifierHandler creates a new identifier handler
func NewIdentifierHandler(
	identifierUsecase *usecases.PropertyIdentifierUsecase,
	mapper *mapper.PropertyIdentifierMapper,
) *PropertyIdentifierHandler {
	return &PropertyIdentifierHandler{
		identifierUsecase: identifierUsecase,
		mapper:            mapper,
	}
}

// RegisterIdentifier handles registerIdentifier request
func (h *PropertyIdentifierHandler) RegisterIdentifier(
	ctx context.Context,
	req *bdspropb.RegisterIdentifierRequest,
) (*bdspropb.RegisterIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Identifier == nil {
		return nil, status.Error(codes.InvalidArgument, "identifier cannot be nil")
	}
	if req.RequestedBy == "" {
		return nil, status.Error(codes.InvalidArgument, "requested_by is required")
	}

	// 2. Proto → DTO
	identifierDTO := h.mapper.FromProtoToDTO(req.Identifier)

	// 3. Set timestamps
	now := _utils.TimeNowPtr()
	identifierDTO.CreatedAt = _utils.FormatTimeToStringCustom(now)
	identifierDTO.UpdatedAt = _utils.FormatTimeToStringCustom(now)

	// 4. DTO → Domain
	identifierDomain := h.mapper.ToDomain(identifierDTO)

	// 5. Call usecase (business logic)
	createdDomain, err := h.identifierUsecase.Create(ctx, identifierDomain, req.RequestedBy)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 6. Domain → DTO → Proto
	createdDTO := h.mapper.ToDTO(createdDomain)
	createdProto := h.mapper.ToProto(createdDTO)

	// 7. Return response
	return &bdspropb.RegisterIdentifierResponse{
		Identifier: createdProto,
		Message:    "Identifier registered successfully",
	}, nil
}

// GetIdentifier handles getIdentifier request
func (h *PropertyIdentifierHandler) GetIdentifier(
	ctx context.Context,
	req *bdspropb.GetIdentifierRequest,
) (*bdspropb.GetIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 2. Call usecase
	identifierDomain, err := h.identifierUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 3. Domain → DTO → Proto
	identifierDTO := h.mapper.ToDTO(identifierDomain)
	identifierProto := h.mapper.ToProto(identifierDTO)

	// 4. Return response
	return &bdspropb.GetIdentifierResponse{
		Identifier: identifierProto,
	}, nil
}

// UpdateIdentifier handles updateIdentifier request
func (h *PropertyIdentifierHandler) UpdateIdentifier(
	ctx context.Context,
	req *bdspropb.UpdateIdentifierRequest,
) (*bdspropb.UpdateIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Identifier == nil {
		return nil, status.Error(codes.InvalidArgument, "identifier cannot be nil")
	}
	if req.RequestedBy == "" {
		return nil, status.Error(codes.InvalidArgument, "requested_by is required")
	}

	// 2. Proto → DTO
	identifierDTO := h.mapper.FromProtoToDTO(req.Identifier)

	// 3. Update timestamp
	now := _utils.TimeNowPtr()
	identifierDTO.UpdatedAt = _utils.FormatTimeToStringCustom(now)

	// 4. DTO → Domain
	identifierDomain := h.mapper.ToDomain(identifierDTO)

	// 5. Call usecase
	updatedDomain, err := h.identifierUsecase.Update(ctx, identifierDomain, req.RequestedBy)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 6. Domain → DTO → Proto
	updatedDTO := h.mapper.ToDTO(updatedDomain)
	updatedProto := h.mapper.ToProto(updatedDTO)

	// 7. Return response
	return &bdspropb.UpdateIdentifierResponse{
		Identifier: updatedProto,
		Message:    "Identifier updated successfully",
	}, nil
}

// DeleteIdentifier handles deleteIdentifier request
func (h *PropertyIdentifierHandler) DeleteIdentifier(
	ctx context.Context,
	req *bdspropb.DeleteIdentifierRequest,
) (*bdspropb.DeleteIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.RequestedBy == "" {
		return nil, status.Error(codes.InvalidArgument, "requested_by is required")
	}

	// 2. Call usecase
	err := h.identifierUsecase.Delete(ctx, req.Id, req.RequestedBy)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 3. Return response
	return &bdspropb.DeleteIdentifierResponse{
		Success: true,
		Message: "Identifier deleted successfully",
	}, nil
}

// SearchIdentifier handles searchIdentifier request
func (h *PropertyIdentifierHandler) SearchIdentifier(
	ctx context.Context,
	req *bdspropb.SearchIdentifierRequest,
) (*bdspropb.SearchIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	// 2. Proto → DTO (search criteria and pagination)
	criteriaDTO := h.mapper.SearchCriteriaFromProto(req.Criteria)
	paginationDTO := h.mapper.PaginationFromProto(req.Pagination)

	// 3. Call usecase
	identifiersDomain, total, err := h.identifierUsecase.Search(ctx, criteriaDTO, paginationDTO)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 4. Domain list → DTO list → Proto list
	identifiersDTO := h.mapper.ToDTOList(identifiersDomain)
	identifiersProto := h.mapper.ToProtoList(identifiersDTO)

	// 5. Return response
	return &bdspropb.SearchIdentifierResponse{
		Identifiers: identifiersProto,
		Total:       total,
		Page:        paginationDTO.Page,
		Limit:       paginationDTO.Limit,
	}, nil
}

// ValidateIdentifier handles validateIdentifier request
func (h *PropertyIdentifierHandler) ValidateIdentifier(
	ctx context.Context,
	req *bdspropb.ValidateIdentifierRequest,
) (*bdspropb.ValidateIdentifierResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	// 2. Call usecase
	identifierDomain, valid, err := h.identifierUsecase.ValidateIdentifier(ctx, req.Code)
	if err != nil {
		return nil, h.mapErrorToGRPC(err)
	}

	// 3. If not valid, return simple response
	if !valid || identifierDomain == nil {
		return &bdspropb.ValidateIdentifierResponse{
			Valid:   false,
			Message: "Identifier does not exist",
		}, nil
	}

	// 4. Domain → DTO → Proto
	identifierDTO := h.mapper.ToDTO(identifierDomain)
	identifierProto := h.mapper.ToProto(identifierDTO)

	// 5. Return response
	return &bdspropb.ValidateIdentifierResponse{
		Valid:      true,
		Message:    "Identifier is valid",
		Identifier: identifierProto,
	}, nil
}

// TransferOwnership handles ownership transfer
// func (h *PropertyIdentifierHandler) TransferOwnership(
// 	ctx context.Context,
// 	req *bdspropb.TransferOwnershipRequest,
// ) (*bdspropb.TransferOwnershipResponse, error) {
// 	// 1. Validate request
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
// 	}
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}
// 	if req.NewOwnerId == "" {
// 		return nil, status.Error(codes.InvalidArgument, "new_owner_id is required")
// 	}
// 	if req.RequestedBy == "" {
// 		return nil, status.Error(codes.InvalidArgument, "requested_by is required")
// 	}

// 	// 2. Call usecase
// 	updatedDomain, err := h.identifierUsecase.TransferOwnership(ctx, req.Id, req.NewOwnerId, req.RequestedBy)
// 	if err != nil {
// 		return nil, h.mapErrorToGRPC(err)
// 	}

// 	// 3. Domain → DTO → Proto
// 	updatedDTO := h.mapper.ToDTO(updatedDomain)
// 	updatedProto := h.mapper.ToProto(updatedDTO)

// 	// 4. Return response
// 	return &bdspropb.TransferOwnershipResponse{
// 		Identifier: updatedProto,
// 		Message:    "Ownership transferred successfully",
// 	}, nil
// }

// // UpdateLegalStatus handles legal status update
// func (h *PropertyIdentifierHandler) UpdateLegalStatus(
// 	ctx context.Context,
// 	req *bdspropb.UpdateLegalStatusRequest,
// ) (*bdspropb.UpdateLegalStatusResponse, error) {
// 	// 1. Validate request
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
// 	}
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}
// 	if req.Status == "" {
// 		return nil, status.Error(codes.InvalidArgument, "status is required")
// 	}
// 	if req.RequestedBy == "" {
// 		return nil, status.Error(codes.InvalidArgument, "requested_by is required")
// 	}

// 	// 2. Call usecase
// 	updatedDomain, err := h.identifierUsecase.UpdateLegalStatus(ctx, req.Id, req.Status, req.RequestedBy)
// 	if err != nil {
// 		return nil, h.mapErrorToGRPC(err)
// 	}

// 	// 3. Domain → DTO → Proto
// 	updatedDTO := h.mapper.ToDTO(updatedDomain)
// 	updatedProto := h.mapper.ToProto(updatedDTO)

// 	// 4. Return response
// 	return &bdspropb.UpdateLegalStatusResponse{
// 		Identifier: updatedProto,
// 		Message:    "Legal status updated successfully",
// 	}, nil
// }

// mapErrorToGRPC maps domain errors to gRPC status codes
func (h *PropertyIdentifierHandler) mapErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrDuplicate):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrConcurrentUpdate):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
