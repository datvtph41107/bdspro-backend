package _middleware

import (
	_callercontext "common/callercontext"
	_enums "common/domain/enum"
	_identitycontext "common/identitycontext"
	_jwt "common/jwt"
	_operationcontext "common/operationcontext"
	_quotasubject "common/quotasubject"
	_requestcontext "common/requestcontext"
	_transporttrust "common/transporttrust"

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
	if requestID, ok := _requestcontext.FromContext(requestCtx); ok {
		md.Set(_requestcontext.RequestIDMetadataKey, requestID)
	}
	if operationID, ok := _requestcontext.OperationIDFromContext(requestCtx); ok {
		md.Set(_requestcontext.OperationIDMetadataKey, operationID)
	}
	if idempotencyKey, ok := _requestcontext.IdempotencyKeyFromContext(requestCtx); ok {
		md.Set(_requestcontext.IdempotencyKeyMetadataKey, idempotencyKey)
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
	if caller, ok := _callercontext.FromContext(r.Context()); ok {
		_callercontext.AppendMetadata(md, caller)
	}
	appendIfExists(md, "x-qhpro-telemetry-key", r.Header.Get("X-QHPro-Telemetry-Key"))
	appendIfExists(md, "if-none-match", r.Header.Get("If-None-Match"))
	appendIfExists(md, "method", r.Method)

	claims, ok := principalFromHTTPRequest(r)
	if !ok {
		return md
	}

	appendActorMetadata(md, _jwt.ActorFromPrincipal(claims))
	if claims.PlanId != nil {
		appendIfExists(md, "planId", *claims.PlanId)
	}
	if claims.PlanFrom != nil {
		appendIfExists(md, "planFrom", *claims.PlanFrom)
	}
	ensureUserCallerWhenActorPresent(md)

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
		_transporttrust.ConfigFromEnv(),
		ctx, req, info, handler,
	)
}

func parseGrpcMetadataContextMiddlewareWithTrust(
	trustConfig _transporttrust.Config,
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
	trustAssessment := _transporttrust.Assess(fullMethod, md, trustConfig)
	if trustAssessment.Verified() {
		bound, bindErr := _transporttrust.BindCaller(ctx, trustAssessment.Caller)
		if bindErr != nil {
			return nil, status.Error(codes.Internal, "transport caller context conflict")
		}
		ctx = bound
	} else if trustAssessment.Reject() {
		return nil, status.Error(codes.Unauthenticated, "trusted transport assertion required")
	} else if trustAssessment.LegacyUnverified() {
		bound, bindErr := _transporttrust.BindCaller(ctx, _transporttrust.Caller{
			ServiceID: "legacy-unverified",
			Level:     _transporttrust.LevelLegacyUnverified,
		})
		if bindErr != nil {
			return nil, status.Error(codes.Internal, "transport caller context conflict")
		}
		ctx = bound
		log.Printf("[gRPC Server] Legacy unverified metadata on %s: %v", fullMethod, trustAssessment.Err)
	}

	actor, hasActor, err := actorFromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	requestCaller, hasRequestCaller, err := _callercontext.FromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !hasRequestCaller && trustAssessment.CanPromotePrivilegedContext() {
		requestCaller = _callercontext.Caller{Kind: _callercontext.KindInternalService}
		if hasActor {
			requestCaller.Kind = _callercontext.KindUser
		}
		hasRequestCaller = true
	}
	if err := validateCallerActor(requestCaller, hasRequestCaller, hasActor); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if hasRequestCaller && trustAssessment.CanPromotePrivilegedContext() {
		bound, bindErr := _callercontext.BindCaller(ctx, requestCaller)
		if bindErr != nil {
			if errors.Is(bindErr, _callercontext.ErrCallerContextConflict) {
				return nil, status.Error(codes.Internal, "caller context conflict")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid caller classification")
		}
		ctx = bound
	}

	quotaSubject, hasQuotaSubject, err := _quotasubject.FromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if hasQuotaSubject {
		if !trustAssessment.CanPromotePrivilegedContext() {
			return nil, status.Error(codes.Unauthenticated, "verified quota subject assertion required")
		}
		bound, bindErr := _quotasubject.BindSubject(ctx, quotaSubject)
		if bindErr != nil {
			if errors.Is(bindErr, _quotasubject.ErrSubjectContextConflict) {
				return nil, status.Error(codes.Internal, "quota subject context conflict")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid quota subject")
		}
		ctx = bound
	}

	operation, hasOperation, err := _operationcontext.FromMetadata(md)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if hasOperation {
		// if trustErr != nil {
		// 	return nil, status.Error(codes.Unauthenticated, "verified operation assertion required")
		// }
		ctx, err = _operationcontext.Bind(ctx, operation)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
	}

	if hasActor {
		if trustAssessment.CanPromotePrivilegedContext() {
			bound, bindErr := _identitycontext.BindActor(ctx, actor)
			if bindErr != nil {
				if errors.Is(bindErr, _identitycontext.ErrActorContextConflict) {
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
		md.Get(_requestcontext.RequestIDMetadataKey),
	)

	ctx, err = bindIncomingOperationID(
		ctx,
		md.Get(_requestcontext.OperationIDMetadataKey),
	)
	if err != nil {
		return nil, err
	}

	ctx, err = bindIncomingIdempotencyKey(
		ctx,
		md.Get(_requestcontext.IdempotencyKeyMetadataKey),
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

	requestID := _requestcontext.NormalizeRequestID(values[0])
	if !_requestcontext.IsValidRequestID(requestID) {
		return ctx
	}
	return _requestcontext.WithRequestID(ctx, requestID)
}

func bindIncomingOperationID(ctx context.Context, values []string) (context.Context, error) {
	if len(values) == 0 {
		return ctx, nil
	}

	operationID, ok := _requestcontext.AcceptOperationID(values)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "invalid operation identity metadata")
	}

	bound, err := _requestcontext.BindOperationID(ctx, operationID)
	if err != nil {
		return ctx, status.Error(codes.Internal, "operation identity context conflict")
	}
	return bound, nil
}

func bindIncomingIdempotencyKey(ctx context.Context, values []string) (context.Context, error) {
	key, ok, err := _requestcontext.ParseIdempotencyKey(values)
	if err != nil {
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency identity metadata")
	}
	if !ok {
		return ctx, nil
	}

	bound, err := _requestcontext.BindIdempotencyKey(ctx, key)
	if err != nil {
		return ctx, status.Error(codes.Internal, "idempotency identity context conflict")
	}
	return bound, nil
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
	allowLegacyPrivilegedPropagation := _transporttrust.CanPropagateLegacyPrivilegedContext(ctx)
	if actor, ok := _identitycontext.FromContext(ctx); ok {
		appendActorMetadata(md, actor)
	} else if allowLegacyPrivilegedPropagation {
		appendUint64IfExists(_enums.ProfileIDKey, _identitycontext.ProfileIDMetadataKey)
		appendUint64IfExists(_enums.OriginIDKey, _identitycontext.OriginIDMetadataKey)
		appendUint64IfExists(_enums.OrganizationIDKey, _identitycontext.OrganizationIDMetadataKey)
		appendUint64IfExists(_enums.AuthIDKey, _identitycontext.AuthIDMetadataKey)
		appendUint64IfExists(_enums.SessionKey, _identitycontext.SessionIDMetadataKey)
		appendStringIfExists(_enums.RoleKey, _identitycontext.RoleMetadataKey)
		appendStringIfExists(_enums.TypeKey, _identitycontext.TokenTypeMetadataKey)
	}

	if caller, ok := _callercontext.FromContext(ctx); ok {
		_callercontext.AppendMetadata(md, caller)
	} else if outgoingMD, ok := metadata.FromOutgoingContext(ctx); !ok || len(outgoingMD.Get(_callercontext.CallerKindMetadataKey)) == 0 {
		defaultKind := _callercontext.KindInternalService
		if hasActorMetadata(md) || hasActorContext(ctx) {
			defaultKind = _callercontext.KindUser
		}
		_callercontext.AppendMetadata(md, _callercontext.Caller{Kind: defaultKind})
	}
	ensureUserCallerWhenActorPresent(md)

	if subject, ok := _quotasubject.FromContext(ctx); ok {
		_quotasubject.AppendMetadata(md, subject)
	}
	if operation, ok := _operationcontext.FromContext(ctx); ok {
		_operationcontext.AppendMetadata(md, operation)
	}

	if allowLegacyPrivilegedPropagation {
		appendUint64IfExists(_enums.PlanIDKey, "planId")
		appendStringIfExists(_enums.APIKeyKey, "x-api-key")
	}
	appendStringIfExists(_enums.DeviceIDKey, "device-id")

	if requestID, ok := _requestcontext.FromContext(ctx); ok {
		md.Append(_requestcontext.RequestIDMetadataKey, requestID)
	}
	if operationID, ok := _requestcontext.OperationIDFromContext(ctx); ok {
		md.Append(_requestcontext.OperationIDMetadataKey, operationID)
	}
	if idempotencyKey, ok := _requestcontext.IdempotencyKeyFromContext(ctx); ok {
		md.Append(_requestcontext.IdempotencyKeyMetadataKey, idempotencyKey)
	}

	if allowLegacyPrivilegedPropagation {
		appendTimeIfExists(_enums.PlanFromKey, "planFrom")
	}

	return md
}

func hasActorContext(ctx context.Context) bool {
	if actor, ok := _identitycontext.FromContext(ctx); ok && !actor.IsZero() {
		return true
	}
	if val, ok := ctx.Value(_enums.ProfileIDKey).(uint64); ok && val > 0 {
		return true
	}
	if val, ok := ctx.Value(_enums.AuthIDKey).(uint64); ok && val > 0 {
		return true
	}
	return false
}

func hasActorMetadata(md metadata.MD) bool {
	_, hasActor, err := actorFromMetadata(md)
	return err == nil && hasActor
}

/**
 * ensureUserCallerWhenActorPresent upgrades caller kind when actor identity is present.
 */
func ensureUserCallerWhenActorPresent(md metadata.MD) {
	if md == nil || !hasActorMetadata(md) {
		return
	}
	caller, hasCaller, err := _callercontext.FromMetadata(md)
	if err != nil {
		return
	}
	if !hasCaller || caller.Kind == _callercontext.KindAnonymous || caller.Kind == _callercontext.KindInternalService {
		_callercontext.AppendMetadata(md, _callercontext.Caller{Kind: _callercontext.KindUser})
	}
}

func validateCallerActor(caller _callercontext.Caller, hasCaller bool, hasActor bool) error {
	if !hasCaller {
		return nil
	}
	switch caller.Kind {
	case _callercontext.KindUser:
		if !hasActor {
			return errors.New("user caller requires actor identity")
		}
	case _callercontext.KindAnonymous, _callercontext.KindAPIKey, _callercontext.KindInternalService:
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
	actor, ok := _identitycontext.FromContext(ctx)
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

/**
 * appendActorMetadata serializes one canonical actor assertion.
 */
func appendActorMetadata(md metadata.MD, actor _identitycontext.Actor) {
	if !actor.IsValid() {
		return
	}
	if actor.AuthID != 0 {
		md.Set(_identitycontext.AuthIDMetadataKey, strconv.FormatUint(actor.AuthID, 10))
	}
	if actor.ProfileID != 0 {
		md.Set(_identitycontext.ProfileIDMetadataKey, strconv.FormatUint(actor.ProfileID, 10))
	}
	if actor.OriginID != 0 {
		md.Set(_identitycontext.OriginIDMetadataKey, strconv.FormatUint(actor.OriginID, 10))
	}
	if actor.OrganizationID != nil {
		md.Set(_identitycontext.OrganizationIDMetadataKey, strconv.FormatUint(*actor.OrganizationID, 10))
	}
	if actor.SessionID != 0 {
		md.Set(_identitycontext.SessionIDMetadataKey, strconv.FormatUint(actor.SessionID, 10))
	}
	if actor.Role != "" {
		md.Set(_identitycontext.RoleMetadataKey, actor.Role)
	}
	if actor.TokenType != "" {
		md.Set(_identitycontext.TokenTypeMetadataKey, actor.TokenType)
	}
}

/**
 * actorFromMetadata validates one transported actor assertion.
 */
func actorFromMetadata(md metadata.MD) (_identitycontext.Actor, bool, error) {
	var actor _identitycontext.Actor
	var present bool

	setUint64 := func(key string, target *uint64) error {
		value, ok, err := singleUint64Metadata(md, key)
		if err != nil {
			return err
		}
		if ok {
			*target = value
			present = true
		}
		return nil
	}

	if err := setUint64(_identitycontext.AuthIDMetadataKey, &actor.AuthID); err != nil {
		return actor, false, err
	}
	if err := setUint64(_identitycontext.ProfileIDMetadataKey, &actor.ProfileID); err != nil {
		return actor, false, err
	}
	if err := setUint64(_identitycontext.OriginIDMetadataKey, &actor.OriginID); err != nil {
		return actor, false, err
	}
	if value, ok, err := singleUint64Metadata(md, _identitycontext.OrganizationIDMetadataKey); err != nil {
		return actor, false, err
	} else if ok {
		actor.OrganizationID = &value
		present = true
	}
	if err := setUint64(_identitycontext.SessionIDMetadataKey, &actor.SessionID); err != nil {
		return actor, false, err
	}
	if value, ok, err := singleStringMetadata(md, _identitycontext.RoleMetadataKey); err != nil {
		return actor, false, err
	} else if ok {
		actor.Role = value
		present = true
	}
	if value, ok, err := singleStringMetadata(md, _identitycontext.TokenTypeMetadataKey); err != nil {
		return actor, false, err
	} else if ok {
		actor.TokenType = value
		present = true
	}

	if !present {
		return actor, false, nil
	}
	if !actor.IsValid() {
		return actor, false, errors.New("invalid actor identity metadata")
	}
	return actor, true, nil
}

func singleUint64Metadata(md metadata.MD, key string) (uint64, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return 0, false, nil
	}
	if len(values) != 1 || values[0] == "" || strings.TrimSpace(values[0]) != values[0] {
		return 0, false, fmt.Errorf("invalid %s identity metadata", key)
	}
	value, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid %s identity metadata", key)
	}
	return value, true, nil
}

func singleStringMetadata(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 || values[0] == "" || len(values[0]) > 128 || strings.TrimSpace(values[0]) != values[0] {
		return "", false, fmt.Errorf("invalid %s identity metadata", key)
	}
	for _, value := range values[0] {
		if value < 0x20 || value == 0x7f {
			return "", false, fmt.Errorf("invalid %s identity metadata", key)
		}
	}
	return values[0], true, nil
}

func withLegacyActorContext(ctx context.Context, actor _identitycontext.Actor) context.Context {
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
		_transporttrust.ConfigFromEnv(),
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
