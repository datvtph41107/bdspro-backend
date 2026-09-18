package sepay

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/subiz/vietqr"
	domain "payment/internal/domain/payment"
)

func TestGatewayAttemptReturnsSignedHostedCheckout(t *testing.T) {
	t.Parallel()
	p := New("sepay", CheckoutConfig{
		MerchantID: "merchant-1", SecretKey: "secret",
		CheckoutURL: "https://pay-sandbox.sepay.vn/v1/checkout/init",
		BridgeURL:   "http://127.0.0.1:8000/v2/payment/sepay/checkout",
		SuccessURL:  "http://localhost:3001/dashboard/pricing?payment=success",
	})
	order := validOrder()

	attempt, err := p.CreateAttempt(context.Background(), order, domain.PaymentMethodCard, "attempt-1")
	require.NoError(t, err)
	require.Equal(t, domain.NextActionRedirect, attempt.NextAction.Kind)
	u, err := url.Parse(attempt.NextAction.RedirectURL)
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:8000", u.Host)
	require.Equal(t, "https://pay-sandbox.sepay.vn/v1/checkout/init", u.Query().Get("checkout_url"))
	require.Equal(t, "99000", u.Query().Get("order_amount"))
	require.Equal(t, "CARD", u.Query().Get("payment_method"))
	require.Equal(t, order.Reference, u.Query().Get("order_invoice_number"))
	require.NotEmpty(t, u.Query().Get("signature"))
}

func TestUnconfiguredGatewayKeepsLegacyBankWebhookFlow(t *testing.T) {
	t.Parallel()
	attempt, err := New("sepay").CreateAttempt(context.Background(), validOrder(), domain.PaymentMethodBankTransfer, "attempt-1")
	require.NoError(t, err)
	require.Equal(t, domain.NextActionWait, attempt.NextAction.Kind)
	require.False(t, New("sepay").Supports(domain.PaymentMethodCard))
}

func TestVietQRAttemptReturnsCanonicalPayloadForClientRendering(t *testing.T) {
	t.Parallel()
	p := New("sepay", CheckoutConfig{
		// Existing hosted credentials must not override the selected QR flow.
		MerchantID: "merchant", SecretKey: "secret",
		CheckoutURL:       "https://pay-sandbox.sepay.vn/v1/checkout/init",
		BridgeURL:         "http://127.0.0.1:8000/v2/payment/sepay/checkout",
		BankTransferFlow:  BankTransferFlowVietQR,
		BankBIN:           "970422",
		BankAccountNumber: "0123456789",
	})

	attempt, err := p.CreateAttempt(context.Background(), validOrder(), domain.PaymentMethodBankTransfer, "attempt-1")
	require.NoError(t, err)
	require.Equal(t, domain.NextActionDisplayQR, attempt.NextAction.Kind)
	require.Empty(t, attempt.NextAction.RedirectURL)
	require.False(t, strings.HasPrefix(attempt.NextAction.QRPayload, "data:image/"))
	require.Equal(t, vietqr.Generate(99000, "970422", "0123456789", validOrder().Reference), attempt.NextAction.QRPayload)
}

func TestVietQRRejectsInvalidBeneficiaryBeforeAttemptCreation(t *testing.T) {
	t.Parallel()
	p := New("sepay", CheckoutConfig{
		BankTransferFlow:  BankTransferFlowVietQR,
		BankBIN:           "MB",
		BankAccountNumber: "123",
	})
	require.False(t, p.Supports(domain.PaymentMethodBankTransfer))
}

func TestMerchantMethodFlagsControlAdvertisedPaymentMethods(t *testing.T) {
	t.Parallel()
	p := New("sepay", CheckoutConfig{
		MerchantID: "merchant-1", SecretKey: "secret",
		CheckoutURL:  "https://pay-sandbox.sepay.vn/v1/checkout/init",
		BridgeURL:    "http://127.0.0.1:8000/v2/payment/sepay/checkout",
		CardDisabled: true,
	})
	require.True(t, p.Supports(domain.PaymentMethodBankTransfer))
	require.False(t, p.Supports(domain.PaymentMethodCard))

	p = New("sepay", CheckoutConfig{
		MerchantID: "merchant-1", SecretKey: "secret",
		CheckoutURL:          "https://pay-sandbox.sepay.vn/v1/checkout/init",
		BridgeURL:            "http://127.0.0.1:8000/v2/payment/sepay/checkout",
		BankTransferDisabled: true,
	})
	require.False(t, p.Supports(domain.PaymentMethodBankTransfer))
	require.True(t, p.Supports(domain.PaymentMethodCard))
}

func validOrder() domain.Order {
	expires := time.Now().UTC().Add(time.Hour)
	return domain.Order{
		ID: 1, Subject: domain.Subject{Kind: domain.SubjectProfile, ID: "29"},
		Reference: "QHP-0123456789ABCDEF", ExpiresAt: expires,
		Terms: domain.CommercialTerms{Price: domain.Money{Currency: domain.CurrencyVND, AmountMinor: 99000}},
	}
}
