package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestNewRuntimeMaterializesHubOwnership(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("server.tcp_port", 8280)
	viper.Set("server.name", "hub-service")
	viper.Set("rpc.user.address", "user:8201")
	viper.Set("rpc.bdspro.address", "bdspro:8202")
	viper.Set("rpc.notification.address", "notification:8204")
	viper.Set("database.dsn", "postgres://hub")
	viper.Set("redis.host", "redis:6379")
	viper.Set("redis.pass", "secret")
	viper.Set("redis.db", 3)
	viper.Set("security.protected_methods", []string{"/hubpb.VersionService/CreateVersionFromDev"})

	t.Setenv("QHPRO_SERVICE_ID", "hub-service")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "metadata-secret")
	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "enforce")

	runtime, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}

	if runtime.GRPCPort != 8280 {
		t.Fatalf("GRPCPort = %d", runtime.GRPCPort)
	}
	if runtime.ServerName != "hub-service" {
		t.Fatalf("ServerName = %q", runtime.ServerName)
	}
	if runtime.UserRPCTarget != "user:8201" ||
		runtime.BDSProRPCTarget != "bdspro:8202" ||
		runtime.NotificationRPCTarget != "notification:8204" {
		t.Fatalf("unexpected RPC targets: %#v", runtime)
	}
	if runtime.Database.DSN != "postgres://hub" {
		t.Fatalf("Database.DSN = %q", runtime.Database.DSN)
	}
	if runtime.Redis.Address != "redis:6379" || runtime.Redis.Password != "secret" || runtime.Redis.DB != 3 {
		t.Fatalf("Redis = %#v", runtime.Redis)
	}
	if !reflect.DeepEqual(
		runtime.Security.ProtectedMethods,
		[]string{"/hubpb.VersionService/CreateVersionFromDev"},
	) {
		t.Fatalf("ProtectedMethods = %#v", runtime.Security.ProtectedMethods)
	}
	if runtime.Transport.ServiceAssertion.ServiceID != "hub-service" {
		t.Fatalf("transport service id = %q", runtime.Transport.ServiceAssertion.ServiceID)
	}
	if !runtime.Transport.RequireServiceAssertion {
		t.Fatal("transport must require service assertion in enforce mode")
	}
}

func TestNewRuntimeRejectsMissingTopology(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("server.tcp_port", 8280)
	viper.Set("server.name", "hub-service")
	viper.Set("rpc.user.address", "user:8201")
	viper.Set("rpc.bdspro.address", "")
	viper.Set("rpc.notification.address", "notification:8204")
	viper.Set("database.dsn", "postgres://hub")
	viper.Set("redis.host", "redis:6379")

	if _, err := NewRuntime(); err == nil {
		t.Fatal("NewRuntime() expected missing bdspro RPC target error")
	}
}
