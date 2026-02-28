package config

import "context"

type Loader interface {
	Load(ctx context.Context, out interface{}) error
	Watch(ctx context.Context, fn func()) error
}

type Config struct {
	// thiết lập trường chung
}
