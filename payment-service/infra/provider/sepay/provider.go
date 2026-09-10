// Package sepay adapts the existing SePay-style inbound bank-transfer flow to
// the canonical PaymentAttempt capability. It does not classify money truth;
// webhook evidence still enters SettlementService as ProviderEvent.
package sepay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/subiz/vietqr"
	domain "payment/internal/domain/payment"
)

const (
	BankTransferFlowHosted = "hosted_checkout"
	BankTransferFlowVietQR = "vietqr"
)

type Provider struct {
	code     string
	checkout CheckoutConfig
}

type CheckoutConfig struct {
	MerchantID           string
	SecretKey            string
	CheckoutURL          string
	BridgeURL            string
	SuccessURL           string
	ErrorURL             string
	CancelURL            string
	BankTransferDisabled bool
	CardDisabled         bool
	BankTransferFlow     string
	BankBIN              string
	BankAccountNumber    string
}

func New(code string, checkout ...CheckoutConfig) *Provider {
	if strings.TrimSpace(code) == "" {
		code = "sepay"
	}
	var cfg CheckoutConfig
	if len(checkout) > 0 {
		cfg = checkout[0]
	}
	return &Provider{code: strings.TrimSpace(code), checkout: cfg.normalized()}
}

func (p *Provider) Code() string { return p.code }

func (p *Provider) Supports(method domain.PaymentMethod) bool {
	if method == domain.PaymentMethodBankTransfer {
		return !p.checkout.BankTransferDisabled && p.checkout.bankTransferEnabled()
	}
	if method == domain.PaymentMethodCard {
		return p.checkout.enabled() && !p.checkout.CardDisabled
	}
	if method == domain.PaymentMethodQR && !p.checkout.BankTransferDisabled {
		if p.checkout.BankTransferFlow == BankTransferFlowVietQR {
			return p.checkout.bankTransferEnabled()
		}
		return p.checkout.enabled()
	}
	return false
}

func (p *Provider) CreateAttempt(ctx context.Context, order domain.Order, method domain.PaymentMethod, commandKey string) (domain.ProviderAttempt, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProviderAttempt{}, err
	}
	if !p.Supports(method) || order.ID == 0 || strings.TrimSpace(order.Reference) == "" || strings.TrimSpace(commandKey) == "" {
		return domain.ProviderAttempt{}, domain.ErrUnsupportedPaymentMethod
	}
	expires := order.ExpiresAt.UTC()
	if method != domain.PaymentMethodCard && p.checkout.BankTransferFlow == BankTransferFlowVietQR {
		payload, err := p.checkout.vietQRDataURI(order)
		if err != nil {
			return domain.ProviderAttempt{}, domain.ErrProviderUnavailable
		}
		return domain.ProviderAttempt{
			ProviderReference: order.Reference,
			Status:            domain.AttemptPendingAction,
			NextAction: domain.NextAction{
				Kind:      domain.NextActionDisplayQR,
				QRPayload: payload,
				ExpiresAt: &expires,
			},
			ExpiresAt: &expires,
		}, nil
	}
	if p.checkout.enabled() {
		redirectURL, err := p.checkout.redirectURL(order, method)
		if err != nil {
			return domain.ProviderAttempt{}, domain.ErrProviderUnavailable
		}
		return domain.ProviderAttempt{
			ProviderReference: order.Reference,
			Status:            domain.AttemptPendingAction,
			NextAction: domain.NextAction{
				Kind:        domain.NextActionRedirect,
				RedirectURL: redirectURL,
				ExpiresAt:   &expires,
			},
			ExpiresAt: &expires,
		}, nil
	}
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

func (c CheckoutConfig) normalized() CheckoutConfig {
	c.MerchantID = strings.TrimSpace(c.MerchantID)
	c.SecretKey = strings.TrimSpace(c.SecretKey)
	c.CheckoutURL = strings.TrimSpace(c.CheckoutURL)
	c.BridgeURL = strings.TrimSpace(c.BridgeURL)
	c.SuccessURL = strings.TrimSpace(c.SuccessURL)
	c.ErrorURL = strings.TrimSpace(c.ErrorURL)
	c.CancelURL = strings.TrimSpace(c.CancelURL)
	c.BankTransferFlow = strings.ToLower(strings.TrimSpace(c.BankTransferFlow))
	if c.BankTransferFlow == "" {
		c.BankTransferFlow = BankTransferFlowHosted
	}
	c.BankBIN = strings.TrimSpace(c.BankBIN)
	c.BankAccountNumber = strings.TrimSpace(c.BankAccountNumber)
	return c
}

func (c CheckoutConfig) enabled() bool {
	return c.MerchantID != "" && c.SecretKey != "" && c.CheckoutURL != "" && c.BridgeURL != ""
}

func (c CheckoutConfig) bankTransferEnabled() bool {
	if c.BankTransferFlow == BankTransferFlowVietQR {
		return validDigits(c.BankBIN, 6, 6) && validDigits(c.BankAccountNumber, 5, 19)
	}
	return true
}

// vietQRDataURI renders the payment instruction inside Payment Service. Bank
// identity, account number, amount and order reference never leave the process
// to reach a public QR-image API. Provider webhook evidence remains the only
// authority that can settle the order.
func (c CheckoutConfig) vietQRDataURI(order domain.Order) (string, error) {
	if order.Terms.Price.Currency != domain.CurrencyVND || order.Terms.Price.AmountMinor <= 0 ||
		order.Terms.Price.AmountMinor > 9_999_999_999_999 {
		return "", fmt.Errorf("unsupported VietQR amount")
	}
	if !validDigits(c.BankBIN, 6, 6) || !validDigits(c.BankAccountNumber, 5, 19) {
		return "", fmt.Errorf("invalid VietQR beneficiary")
	}
	note := strings.TrimSpace(order.Reference)
	if len(note) == 0 || len(note) > 25 || vietqr.Ascii(note) != note {
		return "", fmt.Errorf("invalid VietQR payment reference")
	}
	png, err := vietqr.GeneratePNG(float64(order.Terms.Price.AmountMinor), c.BankBIN, c.BankAccountNumber, note)
	if err != nil {
		return "", fmt.Errorf("render VietQR: %w", err)
	}
	if len(png) == 0 {
		return "", fmt.Errorf("render VietQR: empty PNG")
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

func validDigits(value string, minimum, maximum int) bool {
	if len(value) < minimum || len(value) > maximum {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// signedFieldOrder is the canonical field order published by SePay. It must
// remain stable: SePay signs the comma-separated field sequence, not a sorted
// map. Optional empty fields are omitted exactly as the official SDK does.
var signedFieldOrder = []string{
	"merchant", "operation", "payment_method", "order_invoice_number",
	"order_amount", "currency", "order_description", "customer_id",
	"success_url", "error_url", "cancel_url", "custom_data",
}

func (c CheckoutConfig) redirectURL(order domain.Order, method domain.PaymentMethod) (string, error) {
	checkoutEndpoint, err := url.Parse(c.CheckoutURL)
	if err != nil || checkoutEndpoint.Scheme != "https" || checkoutEndpoint.Host == "" {
		return "", domain.ErrProviderUnavailable
	}
	bridgeEndpoint, err := url.Parse(c.BridgeURL)
	if err != nil || (bridgeEndpoint.Scheme != "http" && bridgeEndpoint.Scheme != "https") || bridgeEndpoint.Host == "" {
		return "", domain.ErrProviderUnavailable
	}
	values := map[string]string{
		"order_amount":         strconv.FormatInt(order.Terms.Price.AmountMinor, 10),
		"merchant":             c.MerchantID,
		"currency":             string(order.Terms.Price.Currency),
		"operation":            "PURCHASE",
		"order_description":    "Thanh toan don hang " + order.Reference,
		"order_invoice_number": order.Reference,
		"customer_id":          string(order.Subject.Kind) + ":" + order.Subject.ID,
		"payment_method":       sepayPaymentMethod(method),
		"success_url":          c.SuccessURL,
		"error_url":            c.ErrorURL,
		"cancel_url":           c.CancelURL,
	}
	parts := make([]string, 0, len(signedFieldOrder))
	form := make(url.Values)
	for _, field := range signedFieldOrder {
		value := strings.TrimSpace(values[field])
		if value == "" {
			continue
		}
		parts = append(parts, field+"="+value)
		form.Set(field, value)
	}
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, _ = mac.Write([]byte(strings.Join(parts, ",")))
	form.Set("signature", base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	// SePay checkout only accepts an HTML POST. Redirecting the mobile/browser
	// client directly with query parameters is not protocol-equivalent. The
	// Gateway bridge contains no merchant secret: it only renders these signed
	// fields into an auto-submitting POST form.
	query := bridgeEndpoint.Query()
	query.Set("checkout_url", checkoutEndpoint.String())
	for _, field := range signedFieldOrder {
		if value := form.Get(field); value != "" {
			query.Set(field, value)
		}
	}
	query.Set("signature", form.Get("signature"))
	bridgeEndpoint.RawQuery = query.Encode()
	return bridgeEndpoint.String(), nil
}

func sepayPaymentMethod(method domain.PaymentMethod) string {
	switch method {
	case domain.PaymentMethodCard:
		return "CARD"
	case domain.PaymentMethodQR:
		return "NAPAS_BANK_TRANSFER"
	default:
		return "BANK_TRANSFER"
	}
}
