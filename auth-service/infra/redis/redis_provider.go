package redis

import (
	_redis "common/redis"
	"context"
	"fmt"
	"time"
	"log"

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

// MGet lấy nhiều giá trị từ Redis bằng các khóa đã cho.
func (r *RedisProvider) MGet(c context.Context, keys ...string) ([]interface{}, error) {
	if len(keys) == 0 {
		return []interface{}{}, nil
	}

	val, err := r.RedisService.Client.MGet(c, keys...).Result()
	if err != nil {
		log.Printf("Lỗi khi MGet từ Redis: %v", err)
		return nil, fmt.Errorf("lỗi khi MGet từ Redis: %w", err)
	}

	return val, nil
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

	var token string
	var err error
	for i := 0; i < 3; i++ {
		token, err = r.RedisService.Client.Get(c, key).Result()
		if err == nil {
			return token, nil
		}
		if err == redis.Nil {
			return "", fmt.Errorf("token không tồn tại cho user %d", userID)
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return "", fmt.Errorf("lỗi khi lấy token từ Redis sau 3 lần thử: %v", err)
}

// IsTokenValid kiểm tra token có hợp lệ không
func (r *RedisProvider) IsTokenValid(c context.Context, userID uint64, token string) (bool, error) {
	key := fmt.Sprintf("user_token:%d", userID)

	var storedToken string
	var err error
	for i := 0; i < 3; i++ {
		storedToken, err = r.RedisService.Client.Get(c, key).Result()
		if err == nil {
			// So sánh token
			return storedToken == token, nil
		}
		if err == redis.Nil {
			return false, nil // Token không tồn tại
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return false, fmt.Errorf("lỗi khi kiểm tra token từ Redis sau 3 lần thử: %v", err)
}

// DeleteToken xóa token khỏi Redis
func (r *RedisProvider) DeleteToken(c context.Context, userID uint64) error {
	key := fmt.Sprintf("user_token:%d", userID)

	var err error
	for i := 0; i < 3; i++ {
		err = r.RedisService.Client.Del(c, key).Err()
		if err == nil {
			return nil
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return fmt.Errorf("lỗi khi xóa token từ Redis sau 3 lần thử: %v", err)
}

// Close đóng Redis connection
func (r *RedisProvider) Close() error {
	return r.RedisService.Client.Close()
}

// Ping kiểm tra kết nối Redis
func (r *RedisProvider) Ping(c context.Context) error {
	var err error
	for i := 0; i < 3; i++ {
		_, err = r.RedisService.Client.Ping(c).Result()
		if err == nil {
			return nil
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return fmt.Errorf("lỗi kết nối Redis sau 3 lần thử: %v", err)
}

func (r *RedisProvider) Increment(c context.Context, key string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = r.RedisService.Client.Incr(c, key).Err()
		if err == nil {
			return nil
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return fmt.Errorf("lỗi khi increment key từ Redis sau 3 lần thử: %v", err)
}

// SaveAppState lưu trạng thái app (foreground/background) của user vào Redis
func (r *RedisProvider) SaveAppState(c context.Context, userID uint64, state string) error {
	key := fmt.Sprintf("user:app_state:%d", userID)
	timestamp := time.Now().Format(time.RFC3339)
	status := state + "|" + timestamp

	var err error
	for i := 0; i < 3; i++ {
		err = r.RedisService.Client.Set(c, key, status, 0).Err()
		if err == nil {
			return nil
		}
		// Nếu là lỗi connection, thử lại sau 100ms
		if i < 2 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return fmt.Errorf("lỗi khi lưu app state vào Redis sau 3 lần thử: %v", err)
}
