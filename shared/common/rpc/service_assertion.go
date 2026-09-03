package rpc

import (
	"common/identity"
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
	ServiceAssertionCallerMetadataKey    = "x-qhpro-internal-caller"
	ServiceAssertionIssuedAtMetadataKey  = "x-qhpro-internal-issued-at"
	ServiceAssertionSignatureMetadataKey = "x-qhpro-internal-signature"
	ServiceAssertionVersionMetadataKey   = "x-qhpro-internal-version"
	ServiceAssertionVersion              = "v1"
)

var (
	ErrServiceAssertionNotConfigured = errors.New("service assertion is not configured")
	ErrServiceAssertionMissing       = errors.New("service assertion is missing")
	ErrServiceAssertionInvalid       = errors.New("service assertion is invalid")
	ErrServiceAssertionExpired       = errors.New("service assertion is expired")
)

// ServiceAssertionConfig contains the key material and verification window for
// the immediate authenticated service hop.
type ServiceAssertionConfig struct {
	ServiceID                  string
	Secret                     string
	VerificationSecrets        map[string]string
	VerificationKeysConfigured bool
	MaxAge                     time.Duration
	ClockSkew                  time.Duration
	Now                        func() time.Time
}

func serviceAssertionMetadataKeys() []string {
	return []string{
		ServiceAssertionCallerMetadataKey,
		ServiceAssertionIssuedAtMetadataKey,
		ServiceAssertionVersionMetadataKey,
	}
}

// HasPrivilegedMetadata reports whether the canonical RPC envelope contains
// metadata that may only be promoted after service assertion verification.
func HasPrivilegedMetadata(md metadata.MD) bool {
	return qhproMetadataContract.hasPrivilegedMetadata(md)
}

// HasServiceAssertion reports whether any canonical assertion field is present.
func HasServiceAssertion(md metadata.MD) bool {
	return len(md.Get(ServiceAssertionCallerMetadataKey)) > 0 ||
		len(md.Get(ServiceAssertionIssuedAtMetadataKey)) > 0 ||
		len(md.Get(ServiceAssertionSignatureMetadataKey)) > 0 ||
		len(md.Get(ServiceAssertionVersionMetadataKey)) > 0
}

// AllowsUnsignedServiceCall reports whether the canonical RPC contract
// explicitly permits this method to arrive without a service assertion.
func AllowsUnsignedServiceCall(method string) bool {
	return qhproMetadataContract.allowsUnsigned(method)
}

// SignServiceAssertion adds a method-bound HMAC assertion to outgoing metadata.
func SignServiceAssertion(
	method string,
	md metadata.MD,
	cfg ServiceAssertionConfig,
) error {
	if strings.TrimSpace(cfg.ServiceID) == "" ||
		strings.TrimSpace(cfg.Secret) == "" {
		return ErrServiceAssertionNotConfigured
	}

	if !(identity.ServiceCaller{ServiceID: cfg.ServiceID}).IsValid() {
		return fmt.Errorf(
			"%w: invalid service ID",
			ErrServiceAssertionInvalid,
		)
	}

	if strings.TrimSpace(method) == "" || md == nil {
		return fmt.Errorf(
			"%w: method or metadata is empty",
			ErrServiceAssertionInvalid,
		)
	}

	now := serviceAssertionNow(cfg).UTC()
	issuedAt := strconv.FormatInt(now.Unix(), 10)

	md.Set(ServiceAssertionVersionMetadataKey, ServiceAssertionVersion)
	md.Set(ServiceAssertionCallerMetadataKey, cfg.ServiceID)
	md.Set(ServiceAssertionIssuedAtMetadataKey, issuedAt)
	md.Delete(ServiceAssertionSignatureMetadataKey)

	signature := calculateServiceAssertionSignature(
		method,
		md,
		cfg.Secret,
		qhproMetadataContract,
	)
	md.Set(ServiceAssertionSignatureMetadataKey, signature)

	return nil
}

// VerifyServiceAssertion validates a method-bound HMAC assertion and returns
// the authenticated immediate service caller.
func VerifyServiceAssertion(
	method string,
	md metadata.MD,
	cfg ServiceAssertionConfig,
) (identity.ServiceCaller, error) {
	version, err := singleServiceAssertionValue(
		md,
		ServiceAssertionVersionMetadataKey,
	)
	if err != nil || version != ServiceAssertionVersion {
		return identity.ServiceCaller{}, ErrServiceAssertionMissing
	}

	caller, err := singleServiceAssertionValue(
		md,
		ServiceAssertionCallerMetadataKey,
	)
	if err != nil ||
		!(identity.ServiceCaller{ServiceID: caller}).IsValid() {
		return identity.ServiceCaller{}, ErrServiceAssertionInvalid
	}

	verificationSecret, err := serviceAssertionSecretForCaller(
		cfg,
		caller,
	)
	if err != nil {
		return identity.ServiceCaller{}, err
	}

	issuedAtRaw, err := singleServiceAssertionValue(
		md,
		ServiceAssertionIssuedAtMetadataKey,
	)
	if err != nil {
		return identity.ServiceCaller{}, ErrServiceAssertionInvalid
	}

	signature, err := singleServiceAssertionValue(
		md,
		ServiceAssertionSignatureMetadataKey,
	)
	if err != nil || len(signature) != sha256.Size*2 {
		return identity.ServiceCaller{}, ErrServiceAssertionInvalid
	}

	issuedUnix, err := strconv.ParseInt(issuedAtRaw, 10, 64)
	if err != nil {
		return identity.ServiceCaller{}, ErrServiceAssertionInvalid
	}

	issuedAt := time.Unix(issuedUnix, 0).UTC()
	now := serviceAssertionNow(cfg).UTC()

	maxAge := cfg.MaxAge
	if maxAge <= 0 {
		maxAge = 30 * time.Second
	}

	clockSkew := cfg.ClockSkew
	if clockSkew < 0 {
		clockSkew = 0
	}

	if issuedAt.After(now.Add(clockSkew)) ||
		now.Sub(issuedAt) > maxAge+clockSkew {
		return identity.ServiceCaller{}, ErrServiceAssertionExpired
	}

	copyMD := md.Copy()
	copyMD.Delete(ServiceAssertionSignatureMetadataKey)

	expected := calculateServiceAssertionSignature(
		method,
		copyMD,
		verificationSecret,
		qhproMetadataContract,
	)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return identity.ServiceCaller{}, ErrServiceAssertionInvalid
	}

	return identity.ServiceCaller{ServiceID: caller}, nil
}

func calculateServiceAssertionSignature(
	method string,
	md metadata.MD,
	secret string,
	contract metadataContract,
) string {
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

func singleServiceAssertionValue(
	md metadata.MD,
	key string,
) (string, error) {
	values := md.Get(key)
	if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
		return "", ErrServiceAssertionInvalid
	}
	return values[0], nil
}

func serviceAssertionNow(cfg ServiceAssertionConfig) time.Time {
	if cfg.Now != nil {
		return cfg.Now()
	}
	return time.Now()
}

func serviceAssertionSecretForCaller(
	cfg ServiceAssertionConfig,
	caller string,
) (string, error) {
	if cfg.VerificationKeysConfigured ||
		len(cfg.VerificationSecrets) > 0 {
		secret := strings.TrimSpace(
			cfg.VerificationSecrets[caller],
		)
		if secret == "" {
			return "", ErrServiceAssertionInvalid
		}
		return secret, nil
	}

	secret := strings.TrimSpace(cfg.Secret)
	if secret == "" {
		return "", ErrServiceAssertionNotConfigured
	}
	return secret, nil
}
