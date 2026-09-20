package _redis

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"common/logging"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log/slog"
)

// RedisService định nghĩa cấu trúc và phương thức thao tác với Redis
type RedisService struct {
	Client    *redis.Client
	Ctx       context.Context
	connected bool
	mu        sync.RWMutex
}

const userPrefix = "USER_"

// errRedisUnavailable returned when Redis was down at boot / marked disconnected.
// Callers (e.g. HasPermissions) should fall back to DB immediately — do not dial.
var errRedisUnavailable = errors.New("redis unavailable")

// Config is the process-owned technical input needed to create one Redis pool.
type Config struct {
	Address      string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	MaxRetries   int
}

// Open creates one Redis pool from explicit configuration. An unreachable
// Redis instance preserves the historical degraded-start behavior: the service
// is returned disconnected so callers can fall back without opening new pools.
func Open(config Config) (*RedisService, error) {
	address := strings.TrimSpace(config.Address)
	if address == "" {
		return nil, errors.New("redis address is required")
	}

	dialTimeout := config.DialTimeout
	if dialTimeout <= 0 {
		dialTimeout = 500 * time.Millisecond
	}
	readTimeout := config.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = 500 * time.Millisecond
	}
	writeTimeout := config.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = 500 * time.Millisecond
	}
	poolTimeout := config.PoolTimeout
	if poolTimeout <= 0 {
		poolTimeout = 500 * time.Millisecond
	}
	maxRetries := config.MaxRetries
	if maxRetries == 0 {
		maxRetries = 1
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolTimeout:  poolTimeout,
		MaxRetries:   maxRetries,
	})
	svc := &RedisService{
		Client: rdb,
		Ctx:    context.Background(),
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	if _, err := rdb.Ping(pingCtx).Result(); err != nil {
		logging.WithComponent(context.Background(), "redis").Warn("redis unavailable at startup", slog.Any("error", err))
		return svc, nil
	}
	svc.connected = true
	return svc, nil
}

// NewRedisService is the compatibility constructor for services that still
// source Redis config from the shared Viper registry.
func NewRedisService() *RedisService {
	svc, err := Open(Config{
		Address:  viper.GetString("redis.host"),
		Password: viper.GetString("redis.pass"),
		DB:       viper.GetInt("redis.db"),
	})
	if err == nil {
		return svc
	}

	// Keep the old no-error constructor contract. The returned disconnected
	// client makes failures explicit at operation time without hiding a second
	// configuration source.
	logging.WithComponent(context.Background(), "redis").Error("redis configuration invalid", slog.Any("error", err))
	return &RedisService{
		Client: redis.NewClient(&redis.Options{}),
		Ctx:    context.Background(),
	}
}

// Close releases the underlying Redis pool. The process that calls Open owns
// this lifecycle and should call Close exactly once during shutdown.
func (r *RedisService) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}
	return r.Client.Close()
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
		logging.WithComponent(r.Ctx, "redis").Warn("redis connection lost", slog.Any("error", err))
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
		logging.WithComponent(r.Ctx, "redis").Warn("redis reconnect failed", slog.Any("error", err))
		return nil, err
	}

	r.connected = true
	logging.WithComponent(r.Ctx, "redis").Info("redis reconnected")
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
		logging.WithComponent(r.Ctx, "redis").Error("redis get failed", slog.String("key", key), slog.Any("error", err))
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
			logging.WithComponent(r.Ctx, "redis").Error("redis key delete failed", slog.String("key", iter.Val()), slog.Any("error", err))
			r.markDisconnected()
		}
	}
	if err := iter.Err(); err != nil {
		logging.WithComponent(r.Ctx, "redis").Error("redis scan failed", slog.Any("error", err))
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
