package _redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisService định nghĩa cấu trúc và phương thức thao tác với Redis
type RedisService struct {
	Client    *redis.Client
	Ctx       context.Context
	logger    *log.Logger
	connected bool
	mu        sync.RWMutex
}

const userPrefix = "USER_"

// errRedisUnavailable returned when Redis was down at boot / marked disconnected.
// Callers (e.g. HasPermissions) should fall back to DB immediately — do not dial.
var errRedisUnavailable = errors.New("redis unavailable")

// NewRedisService khởi tạo RedisService
func NewRedisService() *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"), // Địa chỉ Redis (ví dụ: "localhost:6379")
		Password:     viper.GetString("redis.pass"),
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		PoolTimeout:  500 * time.Millisecond,
		MaxRetries:   1,
	})
	// Không crash app nếu không connect được, chỉ log warning
	svc := &RedisService{
		Client: rdb,
		Ctx:    context.Background(),
		logger: log.New(os.Stdout, "redis: ", log.LstdFlags),
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping(pingCtx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis at %s: %v. App will continue but Redis operations may fail.", viper.GetString("redis.host"), err)
	} else {
		svc.connected = true
		log.Println("Connected to Redis at ", viper.GetString("redis.host"))
	}
	return svc
}

func (r *RedisService) markDisconnected() {
	r.mu.Lock()
	r.connected = false
	r.mu.Unlock()
}

func (r *RedisService) isConnected() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.connected
}

// Client returns the underlying Redis client
func (r *RedisService) Instance() (*redis.Client, error) {
	r.mu.RLock()
	connected := r.connected
	r.mu.RUnlock()

	if connected {
		pingCtx, cancel := context.WithTimeout(r.Ctx, 500*time.Millisecond)
		_, err := r.Client.Ping(pingCtx).Result()
		cancel()
		if err == nil {
			return r.Client, nil
		}
		r.markDisconnected()
		r.logger.Printf("WARN: Redis connection lost: %v", err)
	}

	// Thử kết nối lại (timeout ngắn — tránh treo API ~25s khi Redis down)
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		return r.Client, nil
	}

	pingCtx, cancel := context.WithTimeout(r.Ctx, 500*time.Millisecond)
	_, err := r.Client.Ping(pingCtx).Result()
	cancel()
	if err != nil {
		r.logger.Printf("WARN: Failed to reconnect to Redis: %v", err)
		return nil, err
	}

	r.connected = true
	r.logger.Printf("Reconnected to Redis successfully")
	return r.Client, nil
}

// SaveToken lưu token với thời gian hết hạn (TTL)
func (r *RedisService) SaveToken(userID uint64, token string, expirationInSeconds uint64) error {
	client, err := r.Instance()
	if err != nil {
		return err
	}
	return client.Set(r.Ctx, userPrefix+strconv.FormatUint(userID, 10), token, time.Duration(expirationInSeconds)*time.Second).Err()
}

// GetToken lấy token từ Redis
func (r *RedisService) GetToken(userID uint64) (string, error) {
	key := userPrefix + strconv.FormatUint(userID, 10)
	return r.Get(key)
}

// IsTokenValid kiểm tra token có hợp lệ không
func (r *RedisService) IsTokenValid(userID uint64, token string) (bool, error) {
	storedToken, err := r.GetToken(userID)
	if err != nil {
		return false, err
	}
	return storedToken == token, nil
}

// DeleteToken xóa token khi đăng xuất
func (r *RedisService) DeleteToken(userID uint64) error {
	return r.Delete(userPrefix + strconv.FormatUint(userID, 10))
}

// Set lưu một key-value vào Redis
func (r *RedisService) Set(key, value string) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.Set(r.Ctx, key, value, 15*time.Minute).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}

// SetWithTime lưu một key-value vào Redis
func (r *RedisService) SetWithTime(key string, value int64, expiration time.Duration) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.Set(r.Ctx, key, value, expiration).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}

// Get lấy giá trị từ Redis. Khi Redis down → lỗi ngay (fallback DB), không dial timeout.
func (r *RedisService) Get(key string) (string, error) {
	if !r.isConnected() {
		return "", errRedisUnavailable
	}
	val, err := r.Client.Get(r.Ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.Printf("[RedisService.Get] Error getting key %s: %v", key, err)
		r.markDisconnected()
	}
	return val, err
}

// Delete xóa một key khỏi Redis
func (r *RedisService) Delete(key string) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.Del(r.Ctx, key).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}

// DeleteAllKeysWithPrefix xóa tất cả các key có prefix
func (r *RedisService) DeleteAllKeysWithPrefix(prefix string) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	iter := r.Client.Scan(r.Ctx, 0, prefix, 0).Iterator()
	for iter.Next(r.Ctx) {
		err := r.Client.Del(r.Ctx, iter.Val()).Err()
		if err != nil {
			fmt.Println("Xóa lỗi key:", iter.Val(), err)
			r.markDisconnected()
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Scan lỗi:", err)
		r.markDisconnected()
		return err
	}

	return nil
}

// Increment tăng giá trị của một key số nguyên trong Redis
func (r *RedisService) Increment(key string) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.Incr(r.Ctx, key).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}

// GetAsInt lấy giá trị từ Redis và chuyển đổi thành số nguyên
func (r *RedisService) GetAsInt(key string) (int, error) {
	value, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// Expire đặt thời gian hết hạn cho một key
func (r *RedisService) Expire(key string, seconds int) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.Expire(r.Ctx, key, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}

func (r *RedisService) IncrementBy(key string, delta int64) error {
	if !r.isConnected() {
		return errRedisUnavailable
	}
	err := r.Client.IncrBy(r.Ctx, key, delta).Err()
	if err != nil {
		r.markDisconnected()
	}
	return err
}
