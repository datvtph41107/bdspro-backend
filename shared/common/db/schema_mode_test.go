package _db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestParseSchemaMode(t *testing.T) {
	t.Parallel()

	tests := map[string]SchemaMode{
		"sql":          SchemaModeSQL,
		" migrate ":    SchemaModeSQL,
		"automigrate":  SchemaModeAutoMigrate,
		"AUTO-MIGRATE": SchemaModeAutoMigrate,
		"gorm":         SchemaModeAutoMigrate,
		"off":          SchemaModeOff,
	}
	for input, want := range tests {
		input, want := input, want
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSchemaMode(input)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
	_, err := ParseSchemaMode("both")
	require.Error(t, err, "concurrent schema owners must fail closed")
}

func TestLoadSchemaPolicyPrecedence(t *testing.T) {
	t.Setenv(SchemaModeEnv, "off")
	t.Setenv("QHPRO_USER_DB_SCHEMA_MODE", "automigrate")

	policy, err := LoadSchemaPolicy("user")
	require.NoError(t, err)
	require.Equal(t, SchemaModeAutoMigrate, policy.Mode)
	require.Equal(t, "QHPRO_USER_DB_SCHEMA_MODE", policy.Source)
}

func TestApplySchemaPolicy(t *testing.T) {
	t.Parallel()

	called := false
	callback := func(*gorm.DB) error {
		called = true
		return nil
	}
	require.NoError(t, ApplySchemaPolicy(nil, SchemaPolicy{Mode: SchemaModeSQL}, callback))
	require.False(t, called, "SQL policy invoked AutoMigrate")

	database := &gorm.DB{}
	require.NoError(t, ApplySchemaPolicy(database, SchemaPolicy{Mode: SchemaModeAutoMigrate}, callback))
	require.True(t, called, "AutoMigrate policy did not invoke registry")

	want := errors.New("migration failed")
	err := ApplySchemaPolicy(database, SchemaPolicy{Mode: SchemaModeAutoMigrate}, func(*gorm.DB) error {
		return want
	})
	require.ErrorIs(t, err, want)
}
