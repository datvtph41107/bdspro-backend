package handler_grpc

import (
	"context"
	"strings"

	"common/identity"
	tqdpb "pb/types/tqd"
	quotausecase "tqd/internal/usecase/quota"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileUsageGrpcHandler struct {
	tqdpb.UnimplementedProfileUsageServiceServer
	service *quotausecase.ProfileUsageService
}

func NewProfileUsageGrpcHandler(service *quotausecase.ProfileUsageService) *ProfileUsageGrpcHandler {
	return &ProfileUsageGrpcHandler{service: service}
}

func (h *ProfileUsageGrpcHandler) GetMyGeneratedReportUsage(ctx context.Context, _ *tqdpb.GetMyGeneratedReportUsageRequest) (*tqdpb.ProfileUsageProjection, error) {
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || !strings.EqualFold(actor.TokenType, "ACCESS") {
		return nil, status.Error(codes.Unauthenticated, "authenticated profile is required")
	}
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Unavailable, "profile usage projection is unavailable")
	}
	projection, err := h.service.GetGeneratedReportUsage(ctx)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "profile usage projection is unavailable")
	}
	decision := projection.Access
	response := &tqdpb.ProfileUsageProjection{
		Operation: string(decision.Operation), Allowed: decision.Allowed, Unlimited: decision.Unlimited,
		FeatureCode: decision.Metering.FeatureCode, MeterCode: string(decision.Metering.MeterCode),
		LimitKnown: decision.UsesQuota(), Limit: decision.Limit, DurableUsed: projection.DurableUsed,
		RuntimeUsed: projection.RuntimeUsed, RuntimeReserved: projection.RuntimeReserved, Remaining: projection.Remaining,
		RuntimeProjectionAvailable: projection.RuntimeAvailable, ReconciliationHealth: projection.Reconciliation,
		SubscriptionId: decision.SubscriptionID, PlanCode: decision.PlanCode, PlanVersion: decision.PlanVersion,
		PolicyVersion: decision.Metering.PolicyVersion,
	}
	if !decision.PeriodStart.IsZero() {
		response.PeriodStart = timestamppb.New(decision.PeriodStart)
	}
	if !decision.PeriodEnd.IsZero() {
		response.PeriodEnd = timestamppb.New(decision.PeriodEnd)
	}
	return response, nil
}
