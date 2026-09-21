package handler_grpc

import (
	"context"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"

	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	"tqd/internal"
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
		return nil, _errors.ReturnError(service.AuthorityRequestRequired)
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
		return nil, _errors.ReturnError(
			service.AuthorityIDRequired,
			_errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "id is required"}),
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
		return nil, _errors.ReturnError(
			service.AuthorityIDRequired,
			_errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "id is required"}),
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
		return nil, _errors.ReturnError(service.AuthorityUpdateFieldsRequired)
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
		return nil, _errors.ReturnError(
			service.AuthorityIDRequired,
			_errors.WithViolations(_errors.FieldViolation{Field: "id", Description: "id is required"}),
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

func qhAuthorityIssuringError(err error) error {
	return _errors.ToGRPC(err)
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
