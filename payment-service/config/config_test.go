package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveConfigFileUsesCanonicalRuntimeProfile(t *testing.T) {
	t.Parallel()

	require.Equal(t, "config/runtime.yml", resolveConfigFile(""))
}

func TestResolveConfigFileAllowsExplicitOperatorOverride(t *testing.T) {
	t.Parallel()

	require.Equal(t, "/run/secrets/payment.yml", resolveConfigFile(" /run/secrets/payment.yml "))
}
