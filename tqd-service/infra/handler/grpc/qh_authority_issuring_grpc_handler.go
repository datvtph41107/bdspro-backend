package handler_grpc

import (
	"context"
	"strings"

	_dto "common/domain/dto"
	"common/fault"
	_utils "common/utils"

	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

// QHAuthorityIssuringGrpcHandler xử lý gRPC QHAuthorityIssuringService.
type QHAuthorityIssuringGrpcHandler struct {
	tqdpb.UnimplementedQHAuthorityIssuringServiceServer
	uc usecase.QHAuthorityIssuringUsecase
}

func NewQHAuthorityIssuringGrpcHandler(uc usecase.QHAuthorityIssuringUsecase) *QHAuthorityIssuringGrpcHandler {
	return &QHAuthorityIssuringGrpcHandler{uc: uc}
}

// CreateAuthorityIssuring — POST /v2/tqd/qh/admin/authority-issuring
func (h *QHAuthorityIssuringGrpcHandler) CreateAuthorityIssuring(ctx context.Context, req *tqdpb.CreateAuthorityIssuringRequest) (*tqdpb.QHAuthorityIssuringResponse, error) {
	if req == nil {
		return nil, qhAuthorityIssuringValidation(
			"tqd.qh_authority_issuring.request_required",
			"request is required",
			"",
		)
	}
	row := &qh_domain.QHAuthorityIssuring{
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Description: req.Description,
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		return nil, qhAuthorityIssuringError(err)
	}
	return toQHAuthorityIssuringPB(out), nil
}

// GetAuthorityIssuring — GET /v2/tqd/qh/admin/authority-issuring/{id}
func (h *QHAuthorityIssuringGrpcHandler) GetAuthorityIssuring(ctx context.Context, req *tqdpb.GetAuthorityIssuringRequest) (*tqdpb.QHAuthorityIssuringResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, qhAuthorityIssuringValidation(
			"tqd.qh_authority_issuring.id_required",
			"id is required",
			"id",
		)
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, qhAuthorityIssuringError(err)
	}
	return toQHAuthorityIssuringPB(out), nil
}

// UpdateAuthorityIssuring — PUT /v2/tqd/qh/admin/authority-issuring/{id}
func (h *QHAuthorityIssuringGrpcHandler) UpdateAuthorityIssuring(ctx context.Context, req *tqdpb.UpdateAuthorityIssuringRequest) (*tqdpb.QHAuthorityIssuringResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, qhAuthorityIssuringValidation(
			"tqd.qh_authority_issuring.id_required",
			"id is required",
			"id",
		)
	}
	in := &usecase.QHAuthorityIssuringUpdateInput{}
	if req.Name != nil {
		in.Name = req.Name
	}
	if req.Code != nil {
		in.Code = req.Code
	}
	if req.Description != nil {
		in.Description = req.Description
	}
	if in.Name == nil && in.Code == nil && in.Description == nil {
		return nil, qhAuthorityIssuringValidation(
			"tqd.qh_authority_issuring.update_fields_required",
			"at least one field to update is required",
			"",
		)
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		return nil, qhAuthorityIssuringError(err)
	}
	return toQHAuthorityIssuringPB(out), nil
}

// DeleteAuthorityIssuring — DELETE /v2/tqd/qh/admin/authority-issuring/{id}
func (h *QHAuthorityIssuringGrpcHandler) DeleteAuthorityIssuring(ctx context.Context, req *tqdpb.DeleteAuthorityIssuringRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, qhAuthorityIssuringValidation(
			"tqd.qh_authority_issuring.id_required",
			"id is required",
			"id",
		)
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, qhAuthorityIssuringError(err)
	}
	return &emptypb.Empty{}, nil
}

// ListAuthorityIssuring — GET /v2/tqd/qh/admin/authority-issuring/list
func (h *QHAuthorityIssuringGrpcHandler) ListAuthorityIssuring(ctx context.Context, req *tqdpb.ListAuthorityIssuringRequest) (*tqdpb.ListAuthorityIssuringResponse, error) {
	page := uint32(0)
	size := uint32(0)
	if req != nil {
		page = req.Page
		size = req.Size
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.List(ctx, pagable)
	if err != nil {
		return nil, qhAuthorityIssuringError(err)
	}
	data := make([]*tqdpb.QHAuthorityIssuringResponse, len(rows))
	for i := range rows {
		data[i] = toQHAuthorityIssuringPB(&rows[i])
	}
	return &tqdpb.ListAuthorityIssuringResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func qhAuthorityIssuringValidation(code, message, field string) error {
	violations := make([]fault.FieldViolation, 0, 1)
	if field != "" {
		violations = append(violations, fault.FieldViolation{Field: field, Description: message})
	}
	return fault.ToGRPC(fault.Validation(code, message, violations...))
}

func qhAuthorityIssuringError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := fault.As(err); ok {
		return fault.ToGRPC(err)
	}
	return fault.ToGRPC(fault.Wrap(
		err,
		fault.KindInternal,
		"tqd.qh_authority_issuring.internal",
		"authority issuring operation failed",
	))
}

func toQHAuthorityIssuringPB(e *qh_domain.QHAuthorityIssuring) *tqdpb.QHAuthorityIssuringResponse {
	if e == nil {
		return nil
	}
	return &tqdpb.QHAuthorityIssuringResponse{
		Id:          e.ID,
		Name:        e.Name,
		Code:        e.Code,
		Description: e.Description,
		CreatedAt:   _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(e.UpdatedAt),
	}
}
