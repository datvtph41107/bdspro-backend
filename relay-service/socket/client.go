package socket

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func PushMessageToQueue(key string, message []byte) error {
	_, err := redisClient.LPush(context.Background(), key, message).Result()
	if err != nil {
		log.Printf("Error pushing message to Redis: %v", err)
		return fmt.Errorf("Error pushing message to Redis: %v", err)
	}
	return nil
}

func GetMessagesFromQueue(key string) ([]string, error) {
	messages, err := redisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		log.Printf("Error fetching messages from Redis: %v", err)
		return nil, fmt.Errorf("Error fetching messages from Redis: %v", err)
	}
	return messages, nil
}

func DeleteQueue(key string) {
	redisClient.Del(context.Background(), key)
}
