package facade

import (
	"context"
	"core/kernel/contract"
)

type configKey struct{}

func WithConfig(ctx context.Context, c contract.Config) context.Context {
	return context.WithValue(ctx, configKey{}, c)
}

func Config(ctx context.Context) contract.Config {
	if v := ctx.Value(configKey{}); v != nil {
		return v.(contract.Config)
	}
	panic("config facade not bound")
}
