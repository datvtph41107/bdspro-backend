package payment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"common/configloader"
	qhprorpc "common/rpc"
	"common/rpcenv"
	paymentpb "pb/types/payment"
	"user/internal/usecase/subscription/checkout"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Commands struct {
	client  paymentpb.InternalCommerceServiceClient
	timeout time.Duration
}

func NewCommands() (*Commands, func(), error) {
	target, err := configloader.RequiredString("rpc.payment.address")
	if err != nil {
		return nil, nil, err
	}
	timeout := 5 * time.Second
	client, cleanup, err := qhprorpc.NewBoundClient(qhprorpc.ClientConfig{
		Target: target, BackoffMaxDelay: 5 * time.Second,
		Credentials: insecure.NewCredentials(),
		Transport:   rpcenv.LoadTransportConfig(),
	}, func(conn grpc.ClientConnInterface) paymentpb.InternalCommerceServiceClient {
		return paymentpb.NewInternalCommerceServiceClient(conn)
	})
	if err != nil {
		return nil, nil, err
	}
	return &Commands{client: client, timeout: timeout}, cleanup, nil
}

func (c *Commands) CreateOrder(ctx context.Context, command checkout.CreatePaymentOrderCommand) (checkout.PaymentOrder, bool, error) {
	if c == nil || c.client == nil || !command.Subject.IsValid() || !command.Terms.IsValid() || strings.TrimSpace(command.CommandKey) == "" {
		return checkout.PaymentOrder{}, false, checkout.ErrInvalidCommand
	}
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	r, err := c.client.CreateCommercialOrder(callCtx, &paymentpb.CreateCommercialOrderRequest{
		SubjectKind: string(command.Subject.Kind), SubjectId: command.Subject.ID,
		ProductCode: command.Terms.ProductCode, PlanCode: command.Terms.PlanCode,
		PlanVersionId: command.Terms.PlanVersionID, PlanVersion: command.Terms.PlanVersion,
		TierRank: command.Terms.TierRank, SubscriptionTermDays: command.Terms.SubscriptionTermDays,
		TermsChecksum: command.Terms.TermsChecksum, Currency: command.Terms.Price.Currency,
		AmountMinor: command.Terms.Price.AmountMinor, CommandKey: command.CommandKey,
	})
	if err != nil {
		return checkout.PaymentOrder{}, false, fmt.Errorf("create payment order: %w", err)
	}
	return checkout.PaymentOrder{ID: r.GetOrderId(), Reference: r.GetReference(), Status: r.GetStatus()}, r.GetCreated(), nil
}

func (c *Commands) CreateAttempt(ctx context.Context, command checkout.CreatePaymentAttemptCommand) (checkout.PaymentAttempt, error) {
	if c == nil || c.client == nil || !command.Subject.IsValid() || command.OrderID == 0 || strings.TrimSpace(command.Method) == "" || strings.TrimSpace(command.CommandKey) == "" {
		return checkout.PaymentAttempt{}, checkout.ErrInvalidCommand
	}
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	r, err := c.client.CreatePaymentAttempt(callCtx, &paymentpb.CreatePaymentAttemptRequest{
		OrderId: command.OrderID, Method: command.Method, CommandKey: command.CommandKey,
		SubjectKind: string(command.Subject.Kind), SubjectId: command.Subject.ID,
	})
	if err != nil {
		return checkout.PaymentAttempt{}, fmt.Errorf("create payment attempt: %w", err)
	}
	out := checkout.PaymentAttempt{ID: r.GetAttemptId(), OrderID: r.GetOrderId(), Method: r.GetMethod(), Provider: r.GetProvider(), ProviderReference: r.GetProviderReference(), Status: r.GetStatus(), Created: r.GetCreated()}
	if n := r.GetNextAction(); n != nil {
		out.NextActionKind = n.GetKind()
		out.RedirectURL = n.GetRedirectUrl()
		out.QRPayload = n.GetQrPayload()
	}
	if r.GetExpiresAt() != "" {
		if t, e := time.Parse(time.RFC3339Nano, r.GetExpiresAt()); e == nil {
			out.ExpiresAt = t
		}
	}
	return out, nil
}

var _ checkout.PaymentOrderPort = (*Commands)(nil)
var _ checkout.PaymentAttemptPort = (*Commands)(nil)
