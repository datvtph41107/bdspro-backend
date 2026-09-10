package providers

import (
	"time"

	"pb/clients"
	userpb "pb/types/user"
	"tqd/config"
	"tqd/infra/client"
	handler_grpc "tqd/infra/handler/grpc"
	"tqd/infra/worker"
	qh_usecase "tqd/internal/usecase/qh"

	_redis "common/redis"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// NewFileConfig builds file-service client config for Wire.
func NewFileConfig(cfg *config.RuntimeConfig) client.FileConfig {
	if cfg != nil {
		return client.FileConfig{
			BaseURL:        cfg.File.BaseURL,
			ServiceAuthKey: cfg.File.ServiceAuthKey,
			Timeout:        cfg.File.Timeout,
		}
	}
	timeout := time.Duration(viper.GetInt("file_service.timeout_seconds")) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return client.FileConfig{
		BaseURL:        viper.GetString("file_service.base_url"),
		ServiceAuthKey: viper.GetString("file_service.service_auth_key"),
		Timeout:        timeout,
	}
}

// NewClassifyPolicy builds classify usecase policy for Wire.
func NewClassifyPolicy(cfg *config.RuntimeConfig) qh_usecase.ClassifyPolicy {
	if cfg != nil {
		return qh_usecase.ClassifyPolicy{
			ReadContent:    cfg.Classify.ReadContent,
			MaxInlineBytes: cfg.Classify.MaxInlineBytes,
		}
	}
	return qh_usecase.ClassifyPolicy{
		ReadContent:    viper.GetBool("qh_planning_classify.read_content"),
		MaxInlineBytes: viper.GetInt("qh_planning_classify.max_inline_bytes"),
	}
}

// NewClassifyWorkerConfig builds classify worker config for Wire.
func NewClassifyWorkerConfig(cfg *config.RuntimeConfig) worker.ClassifyWorkerConfig {
	if cfg != nil {
		return worker.ClassifyWorkerConfig{
			Enabled:   cfg.Classify.Enabled,
			Interval:  cfg.Classify.Interval,
			BatchSize: cfg.Classify.BatchSize,
		}
	}
	interval := time.Duration(viper.GetInt("qh_planning_classify.interval_sec")) * time.Second
	if interval <= 0 {
		interval = 15 * time.Second
	}
	return worker.ClassifyWorkerConfig{
		Enabled:   viper.GetBool("qh_planning_classify.enabled"),
		Interval:  interval,
		BatchSize: viper.GetInt("qh_planning_classify.batch_size"),
	}
}

// NewRedisClient exposes the Redis go-redis client from the process RedisService.
func NewRedisClient(service *_redis.RedisService) *redis.Client {
	if service == nil {
		return nil
	}
	return service.Client
}

// NewInternalAccessClient extracts InternalAccess from the User gRPC client.
func NewInternalAccessClient(rpcClient *clients.UserGrpcClient) userpb.InternalAccessServiceClient {
	if rpcClient == nil {
		return nil
	}
	return rpcClient.AccessClient
}

// NewUsagePermissionAuthorizer adapts Auth gRPC client for admin usage RPCs.
func NewUsagePermissionAuthorizer(rpcClient *clients.AuthGrpcClient) handler_grpc.UsagePermissionAuthorizer {
	return rpcClient
}
