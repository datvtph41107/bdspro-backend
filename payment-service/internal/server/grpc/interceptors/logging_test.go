package interceptors

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

func captureDefaultLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()

	var buffer bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))

	return &buffer, func() {
		slog.SetDefault(original)
	}
}

func decodeSingleLogRecord(t *testing.T, buffer *bytes.Buffer) map[string]any {
	t.Helper()

	var record map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buffer.Bytes()), &record))
	return record
}

func TestUnaryLoggerInterceptorProjectsCorrelation(t *testing.T) {
	buffer, restore := captureDefaultLogger(t)
	defer restore()

	requestID := request.NewRequestID()
	ctx := request.WithRequestID(context.Background(), requestID)

	interceptor := UnaryLoggerInterceptor()
	response, err := interceptor(
		ctx,
		"request",
		&grpc.UnaryServerInfo{FullMethod: "/payment.PaymentService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			return "response", nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, "response", response)

	record := decodeSingleLogRecord(t, buffer)
	require.Equal(t, "unary gRPC completed", record["msg"])
	require.Equal(t, requestID, record["request_id"])
	require.Equal(t, "grpc", record["component"])
	require.Equal(t, "/payment.PaymentService/Get", record["method"])
	require.Equal(t, "INFO", record["level"])
}

func TestUnaryRecoveryInterceptorProjectsCanonicalError(t *testing.T) {
	buffer, restore := captureDefaultLogger(t)
	defer restore()

	requestID := request.NewRequestID()
	ctx := request.WithRequestID(context.Background(), requestID)

	interceptor := UnaryRecoveryInterceptor()
	_, err := interceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/payment.PaymentService/Panic"},
		func(context.Context, interface{}) (interface{}, error) {
			panic("boom")
		},
	)

	require.Equal(t, codes.Internal, status.Code(err))

	record := decodeSingleLogRecord(t, buffer)
	require.Equal(t, "panic recovered in unary interceptor", record["msg"])
	require.Equal(t, requestID, record["request_id"])
	require.Equal(t, "grpc", record["component"])
	require.Equal(t, "/payment.PaymentService/Panic", record["method"])
	require.Equal(t, "ERROR", record["level"])
}
