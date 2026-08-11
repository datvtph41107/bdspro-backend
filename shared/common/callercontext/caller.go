package callercontext

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	CallerKindMetadataKey = "x-qhpro-caller-kind"
	APIKeyIDMetadataKey   = "x-qhpro-api-key-id"
	APIKeyAppMetadataKey  = "x-qhpro-api-key-app"
)

type Kind string

const (
	KindAnonymous       Kind = "anonymous"
	KindUser            Kind = "user"
	KindAPIKey          Kind = "api_key"
	KindInternalService Kind = "internal_service"
)

/**
 * Caller classifies the logical request source.
 */
type Caller struct {
	Kind      Kind
	APIKeyID  uint64
	APIKeyApp string
}

type callerContextKey struct{}

var (
	ErrInvalidCaller         = errors.New("caller classification is invalid")
	ErrCallerContextConflict = errors.New("caller classification conflicts with the existing context")
)

/**
 * IsValid reports whether the caller classification is internally consistent.
 */
func (c Caller) IsValid() bool {
	switch c.Kind {
	case KindAnonymous, KindInternalService:
		return c.APIKeyID == 0 && strings.TrimSpace(c.APIKeyApp) == ""
	case KindUser:
		return (c.APIKeyID == 0 && strings.TrimSpace(c.APIKeyApp) == "") ||
			(c.APIKeyID > 0 && validText(c.APIKeyApp))
	case KindAPIKey:
		return c.APIKeyID > 0 && validText(c.APIKeyApp)
	default:
		return false
	}
}

/**
 * BindCaller validates and immutably binds one logical caller.
 */
func BindCaller(ctx context.Context, caller Caller) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !caller.IsValid() {
		return ctx, ErrInvalidCaller
	}
	caller.APIKeyApp = strings.TrimSpace(caller.APIKeyApp)
	if existing, ok := FromContext(ctx); ok {
		if existing == caller {
			return ctx, nil
		}
		return ctx, ErrCallerContextConflict
	}
	return context.WithValue(ctx, callerContextKey{}, caller), nil
}

/**
 * FromContext returns the caller classification.
 */
func FromContext(ctx context.Context) (Caller, bool) {
	if ctx == nil {
		return Caller{}, false
	}
	caller, ok := ctx.Value(callerContextKey{}).(Caller)
	return caller, ok && caller.IsValid()
}

/**
 * AppendMetadata writes the caller classification to outgoing metadata.
 */
func AppendMetadata(md metadata.MD, caller Caller) {
	if md == nil || !caller.IsValid() {
		return
	}
	md.Set(CallerKindMetadataKey, string(caller.Kind))
	if caller.APIKeyID > 0 {
		md.Set(APIKeyIDMetadataKey, strconv.FormatUint(caller.APIKeyID, 10))
		md.Set(APIKeyAppMetadataKey, strings.TrimSpace(caller.APIKeyApp))
	} else {
		md.Delete(APIKeyIDMetadataKey)
		md.Delete(APIKeyAppMetadataKey)
	}
}

/**
 * FromMetadata parses one signed caller classification.
 */
func FromMetadata(md metadata.MD) (Caller, bool, error) {
	if md == nil {
		return Caller{}, false, nil
	}
	kindValue, hasKind, err := singleOptional(md, CallerKindMetadataKey)
	if err != nil {
		return Caller{}, false, err
	}
	idValue, hasID, err := singleOptional(md, APIKeyIDMetadataKey)
	if err != nil {
		return Caller{}, false, err
	}
	appValue, hasApp, err := singleOptional(md, APIKeyAppMetadataKey)
	if err != nil {
		return Caller{}, false, err
	}
	if !hasKind && !hasID && !hasApp {
		return Caller{}, false, nil
	}
	if !hasKind {
		return Caller{}, false, fmt.Errorf("caller kind is missing")
	}

	caller := Caller{Kind: Kind(kindValue)}
	if hasID != hasApp {
		return Caller{}, false, fmt.Errorf("api key caller metadata is incomplete")
	}
	if hasID {
		apiKeyID, parseErr := strconv.ParseUint(idValue, 10, 64)
		if parseErr != nil {
			return Caller{}, false, fmt.Errorf("invalid api key caller ID")
		}
		caller.APIKeyID = apiKeyID
		caller.APIKeyApp = appValue
	}

	if !caller.IsValid() {
		return Caller{}, false, fmt.Errorf("invalid caller classification")
	}
	return caller, true, nil
}

func singleOptional(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 {
		return "", false, fmt.Errorf("duplicate %s metadata", key)
	}
	raw := values[0]
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, fmt.Errorf("empty %s metadata", key)
	}
	if value != raw {
		return "", false, fmt.Errorf("invalid %s metadata", key)
	}
	return value, true, nil
}

func validText(value string) bool {
	if value == "" || strings.TrimSpace(value) != value || len(value) > 128 {
		return false
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}
