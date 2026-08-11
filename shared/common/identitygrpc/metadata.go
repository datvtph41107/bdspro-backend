package identitygrpc

import (
	"common/callercontext"
	"common/identitycontext"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/**
 * AppendFromContext đưa Actor và Caller đã được xác thực vào metadata nội bộ.
 */
func AppendFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.MD{}
	}

	if actor, ok := identitycontext.FromContext(ctx); ok {
		appendActor(md, actor)
	}
	if caller, ok := callercontext.FromContext(ctx); ok {
		callercontext.AppendMetadata(md, caller)
	}

	return md
}

/**
 * RestoreFromMetadata bind Actor và Caller sau khi transport signature đã được verify.
 */
func RestoreFromMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	actor, hasActor, err := actorFromMetadata(md)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	caller, hasCaller, err := callercontext.FromMetadata(md)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}
	if err := validateCaller(caller, hasCaller, hasActor); err != nil {
		return ctx, status.Error(codes.Unauthenticated, err.Error())
	}

	if hasActor {
		bound, bindErr := identitycontext.BindActor(ctx, actor)
		if bindErr != nil {
			if errors.Is(bindErr, identitycontext.ErrActorContextConflict) {
				return ctx, status.Error(codes.Internal, "actor context conflict")
			}
			return ctx, status.Error(codes.Unauthenticated, "invalid actor")
		}
		ctx = bound
	}

	if hasCaller {
		bound, bindErr := callercontext.BindCaller(ctx, caller)
		if bindErr != nil {
			if errors.Is(bindErr, callercontext.ErrCallerContextConflict) {
				return ctx, status.Error(codes.Internal, "caller context conflict")
			}
			return ctx, status.Error(codes.Unauthenticated, "invalid caller")
		}
		ctx = bound
	}

	return ctx, nil
}

func appendActor(md metadata.MD, actor identitycontext.Actor) {
	setUint64(md, identitycontext.AuthIDMetadataKey, actor.AuthID)
	setUint64(md, identitycontext.ProfileIDMetadataKey, actor.ProfileID)
	setUint64(md, identitycontext.OriginIDMetadataKey, actor.OriginID)
	setUint64(md, identitycontext.SessionIDMetadataKey, actor.SessionID)
	if actor.OrganizationID != nil {
		setUint64(md, identitycontext.OrganizationIDMetadataKey, *actor.OrganizationID)
	}
	setText(md, identitycontext.RoleMetadataKey, actor.Role)
	setText(md, identitycontext.TokenTypeMetadataKey, actor.TokenType)
}

func actorFromMetadata(md metadata.MD) (identitycontext.Actor, bool, error) {
	if md == nil {
		return identitycontext.Actor{}, false, nil
	}

	actor := identitycontext.Actor{}
	hasValue := false

	values := []struct {
		key string
		set func(uint64)
	}{
		{identitycontext.AuthIDMetadataKey, func(v uint64) { actor.AuthID = v }},
		{identitycontext.ProfileIDMetadataKey, func(v uint64) { actor.ProfileID = v }},
		{identitycontext.OriginIDMetadataKey, func(v uint64) { actor.OriginID = v }},
		{identitycontext.SessionIDMetadataKey, func(v uint64) { actor.SessionID = v }},
	}

	for _, item := range values {
		value, found, err := readUint64(md, item.key)
		if err != nil {
			return identitycontext.Actor{}, false, err
		}
		if found {
			hasValue = true
			item.set(value)
		}
	}

	organizationID, found, err := readUint64(md, identitycontext.OrganizationIDMetadataKey)
	if err != nil {
		return identitycontext.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.OrganizationID = &organizationID
	}

	role, found, err := readText(md, identitycontext.RoleMetadataKey)
	if err != nil {
		return identitycontext.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.Role = role
	}

	tokenType, found, err := readText(md, identitycontext.TokenTypeMetadataKey)
	if err != nil {
		return identitycontext.Actor{}, false, err
	}
	if found {
		hasValue = true
		actor.TokenType = tokenType
	}

	if !hasValue {
		return identitycontext.Actor{}, false, nil
	}
	if !actor.IsValid() {
		return identitycontext.Actor{}, false, errors.New("actor metadata is invalid")
	}
	return actor, true, nil
}

func validateCaller(caller callercontext.Caller, hasCaller, hasActor bool) error {
	if !hasCaller {
		if hasActor {
			return errors.New("actor requires caller metadata")
		}
		return nil
	}

	switch caller.Kind {
	case callercontext.KindUser:
		if !hasActor {
			return errors.New("user caller requires actor")
		}
	case callercontext.KindAnonymous, callercontext.KindAPIKey, callercontext.KindInternalService:
		if hasActor {
			return errors.New("caller type does not allow actor")
		}
	default:
		return errors.New("caller type is invalid")
	}
	return nil
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
	if value == "" || value != values[0] {
		return "", false, fmt.Errorf("%s metadata is invalid", key)
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
