package redis

import (
	_redis "common/redis"
	"context"
	"fmt"
	"log"
	"time"
	"user/internal/interface/providers"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	*_redis.RedisService
}

func NewRedisProvider(redisService *_redis.RedisService) providers.CacheProvider {
	return &RedisProvider{RedisService: redisService}
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

func (r *RedisProvider) Set(c context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.RedisService.Client.Set(c, key, value, expiration).Err()
}

func (r *RedisProvider) Get(c context.Context, key string) (string, error) {
	val, err := r.RedisService.Client.Get(c, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisProvider) Delete(c context.Context, key string) error {
	return r.RedisService.Client.Del(c, key).Err()
}

func (r *RedisProvider) DeleteAllKeysWithPrefix(c context.Context, prefix string) error {
	iter := r.RedisService.Client.Scan(c, 0, prefix+"*", 0).Iterator()
	for iter.Next(c) {
		if err := r.RedisService.Client.Del(c, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (r *RedisProvider) Expire(c context.Context, key string, expiration time.Duration) error {
	return r.RedisService.Client.Expire(c, key, expiration).Err()
}

func (r *RedisProvider) IncrementBy(c context.Context, key string, delta int64) error {
	return r.RedisService.Client.IncrBy(c, key, delta).Err()
}

func (r *RedisProvider) GetAsInt(c context.Context, key string) (int, error) {
	val, err := r.Get(c, key)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	var result int
	_, err = fmt.Sscanf(val, "%d", &result)
	return result, err
}
