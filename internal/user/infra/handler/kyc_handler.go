package handler

import (
	_dto "common/domain/dto"
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/dto"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KYCHandler struct {
	userpb.UnimplementedKYCServiceServer
	kycUsecase      usecases.IKYCUsecase
	adminKYCUsecase usecases.IAdminKYCUsecase
	mapper          *mapper.KYCMapper
}

// @bind: internal/usecases.IKYCUsecase
// @bind: internal/usecases.IAdminKYCUsecase
func NewKYCHandler(
	kycUsecase usecases.IKYCUsecase,
	adminKYCUsecase usecases.IAdminKYCUsecase,
	mapper *mapper.KYCMapper,
) *KYCHandler {
	return &KYCHandler{
		kycUsecase:      kycUsecase,
		adminKYCUsecase: adminKYCUsecase,
		mapper:          mapper,
	}
}

// SubmitKYC submit KYC request
// @Summary Submit KYC request
// @Description User submit KYC verification request
// @Tags KYC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.SubmitKYCRequest true "Thông tin KYC"
// @Success 200 {object} sharepb.KYCV3Proto
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/kyc/submit [post]
func (h *KYCHandler) SubmitKYC(ctx context.Context, req *userpb.SubmitKYCRequest) (*sharepb.KYCV3Proto, error) {
	submitDTO := dto.SubmitKYCRequest{
		FullName:     req.FullName,
		IdentityCard: req.IdentityCard,
		FrontImage:   req.FrontImage,
		BackImage:    req.BackImage,
		SelfieImage:  req.SelfieImage,
	}

	// Thêm 4 trường mới nếu có
	if req.IdNumber != nil {
		submitDTO.IDNumber = *req.IdNumber
	}
	if req.DateOfBirth != nil {
		submitDTO.DateOfBirth = *req.DateOfBirth
	}
	if req.ExpiryDate != nil {
		submitDTO.ExpiryDate = *req.ExpiryDate
	}

	kyc, err := h.kycUsecase.SubmitKYC(ctx, submitDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	// SubmitKYC luôn trả về cho owner (chính user submit)
	return h.mapper.ToProtoWithOwnerCheck(kyc, true), nil
}

// ListKYCs lấy danh sách KYC requests (Admin)
// @Summary Lấy danh sách KYC requests
// @Description Lấy danh sách tất cả KYC requests với phân trang và tìm kiếm (Admin only)
// @Tags Admin KYC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại (mặc định: 1)"
// @Param size query int false "Số lượng item trên mỗi trang (mặc định: 20)"
// @Param search query string false "Tìm kiếm theo tên, số CMND"
// @Param status query int false "Lọc theo trạng thái: 10: Pending, 20: Approved, 30: Rejected"
// @Param sortBy query string false "Sắp xếp theo trường"
// @Param sortOrder query string false "Thứ tự sắp xếp (asc, desc)"
// @Success 200 {object} userpb.ListKYCsResponse
// @Failure 500 {object} map[string]string
// @Router /v2/user/admin/kycs [get]
func (h *KYCHandler) ListKYCs(ctx context.Context, req *userpb.ListKYCsRequest) (*userpb.ListKYCsResponse, error) {
	listDTO := dto.ListKYCsRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Search:    req.Search,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	entities, total, err := h.adminKYCUsecase.ListKYCs(ctx, listDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &userpb.ListKYCsResponse{
		Data:  h.mapper.ToProtoList(entities),
		Total: total,
	}, nil
}

// ApproveKYC duyệt KYC request (Admin)
// @Summary Duyệt KYC request
// @Description Admin duyệt yêu cầu KYC của user
// @Tags Admin KYC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "KYC ID"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/user/admin/kyc/approve/{id} [post]
func (h *KYCHandler) ApproveKYC(ctx context.Context, req *userpb.ApproveKYCRequest) (*sharepb.Empty, error) {
	err := h.adminKYCUsecase.ApproveKYC(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &sharepb.Empty{}, nil
}

// RejectKYC từ chối KYC request (Admin)
// @Summary Từ chối KYC request
// @Description Admin từ chối yêu cầu KYC của user
// @Tags Admin KYC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "KYC ID"
// @Param request body userpb.RejectKYCRequest true "Lý do từ chối"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/user/admin/kyc/reject/{id} [post]
func (h *KYCHandler) RejectKYC(ctx context.Context, req *userpb.RejectKYCRequest) (*sharepb.Empty, error) {
	err := h.adminKYCUsecase.RejectKYC(ctx, req.Id, req.Reason)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &sharepb.Empty{}, nil
}
