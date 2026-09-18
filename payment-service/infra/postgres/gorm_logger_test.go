package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"common/request"

	"github.com/stretchr/testify/require"
)

func TestSlogGormLoggerProjectsSlowQuery(t *testing.T) {
	var buffer bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	defer slog.SetDefault(original)

	requestID := request.NewRequestID()
	ctx := request.WithRequestID(context.Background(), requestID)

	logger := newGormLogger()
	logger.Trace(
		ctx,
		time.Now().Add(-500*time.Millisecond),
		func() (string, int64) {
			return "SELECT 1", 1
		},
		nil,
	)

	var record map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(buffer.Bytes()), &record))
	require.Equal(t, "slow database query", record["msg"])
	require.Equal(t, requestID, record["request_id"])
	require.Equal(t, "postgres", record["component"])
	require.Equal(t, "WARN", record["level"])
	require.Equal(t, "SELECT 1", record["sql"])
	require.Equal(t, float64(1), record["rows"])
}
