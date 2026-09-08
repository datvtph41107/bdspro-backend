package grpc

import (
	"common/identity"
	"context"
	userpb "pb/types/user"

	"user/internal/models"
	"user/internal/usecase/subscription/settlement"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionSettlementHandler struct {
	userpb.UnimplementedInternalSubscriptionServiceServer
	service *settlement.Service
}

func NewSubscriptionSettlementHandler(s *settlement.Service) *SubscriptionSettlementHandler {
	return &SubscriptionSettlementHandler{service: s}
}
func (h *SubscriptionSettlementHandler) ApplySettlement(ctx context.Context, req *userpb.ApplySettlementRequest) (*userpb.ApplySettlementResponse, error) {
	caller, ok := identity.ServiceCallerFromContext(ctx)
	if !ok || caller.ServiceID != "payment-service" {
		return nil, status.Error(codes.PermissionDenied, "verified payment-service caller required")
	}
	if h == nil || h.service == nil || req == nil || req.GetOccurredAt() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid settlement effect")
	}
	e := settlement.Effect{EffectKey: req.GetEffectKey(), OrderID: req.GetOrderId(), SubjectKind: models.SubscriptionSubjectKind(req.GetSubjectKind()), SubjectID: req.GetSubjectId(), ProductCode: req.GetProductCode(), PlanCode: req.GetPlanCode(), PlanVersionID: req.GetPlanVersionId(), PlanVersion: req.GetPlanVersion(), TierRank: req.GetTierRank(), SubscriptionTermDays: req.GetSubscriptionTermDays(), TermsChecksum: req.GetTermsChecksum(), OccurredAt: req.GetOccurredAt().AsTime()}
	r, err := h.service.ApplySettlement(ctx, e)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &userpb.ApplySettlementResponse{SubscriptionId: r.SubscriptionID, Action: string(r.Action), Replayed: r.Replayed}, nil
}
