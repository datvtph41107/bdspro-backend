package config

import (
	"common/configloader"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address         string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

type CommerceConfig struct {
	OrderTTL         time.Duration
	FulfillmentRetry time.Duration
	FulfillmentLease time.Duration
	FulfillmentPoll  time.Duration
}

type RabbitConfig struct {
	URL            string
	Exchange       string
	PublisherLease time.Duration
	PublisherPoll  time.Duration
	PublisherRetry time.Duration
}

type ProviderConfig struct {
	// SepayAPIKey authenticates the legacy bank-transaction webhook. It is
	// deliberately separate from Payment Gateway merchant credentials.
	SepayAPIKey string
	Sepay       SepayGatewayConfig
}

type SepayGatewayConfig struct {
	Environment         string
	BankTransferFlow    string
	WebhookEnvironment  string
	MerchantID          string
	SecretKey           string
	CheckoutURL         string
	BridgeURL           string
	SuccessURL          string
	ErrorURL            string
	CancelURL           string
	IPNSecret           string
	BankTransferEnabled bool
	CardEnabled         bool
	BankBIN             string
	BankAccountNumber   string
}

const (
	SepayBankTransferHosted = "hosted_checkout"
	SepayBankTransferVietQR = "vietqr"
	SepayWebhookTest        = "test"
	SepayWebhookLive        = "live"
)

type Config struct {
	Environment  string
	Server       ServerConfig
	Database     DatabaseConfig
	User         ClientConfig
	Auth         ClientConfig
	Notification ClientConfig
	Commerce     CommerceConfig
	Rabbit       RabbitConfig
	Provider     ProviderConfig
}

// Load reads process config from the environment YAML contract:
// runtime.yml merged with optional config/local.yml (workstation), or
// CONFIG_FILE / QHPRO_CONFIG_FILE when deploy mounts an operator file.
// Database DSN is owned by that YAML; PAYMENT_DATABASE_URL must not fight it.
func Load() (Config, error) {
	selection, err := configloader.LoadRuntimeYML()
	if err != nil {
		return Config{}, err
	}

	host := setting("PAYMENT_GRPC_HOST", "host", "0.0.0.0")
	port := setting("PAYMENT_GRPC_PORT", "port", "8205")
	address := firstNonBlank(
		os.Getenv("PAYMENT_GRPC_ADDRESS"),
		host+":"+port,
	)

	// Database DSN is owned by the environment YAML:
	// - workstation: config/runtime.yml merged with config/local.yml
	// - deploy: CONFIG_FILE / QHPRO_CONFIG_FILE
	// Accept legacy database.url overlays that predate database.dsn.
	dsn := firstNonBlank(viper.GetString("database.dsn"), viper.GetString("database.url"))
	if strings.TrimSpace(dsn) == "" || strings.Contains(strings.ToUpper(dsn), "CHANGE_ME") {
		return Config{}, errors.New("database.dsn is required in the environment config YAML")
	}

	cfg := Config{
		Environment: selection.Environment,
		Server: ServerConfig{
			Address:         address,
			ShutdownTimeout: durationSetting("PAYMENT_SHUTDOWN_TIMEOUT", "shutdown_timeout", 15*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             dsn,
			MaxOpenConns:    intSetting("PAYMENT_DATABASE_MAX_OPEN_CONNS", "database.max_open_conns", 30),
			MaxIdleConns:    intSetting("PAYMENT_DATABASE_MAX_IDLE_CONNS", "database.max_idle_conns", 10),
			ConnMaxLifetime: durationSetting("PAYMENT_DATABASE_CONN_MAX_LIFETIME", "database.conn_max_lifetime", 5*time.Minute),
		},
		User: ClientConfig{
			Address: setting("QHPRO_USER_GRPC_ADDRESS", "rpc.user.address", "localhost:8201"),
			Timeout: durationSetting("PAYMENT_USER_RPC_TIMEOUT", "rpc.user.timeout", 5*time.Second),
		},
		Auth: ClientConfig{
			Address: setting("QHPRO_AUTH_GRPC_ADDRESS", "rpc.auth.address", "localhost:8216"),
			Timeout: durationSetting("PAYMENT_AUTH_RPC_TIMEOUT", "rpc.auth.timeout", 5*time.Second),
		},
		Notification: ClientConfig{
			Address: notificationAddress(),
			Timeout: durationSetting("PAYMENT_NOTIFICATION_RPC_TIMEOUT", "rpc.notification.timeout", 5*time.Second),
		},
		Commerce: CommerceConfig{
			OrderTTL:         durationSetting("PAYMENT_ORDER_TTL", "commerce.order_ttl", 30*time.Minute),
			FulfillmentRetry: durationSetting("PAYMENT_FULFILLMENT_RETRY_DELAY", "commerce.fulfillment_retry", 5*time.Second),
			FulfillmentLease: durationSetting("PAYMENT_FULFILLMENT_CLAIM_LEASE", "commerce.fulfillment_lease", 30*time.Second),
			FulfillmentPoll:  durationSetting("PAYMENT_FULFILLMENT_POLL_INTERVAL", "commerce.fulfillment_poll", time.Second),
		},
		Rabbit: RabbitConfig{
			URL:            setting("PAYMENT_RABBIT_URL", "rabbit.url", ""),
			Exchange:       setting("PAYMENT_EVENT_EXCHANGE", "rabbit.exchange", "qhpro.payment.events"),
			PublisherLease: durationSetting("PAYMENT_OUTBOX_CLAIM_LEASE", "rabbit.publisher_lease", 30*time.Second),
			PublisherPoll:  durationSetting("PAYMENT_OUTBOX_POLL_INTERVAL", "rabbit.publisher_poll", time.Second),
			PublisherRetry: durationSetting("PAYMENT_OUTBOX_RETRY_DELAY", "rabbit.publisher_retry", 5*time.Second),
		},
		Provider: ProviderConfig{
			SepayAPIKey: setting("SEPAY_API_KEY", "sepay.api_key", ""),
			Sepay: SepayGatewayConfig{
				Environment:        setting("SEPAY_GATEWAY_ENVIRONMENT", "sepay.gateway_environment", "sandbox"),
				BankTransferFlow:   strings.ToLower(setting("SEPAY_BANK_TRANSFER_FLOW", "sepay.bank_transfer_flow", SepayBankTransferHosted)),
				WebhookEnvironment: strings.ToLower(setting("SEPAY_WEBHOOK_ENVIRONMENT", "sepay.webhook_environment", SepayWebhookLive)),
				MerchantID:         setting("SEPAY_MERCHANT_ID", "sepay.merchant_id", ""),
				SecretKey:          setting("SEPAY_MERCHANT_SECRET_KEY", "sepay.merchant_secret_key", ""),
				CheckoutURL:        setting("SEPAY_CHECKOUT_URL", "sepay.checkout_url", ""),
				BridgeURL: setting(
					"SEPAY_CHECKOUT_BRIDGE_URL",
					"sepay.checkout_bridge_url",
					"http://127.0.0.1:8000/v2/payment/sepay/checkout",
				),
				SuccessURL: setting("SEPAY_SUCCESS_URL", "sepay.success_url", ""),
				ErrorURL:   setting("SEPAY_ERROR_URL", "sepay.error_url", ""),
				CancelURL:  setting("SEPAY_CANCEL_URL", "sepay.cancel_url", ""),
				IPNSecret:  setting("SEPAY_IPN_SECRET_KEY", "sepay.ipn_secret_key", ""),
				// Availability is operator-owned. Bank transfer keeps the
				// historical default; international cards require an explicit
				// opt-in after the SePay merchant has been approved for CARD.
				BankTransferEnabled: boolSetting("SEPAY_BANK_TRANSFER_ENABLED", "sepay.bank_transfer_enabled", true),
				CardEnabled:         boolSetting("SEPAY_CARD_ENABLED", "sepay.card_enabled", false),
				BankBIN:             setting("SEPAY_QR_BANK_BIN", "sepay.qr_bank_bin", ""),
				BankAccountNumber:   setting("SEPAY_QR_ACCOUNT_NUMBER", "sepay.qr_account_number", ""),
			},
		},
	}
	if cfg.Server.Address == "" || cfg.Database.MaxOpenConns <= 0 || cfg.Database.MaxIdleConns < 0 ||
		cfg.Commerce.OrderTTL <= 0 || cfg.Commerce.FulfillmentLease <= 0 || cfg.Commerce.FulfillmentPoll <= 0 {
		return Config{}, fmt.Errorf("invalid Payment runtime configuration")
	}
	if err := validateProvider(cfg.Provider); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func notificationAddress() string {
	if value := setting("QHPRO_NOTIFICATION_GRPC_ADDRESS", "rpc.notification.address", ""); value != "" {
		return value
	}
	// Legacy overlays used grpc.notification.host/port instead of rpc.notification.address.
	host := strings.TrimSpace(viper.GetString("grpc.notification.host"))
	port := strings.TrimSpace(viper.GetString("grpc.notification.port"))
	if host != "" && port != "" {
		return host + ":" + port
	}
	return "localhost:8204"
}

func validateProvider(cfg ProviderConfig) error {
	if err := validateSepayGateway(cfg.Sepay); err != nil {
		return err
	}
	if strings.EqualFold(strings.TrimSpace(cfg.Sepay.BankTransferFlow), SepayBankTransferVietQR) {
		webhookEnvironment := strings.ToLower(strings.TrimSpace(cfg.Sepay.WebhookEnvironment))
		if webhookEnvironment == "" {
			webhookEnvironment = SepayWebhookLive
		}
		if webhookEnvironment != SepayWebhookTest && webhookEnvironment != SepayWebhookLive {
			return errors.New("SEPAY_WEBHOOK_ENVIRONMENT must be test or live")
		}
		// Bank Webhooks Test mode is independent of Payment Gateway sandbox.
		// Require an explicit selection before treating an inline QR as a test.
		sandbox := strings.EqualFold(strings.TrimSpace(cfg.Sepay.Environment), "sandbox")
		if sandbox != (webhookEnvironment == SepayWebhookTest) {
			return errors.New("VietQR requires sandbox + SEPAY_WEBHOOK_ENVIRONMENT=test or production + SEPAY_WEBHOOK_ENVIRONMENT=live")
		}
		if strings.TrimSpace(cfg.SepayAPIKey) == "" {
			return errors.New("SEPAY_API_KEY is required for authenticated bank webhook settlement when SEPAY_BANK_TRANSFER_FLOW=vietqr")
		}
	}
	return nil
}

func validateSepayGateway(cfg SepayGatewayConfig) error {
	environment := strings.ToLower(strings.TrimSpace(cfg.Environment))
	if environment != "sandbox" && environment != "production" {
		return fmt.Errorf("SEPAY_GATEWAY_ENVIRONMENT must be sandbox or production")
	}
	flow := strings.ToLower(strings.TrimSpace(cfg.BankTransferFlow))
	if flow == "" {
		flow = SepayBankTransferHosted
	}
	if flow != SepayBankTransferHosted && flow != SepayBankTransferVietQR {
		return fmt.Errorf("SEPAY_BANK_TRANSFER_FLOW must be hosted_checkout or vietqr")
	}

	configured := cfg.MerchantID != "" || cfg.SecretKey != "" || cfg.IPNSecret != ""
	usesHostedCheckout := flow == SepayBankTransferHosted || cfg.CardEnabled
	if usesHostedCheckout && (configured || cfg.CardEnabled) {
		if cfg.MerchantID == "" || cfg.SecretKey == "" || cfg.IPNSecret == "" {
			return fmt.Errorf("SEPAY_MERCHANT_ID, SEPAY_MERCHANT_SECRET_KEY and SEPAY_IPN_SECRET_KEY must be configured together")
		}
		bridge, err := url.Parse(strings.TrimSpace(cfg.BridgeURL))
		if err != nil || bridge.Scheme == "" || bridge.Host == "" {
			return fmt.Errorf("SEPAY_CHECKOUT_BRIDGE_URL must be an absolute HTTP(S) URL")
		}
		if environment == "production" && bridge.Scheme != "https" {
			return fmt.Errorf("SEPAY_CHECKOUT_BRIDGE_URL must use HTTPS in production")
		}
	}

	if flow == SepayBankTransferVietQR {
		if !cfg.BankTransferEnabled {
			return fmt.Errorf("SEPAY_BANK_TRANSFER_ENABLED must be true when SEPAY_BANK_TRANSFER_FLOW=vietqr")
		}
		if !fixedDigits(cfg.BankBIN, 6) {
			return fmt.Errorf("SEPAY_QR_BANK_BIN must contain exactly 6 digits")
		}
		if !digitsBetween(cfg.BankAccountNumber, 5, 19) {
			return fmt.Errorf("SEPAY_QR_ACCOUNT_NUMBER must contain 5 to 19 digits")
		}
	}
	return nil
}

func fixedDigits(value string, length int) bool {
	return digitsBetween(value, length, length)
}

func digitsBetween(value string, minimum, maximum int) bool {
	value = strings.TrimSpace(value)
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

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return ""
}

// setting prefers an explicit process env override, then the YAML key, then fallback.
// AutomaticEnv already maps many YAML keys; the env key here covers historical PAYMENT_*/QHPRO_* names.
func setting(envKey, yamlKey, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		return value
	}
	if value := strings.TrimSpace(viper.GetString(yamlKey)); value != "" && !strings.Contains(strings.ToUpper(value), "CHANGE_ME") {
		return value
	}
	return fallback
}

func durationSetting(envKey, yamlKey string, fallback time.Duration) time.Duration {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	if raw := strings.TrimSpace(viper.GetString(yamlKey)); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func intSetting(envKey, yamlKey string, fallback int) int {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			return parsed
		}
	}
	if viper.IsSet(yamlKey) {
		if parsed := viper.GetInt(yamlKey); parsed >= 0 {
			return parsed
		}
	}
	return fallback
}

func boolSetting(envKey, yamlKey string, fallback bool) bool {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	if viper.IsSet(yamlKey) {
		return viper.GetBool(yamlKey)
	}
	return fallback
}
