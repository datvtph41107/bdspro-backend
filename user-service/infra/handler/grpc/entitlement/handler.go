package accessgrpc

import (
	"common/identity"
	"common/operation"
	"context"
	"errors"
	"strings"
	"time"
	"user/internal/domain/entitlement"
	internalaccess "user/internal/usecase/entitlement"

	userpb "pb/types/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * Handler expose access qua gRPC nội bộ.
 *
 * Profile ID được lấy từ Actor đã có trong context, không nhận từ request.
 * Nhờ vậy service gọi không thể tự khai một profile khác để hỏi quyền.
 */
type Handler struct {
	userpb.UnimplementedInternalAccessServiceServer
	service *internalaccess.Service
	now     func() time.Time
}

func NewHandler(service *internalaccess.Service) *Handler {
	return &Handler{
		service: service,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

/**
 * GetAccess trả quyền dùng operation cho profile đang gọi request.
 */
func (h *Handler) GetAccess(
	ctx context.Context,
	req *userpb.GetAccessRequest,
) (*userpb.GetAccessResponse, error) {
	if h == nil || h.service == nil {
		return nil, status.Error(codes.Internal, "access service is not configured")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 {
		return nil, status.Error(codes.Unauthenticated, "profile identity is required")
	}

	operationCode, err := operation.Parse(strings.TrimSpace(req.GetOperation()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "operation is invalid")
	}

	subject, err := internalaccess.SubjectFromActor(actor)
	if err != nil {
		if errors.Is(err, internalaccess.ErrActorProfileMissing) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, "cannot resolve access subject")
	}

	result, err := h.service.GetAccess(ctx, subject, operationCode, h.now())
	if err != nil {
		switch {
		case errors.Is(err, internalaccess.ErrInvalidSubject),
			errors.Is(err, internalaccess.ErrInvalidOperation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "cannot load access")
		}
	}

	return toProto(result), nil
}

func toProto(result access.Result) *userpb.GetAccessResponse {
	response := &userpb.GetAccessResponse{
		SubjectType:    string(result.Subject.Type),
		SubjectId:      result.Subject.ID,
		Operation:      string(result.Operation),
		FeatureCode:    result.Metering.FeatureCode,
		MeterCode:      string(result.Metering.MeterCode),
		UnitsPerAction: result.Metering.UnitsPerAction,
		PolicyVersion:  result.Metering.PolicyVersion,
		Allowed:        result.Allowed,
		Unlimited:      result.Unlimited,
		Limit:          result.Limit,
		Period:         string(result.Period),
		SubscriptionId: result.SubscriptionID,
		PlanCode:       result.PlanCode,
		PlanVersion:    result.PlanVersion,
	}

	if !result.PeriodStart.IsZero() {
		response.PeriodStartUnix = result.PeriodStart.Unix()
	}
	if !result.PeriodEnd.IsZero() {
		response.PeriodEndUnix = result.PeriodEnd.Unix()
	}

	return response
}
