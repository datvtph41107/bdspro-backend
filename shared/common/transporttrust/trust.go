// Package transporttrust is the compatibility surface for the historical QHPRO
// trusted-metadata contract. New transport code should depend on common/trust
// through the canonical rpc package once F01-I1B-4 is complete.
package transporttrust

import (
	_rpc "common/rpc"
	_trust "common/trust"

	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	CallerMetadataKey    = _trust.CallerMetadataKey
	IssuedAtMetadataKey  = _trust.IssuedAtMetadataKey
	SignatureMetadataKey = _trust.SignatureMetadataKey
	VersionMetadataKey   = _trust.VersionMetadataKey
	SignatureVersion     = _trust.SignatureVersion
)

type Mode = _trust.Mode

const (
	ModeAudit   = _trust.ModeAudit
	ModeEnforce = _trust.ModeEnforce
)

type Level = _trust.Level

const (
	LevelLegacyUnverified = _trust.LevelLegacyUnverified
	LevelVerified         = _trust.LevelVerified
)

var (
	ErrNotConfigured    = _trust.ErrNotConfigured
	ErrMissingAssertion = _trust.ErrMissingAssertion
	ErrInvalidAssertion = _trust.ErrInvalidAssertion
	ErrExpiredAssertion = _trust.ErrExpiredAssertion
	ErrCallerConflict   = _trust.ErrCallerConflict
)

type Config = _trust.Config
type Caller = _trust.Caller
type Assessment = _trust.Assessment

// The protected QHPRO metadata set is now owned by common/rpc.
// ConfigFromEnv preserves the legacy runtime configuration surface while
// canonical trust remains process-neutral and explicit.
func ConfigFromEnv() Config {
	mode := Mode(strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_TRUSTED_METADATA_MODE"))))
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
		),
		Secret:                     strings.TrimSpace(os.Getenv("QHPRO_INTERNAL_METADATA_SECRET")),
		VerificationSecrets:        verificationSecrets,
		VerificationKeysConfigured: verificationKeysConfigured,
		Mode:                       mode,
		MaxAge:                     durationFromEnv("QHPRO_INTERNAL_METADATA_MAX_AGE", 30*time.Second),
		ClockSkew:                  durationFromEnv("QHPRO_INTERNAL_METADATA_CLOCK_SKEW", 5*time.Second),
	}
}

func BindCaller(ctx context.Context, caller Caller) (context.Context, error) {
	return _trust.BindCaller(ctx, caller)
}

func FromContext(ctx context.Context) (Caller, bool) {
	return _trust.FromContext(ctx)
}

func IsVerified(ctx context.Context) bool {
	return _trust.IsVerified(ctx)
}

func HasPrivilegedMetadata(md metadata.MD) bool {
	return _trust.HasPrivilegedMetadata(md, _rpc.QHPROMetadataContract())
}

func HasAssertion(md metadata.MD) bool {
	return _trust.HasAssertion(md)
}

func Sign(method string, md metadata.MD, cfg Config) error {
	return _trust.Sign(method, md, cfg, _rpc.QHPROMetadataContract())
}

func Verify(method string, md metadata.MD, cfg Config) (Caller, error) {
	return _trust.Verify(method, md, cfg, _rpc.QHPROMetadataContract())
}

func Assess(method string, md metadata.MD, cfg Config) Assessment {
	return _trust.Assess(method, md, cfg, _rpc.QHPROMetadataContract())
}

func CanPropagateLegacyPrivilegedContext(ctx context.Context) bool {
	return _trust.CanPropagateLegacyPrivilegedContext(ctx)
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
