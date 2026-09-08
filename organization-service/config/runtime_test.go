package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestInitConfigSelectsHostTopologyAndEnvOverrides(t *testing.T) {
	oldWorkingDirectory, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(".."))
	t.Cleanup(func() { require.NoError(t, os.Chdir(oldWorkingDirectory)) })

	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("QHPRO_ENVIRONMENT", "development")
	t.Setenv("QHPRO_EXECUTION_MODE", "host")
	t.Setenv("ORGANIZATION_DATABASE_URL", "host=db.test dbname=organization")
	t.Setenv("ORGANIZATION_USER_GRPC_ADDRESS", "user.test:8201")

	require.NoError(t, InitConfig())
	require.Equal(t, "host=db.test dbname=organization", viper.GetString("database.dsn"))
	require.Equal(t, "user.test:8201", viper.GetString("rpc.user.address"))
}
