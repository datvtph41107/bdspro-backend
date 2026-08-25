package transporttrust

import (
	_callercontext "common/callercontext"
	_operationcontext "common/operationcontext"
	_quotasubject "common/quotasubject"

	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

const (
	CallerMetadataKey    = "x-qhpro-internal-caller"
	IssuedAtMetadataKey  = "x-qhpro-internal-issued-at"
	SignatureMetadataKey = "x-qhpro-internal-signature"
	VersionMetadataKey   = "x-qhpro-internal-version"
	SignatureVersion     = "v1"
)

type Mode string

const (
	ModeAudit   Mode = "audit"
	ModeEnforce Mode = "enforce"
)

type Level string

const (
	LevelLegacyUnverified Level = "legacy_unverified"
	LevelVerified         Level = "verified"
)

var (
	ErrNotConfigured    = errors.New("transport trust is not configured")
	ErrMissingAssertion = errors.New("trusted metadata assertion is missing")
	ErrInvalidAssertion = errors.New("trusted metadata assertion is invalid")
	ErrExpiredAssertion = errors.New("trusted metadata assertion is expired")
	ErrCallerConflict   = errors.New("transport caller conflicts with the existing context")
)

/**
 * Config controls internal metadata signing and verification.
 */
type Config struct {
	ServiceID                  string
	Secret                     string
	VerificationSecrets        map[string]string
	VerificationKeysConfigured bool
	Mode                       Mode
	MaxAge                     time.Duration
	ClockSkew                  time.Duration
	Now                        func() time.Time
}

/**
 * Caller describes the verified internal transport caller.
 */
type Caller struct {
	ServiceID string
	Level     Level
	IssuedAt  time.Time
}

type callerContextKey struct{}

/**
 * ConfigFromEnv loads trust config from environment, then YAML (viper) fallback.
 * Env always wins so production can inject secrets without committing them.
 *
 * YAML keys (after common.LoadYMLConfigFile):
 *   trusted_metadata.service_id
 *   trusted_metadata.secret
 *   trusted_metadata.mode
 *   trusted_metadata.max_age
 *   trusted_metadata.clock_skew
 * Service ID also falls back to server.name when unset.
 */
func ConfigFromEnv() Config {
	mode := Mode(strings.ToLower(strings.TrimSpace(firstNonEmpty(
		os.Getenv("QHPRO_TRUSTED_METADATA_MODE"),
		viper.GetString("trusted_metadata.mode"),
	))))
	if mode != ModeEnforce {
		mode = ModeAudit
	}

	verificationSecrets, verificationKeysConfigured := verificationSecretsFromEnv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
	)

	return Config{
		ServiceID: firstNonEmpty(
			os.Getenv("QHPRO_SERVICE_ID"),
			os.Getenv("SERVICE_NAME"),
			viper.GetString("trusted_metadata.service_id"),
			viper.GetString("server.name"),
		),
		Secret: firstNonEmpty(
			os.Getenv("QHPRO_INTERNAL_METADATA_SECRET"),
			viper.GetString("trusted_metadata.secret"),
		),
		VerificationSecrets:        verificationSecrets,
		VerificationKeysConfigured: verificationKeysConfigured,
		Mode:                       mode,
		MaxAge: durationFromEnvOrViper(
			"QHPRO_INTERNAL_METADATA_MAX_AGE",
			"trusted_metadata.max_age",
			30*time.Second,
		),
		ClockSkew: durationFromEnvOrViper(
			"QHPRO_INTERNAL_METADATA_CLOCK_SKEW",
			"trusted_metadata.clock_skew",
			5*time.Second,
		),
	}
}

/**
 * BindCaller immutably binds the verified transport provenance.
 */
func BindCaller(ctx context.Context, caller Caller) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !caller.IsValid() {
		return ctx, ErrInvalidAssertion
	}
	if existing, ok := FromContext(ctx); ok {
		if existing == caller {
			return ctx, nil
		}
		return ctx, ErrCallerConflict
	}
	return context.WithValue(ctx, callerContextKey{}, caller), nil
}

/**
 * FromContext returns the transport caller.
 */
func FromContext(ctx context.Context) (Caller, bool) {
	if ctx == nil {
		return Caller{}, false
	}
	caller, ok := ctx.Value(callerContextKey{}).(Caller)
	return caller, ok && caller.IsValid()
}

func (c Caller) IsValid() bool {
	if !validServiceID(c.ServiceID) {
		return false
	}
	switch c.Level {
	case LevelVerified:
		return !c.IssuedAt.IsZero()
	case LevelLegacyUnverified:
		return c.ServiceID == "legacy-unverified" && c.IssuedAt.IsZero()
	default:
		return false
	}
}

/**
 * IsVerified reports whether metadata came from a verified internal caller.
 */
func IsVerified(ctx context.Context) bool {
	caller, ok := FromContext(ctx)
	return ok && caller.Level == LevelVerified
}

/**
 * HasPrivilegedMetadata detects identity and entitlement assertions.
 */
func HasPrivilegedMetadata(md metadata.MD) bool {
	for _, key := range privilegedKeys {
		if len(md.Get(key)) > 0 {
			return true
		}
	}
	return false
}

/**
 * HasAssertion detects transport trust metadata.
 */
func HasAssertion(md metadata.MD) bool {
	return len(md.Get(CallerMetadataKey)) > 0 ||
		len(md.Get(IssuedAtMetadataKey)) > 0 ||
		len(md.Get(SignatureMetadataKey)) > 0 ||
		len(md.Get(VersionMetadataKey)) > 0
}

/**
 * Sign adds a method-bound HMAC assertion to outgoing metadata.
 */
func Sign(method string, md metadata.MD, cfg Config) error {
	if strings.TrimSpace(cfg.ServiceID) == "" || strings.TrimSpace(cfg.Secret) == "" {
		return ErrNotConfigured
	}
	if !validServiceID(cfg.ServiceID) {
		return fmt.Errorf("%w: invalid service ID", ErrInvalidAssertion)
	}
	if strings.TrimSpace(method) == "" {
		return fmt.Errorf("%w: method is empty", ErrInvalidAssertion)
	}

	now := configNow(cfg).UTC()
	issuedAt := strconv.FormatInt(now.Unix(), 10)

	md.Set(VersionMetadataKey, SignatureVersion)
	md.Set(CallerMetadataKey, cfg.ServiceID)
	md.Set(IssuedAtMetadataKey, issuedAt)
	md.Delete(SignatureMetadataKey)

	signature := calculateSignature(method, md, cfg.Secret)
	md.Set(SignatureMetadataKey, signature)
	return nil
}

/**
 * Verify validates a method-bound HMAC assertion.
 */
func Verify(method string, md metadata.MD, cfg Config) (Caller, error) {
	version, err := single(md, VersionMetadataKey)
	if err != nil || version != SignatureVersion {
		return Caller{}, ErrMissingAssertion
	}
	caller, err := single(md, CallerMetadataKey)
	if err != nil || !validServiceID(caller) {
		return Caller{}, ErrInvalidAssertion
	}
	verificationSecret, err := secretForCaller(cfg, caller)
	if err != nil {
		return Caller{}, err
	}
	issuedAtRaw, err := single(md, IssuedAtMetadataKey)
	if err != nil {
		return Caller{}, ErrInvalidAssertion
	}
	signature, err := single(md, SignatureMetadataKey)
	if err != nil || len(signature) != sha256.Size*2 {
		return Caller{}, ErrInvalidAssertion
	}

	issuedUnix, err := strconv.ParseInt(issuedAtRaw, 10, 64)
	if err != nil {
		return Caller{}, ErrInvalidAssertion
	}
	issuedAt := time.Unix(issuedUnix, 0).UTC()
	now := configNow(cfg).UTC()
	maxAge := cfg.MaxAge
	if maxAge <= 0 {
		maxAge = 30 * time.Second
	}
	clockSkew := cfg.ClockSkew
	if clockSkew < 0 {
		clockSkew = 0
	}
	if issuedAt.After(now.Add(clockSkew)) || now.Sub(issuedAt) > maxAge+clockSkew {
		return Caller{}, ErrExpiredAssertion
	}

	copyMD := md.Copy()
	copyMD.Delete(SignatureMetadataKey)
	expected := calculateSignature(method, copyMD, verificationSecret)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return Caller{}, ErrInvalidAssertion
	}

	return Caller{ServiceID: caller, Level: LevelVerified, IssuedAt: issuedAt}, nil
}

var privilegedKeys = []string{
	_callercontext.CallerKindMetadataKey,
	_callercontext.APIKeyIDMetadataKey,
	_callercontext.APIKeyAppMetadataKey,
	_quotasubject.KindMetadataKey,
	_quotasubject.IDMetadataKey,
	_quotasubject.SourceMetadataKey,
	_operationcontext.CodeMetadataKey,
	_operationcontext.BindingMetadataKey,
	_operationcontext.SourceMetadataKey,
	"x-api-key",
	"authid",
	"profileid",
	"originid",
	"organizationid",
	"session",
	"role",
	"type",
	"planid",
	"planfrom",
}

var signedKeys = []string{
	_callercontext.CallerKindMetadataKey,
	_callercontext.APIKeyIDMetadataKey,
	_callercontext.APIKeyAppMetadataKey,
	_quotasubject.KindMetadataKey,
	_quotasubject.IDMetadataKey,
	_quotasubject.SourceMetadataKey,
	_operationcontext.CodeMetadataKey,
	_operationcontext.BindingMetadataKey,
	_operationcontext.SourceMetadataKey,
	"x-api-key",
	"authid",
	"profileid",
	"originid",
	"organizationid",
	"session",
	"role",
	"type",
	"planid",
	"planfrom",
	"x-request-id",
	"x-operation-id",
	"idempotency-key",
	CallerMetadataKey,
	IssuedAtMetadataKey,
	VersionMetadataKey,
}

func calculateSignature(method string, md metadata.MD, secret string) string {
	keys := append([]string(nil), signedKeys...)
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString("method=")
	builder.WriteString(method)
	builder.WriteByte('\n')
	for _, key := range keys {
		values := append([]string(nil), md.Get(key)...)
		sort.Strings(values)
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(strings.Join(values, "\x1f"))
		builder.WriteByte('\n')
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(builder.String()))
	return hex.EncodeToString(mac.Sum(nil))
}

func single(md metadata.MD, key string) (string, error) {
	values := md.Get(key)
	if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
		return "", ErrInvalidAssertion
	}
	return values[0], nil
}

func validServiceID(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == ':' {
			continue
		}
		return false
	}
	return true
}

func configNow(cfg Config) time.Time {
	if cfg.Now != nil {
		return cfg.Now()
	}
	return time.Now()
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
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

func durationFromEnvOrViper(envKey, viperKey string, fallback time.Duration) time.Duration {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	if value := strings.TrimSpace(viper.GetString(viperKey)); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func secretForCaller(cfg Config, caller string) (string, error) {
	if cfg.VerificationKeysConfigured || len(cfg.VerificationSecrets) > 0 {
		secret := strings.TrimSpace(cfg.VerificationSecrets[caller])
		if secret == "" {
			return "", ErrInvalidAssertion
		}
		return secret, nil
	}
	secret := strings.TrimSpace(cfg.Secret)
	if secret == "" {
		return "", ErrNotConfigured
	}
	return secret, nil
}

func verificationSecretsFromEnv(key string) (map[string]string, bool) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil, false
	}

	var values map[string]string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, true
	}

	result := make(map[string]string, len(values))
	for serviceID, secret := range values {
		serviceID = strings.TrimSpace(serviceID)
		secret = strings.TrimSpace(secret)
		if validServiceID(serviceID) && secret != "" {
			result[serviceID] = secret
		}
	}
	if len(result) == 0 {
		return nil, true
	}
	return result, true
}
