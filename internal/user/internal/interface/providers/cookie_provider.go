package providers

import (
	"context"
)

type CookieProvider interface {
	// SetSessionBrowser(c context.Context, refreshToken string)

	// Đặt refresh token trong cookie
	SetRefreshToken(c context.Context, refreshToken string)
	DeleteRefreshToken(c context.Context)

	// Set cookie
	Get(c context.Context, key string) (string, error)
	Set(c context.Context, key string, value string)
	// GetAsInt(c context.Context, key string) (int, error)
	Expire(c context.Context, key string, seconds int) error
	ResetAll(c context.Context)
}
