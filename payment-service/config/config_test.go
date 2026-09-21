package config

import (
	"testing"
	"time"

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

func TestRabbitTimingConfigPreservesLegacyRetryBridge(t *testing.T) {
	t.Setenv("PAYMENT_OUTBOX_RETRY_DELAY", "7s")
	t.Setenv("PAYMENT_OUTBOX_RETRY_BASE", "")
	t.Setenv("PAYMENT_OUTBOX_RETRY_MAX", "")
	t.Setenv("PAYMENT_RABBIT_RECONNECT_BASE", "")
	t.Setenv("PAYMENT_RABBIT_RECONNECT_MAX", "")

	outboxBase, outboxMax, reconnectBase, reconnectMax := rabbitTimingConfig()

	require.Equal(t, 7*time.Second, outboxBase)
	require.Equal(t, time.Minute, outboxMax)
	require.Equal(t, 7*time.Second, reconnectBase)
	require.Equal(t, 30*time.Second, reconnectMax)
}

func TestRabbitTimingConfigLetsCanonicalSettingsOverrideLegacyBridge(t *testing.T) {
	t.Setenv("PAYMENT_OUTBOX_RETRY_DELAY", "9s")
	t.Setenv("PAYMENT_OUTBOX_RETRY_BASE", "3s")
	t.Setenv("PAYMENT_OUTBOX_RETRY_MAX", "45s")
	t.Setenv("PAYMENT_RABBIT_RECONNECT_BASE", "2s")
	t.Setenv("PAYMENT_RABBIT_RECONNECT_MAX", "20s")

	outboxBase, outboxMax, reconnectBase, reconnectMax := rabbitTimingConfig()

	require.Equal(t, 3*time.Second, outboxBase)
	require.Equal(t, 45*time.Second, outboxMax)
	require.Equal(t, 2*time.Second, reconnectBase)
	require.Equal(t, 20*time.Second, reconnectMax)
}
