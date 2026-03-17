package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	chatpb "pb/types/chat"

	"chat/internal/constants"

	"github.com/hyperledger/fabric/common/flogging"
)

type redisClientWorker struct {
	parent              *chatHandler
	redisClient         redisClient
	requestRedisTimeout time.Duration
	ctx                 context.Context
}

func NewRedisClientWorker(parent *chatHandler, redisClient redisClient, logger *flogging.FabricLogger) *redisClientWorker {
	return &redisClientWorker{
		parent:              parent,
		redisClient:         redisClient,
		requestRedisTimeout: 2 * time.Second,
		ctx:                 context.Background(),
	}
}

func (r *redisClientWorker) serve() {
	for event := range r.parent.eventChannel {
		r.parent.logger.Infof("Received event: %v", event)

		switch v := event.(type) {

		case *chatpb.UpdateRoomMessage:
			members, err := r.parent.participantUsecases.InternalGetListParticipant(r.ctx, v.Data.RoomId)
			if err != nil {
				r.parent.logger.Errorf("Failed to get list participant: %v", err)
				continue
			}

			// Gán UserIds nếu field có trong protobuf
			if v.Data != nil {
				userIds := make([]string, 0, len(members))
				for _, member := range members {
					userIds = append(userIds, strconv.FormatUint(member.UserID, 10))
				}
				v.Data.Members = userIds
			}

			r.publishToRedis(v)

		default:
			r.publishToRedis(v)
		}
	}
}

func (s *redisClientWorker) publishToRedis(event any) {
	data, err := json.Marshal(event)
	if err != nil {
		s.parent.logger.Errorf("Failed to marshal message: %v", err)
	}

	redisCtx, cancel := context.WithTimeout(context.Background(), s.requestRedisTimeout)
	defer cancel()

	err = s.redisClient.Publish(redisCtx, constants.API_WS_CHANNEL, string(data))
	if err != nil {
		s.parent.logger.Errorf("Failed to publish message to Redis: %v", err)
	}
}
