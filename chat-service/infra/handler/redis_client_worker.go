package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	chatpb "pb/types/chat"

	"chat/infra/redis"
	"chat/internal/constants"

	"common/logging"
)

type redisClientWorker struct {
	parent              *chatHandler
	redisClient         *redis.RedisClient
	requestRedisTimeout time.Duration
	ctx                 context.Context
}

func NewRedisClientWorker(ctx context.Context, parent *chatHandler, redisClient *redis.RedisClient) *redisClientWorker {
	return &redisClientWorker{
		parent:              parent,
		redisClient:         redisClient,
		requestRedisTimeout: 2 * time.Second,
		ctx:                 ctx,
	}
}

func (r *redisClientWorker) serve() {
	for {
		select {
		case event := <-r.parent.eventChannel:
			r.handleEvent(event)

		case <-r.ctx.Done():
			logging.WithComponent(r.ctx, "redis_client_worker").Info("Redis worker stopped")
			return
		}
	}
}

func (r *redisClientWorker) handleEvent(event any) {
	switch v := event.(type) {

	case *chatpb.UpdateRoomMessage:
		members, err := r.parent.participantUsecases.InternalGetListParticipant(r.ctx, v.Data.RoomId)
		if err != nil {
			logging.WithComponent(r.ctx, "redis_client_worker").Error(
				fmt.Sprintf("Get participants error: %v", err),
			)
			return
		}

		userIds := make([]string, 0, len(members))
		for _, m := range members {
			userIds = append(userIds, strconv.FormatUint(m.UserID, 10))
		}
		v.Data.Members = userIds

		r.publish(v)
	default:
		fmt.Println("EVENT IS SendMessage")
		r.publish(v)
	}
}

func (r *redisClientWorker) publish(event any) {
	data, err := json.Marshal(event)
	if err != nil {
		logging.WithComponent(r.ctx, "redis_client_worker").Error(
			fmt.Sprintf("Marshal error: %v", err),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.requestRedisTimeout)
	defer cancel()

	if err := r.redisClient.Publish(ctx, constants.API_WS_CHANNEL, string(data)); err != nil {
		logging.WithComponent(r.ctx, "redis_client_worker").Error(
			fmt.Sprintf("Redis publish error: %v", err),
		)
	}
}
