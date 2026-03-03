package config

import "context"

type Config interface {
	Load(ctx context.Context) error

	// Bind toàn bộ config vào struct (recommended)
	Scan(out interface{}) error

	// Bind theo namespace (fallback)
	Get(path string, out interface{}) error

	// Watch change
	OnChange(path string, fn func()) error
}
