package _middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	_errors "common/errors"
	"common/request"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func captureErrorBoundaryLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()

	var buffer bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	return &buffer, func() { slog.SetDefault(original) }
}

func decodeBoundaryRecord(t *testing.T, buffer *bytes.Buffer) map[string]any {
	t.Helper()
	var record map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buffer.Bytes()), &record))
	return record
}

func TestUnaryErrorInterceptorLogsCanonicalOperationalCauseAndSanitizesWire(t *testing.T) {
	buffer, restore := captureErrorBoundaryLogger(t)
	defer restore()

	spec := _errors.MustSpec(
		990001,
		"PROFILE_UNAVAILABLE",
		"profile is temporarily unavailable",
		codes.Unavailable,
	)
	cause := errors.New("postgres dial failed: secret-internal-detail")
	requestID := request.NewRequestID()
	ctx := request.WithRequestID(context.Background(), requestID)

	interceptor := UnaryErrorInterceptor()
	_, err := interceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/profile.ProfileService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			return nil, _errors.ReturnError(spec, _errors.WithCause(cause))
		},
	)

	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unavailable, grpcStatus.Code())
	require.Equal(t, "profile is temporarily unavailable", grpcStatus.Message())
	require.NotContains(t, grpcStatus.Message(), "secret-internal-detail")

	record := decodeBoundaryRecord(t, buffer)
	require.Equal(t, "ERROR", record["level"])
	require.Equal(t, "canonical application failure", record["msg"])
	require.Equal(t, "grpc.application.error", record["event_name"])
	require.Equal(t, requestID, record["request_id"])
	require.Equal(t, "/profile.ProfileService/Get", record["grpc.method"])
	require.Equal(t, float64(990001), record["error_code"])
	require.Equal(t, "PROFILE_UNAVAILABLE", record["error_reason"])
	require.Equal(t, codes.Unavailable.String(), record["grpc_code"])
	require.Contains(t, buffer.String(), "secret-internal-detail")
}

func TestUnaryErrorInterceptorDoesNotPromoteExpectedCanonicalFailureToOperationalError(t *testing.T) {
	buffer, restore := captureErrorBoundaryLogger(t)
	defer restore()

	spec := _errors.MustSpec(
		990002,
		"PROFILE_NOT_FOUND",
		"profile not found",
		codes.NotFound,
	)
	interceptor := UnaryErrorInterceptor()
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/profile.ProfileService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			return nil, _errors.ReturnError(spec, _errors.WithCause(errors.New("row missing")))
		},
	)

	require.Equal(t, codes.NotFound, status.Code(err))
	require.Empty(t, strings.TrimSpace(buffer.String()))
}

func TestUnaryErrorInterceptorLogsUnknownTechnicalFailureAndSanitizesWire(t *testing.T) {
	buffer, restore := captureErrorBoundaryLogger(t)
	defer restore()

	interceptor := UnaryErrorInterceptor()
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/profile.ProfileService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			return nil, errors.New("sql: connection reset by peer")
		},
	)

	require.Equal(t, codes.Internal, status.Code(err))
	require.Equal(t, "internal server error", status.Convert(err).Message())
	require.Contains(t, buffer.String(), "sql: connection reset by peer")
	require.Contains(t, buffer.String(), "grpc.handler.error")
}
