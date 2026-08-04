package transporttrust

import (
	_callercontext "common/callercontext"
	_operationcontext "common/operationcontext"
	_quotasubject "common/quotasubject"

	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
)

/**
 * Config controls internal metadata signing and verification.
 */
type Config struct {
	ServiceID string
	Secret    string
	Mode      Mode
	MaxAge    time.Duration
	ClockSkew time.Duration
	Now       func() time.Time
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
 * ConfigFromEnv loads non-secret policy and secret material from environment.
 */
func ConfigFromEnv() Config {
	mode := Mode(strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_TRUSTED_METADATA_MODE"))))
	if mode != ModeEnforce {
		mode = ModeAudit
	}

	return Config{
		ServiceID: firstNonEmpty(
			os.Getenv("QHPRO_SERVICE_ID"),
			os.Getenv("SERVICE_NAME"),
		),
		Secret:    strings.TrimSpace(os.Getenv("QHPRO_INTERNAL_METADATA_SECRET")),
		Mode:      mode,
		MaxAge:    durationFromEnv("QHPRO_INTERNAL_METADATA_MAX_AGE", 30*time.Second),
		ClockSkew: durationFromEnv("QHPRO_INTERNAL_METADATA_CLOCK_SKEW", 5*time.Second),
	}
}

/**
 * WithCaller stores the transport caller in context.
 */
func WithCaller(ctx context.Context, caller Caller) context.Context {
	if ctx == nil || caller.ServiceID == "" || caller.Level == "" {
		return ctx
	}
	return context.WithValue(ctx, callerContextKey{}, caller)
}

/**
 * FromContext returns the transport caller.
 */
func FromContext(ctx context.Context) (Caller, bool) {
	if ctx == nil {
		return Caller{}, false
	}
	caller, ok := ctx.Value(callerContextKey{}).(Caller)
	return caller, ok && caller.ServiceID != "" && caller.Level != ""
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
	if strings.TrimSpace(cfg.Secret) == "" {
		return Caller{}, ErrNotConfigured
	}

	version, err := single(md, VersionMetadataKey)
	if err != nil || version != SignatureVersion {
		return Caller{}, ErrMissingAssertion
	}
	caller, err := single(md, CallerMetadataKey)
	if err != nil || !validServiceID(caller) {
		return Caller{}, ErrInvalidAssertion
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
	expected := calculateSignature(method, copyMD, cfg.Secret)
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
