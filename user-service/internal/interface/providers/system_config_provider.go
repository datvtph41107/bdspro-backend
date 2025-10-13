package providers

import "context"

type SystemConfigProvider interface {
	GetSystemConfig(key string) string
	SetSystemConfig(key, value string)
	InitSystemConfig(ctx context.Context) error
	ReloadSystemConfig(ctx context.Context) (int, error)
}
