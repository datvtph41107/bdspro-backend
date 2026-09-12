package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSepayGatewayAllowsUnconfiguredProvider(t *testing.T) {
	t.Parallel()
	require.NoError(t, validateSepayGateway(SepayGatewayConfig{Environment: "sandbox"}))
}

func TestValidateSepayGatewayRequiresCompleteCredentials(t *testing.T) {
	t.Parallel()
	err := validateSepayGateway(SepayGatewayConfig{Environment: "sandbox", MerchantID: "merchant"})
	require.ErrorContains(t, err, "must be configured together")
}

func TestValidateSepayGatewayRejectsUnknownEnvironment(t *testing.T) {
	t.Parallel()
	err := validateSepayGateway(SepayGatewayConfig{Environment: "staging"})
	require.ErrorContains(t, err, "sandbox or production")
}

func TestValidateSepayGatewayRequiresHTTPSBridgeInProduction(t *testing.T) {
	t.Parallel()
	err := validateSepayGateway(SepayGatewayConfig{
		Environment: "production", MerchantID: "merchant", SecretKey: "secret",
		IPNSecret: "ipn-secret", BridgeURL: "http://localhost:8000/v2/payment/sepay/checkout",
	})
	require.ErrorContains(t, err, "must use HTTPS")
}

func TestValidateSepayGatewayAcceptsCompleteVietQRBeneficiary(t *testing.T) {
	t.Parallel()
	require.NoError(t, validateSepayGateway(SepayGatewayConfig{
		Environment: "production", BankTransferFlow: SepayBankTransferVietQR,
		BankTransferEnabled: true, BankBIN: "970422", BankAccountNumber: "0123456789",
	}))
}

func TestValidateSepayGatewayRejectsIncompleteVietQRBeneficiary(t *testing.T) {
	t.Parallel()
	err := validateSepayGateway(SepayGatewayConfig{
		Environment: "production", BankTransferFlow: SepayBankTransferVietQR,
		BankTransferEnabled: true, BankBIN: "MB", BankAccountNumber: "123",
	})
	require.ErrorContains(t, err, "SEPAY_QR_BANK_BIN")
}

func TestValidateProviderRejectsRealBankQRInSandbox(t *testing.T) {
	t.Parallel()
	err := validateProvider(ProviderConfig{
		SepayAPIKey: "webhook-key",
		Sepay: SepayGatewayConfig{
			Environment: "sandbox", BankTransferFlow: SepayBankTransferVietQR,
			BankTransferEnabled: true, BankBIN: "970422", BankAccountNumber: "0123456789",
		},
	})
	require.ErrorContains(t, err, "SEPAY_WEBHOOK_ENVIRONMENT=test")
}

func TestValidateProviderAllowsExplicitBankWebhookTestMode(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		gateway, webhook string
		valid            bool
	}{
		{"sandbox", "test", true},
		{"production", "live", true},
		{"production", "test", false},
		{"sandbox", "live", false},
		{"sandbox", "unknown", false},
	} {
		t.Run(tc.gateway+"/"+tc.webhook, func(t *testing.T) {
			cfg := ProviderConfig{SepayAPIKey: "webhook-key", Sepay: SepayGatewayConfig{
				Environment: tc.gateway, WebhookEnvironment: tc.webhook,
				BankTransferFlow: SepayBankTransferVietQR, BankTransferEnabled: true,
				BankBIN: "970422", BankAccountNumber: "0123456789",
			}}
			if tc.valid {
				require.NoError(t, validateProvider(cfg))
				cfg.SepayAPIKey = ""
				require.ErrorContains(t, validateProvider(cfg), "SEPAY_API_KEY")
			} else {
				require.Error(t, validateProvider(cfg))
			}
		})
	}
}

func TestValidateProviderRequiresWebhookAuthenticationForVietQR(t *testing.T) {
	t.Parallel()
	err := validateProvider(ProviderConfig{Sepay: SepayGatewayConfig{
		Environment: "production", BankTransferFlow: SepayBankTransferVietQR,
		BankTransferEnabled: true, BankBIN: "970422", BankAccountNumber: "0123456789",
	}})
	require.ErrorContains(t, err, "SEPAY_API_KEY")
}

func TestResolveConfigFileUsesCanonicalRuntimeProfile(t *testing.T) {
	t.Parallel()

	require.Equal(t, "config/runtime.yml", resolveConfigFile(""))
}

func TestResolveConfigFileAllowsExplicitOperatorOverride(t *testing.T) {
	t.Parallel()

	require.Equal(t, "/run/secrets/payment.yml", resolveConfigFile(" /run/secrets/payment.yml "))
}
