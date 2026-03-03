package config

import "context"

type Provider interface {
	Name() string

	// Load toàn bộ config
	Load(ctx context.Context) (map[string]interface{}, error)

	// Watch config change (optional)
	Watch(ctx context.Context, onChange func()) error

	Priority() int // càng cao càng override
}

type SourceOption func(Provider)
