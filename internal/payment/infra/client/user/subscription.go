package usergrpc

import (
	qhprorpc "common/rpc"
	"common/rpcenv"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	domain "payment/internal/domain/payment"
	userpb "pb/types/user"
	"time"
)

type Client struct {
	client  userpb.InternalSubscriptionServiceClient
	timeout time.Duration
}

func New(target string, timeout time.Duration) (*Client, func(), error) {
	c, cleanup, err := qhprorpc.NewBoundClient(qhprorpc.ClientConfig{Target: target, BackoffMaxDelay: 5 * time.Second, Credentials: insecure.NewCredentials(), Transport: rpcenv.LoadTransportConfig()}, func(conn grpc.ClientConnInterface) userpb.InternalSubscriptionServiceClient {
		return userpb.NewInternalSubscriptionServiceClient(conn)
	})
	if err != nil {
		return nil, nil, err
	}
	return &Client{client: c, timeout: timeout}, cleanup, nil
}
func (c *Client) ApplySettlement(ctx context.Context, e domain.SettlementEffect) error {
	cc, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.client.ApplySettlement(cc, &userpb.ApplySettlementRequest{EffectKey: e.EffectKey, OrderId: e.OrderID, SubjectKind: string(e.Subject.Kind), SubjectId: e.Subject.ID, ProductCode: e.ProductCode, PlanCode: e.PlanCode, PlanVersionId: e.PlanVersionID, PlanVersion: e.PlanVersion, TierRank: e.TierRank, SubscriptionTermDays: e.SubscriptionTermDays, TermsChecksum: e.TermsChecksum, OccurredAt: timestamppb.New(e.OccurredAt)})
	if err == nil {
		return nil
	}
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
		return fmt.Errorf("%w: %v", domain.ErrSubscriptionTransient, err)
	case codes.AlreadyExists, codes.Aborted:
		return fmt.Errorf("%w: %v", domain.ErrSubscriptionEffectConflict, err)
	case codes.InvalidArgument, codes.FailedPrecondition, codes.PermissionDenied:
		return fmt.Errorf("%w: %v", domain.ErrSubscriptionPermanent, err)
	default:
		if errors.Is(err, context.Canceled) {
			return err
		}
		return fmt.Errorf("%w: %v", domain.ErrSubscriptionTransient, err)
	}
}
