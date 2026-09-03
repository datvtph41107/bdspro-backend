package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentEventExchangeMatchesPublishedContract(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile("../.env")
	require.NoError(t, err)
	require.Contains(t, string(contents),
		"QHPRO_BUSINESS_EVENTS_EXCHANGE="+defaultBusinessEventsExchange)
}

func TestParseRedisDBOverride(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		fallback int
		want     int
		wantErr  bool
	}{
		{name: "fallback", fallback: 2, want: 2},
		{name: "zero", raw: "0", fallback: 2, want: 0},
		{name: "highest logical database", raw: "15", want: 15},
		{name: "negative", raw: "-1", wantErr: true},
		{name: "too high", raw: "16", wantErr: true},
		{name: "not a number", raw: "two", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseRedisDBOverride(tt.raw, tt.fallback)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
