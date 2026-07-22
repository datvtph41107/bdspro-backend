package grpc_interceptor

import (
	reqctx "bdspro/core/context"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataInterceptor extracts metadata from gRPC context and injects into Go context.
func MetadataInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if userID := getFirst(md, "user-id"); userID != "" {
		ctx = reqctx.WithUserID(ctx, userID)
	}
	if traceID := getFirst(md, "trace-id"); traceID != "" {
		ctx = reqctx.WithTraceID(ctx, traceID)
	}
	if reqID := getFirst(md, "request-id"); reqID != "" {
		ctx = reqctx.WithRequestID(ctx, reqID)
	}
	if lang := getFirst(md, "lang"); lang != "" {
		ctx = reqctx.WithLang(ctx, lang)
	}

	return handler(ctx, req)
}

func getFirst(md metadata.MD, key string) string {
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}
