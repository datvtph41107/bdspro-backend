package handler

import (
	"context"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/dto"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MembershipHandler struct {
	userpb.UnimplementedMembershipServiceServer
	membershipUsecase usecases.IMembershipUsecase
	mapper            *mapper.MembershipMapper
}

func NewMembershipHandler(
	membershipUsecase usecases.IMembershipUsecase,
	mapper *mapper.MembershipMapper,
) *MembershipHandler {
	return &MembershipHandler{
		membershipUsecase: membershipUsecase,
		mapper:            mapper,
	}
}

// CreateOrder creates a new membership order
// @Summary Tạo đơn hàng mua gói membership
// @Description Tạo đơn hàng để mua gói membership và nhận link thanh toán
// @Tags Membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreateOrderRequest true "Thông tin đơn hàng"
// @Success 200 {object} pb.CreateOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/membership/order [post]
func (h *MembershipHandler) CreateOrder(ctx context.Context, req *userpb.CreateOrderRequest) (*userpb.CreateOrderResponse, error) {
	membership, paymentURL, err := h.membershipUsecase.CreateOrder(ctx, req.PlanId)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &userpb.CreateOrderResponse{
		Membership: h.mapper.ToProto(membership),
		PaymentUrl: paymentURL,
	}, nil
}

// VerifyOrder verifies the payment and activates membership
// @Summary Xác nhận thanh toán từ payment service
// @Description Webhook từ payment service để xác nhận thanh toán và kích hoạt membership
// @Tags Membership
// @Accept json
// @Produce json
// @Param request body pb.VerifyOrderRequest true "Thông tin xác nhận thanh toán"
// @Success 200 {object} pb.VerifyOrderResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user/membership/verify [post]
func (h *MembershipHandler) VerifyOrder(ctx context.Context, req *userpb.VerifyOrderRequest) (*userpb.VerifyOrderResponse, error) {
	verifyDTO := dto.VerifyOrderDTO{
		TargetID:  req.TargetId,
		Type:      dto.OrderType(req.Type),
		Signature: req.Signature,
	}

	if req.PaidAt != nil {
		paidAt := req.PaidAt.AsTime()
		verifyDTO.PaidAt = &paidAt
	}

	membership, err := h.membershipUsecase.VerifyOrder(ctx, verifyDTO)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &userpb.VerifyOrderResponse{
		Membership: h.mapper.ToProto(membership),
		Success:    true,
	}, nil
}

// GetMyMembership gets the current user's membership
// @Summary Lấy thông tin membership của user hiện tại
// @Description Lấy thông tin gói membership đang sử dụng của user
// @Tags Membership
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pb.GetMyMembershipResponse
// @Failure 404 {object} map[string]string
// @Router /user/membership/me [get]
func (h *MembershipHandler) GetMyMembership(ctx context.Context, req *userpb.GetMyMembershipRequest) (*userpb.GetMyMembershipResponse, error) {
	membership, err := h.membershipUsecase.GetMyMembership(ctx)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return &userpb.GetMyMembershipResponse{
		Membership: h.mapper.ToProto(membership),
	}, nil
}
