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

func Load() (Config, error) {
	selection, err := configloader.ResolveRuntimeSelection()
	if err != nil {
		return Config{}, err
	}
	v := viper.New()
	v.SetConfigType("yaml")
	file := resolveConfigFile(os.Getenv("CONFIG_FILE"))
	v.SetConfigFile(file)
	// YAML là default không nhạy cảm; environment là override/secret. Thiếu hoặc
	// lỗi YAML phải fail sớm để dev và container dùng cùng một contract.
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("load Payment config %s: %w", file, err)
	}

	host := envOr("PAYMENT_GRPC_HOST", v.GetString("host"), "0.0.0.0")
	port := envOr("PAYMENT_GRPC_PORT", v.GetString("port"), "8205")
	address := envOr("PAYMENT_GRPC_ADDRESS", "", host+":"+port)
	dsn := firstNonBlank(os.Getenv("PAYMENT_DATABASE_URL"), os.Getenv("PAYMENT_DATABASE_DSN"), v.GetString("database.url"))
	if strings.TrimSpace(dsn) == "" || strings.Contains(strings.ToUpper(dsn), "CHANGE_ME") {
		return Config{}, errors.New("PAYMENT_DATABASE_URL is required")
	}

	cfg := Config{
		Environment: selection.Environment,
		Server: ServerConfig{
			Address:         address,
			ShutdownTimeout: durationEnv("PAYMENT_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             dsn,
			MaxOpenConns:    intEnv("PAYMENT_DATABASE_MAX_OPEN_CONNS", 30),
			MaxIdleConns:    intEnv("PAYMENT_DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: durationEnv("PAYMENT_DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		User: ClientConfig{
			Address: envOr("QHPRO_USER_GRPC_ADDRESS", "", "localhost:8201"),
			Timeout: durationEnv("PAYMENT_USER_RPC_TIMEOUT", 5*time.Second),
		},
		Auth: ClientConfig{
			Address: envOr("QHPRO_AUTH_GRPC_ADDRESS", "", "localhost:8216"),
			Timeout: durationEnv("PAYMENT_AUTH_RPC_TIMEOUT", 5*time.Second),
		},
		Notification: ClientConfig{
			Address: envOr("QHPRO_NOTIFICATION_GRPC_ADDRESS", "", "localhost:8204"),
			Timeout: durationEnv("PAYMENT_NOTIFICATION_RPC_TIMEOUT", 5*time.Second),
		},
		Commerce: CommerceConfig{
			OrderTTL:         durationEnv("PAYMENT_ORDER_TTL", 30*time.Minute),
			FulfillmentRetry: durationEnv("PAYMENT_FULFILLMENT_RETRY_DELAY", 5*time.Second),
			FulfillmentLease: durationEnv("PAYMENT_FULFILLMENT_CLAIM_LEASE", 30*time.Second),
			FulfillmentPoll:  durationEnv("PAYMENT_FULFILLMENT_POLL_INTERVAL", time.Second),
		},
		Rabbit: RabbitConfig{
			URL:            strings.TrimSpace(os.Getenv("PAYMENT_RABBIT_URL")),
			Exchange:       envOr("PAYMENT_EVENT_EXCHANGE", "", "qhpro.payment.events"),
			PublisherLease: durationEnv("PAYMENT_OUTBOX_CLAIM_LEASE", 30*time.Second),
			PublisherPoll:  durationEnv("PAYMENT_OUTBOX_POLL_INTERVAL", time.Second),
			PublisherRetry: durationEnv("PAYMENT_OUTBOX_RETRY_DELAY", 5*time.Second),
		},
		Provider: ProviderConfig{
			SepayAPIKey: strings.TrimSpace(os.Getenv("SEPAY_API_KEY")),
			Sepay: SepayGatewayConfig{
				Environment:        envOr("SEPAY_GATEWAY_ENVIRONMENT", "", "sandbox"),
				BankTransferFlow:   strings.ToLower(envOr("SEPAY_BANK_TRANSFER_FLOW", "", SepayBankTransferHosted)),
				WebhookEnvironment: strings.ToLower(envOr("SEPAY_WEBHOOK_ENVIRONMENT", "", SepayWebhookLive)),
				MerchantID:         strings.TrimSpace(os.Getenv("SEPAY_MERCHANT_ID")),
				SecretKey:          strings.TrimSpace(os.Getenv("SEPAY_MERCHANT_SECRET_KEY")),
				CheckoutURL:        strings.TrimSpace(os.Getenv("SEPAY_CHECKOUT_URL")),
				BridgeURL: envOr(
					"SEPAY_CHECKOUT_BRIDGE_URL",
					"",
					"http://127.0.0.1:8000/v2/payment/sepay/checkout",
				),
				SuccessURL: strings.TrimSpace(os.Getenv("SEPAY_SUCCESS_URL")),
				ErrorURL:   strings.TrimSpace(os.Getenv("SEPAY_ERROR_URL")),
				CancelURL:  strings.TrimSpace(os.Getenv("SEPAY_CANCEL_URL")),
				IPNSecret:  strings.TrimSpace(os.Getenv("SEPAY_IPN_SECRET_KEY")),
				// Availability is operator-owned. Bank transfer keeps the
				// historical default; international cards require an explicit
				// opt-in after the SePay merchant has been approved for CARD.
				BankTransferEnabled: boolEnv("SEPAY_BANK_TRANSFER_ENABLED", true),
				CardEnabled:         boolEnv("SEPAY_CARD_ENABLED", false),
				BankBIN:             strings.TrimSpace(os.Getenv("SEPAY_QR_BANK_BIN")),
				BankAccountNumber:   strings.TrimSpace(os.Getenv("SEPAY_QR_ACCOUNT_NUMBER")),
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

// resolveConfigFile uses one committed runtime default. CONFIG_FILE is the
// explicit operator escape hatch for an externally mounted configuration.
func resolveConfigFile(override string) string {
	if file := strings.TrimSpace(override); file != "" {
		return file
	}
	return "config/runtime.yml"
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return ""
}

func envOr(key, compatibility, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	if value := strings.TrimSpace(compatibility); value != "" && !strings.Contains(strings.ToUpper(value), "CHANGE_ME") {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func intEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

func boolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
