package useraccess

import (
	"common/access"
	commonmetering "common/metering"
	"common/operation"
	"context"
	"errors"
	"strings"
	"time"

	userpb "pb/types/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAccessUnavailable = errors.New("entitlement access is unavailable")

/**
 * Client lấy access từ user-service qua gRPC nội bộ.
 *
 * Không truyền profile ID. user-service lấy Actor từ trusted context.
 */
type Client struct {
	grpc userpb.InternalAccessServiceClient
}

func NewClient(grpcClient userpb.InternalAccessServiceClient) *Client {
	return &Client{grpc: grpcClient}
}

/**
 * GetAccess resolves one operation against trusted caller context.
 */
func (c *Client) GetAccess(
	ctx context.Context,
	operationCode operation.Code,
) (access.Result, error) {
	if c == nil || c.grpc == nil {
		return access.Result{}, ErrAccessUnavailable
	}
	if !operationCode.IsValid() {
		return access.Result{}, ErrAccessUnavailable
	}

	response, err := c.grpc.GetAccess(ctx, &userpb.GetAccessRequest{
		Operation: string(operationCode),
	})
	if err != nil {
		if status.Code(err) == codes.PermissionDenied || status.Code(err) == codes.Unauthenticated {
			return access.Result{}, err
		}
		return access.Result{}, errors.Join(ErrAccessUnavailable, err)
	}

	return fromProto(response)
}

func fromProto(response *userpb.GetAccessResponse) (access.Result, error) {
	if response == nil {
		return access.Result{}, ErrAccessUnavailable
	}

	operationCode, err := operation.Parse(strings.TrimSpace(response.GetOperation()))
	if err != nil {
		return access.Result{}, errors.Join(ErrAccessUnavailable, err)
	}

	result := access.Result{
		Subject: access.Subject{
			Type: access.SubjectType(response.GetSubjectType()),
			ID:   response.GetSubjectId(),
		},
		Operation: operationCode,
		Metering: access.Metering{
			FeatureCode:    response.GetFeatureCode(),
			MeterCode:      commonmetering.Code(response.GetMeterCode()),
			UnitsPerAction: response.GetUnitsPerAction(),
			PolicyVersion:  response.GetPolicyVersion(),
		},
		Allowed:        response.GetAllowed(),
		Unlimited:      response.GetUnlimited(),
		Limit:          response.GetLimit(),
		Period:         access.Period(response.GetPeriod()),
		SubscriptionID: response.GetSubscriptionId(),
		PlanCode:       response.GetPlanCode(),
		PlanVersion:    response.GetPlanVersion(),
	}

	if response.GetPeriodStartUnix() > 0 {
		result.PeriodStart = time.Unix(response.GetPeriodStartUnix(), 0).UTC()
	}
	if response.GetPeriodEndUnix() > 0 {
		result.PeriodEnd = time.Unix(response.GetPeriodEndUnix(), 0).UTC()
	}

	if !result.IsValid() {
		return access.Result{}, ErrAccessUnavailable
	}

	return result, nil
}
