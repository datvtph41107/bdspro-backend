package _middleware

import (
	_enums "common/domain/enum"
	_jwt "common/jwt"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthError struct {
	Message string
}
type authErrorKeyType string

const authErrorKey authErrorKeyType = "auth_error"

func InjectGrpcMetadataContextMiddleware(ctx context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}

	// Luôn inject IP và User-Agent vào metadata (dù có token hay không)
	appendIfExists(md, "client-ip", getClientIP(r))
	appendIfExists(md, "user-agent", r.Header.Get("User-Agent"))

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		ctx = context.WithValue(ctx, authErrorKey, AuthError{Message: "missing token"})
		return md
	}

	// Tách token từ "Bearer ..."
	parts := strings.SplitN(tokenString, "Bearer ", 2)
	tokenStr := strings.TrimSpace(parts[len(parts)-1]) // lấy phần cuối bất kể có "Bearer " hay không

	claims, err := _jwt.ParseJWT(tokenStr)
	if err != nil || claims == nil {
		ctx = context.WithValue(ctx, authErrorKey, AuthError{Message: "invalid or expired token"})
		return md
	}

	appendIfExists(md, "profileId", claims.ProfileId)
	if claims.OrganizationId != nil {
		appendIfExists(md, "organizationId", *claims.OrganizationId)
	}
	appendIfExists(md, "session", claims.Session)
	appendIfExists(md, "role", claims.Role)
	appendIfExists(md, "type", claims.Type)
	appendIfExists(md, "authId", claims.AuthID)
	if claims.PlanId != nil {
		appendIfExists(md, "planId", *claims.PlanId)
	}
	if claims.PlanFrom != nil {
		appendIfExists(md, "planFrom", *claims.PlanFrom)
	}
	appendIfExists(md, "method", r.Method)

	// Thêm writer vào context để có thể truy cập từ gRPC server
	if writer, ok := r.Context().Value("response-writer").(http.ResponseWriter); ok {
		// Lưu writer vào context để có thể truy cập từ gRPC server
		ctx = context.WithValue(ctx, "response-writer", writer)
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	ctx = parseUint64Metadata(ctx, md, "profileId", _enums.ProfileIDKey)
	ctx = parseUint64Metadata(ctx, md, "organizationId", _enums.OrganizationIDKey)
	ctx = parseUint64Metadata(ctx, md, "authId", _enums.AuthIDKey)
	ctx = parseUint64Metadata(ctx, md, "session", _enums.SessionKey)
	ctx = parseUint64Metadata(ctx, md, "planId", _enums.PlanIDKey)
	ctx = parseStringMetadata(ctx, md, "role", _enums.RoleKey)
	ctx = parseStringMetadata(ctx, md, "type", _enums.TypeKey)
	ctx = parseTimeMetadata(ctx, md, "planFrom", _enums.PlanFromKey)

	// Parse IP address và User-Agent từ metadata
	ctx = parseStringMetadata(ctx, md, "client-ip", _enums.ClientIPKey)
	ctx = parseStringMetadata(ctx, md, "user-agent", _enums.UserAgentKey)

	// Parse writer từ context nếu có
	if writer, ok := ctx.Value("response-writer").(http.ResponseWriter); ok {
		// Sử dụng appcontext để lưu writer vào context
		ctx = context.WithValue(ctx, "response-writer", writer)
	}

	return handler(ctx, req)
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

	// Các key giống với middleware ở server
	appendUint64IfExists(_enums.ProfileIDKey, "profileId")
	appendUint64IfExists(_enums.OrganizationIDKey, "organizationId")
	appendUint64IfExists(_enums.AuthIDKey, "authId")
	appendUint64IfExists(_enums.SessionKey, "session")
	appendUint64IfExists(_enums.PlanIDKey, "planId")

	appendStringIfExists(_enums.RoleKey, "role")
	appendStringIfExists(_enums.TypeKey, "type")

	appendTimeIfExists(_enums.PlanFromKey, "planFrom")

	return md
}

// PrincipalFromContext đọc toàn bộ metadata và sinh ra struct Principal.
func PrincipalFromContext(ctx context.Context) (*_jwt.Principal, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no grpc metadata found in context")
	}

	p := &_jwt.Principal{}
	// ─── helper nội bộ ────────────────────────────────────────────────────────────
	parseUint64 := func(key string) (uint64, bool) {
		if v := md.Get(key); len(v) > 0 {
			if num, err := strconv.ParseUint(v[0], 10, 64); err == nil {
				return num, true
			}
		}
		return 0, false
	}

	parseUintPtr := func(key string) *uint64 {
		if num, ok := parseUint64(key); ok {
			return &num
		}
		return nil
	}

	parseString := func(key string) string {
		if v := md.Get(key); len(v) > 0 {
			return v[0]
		}
		return ""
	}

	// ─── map metadata→struct ─────────────────────────────────────────────────────
	if v, ok := parseUint64("authId"); ok {
		p.AuthID = v
	}
	if v, ok := parseUint64("profileId"); ok {
		p.ProfileId = v
	}

	p.OrganizationId = parseUintPtr("organizationId")
	p.PlanId = parseUintPtr("planId")

	p.Role = parseString("role")
	p.Type = parseString("type")

	if v, ok := parseUint64("session"); ok {
		p.Session = v
	}

	// planFrom (RFC3339 → jwt.NumericDate)
	if ts := parseString("planFrom"); ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			p.PlanFrom = jwt.NewNumericDate(t)
		}
	}

	// RegisteredClaims (nếu bạn muốn copy thêm fields standard như iss/exp…)
	// vd: p.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24*time.Hour))

	return p, nil
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
