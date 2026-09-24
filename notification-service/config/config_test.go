package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultBusinessEventsExchange(t *testing.T) {
	t.Parallel()
	require.Equal(t, "qhpro.payment.events", defaultBusinessEventsExchange)
}

func TestParseDurationOr(t *testing.T) {
	t.Parallel()

	require.Equal(t, time.Second, parseDurationOr("", time.Second))
	require.Equal(t, 5*time.Second, parseDurationOr("5s", time.Second))
	require.Equal(t, time.Second, parseDurationOr("bad", time.Second))
	require.Equal(t, time.Second, parseDurationOr("0s", time.Second))
}

func TestValidateRedisDB(t *testing.T) {
	t.Parallel()

	require.NoError(t, validateRedisDB(0))
	require.NoError(t, validateRedisDB(15))
	require.Error(t, validateRedisDB(-1))
	require.Error(t, validateRedisDB(16))
}

func TestDatabaseDSNPrefersYAMLField(t *testing.T) {
	AppProperties = UserProperties{}
	AppProperties.Database.DSN = "host=yml user=u password=p dbname=db port=5432 sslmode=disable"
	require.Equal(t, AppProperties.Database.DSN, databaseDSN())
}
