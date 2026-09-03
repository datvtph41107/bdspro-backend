// Package sepay adapts the existing SePay-style inbound bank-transfer flow to
// the canonical PaymentAttempt capability. It does not classify money truth;
// webhook evidence still enters SettlementService as ProviderEvent.
package sepay

import (
	"context"
	"strings"

	domain "payment/internal/domain/payment"
)

type Provider struct {
	code string
}

func New(code string) *Provider {
	if strings.TrimSpace(code) == "" {
		code = "sepay"
	}
	return &Provider{code: strings.TrimSpace(code)}
}

func (p *Provider) Code() string { return p.code }

func (p *Provider) Supports(method domain.PaymentMethod) bool {
	return method == domain.PaymentMethodBankTransfer
}

func (p *Provider) CreateAttempt(ctx context.Context, order domain.Order, method domain.PaymentMethod, commandKey string) (domain.ProviderAttempt, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProviderAttempt{}, err
	}
	if !p.Supports(method) || order.ID == 0 || strings.TrimSpace(order.Reference) == "" || strings.TrimSpace(commandKey) == "" {
		return domain.ProviderAttempt{}, domain.ErrUnsupportedPaymentMethod
	}
	expires := order.ExpiresAt.UTC()
	return domain.ProviderAttempt{
		ProviderReference: order.Reference,
		Status:            domain.AttemptPendingAction,
		NextAction: domain.NextAction{
			Kind:      domain.NextActionWait,
			ExpiresAt: &expires,
		},
		ExpiresAt: &expires,
	}, nil
}
