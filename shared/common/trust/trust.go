// Package trust owns process-neutral service-hop provenance and cryptographic
// verification. It deliberately does not know Gateway, business services,
// request identity semantics, quota semantics, or environment configuration.
package trust

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

var assertionMetadataKeys = []string{
	CallerMetadataKey,
	IssuedAtMetadataKey,
	VersionMetadataKey,
}

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

// Config contains explicit trust policy and key material. Environment loading
// belongs to the process/configuration adapter, not to this package.
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

// Caller describes the immediate service-hop provenance.
type Caller struct {
	ServiceID string
	Level     Level
	IssuedAt  time.Time
}

type callerContextKey struct{}

// BindCaller immutably binds one transport provenance value.
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

// FromContext returns the canonical immediate service-hop provenance.
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

// IsVerified reports whether the context has verified service-hop provenance.
func IsVerified(ctx context.Context) bool {
	caller, ok := FromContext(ctx)
	return ok && caller.Level == LevelVerified
}

// HasPrivilegedMetadata checks the application-specific privileged key set
// supplied by the transport adapter.
func HasPrivilegedMetadata(md metadata.MD, contract MetadataContract) bool {
	return contract.hasPrivilegedMetadata(md)
}

// HasAssertion detects canonical transport trust assertion fields.
func HasAssertion(md metadata.MD) bool {
	return len(md.Get(CallerMetadataKey)) > 0 ||
		len(md.Get(IssuedAtMetadataKey)) > 0 ||
		len(md.Get(SignatureMetadataKey)) > 0 ||
		len(md.Get(VersionMetadataKey)) > 0
}

// Sign adds a method-bound HMAC assertion to outgoing metadata.
func Sign(method string, md metadata.MD, cfg Config, contract MetadataContract) error {
	if strings.TrimSpace(cfg.ServiceID) == "" || strings.TrimSpace(cfg.Secret) == "" {
		return ErrNotConfigured
	}
	if !validServiceID(cfg.ServiceID) {
		return fmt.Errorf("%w: invalid service ID", ErrInvalidAssertion)
	}
	if strings.TrimSpace(method) == "" || md == nil {
		return fmt.Errorf("%w: method or metadata is empty", ErrInvalidAssertion)
	}

	now := configNow(cfg).UTC()
	issuedAt := strconv.FormatInt(now.Unix(), 10)

	md.Set(VersionMetadataKey, SignatureVersion)
	md.Set(CallerMetadataKey, cfg.ServiceID)
	md.Set(IssuedAtMetadataKey, issuedAt)
	md.Delete(SignatureMetadataKey)

	signature := calculateSignature(method, md, cfg.Secret, contract)
	md.Set(SignatureMetadataKey, signature)
	return nil
}

// Verify validates a method-bound HMAC assertion.
func Verify(method string, md metadata.MD, cfg Config, contract MetadataContract) (Caller, error) {
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
	expected := calculateSignature(method, copyMD, verificationSecret, contract)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return Caller{}, ErrInvalidAssertion
	}

	return Caller{ServiceID: caller, Level: LevelVerified, IssuedAt: issuedAt}, nil
}

// CanPropagateLegacyPrivilegedContext prevents an explicitly unverified inbound
// context from being re-signed as trusted by a downstream client.
func CanPropagateLegacyPrivilegedContext(ctx context.Context) bool {
	caller, ok := FromContext(ctx)
	return !ok || caller.Level == LevelVerified
}

func calculateSignature(method string, md metadata.MD, secret string, contract MetadataContract) string {
	keys := append([]string(nil), contract.signedKeys...)
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
