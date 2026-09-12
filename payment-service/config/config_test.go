package config

import (
	"os"
	"path/filepath"
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

func TestLoadDatabaseDSNFromLocalYMLIgnoresEnv(t *testing.T) {
	serviceRoot, err := filepath.Abs("..")
	require.NoError(t, err)
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(serviceRoot))
	t.Cleanup(func() { _ = os.Chdir(wd) })

	for _, key := range []string{"QHPRO_CONFIG_FILE", "CONFIG_FILE"} {
		t.Setenv(key, "")
	}
	const envOverride = "postgres://env-override-must-not-win@127.0.0.1:65432/env_payment?sslmode=disable"
	t.Setenv("PAYMENT_DATABASE_URL", envOverride)
	t.Setenv("PAYMENT_DATABASE_DSN", envOverride)

	if _, err := os.Stat("config/local.yml"); err != nil {
		t.Skip("config/local.yml missing; workstation overlay required for this evidence")
	}

	cfg, err := Load()
	require.NoError(t, err)
	require.NotContains(t, cfg.Database.DSN, "env-override-must-not-win", "DSN must come from YAML, not PAYMENT_DATABASE_URL")
	require.NotEmpty(t, cfg.Database.DSN)
}

func TestLoadUsesMountedConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payment.yml")
	require.NoError(t, os.WriteFile(path, []byte(`
host: 127.0.0.1
port: "9205"
database:
  dsn: "postgres://payment:secret@db:5432/db_payment?sslmode=disable"
sepay:
  gateway_environment: sandbox
  bank_transfer_flow: hosted_checkout
`), 0o600))

	t.Setenv("CONFIG_FILE", path)
	t.Setenv("QHPRO_CONFIG_FILE", "")
	t.Setenv("PAYMENT_DATABASE_URL", "postgres://ignored@localhost/payment_service")
	t.Setenv("PAYMENT_GRPC_ADDRESS", "")
	t.Setenv("PAYMENT_GRPC_HOST", "")
	t.Setenv("PAYMENT_GRPC_PORT", "")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "postgres://payment:secret@db:5432/db_payment?sslmode=disable", cfg.Database.DSN)
	require.Equal(t, "127.0.0.1:9205", cfg.Server.Address)
}
