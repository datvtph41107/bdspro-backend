package appcontext

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	ResponseWriterKey ContextKey = "response-writer"
	RequestKey        ContextKey = "request"
)

func WithResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, ResponseWriterKey, w)
}

func GetResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter)
	return w, ok
}

// GetResponseWriterFromGrpcContext lấy writer từ gRPC context (được inject từ gateway)
func GetResponseWriterFromGrpcContext(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value("response-writer").(http.ResponseWriter)
	return w, ok
}

func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, RequestKey, r)
}

func GetRequest(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(RequestKey).(*http.Request)
	return r, ok
}
