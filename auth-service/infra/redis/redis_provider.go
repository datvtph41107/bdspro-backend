package redis

import (
	_redis "common/redis"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// @bind: auth/internal/interface/provider.CacheProvider
type RedisProvider struct {
	*_redis.RedisService
}

func NewRedisProvider() *RedisProvider {
	return &RedisProvider{
		RedisService: _redis.NewRedisService(),
	}
}

// SaveToken lưu token vào Redis với expiration
func (r *RedisProvider) SaveToken(c context.Context, userID uint64, token string, expirationInSeconds uint64) error {
	key := fmt.Sprintf("user_token:%d", userID)

	println("SaveToken: ", key, token, expirationInSeconds)

	// Lưu token với expiration
	err := r.RedisService.Client.Set(c, key, token, time.Duration(expirationInSeconds)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi lưu token vào Redis: %v", err)
	}

	return nil
}

// GetToken lấy token từ Redis theo userID
func (r *RedisProvider) GetToken(c context.Context, userID uint64) (string, error) {
	key := fmt.Sprintf("user_token:%d", userID)

	token, err := r.RedisService.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("token không tồn tại cho user %d", userID)
		}
		return "", fmt.Errorf("lỗi khi lấy token từ Redis: %v", err)
	}

	return token, nil
}

// IsTokenValid kiểm tra token có hợp lệ không
func (r *RedisProvider) IsTokenValid(c context.Context, userID uint64, token string) (bool, error) {
	key := fmt.Sprintf("user_token:%d", userID)

	storedToken, err := r.RedisService.Client.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Token không tồn tại
		}
		return false, fmt.Errorf("lỗi khi kiểm tra token từ Redis: %v", err)
	}

	// So sánh token
	return storedToken == token, nil
}

// DeleteToken xóa token khỏi Redis
func (r *RedisProvider) DeleteToken(c context.Context, userID uint64) error {
	key := fmt.Sprintf("user_token:%d", userID)

	err := r.RedisService.Client.Del(c, key).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi xóa token từ Redis: %v", err)
	}

	return nil
}

// Close đóng Redis connection
func (r *RedisProvider) Close() error {
	return r.RedisService.Client.Close()
}

// Ping kiểm tra kết nối Redis
func (r *RedisProvider) Ping(c context.Context) error {
	_, err := r.RedisService.Client.Ping(c).Result()
	if err != nil {
		return fmt.Errorf("lỗi kết nối Redis: %v", err)
	}
	return nil
}

func (r *RedisProvider) Increment(c context.Context, key string) error {
	return r.RedisService.Client.Incr(c, key).Err()
}

// SaveAppState lưu trạng thái app (foreground/background) của user vào Redis
func (r *RedisProvider) SaveAppState(c context.Context, userID uint64, state string) error {
	key := fmt.Sprintf("user:app_state:%d", userID)
	timestamp := time.Now().Format(time.RFC3339)
	status := state + "|" + timestamp

	err := r.RedisService.Client.Set(c, key, status, 0).Err()
	if err != nil {
		return fmt.Errorf("lỗi khi lưu app state vào Redis: %v", err)
	}

	return nil
}
