package rpc

import (
	"common/identity"
	"common/request"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AppendRequestFromContext serializes canonical request-scoped identities into
// gRPC metadata. It never originates a new Request-ID, Operation-ID, or
// Idempotency-Key.
func AppendRequestFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.MD{}
	}
	if requestID, ok := request.RequestIDFromContext(ctx); ok {
		md.Set(RequestIDMetadataKey, requestID)
	}
	if operationID, ok := request.OperationIDFromContext(ctx); ok {
		md.Set(OperationIDMetadataKey, operationID)
	}
	if key, ok := request.IdempotencyKeyFromContext(ctx); ok {
		md.Set(IdempotencyKeyMetadataKey, key)
	}
	return md
}

// RestoreRequestFromMetadata restores canonical request identities. Request-ID
// remains trace metadata and malformed values are ignored; Operation-ID and
// Idempotency-Key are strict request contracts.
func RestoreRequestFromMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	if ctx == nil {
		return nil, status.Error(codes.Internal, "context is nil")
	}
	if md == nil {
		return ctx, nil
	}

	ctx = restoreRequestID(ctx, md.Get(RequestIDMetadataKey))

	var err error
	ctx, err = restoreOperationID(ctx, md.Get(OperationIDMetadataKey))
	if err != nil {
		return ctx, err
	}
	ctx, err = restoreIdempotencyKey(ctx, md.Get(IdempotencyKeyMetadataKey))
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func restoreRequestID(ctx context.Context, values []string) context.Context {
	if len(values) != 1 {
		return ctx
	}
	requestID := request.NormalizeRequestID(values[0])
	if !request.IsValidRequestID(requestID) {
		return ctx
	}
	return request.WithRequestID(ctx, requestID)
}

func restoreOperationID(ctx context.Context, values []string) (context.Context, error) {
	if len(values) == 0 {
		return ctx, nil
	}
	operationID, ok := request.AcceptOperationID(values)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "invalid operation-id metadata")
	}
	bound, err := request.BindOperationID(ctx, operationID)
	if err != nil {
		if errors.Is(err, request.ErrOperationIDContextConflict) {
			return ctx, status.Error(codes.Internal, "operation-id context conflict")
		}
		return ctx, status.Error(codes.InvalidArgument, "invalid operation-id metadata")
	}
	return bound, nil
}

func restoreIdempotencyKey(ctx context.Context, values []string) (context.Context, error) {
	key, ok, err := request.ParseIdempotencyKey(values)
	if err != nil {
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency-key metadata")
	}
	if !ok {
		return ctx, nil
	}
	bound, err := request.BindIdempotencyKey(ctx, key)
	if err != nil {
		if errors.Is(err, request.ErrIdempotencyKeyContextConflict) {
			return ctx, status.Error(codes.Internal, "idempotency-key context conflict")
		}
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency-key metadata")
	}
	return bound, nil
}

// AppendIdentityFromContext serializes the canonical Actor and logical Caller.
// Local service-originated calls receive internal_service when no logical
// caller exists. If an Actor exists without an explicit Caller, the compatible
// logical caller is user.
func AppendIdentityFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.MD{}
	}

	actor, hasActor := identity.ActorFromContext(ctx)
	if hasActor {
		appendActor(md, actor)
	}

	caller, hasCaller := identity.CallerFromContext(ctx)
	if !hasCaller {
		if hasActor {
			caller = identity.Caller{Kind: identity.CallerUser}
		} else {
			caller = identity.Caller{Kind: identity.CallerInternalService}
		}
	}
	appendCaller(md, caller)
	return md
}

// RestoreIdentityFromMetadata validates the transported Actor/Caller pair. If
// promote is false, syntax and consistency are still checked but privileged
// identity is not bound into canonical context.
func RestoreIdentityFromMetadata(
	ctx context.Context,
	md metadata.MD,
	promote bool,
) (context.Context, error) {
	return restoreIdentityFromMetadata(ctx, md, promote, false)
}

func restoreIdentityFromMetadata(
	ctx context.Context,
	md metadata.MD,
	promote bool,
	inferCaller bool,
) (context.Context, error) {
	if ctx == nil {
		return nil, status.Error(codes.Internal, "context is nil")
	}

	actor, hasActor, err := actorFromMetadata(md)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}
	caller, hasCaller, err := callerFromMetadata(md)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	if promote && inferCaller && !hasCaller {
		if hasActor {
			caller = identity.Caller{Kind: identity.CallerUser}
		} else {
			caller = identity.Caller{Kind: identity.CallerInternalService}
		}
		hasCaller = true
	}

	if err := validateCallerActor(caller, hasCaller, hasActor); err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}
	if !promote {
		return ctx, nil
	}

	if hasActor {
		bound, bindErr := identity.BindActor(ctx, actor)
		if bindErr != nil {
			if errors.Is(bindErr, identity.ErrActorContextConflict) {
				return ctx, status.Error(codes.Internal, "actor context conflict")
			}
			return ctx, status.Error(codes.Unauthenticated, "invalid actor")
		}
		ctx = bound
	}
	if hasCaller {
		bound, bindErr := identity.BindCaller(ctx, caller)
		if bindErr != nil {
			if errors.Is(bindErr, identity.ErrCallerContextConflict) {
				return ctx, status.Error(codes.Internal, "caller context conflict")
			}
			return ctx, status.Error(codes.Unauthenticated, "invalid caller")
		}
		ctx = bound
	}
	return ctx, nil
}

// AppendActorMetadata serializes one validated Actor into canonical QHPRO
// gRPC metadata. It does not bind or originate identity.
func AppendActorMetadata(md metadata.MD, actor identity.Actor) {
	appendActor(md, actor)
}

// ActorFromMetadata parses one canonical Actor from QHPRO gRPC metadata
// without binding it into context.
func ActorFromMetadata(md metadata.MD) (identity.Actor, bool, error) {
	return actorFromMetadata(md)
}

func appendActor(md metadata.MD, actor identity.Actor) {
	setUint64(md, AuthIDMetadataKey, actor.AuthID)
	setUint64(md, ProfileIDMetadataKey, actor.ProfileID)
	setUint64(md, OriginIDMetadataKey, actor.OriginID)
	setUint64(md, SessionIDMetadataKey, actor.SessionID)
	if actor.OrganizationID != nil {
		setUint64(md, OrganizationIDMetadataKey, *actor.OrganizationID)
	}
	setText(md, RoleMetadataKey, actor.Role)
	setText(md, TokenTypeMetadataKey, actor.TokenType)
}

func actorFromMetadata(md metadata.MD) (identity.Actor, bool, error) {
	if md == nil {
		return identity.Actor{}, false, nil
	}
	actor := identity.Actor{}
	hasValue := false
	values := []struct {
		key string
		set func(uint64)
	}{
		{AuthIDMetadataKey, func(v uint64) { actor.AuthID = v }},
		{ProfileIDMetadataKey, func(v uint64) { actor.ProfileID = v }},
		{OriginIDMetadataKey, func(v uint64) { actor.OriginID = v }},
		{SessionIDMetadataKey, func(v uint64) { actor.SessionID = v }},
	}
	for _, item := range values {
		value, found, err := readActorUint64(md, item.key)
		if err != nil {
			return identity.Actor{}, false, err
		}
		if found {
			hasValue = true
			item.set(value)
		}
	}
	organizationID, found, err := readActorUint64(md, OrganizationIDMetadataKey)
	if err != nil {
		return identity.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.OrganizationID = &organizationID
	}
	role, found, err := readActorText(md, RoleMetadataKey)
	if err != nil {
		return identity.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.Role = role
	}
	tokenType, found, err := readActorText(md, TokenTypeMetadataKey)
	if err != nil {
		return identity.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.TokenType = tokenType
	}
	if !hasValue {
		return identity.Actor{}, false, nil
	}
	if !actor.IsValid() {
		return identity.Actor{}, false, errors.New("invalid actor identity metadata")
	}
	return actor, true, nil
}

// AppendCallerMetadata serializes one validated logical caller into canonical
// QHPRO gRPC metadata. It does not bind or originate caller identity.
func AppendCallerMetadata(md metadata.MD, caller identity.Caller) {
	appendCaller(md, caller)
}

// CallerFromMetadata parses one canonical logical caller from QHPRO gRPC
// metadata without binding it into context.
func CallerFromMetadata(md metadata.MD) (identity.Caller, bool, error) {
	return callerFromMetadata(md)
}

func appendCaller(md metadata.MD, caller identity.Caller) {
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

func callerFromMetadata(md metadata.MD) (identity.Caller, bool, error) {
	if md == nil {
		return identity.Caller{}, false, nil
	}
	kindValue, hasKind, err := readText(md, CallerKindMetadataKey)
	if err != nil {
		return identity.Caller{}, false, err
	}
	idValue, hasID, err := readText(md, APIKeyIDMetadataKey)
	if err != nil {
		return identity.Caller{}, false, err
	}
	appValue, hasApp, err := readText(md, APIKeyAppMetadataKey)
	if err != nil {
		return identity.Caller{}, false, err
	}
	if !hasKind && !hasID && !hasApp {
		return identity.Caller{}, false, nil
	}
	if !hasKind {
		return identity.Caller{}, false, errors.New("caller kind is missing")
	}
	caller := identity.Caller{Kind: identity.CallerKind(kindValue)}
	if hasID != hasApp {
		return identity.Caller{}, false, errors.New("api key caller metadata is incomplete")
	}
	if hasID {
		apiKeyID, parseErr := strconv.ParseUint(idValue, 10, 64)
		if parseErr != nil || apiKeyID == 0 {
			return identity.Caller{}, false, errors.New("invalid api key caller ID")
		}
		caller.APIKeyID = apiKeyID
		caller.APIKeyApp = appValue
	}
	if !caller.IsValid() {
		return identity.Caller{}, false, errors.New("invalid caller classification")
	}
	return caller, true, nil
}

func validateCallerActor(caller identity.Caller, hasCaller, hasActor bool) error {
	if !hasCaller {
		if hasActor {
			return errors.New("actor requires caller metadata")
		}
		return nil
	}
	switch caller.Kind {
	case identity.CallerUser:
		if !hasActor {
			return errors.New("user caller requires actor")
		}
	case identity.CallerAnonymous, identity.CallerAPIKey, identity.CallerInternalService:
		if hasActor {
			return errors.New("caller type does not allow actor")
		}
	default:
		return errors.New("caller type is invalid")
	}
	return nil
}

func readActorUint64(md metadata.MD, key string) (uint64, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return 0, false, nil
	}
	if len(values) != 1 ||
		values[0] == "" ||
		strings.TrimSpace(values[0]) != values[0] {
		return 0, false, fmt.Errorf("invalid %s identity metadata", key)
	}

	value, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s identity metadata", key)
	}

	return value, true, nil
}

func readActorText(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 ||
		values[0] == "" ||
		len(values[0]) > identity.MaxActorTextLength ||
		strings.TrimSpace(values[0]) != values[0] {
		return "", false, fmt.Errorf("invalid %s identity metadata", key)
	}

	for _, value := range values[0] {
		if value < 0x20 || value == 0x7f {
			return "", false, fmt.Errorf("invalid %s identity metadata", key)
		}
	}

	return values[0], true, nil
}

func readUint64(md metadata.MD, key string) (uint64, bool, error) {
	text, found, err := readText(md, key)
	if err != nil || !found {
		return 0, found, err
	}
	value, parseErr := strconv.ParseUint(text, 10, 64)
	if parseErr != nil || value == 0 {
		return 0, false, fmt.Errorf("%s metadata is invalid", key)
	}
	return value, true, nil
}

func readText(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 {
		return "", false, fmt.Errorf("duplicate %s metadata", key)
	}
	value := strings.TrimSpace(values[0])
	if value == "" || value != values[0] || len(value) > 128 {
		return "", false, fmt.Errorf("%s metadata is invalid", key)
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return "", false, fmt.Errorf("%s metadata is invalid", key)
		}
	}
	return value, true, nil
}

func setUint64(md metadata.MD, key string, value uint64) {
	if value > 0 {
		md.Set(key, strconv.FormatUint(value, 10))
	}
}

func setText(md metadata.MD, key, value string) {
	if value = strings.TrimSpace(value); value != "" {
		md.Set(key, value)
	}
}
