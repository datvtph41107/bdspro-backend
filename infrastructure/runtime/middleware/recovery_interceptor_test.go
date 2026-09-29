package _middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"common/request"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryRecoveryInterceptorProjectsCanonicalError(t *testing.T) {
	var buffer bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	defer slog.SetDefault(original)

	requestID := request.NewRequestID()
	ctx := request.WithRequestID(context.Background(), requestID)

	interceptor := UnaryRecoveryInterceptor()
	_, err := interceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Panic"},
		func(context.Context, interface{}) (interface{}, error) {
			panic("boom")
		},
	)

	require.Equal(t, codes.Internal, status.Code(err))

	var record map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buffer.Bytes()), &record))
	require.Equal(t, "panic recovered in unary interceptor", record["msg"])
	require.Equal(t, requestID, record["request_id"])
	require.Equal(t, "grpc", record["component"])
	require.Equal(t, "/test.Service/Panic", record["method"])
	require.Equal(t, "ERROR", record["level"])
}

func TestUnaryRecoveryInterceptorPassesThroughHandler(t *testing.T) {
	interceptor := UnaryRecoveryInterceptor()
	response, err := interceptor(
		context.Background(),
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/OK"},
		func(context.Context, interface{}) (interface{}, error) {
			return "response", nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, "response", response)
}
