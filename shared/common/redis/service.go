package _redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisService định nghĩa cấu trúc và phương thức thao tác với Redis
type RedisService struct {
	Client *redis.Client
	Ctx    context.Context
}

const userPrefix = "USER_"

// NewRedisService khởi tạo RedisService
func NewRedisService() *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"), // Địa chỉ Redis (ví dụ: "localhost:6379")
		Password: viper.GetString("redis.pass"),
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis at ", viper.GetString("redis.host"))
	return &RedisService{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

// SaveToken lưu token với thời gian hết hạn (TTL)
func (r *RedisService) SaveToken(userID uint64, token string, expirationInSeconds uint64) error {
	return r.Client.Set(r.Ctx, userPrefix+strconv.FormatUint(userID, 10), token, time.Duration(expirationInSeconds)*time.Second).Err()
}

// GetToken lấy token từ Redis
func (r *RedisService) GetToken(userID uint64) (string, error) {
	key := userPrefix + strconv.FormatUint(userID, 10)
	return r.Client.Get(r.Ctx, key).Result()
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
	return r.Client.Del(r.Ctx, userPrefix+strconv.FormatUint(userID, 10)).Err()
}

// Set lưu một key-value vào Redis
func (r *RedisService) Set(key, value string) error {
	return r.Client.Set(r.Ctx, key, value, 15*time.Minute).Err()
}

// Get lấy giá trị từ Redis
func (r *RedisService) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

// Delete xóa một key khỏi Redis
func (r *RedisService) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

// DeleteAllKeysWithPrefix xóa tất cả các key có prefix
func (r *RedisService) DeleteAllKeysWithPrefix(prefix string) error {
	iter := r.Client.Scan(r.Ctx, 0, prefix, 0).Iterator()
	for iter.Next(r.Ctx) {
		err := r.Client.Del(r.Ctx, iter.Val()).Err()
		if err != nil {
			fmt.Println("Xóa lỗi key:", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Scan lỗi:", err)
	}

	return nil
}

// Increment tăng giá trị của một key số nguyên trong Redis
func (r *RedisService) Increment(key string) error {
	return r.Client.Incr(r.Ctx, key).Err()
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
	return r.Client.Expire(r.Ctx, key, time.Duration(seconds)*time.Second).Err()
}
