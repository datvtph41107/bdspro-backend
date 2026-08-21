package callercontext

import (
	"common/identity"
	"context"
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

// Deprecated: use identity.CallerKind.
type Kind = identity.CallerKind

const (
	KindAnonymous       = identity.CallerAnonymous
	KindUser            = identity.CallerUser
	KindAPIKey          = identity.CallerAPIKey
	KindInternalService = identity.CallerInternalService
)

// Deprecated: use identity.Caller. This alias remains while callers migrate.
type Caller = identity.Caller

var (
	ErrInvalidCaller         = identity.ErrInvalidCaller
	ErrCallerContextConflict = identity.ErrCallerContextConflict
)

// Deprecated: use identity.BindCaller.
func BindCaller(ctx context.Context, caller Caller) (context.Context, error) {
	return identity.BindCaller(ctx, caller)
}

// Deprecated: use identity.CallerFromContext.
func FromContext(ctx context.Context) (Caller, bool) {
	return identity.CallerFromContext(ctx)
}

// AppendMetadata is a transport compatibility helper. It remains here until
// identitygrpc/caller metadata is absorbed by the canonical rpc package.
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

// FromMetadata parses one signed caller classification. It remains a transport
// compatibility helper until the canonical rpc package owns gRPC metadata.
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
