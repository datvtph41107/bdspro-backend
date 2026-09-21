package interceptors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	_errors "common/errors"
	"common/request"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func TestUnaryLoggerInterceptorProjectsSuccessfulCompletion(t *testing.T) {
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
	require.Equal(t, "success", record["outcome"])
	require.Equal(t, "INFO", record["level"])
}

func TestUnaryLoggerInterceptorProjectsCanonicalFailureWithoutOwningSeverity(t *testing.T) {
	buffer, restore := captureDefaultLogger(t)
	defer restore()

	spec := _errors.MustSpec(
		599991,
		"PAYMENT_TEST_UNAVAILABLE",
		"payment temporarily unavailable",
		codes.Unavailable,
	)
	want := _errors.ReturnError(spec, _errors.WithCause(errors.New("database unavailable")))

	interceptor := UnaryLoggerInterceptor()
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/payment.PaymentService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			return nil, want
		},
	)

	require.ErrorIs(t, err, want)

	record := decodeSingleLogRecord(t, buffer)
	require.Equal(t, "unary gRPC completed", record["msg"])
	require.Equal(t, "error", record["outcome"])
	require.Equal(t, "INFO", record["level"])
	require.Equal(t, float64(599991), record["error_code"])
	require.Equal(t, "PAYMENT_TEST_UNAVAILABLE", record["error_reason"])
	require.Equal(t, codes.Unavailable.String(), record["grpc_code"])
	require.NotContains(t, buffer.String(), "database unavailable")
}
