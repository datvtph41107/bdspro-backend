package configloader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveRuntimeSelectionCanonicalProfiles(t *testing.T) {
	t.Setenv(EnvironmentKey, "development")
	t.Setenv(ExecutionModeKey, "host")
	t.Setenv(legacyProfileKey, "production")

	selection, err := ResolveRuntimeSelection()
	require.NoError(t, err)
	require.Equal(t, "development", selection.Environment)
	require.Equal(t, "host", selection.ExecutionMode)

	t.Setenv(ExecutionModeKey, "container")
	selection, err = ResolveRuntimeSelection()
	require.NoError(t, err)
	require.Equal(t, "container", selection.ExecutionMode)
}

func TestResolveRuntimeSelectionProductionHasOneProfile(t *testing.T) {
	t.Setenv(EnvironmentKey, "production")
	t.Setenv(ExecutionModeKey, "container")

	selection, err := ResolveRuntimeSelection()
	require.NoError(t, err)
	require.Equal(t, "production", selection.Environment)
}

func TestResolveRuntimeSelectionRejectsInvalidMode(t *testing.T) {
	t.Setenv(EnvironmentKey, "development")
	t.Setenv(ExecutionModeKey, "vm")

	_, err := ResolveRuntimeSelection()
	require.ErrorContains(t, err, ExecutionModeKey)
}

func TestResolveRuntimeSelectionKeepsLegacyCompatibility(t *testing.T) {
	t.Setenv(EnvironmentKey, "")
	t.Setenv(ExecutionModeKey, "")
	t.Setenv(legacyProfileKey, "develop")

	selection, err := ResolveRuntimeSelection()
	require.NoError(t, err)
	require.Equal(t, "development", selection.Environment)
	require.Equal(t, "container", selection.ExecutionMode)
}
