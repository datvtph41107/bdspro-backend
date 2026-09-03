package _middleware

import (
	_enums "common/domain/enum"
	_identity "common/identity"
	_jwt "common/jwt"
	_request "common/request"
	_rpc "common/rpc"
	_rpcenv "common/rpcenv"

	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func InjectGrpcMetadataContextMiddleware(_ context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}
	if r == nil {
		return md
	}

	requestCtx := r.Context()
	if requestID, ok := _request.RequestIDFromContext(requestCtx); ok {
		md.Set(_request.RequestIDMetadataKey, requestID)
	}
	if operationID, ok := _request.OperationIDFromContext(requestCtx); ok {
		md.Set(_request.OperationIDMetadataKey, operationID)
	}
	if idempotencyKey, ok := _request.IdempotencyKeyFromContext(requestCtx); ok {
		md.Set(_request.IdempotencyKeyMetadataKey, idempotencyKey)
	}
	// Luôn inject IP và User-Agent vào metadata (dù có token hay không)
	appendIfExists(md, "client-ip", getClientIP(r))
	appendIfExists(md, "user-agent", r.Header.Get("User-Agent"))
	if apiKey, ok, err := _jwt.APIKeyFromRequest(r); err == nil && ok {
		appendIfExists(md, "x-api-key", apiKey)
	}
	if deviceID := r.Header.Get("device-id"); deviceID != "" {
		appendIfExists(md, "device-id", deviceID)
	} else {
		appendIfExists(md, "device-id", r.Header.Get("device-id"))
	}
	if caller, ok := _identity.CallerFromContext(r.Context()); ok {
		_rpc.AppendCallerMetadata(md, caller)
	}
	appendIfExists(md, "x-qhpro-telemetry-key", r.Header.Get("X-QHPro-Telemetry-Key"))
	appendIfExists(md, "if-none-match", r.Header.Get("If-None-Match"))
	appendIfExists(md, "method", r.Method)

	claims, ok := principalFromHTTPRequest(r)
	if !ok {
		return md
	}

	_rpc.AppendActorMetadata(md, _jwt.ActorFromPrincipal(claims))
	if claims.PlanId != nil {
		appendIfExists(md, "planId", *claims.PlanId)
	}
	if claims.PlanFrom != nil {
		appendIfExists(md, "planFrom", *claims.PlanFrom)
	}

	return md
}

// getClientIP lấy IP address thực của client từ request
func getClientIP(r *http.Request) string {
	// Ưu tiên lấy từ X-Forwarded-For (khi đi qua proxy/load balancer)
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For có thể chứa nhiều IP, lấy IP đầu tiên
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Nếu không có X-Forwarded-For, thử lấy từ X-Real-IP
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// Cuối cùng lấy từ RemoteAddr
	// RemoteAddr có format "IP:Port", cần tách IP ra
	if remoteAddr := r.RemoteAddr; remoteAddr != "" {
		// Tách IP khỏi port
		if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
			return remoteAddr[:idx]
		}
		return remoteAddr
	}

	return ""
}

func ParseGrpcMetadataContextMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return parseGrpcMetadataContextMiddlewareWithTrust(
		_rpcenv.LoadTransportConfig(),
		ctx, req, info, handler,
	)
}

func parseGrpcMetadataContextMiddlewareWithTrust(
	transportConfig _rpc.TransportConfig,
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		md = metadata.MD{}
	}

	fullMethod := ""
	if info != nil {
		fullMethod = info.FullMethod
	}
	transportCaller, verifyErr := _rpc.VerifyServiceAssertion(
		fullMethod,
		md,
		transportConfig.ServiceAssertion,
	)

	verified := verifyErr == nil
	hasPrivilegedMetadata := _rpc.HasPrivilegedMetadata(md)
	hasServiceAssertion := _rpc.HasServiceAssertion(md)
	unsignedAllowed := _rpc.AllowsUnsignedServiceCall(fullMethod) &&
		!hasPrivilegedMetadata &&
		!hasServiceAssertion

	switch {
	case verified:
		bound, bindErr := _identity.BindServiceCaller(ctx, transportCaller)
		if bindErr != nil {
			return nil, status.Error(codes.Internal, "transport caller context conflict")
		}
		ctx = bound

	case unsignedAllowed:
		// Explicit infrastructure exception. No trusted caller is promoted.

	case transportConfig.RequireServiceAssertion:
		return nil, status.Error(codes.Unauthenticated, "trusted transport assertion required")

	case !transportConfig.RequireServiceAssertion &&
		(hasPrivilegedMetadata || hasServiceAssertion):
		ctx = markInboundMetadataUnverified(ctx)
		log.Printf(
			"[gRPC Server] Unverified transport metadata on %s: %v",
			fullMethod,
			verifyErr,
		)
	}

	actor, hasActor, err := _rpc.ActorFromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	requestCaller, hasRequestCaller, err := _rpc.CallerFromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !hasRequestCaller && verified {
		requestCaller = _identity.Caller{Kind: _identity.CallerInternalService}
		if hasActor {
			requestCaller.Kind = _identity.CallerUser
		}
		hasRequestCaller = true
	}
	if err := validateCallerActor(requestCaller, hasRequestCaller, hasActor); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if hasRequestCaller && verified {
		bound, bindErr := _identity.BindCaller(ctx, requestCaller)
		if bindErr != nil {
			if errors.Is(bindErr, _identity.ErrCallerContextConflict) {
				return nil, status.Error(codes.Internal, "caller context conflict")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid caller classification")
		}
		ctx = bound
	}

	if hasActor {
		if verified {
			bound, bindErr := _identity.BindActor(ctx, actor)
			if bindErr != nil {
				if errors.Is(bindErr, _identity.ErrActorContextConflict) {
					return nil, status.Error(codes.Internal, "actor identity context conflict")
				}
				return nil, status.Error(codes.Unauthenticated, "invalid actor identity")
			}
			ctx = bound
		}
		ctx = withLegacyActorContext(ctx, actor)
	}

	ctx = parseUint64Metadata(ctx, md, "planId", _enums.PlanIDKey)
	ctx = parseTimeMetadata(ctx, md, "planFrom", _enums.PlanFromKey)
	ctx = parseStringMetadata(ctx, md, "x-api-key", _enums.APIKeyKey)

	ctx = acceptIncomingRequestID(
		ctx,
		md.Get(_request.RequestIDMetadataKey),
	)

	ctx, err = bindIncomingOperationID(
		ctx,
		md.Get(_request.OperationIDMetadataKey),
	)
	if err != nil {
		return nil, err
	}

	ctx, err = bindIncomingIdempotencyKey(
		ctx,
		md.Get(_request.IdempotencyKeyMetadataKey),
	)
	if err != nil {
		return nil, err
	}

	// Parse IP address và User-Agent từ metadata
	ctx = parseStringMetadata(ctx, md, "client-ip", _enums.ClientIPKey)
	ctx = parseStringMetadata(ctx, md, "user-agent", _enums.UserAgentKey)
	ctx = parseStringMetadata(ctx, md, "device-id", _enums.DeviceIDKey)

	// Parse writer từ context nếu có
	if writer, ok := ctx.Value("response-writer").(http.ResponseWriter); ok {
		// Sử dụng appcontext để lưu writer vào context
		ctx = context.WithValue(ctx, "response-writer", writer)
	}

	return handler(ctx, req)
}

func acceptIncomingRequestID(ctx context.Context, values []string) context.Context {
	if len(values) != 1 {
		return ctx
	}

	requestID := _request.NormalizeRequestID(values[0])
	if !_request.IsValidRequestID(requestID) {
		return ctx
	}
	return _request.WithRequestID(ctx, requestID)
}

func bindIncomingOperationID(ctx context.Context, values []string) (context.Context, error) {
	if len(values) == 0 {
		return ctx, nil
	}

	operationID, ok := _request.AcceptOperationID(values)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "invalid operation identity metadata")
	}

	bound, err := _request.BindOperationID(ctx, operationID)
	if err != nil {
		return ctx, status.Error(codes.Internal, "operation identity context conflict")
	}
	return bound, nil
}

func bindIncomingIdempotencyKey(ctx context.Context, values []string) (context.Context, error) {
	key, ok, err := _request.ParseIdempotencyKey(values)
	if err != nil {
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency identity metadata")
	}
	if !ok {
		return ctx, nil
	}

	bound, err := _request.BindIdempotencyKey(ctx, key)
	if err != nil {
		return ctx, status.Error(codes.Internal, "idempotency identity context conflict")
	}
	return bound, nil
}

type inboundMetadataUnverifiedKey struct{}

func markInboundMetadataUnverified(ctx context.Context) context.Context {
	return context.WithValue(ctx, inboundMetadataUnverifiedKey{}, true)
}

func canForwardPrivilegedMetadata(ctx context.Context) bool {
	unverified, _ := ctx.Value(inboundMetadataUnverifiedKey{}).(bool)
	return !unverified
}

func GrpcClientMetadataFromContext(ctx context.Context) metadata.MD {
	md := metadata.MD{}

	appendUint64IfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(uint64); ok {
			md.Append(mdKey, strconv.FormatUint(val, 10))
		}
	}

	appendStringIfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(string); ok && val != "" {
			md.Append(mdKey, val)
		}
	}

	appendTimeIfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(time.Time); ok && !val.IsZero() {
			md.Append(mdKey, val.Format(time.RFC3339))
		}
	}
	forwardPrivilegedMetadata := canForwardPrivilegedMetadata(ctx)
	if actor, ok := _identity.ActorFromContext(ctx); ok {
		_rpc.AppendActorMetadata(md, actor)
	} else if forwardPrivilegedMetadata {
		appendUint64IfExists(_enums.ProfileIDKey, _rpc.ProfileIDMetadataKey)
		appendUint64IfExists(_enums.OriginIDKey, _rpc.OriginIDMetadataKey)
		appendUint64IfExists(_enums.OrganizationIDKey, _rpc.OrganizationIDMetadataKey)
		appendUint64IfExists(_enums.AuthIDKey, _rpc.AuthIDMetadataKey)
		appendUint64IfExists(_enums.SessionKey, _rpc.SessionIDMetadataKey)
		appendStringIfExists(_enums.RoleKey, _rpc.RoleMetadataKey)
		appendStringIfExists(_enums.TypeKey, _rpc.TokenTypeMetadataKey)
	}

	if caller, ok := _identity.CallerFromContext(ctx); ok {
		_rpc.AppendCallerMetadata(md, caller)
	} else if outgoingMD, ok := metadata.FromOutgoingContext(ctx); !ok || len(outgoingMD.Get(_rpc.CallerKindMetadataKey)) == 0 {
		_rpc.AppendCallerMetadata(md, _identity.Caller{Kind: _identity.CallerInternalService})
	}

	if forwardPrivilegedMetadata {
		appendUint64IfExists(_enums.PlanIDKey, "planId")
		appendStringIfExists(_enums.APIKeyKey, "x-api-key")
	}
	appendStringIfExists(_enums.DeviceIDKey, "device-id")

	if requestID, ok := _request.RequestIDFromContext(ctx); ok {
		md.Append(_request.RequestIDMetadataKey, requestID)
	}
	if operationID, ok := _request.OperationIDFromContext(ctx); ok {
		md.Append(_request.OperationIDMetadataKey, operationID)
	}
	if idempotencyKey, ok := _request.IdempotencyKeyFromContext(ctx); ok {
		md.Append(_request.IdempotencyKeyMetadataKey, idempotencyKey)
	}

	if forwardPrivilegedMetadata {
		appendTimeIfExists(_enums.PlanFromKey, "planFrom")
	}

	return md
}

func validateCallerActor(caller _identity.Caller, hasCaller bool, hasActor bool) error {
	if !hasCaller {
		return nil
	}
	switch caller.Kind {
	case _identity.CallerUser:
		if !hasActor {
			return errors.New("user caller requires actor identity")
		}
	case _identity.CallerAnonymous, _identity.CallerAPIKey, _identity.CallerInternalService:
		if hasActor {
			return errors.New("caller kind does not allow actor identity")
		}
	default:
		return errors.New("invalid caller kind")
	}
	return nil
}

/**
 * PrincipalFromContext projects the canonical actor for legacy callers.
 */
func PrincipalFromContext(ctx context.Context) (*_jwt.Principal, error) {
	actor, ok := _identity.ActorFromContext(ctx)
	if !ok {
		return nil, errors.New("actor identity not found in context")
	}

	principal := &_jwt.Principal{
		AuthID:         actor.AuthID,
		ProfileId:      actor.ProfileID,
		OriginId:       actor.OriginID,
		OrganizationId: cloneUint64Pointer(actor.OrganizationID),
		Session:        actor.SessionID,
		Role:           actor.Role,
		Type:           actor.TokenType,
	}

	if planID, ok := ctx.Value(_enums.PlanIDKey).(uint64); ok {
		principal.PlanId = &planID
	}
	if planFrom, ok := ctx.Value(_enums.PlanFromKey).(time.Time); ok && !planFrom.IsZero() {
		principal.PlanFrom = jwt.NewNumericDate(planFrom)
	}

	return principal, nil
}

/**
 * principalFromHTTPRequest reuses the principal verified by HTTP auth.
 */
func principalFromHTTPRequest(r *http.Request) (*_jwt.Principal, bool) {
	if r == nil {
		return nil, false
	}
	return _jwt.PrincipalFromRequestContext(r.Context())
}

func withLegacyActorContext(ctx context.Context, actor _identity.Actor) context.Context {
	if actor.ProfileID != 0 {
		ctx = context.WithValue(ctx, _enums.ProfileIDKey, actor.ProfileID)
	}
	if actor.OriginID != 0 {
		ctx = context.WithValue(ctx, _enums.OriginIDKey, actor.OriginID)
	}
	if actor.OrganizationID != nil {
		ctx = context.WithValue(ctx, _enums.OrganizationIDKey, *actor.OrganizationID)
	}
	if actor.AuthID != 0 {
		ctx = context.WithValue(ctx, _enums.AuthIDKey, actor.AuthID)
	}
	if actor.SessionID != 0 {
		ctx = context.WithValue(ctx, _enums.SessionKey, actor.SessionID)
	}
	if actor.Role != "" {
		ctx = context.WithValue(ctx, _enums.RoleKey, actor.Role)
	}
	if actor.TokenType != "" {
		ctx = context.WithValue(ctx, _enums.TypeKey, actor.TokenType)
	}
	return ctx
}

func cloneUint64Pointer(value *uint64) *uint64 {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func parseUint64Metadata(ctx context.Context, md metadata.MD, key string, ctxKey _enums.ContextKey) context.Context {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ctx
	}
	if val, err := strconv.ParseUint(vals[0], 10, 64); err == nil {
		return context.WithValue(ctx, ctxKey, val)
	}
	return ctx
}

func parseStringMetadata(ctx context.Context, md metadata.MD, key string, ctxKey _enums.ContextKey) context.Context {
	vals := md.Get(key)
	if len(vals) > 0 && vals[0] != "" {
		return context.WithValue(ctx, ctxKey, vals[0])
	}
	return ctx
}

func parseTimeMetadata(ctx context.Context, md metadata.MD, key string, ctxKey _enums.ContextKey) context.Context {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ctx
	}
	if t, err := time.Parse(time.RFC3339, vals[0]); err == nil {
		return context.WithValue(ctx, ctxKey, t)
	}
	return ctx
}

func appendIfExists(md metadata.MD, key string, value any) {
	if value == nil {
		return
	}

	// v := reflect.ValueOf(value)

	// Nếu là con trỏ, kiểm tra nil và dereference
	// if v.Kind() == reflect.Ptr {
	// 	if v.IsNil() {
	// 		return
	// 	}
	// 	v = v.Elem()
	// 	value = v.Interface()
	// }

	switch v := value.(type) {
	case string:
		if v != "" {
			md.Append(key, v)
		}
	case *string:
		if v != nil && *v != "" {
			md.Append(key, *v)
		}
	case int, int64, uint64, float64:
		md.Append(key, fmt.Sprintf("%v", v))
	case *int, *int64, *uint64, *float64:
		if v != nil {
			md.Append(key, fmt.Sprintf("%v", v))
		}
	case *time.Time:
		if v != nil && !v.IsZero() {
			md.Append(key, v.Format(time.RFC3339))
		}
	case jwt.NumericDate:
		if !v.Time.IsZero() {
			md.Append(key, v.Time.Format(time.RFC3339))
		}
	case *jwt.NumericDate:
		if v != nil && !v.Time.IsZero() {
			md.Append(key, v.Time.Format(time.RFC3339))
		}
	}
}

/**
 * ParseGrpcMetadataContextStreamMiddleware applies the unary metadata contract to streams.
 */
func ParseGrpcMetadataContextStreamMiddleware(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if stream == nil {
		return status.Error(codes.Internal, "server stream is nil")
	}

	parsedCtx := stream.Context()
	fullMethod := ""
	if info != nil {
		fullMethod = info.FullMethod
	}
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		_rpcenv.LoadTransportConfig(),
		stream.Context(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: fullMethod},
		func(handlerCtx context.Context, _ interface{}) (interface{}, error) {
			parsedCtx = handlerCtx
			return nil, nil
		},
	)
	if err != nil {
		return err
	}

	return handler(srv, &serverStreamWithContext{ServerStream: stream, ctx: parsedCtx})
}

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamWithContext) Context() context.Context {
	return s.ctx
}
