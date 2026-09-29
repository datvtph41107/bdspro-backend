package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRedisDBOverride(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		fallback int
		want     int
		wantErr  bool
	}{
		{name: "fallback", fallback: 1, want: 1},
		{name: "zero", raw: "0", fallback: 1, want: 0},
		{name: "highest logical database", raw: "15", want: 15},
		{name: "negative", raw: "-1", wantErr: true},
		{name: "too high", raw: "16", wantErr: true},
		{name: "not a number", raw: "one", wantErr: true},
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

func TestParsePortOverride(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int
		wantErr bool
	}{
		{name: "yaml fallback", raw: "", want: 8219},
		{name: "environment override", raw: " 9220 ", want: 9220},
		{name: "not a number", raw: "grpc", wantErr: true},
		{name: "zero", raw: "0", wantErr: true},
		{name: "out of range", raw: "65536", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parsePortOverride(test.raw, 8219)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
