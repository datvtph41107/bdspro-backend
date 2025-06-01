package providers

import (
	"context"
	"sync"
	"time"
)

type ListCacheProvider[T any, REQ any] struct {
	Data      []T
	Total     int64
	ExpiresAt time.Time
	Lock      sync.RWMutex
	Ttl       time.Duration
	// Loader    func() ([]T, error)
}

func (c *ListCacheProvider[T, REQ]) Setup() {
	c.Ttl = time.Duration(20 * time.Second)
	c.ExpiresAt = time.Now().Add(c.Ttl)
	c.Lock = sync.RWMutex{}
}

func (c *ListCacheProvider[T, REQ]) GetDataWithCache(ctx context.Context, dto *REQ, f func(context.Context, *REQ) ([]T, int64, error)) ([]T, int64, error) {
	c.Lock.RLock()
	if time.Now().Before(c.ExpiresAt) && c.Data != nil {
		data := c.Data
		total := c.Total
		c.Lock.RUnlock()
		return data, total, nil
	}
	c.Lock.RUnlock()

	// Upgrade to write lock
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// Double check (tránh race condition)
	if time.Now().Before(c.ExpiresAt) && c.Data != nil {
		return c.Data, c.Total, nil
	}

	// Load lại dữ liệu
	data, total, err := f(ctx, dto)
	if err != nil {
		return nil, 0, err
	}
	c.Data = data
	c.Total = total
	c.ExpiresAt = time.Now().Add(c.Ttl)

	return data, total, nil
}
