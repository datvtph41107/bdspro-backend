package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	for {
		select {
		case event := <-r.parent.eventChannel:
			r.handleEvent(event)

		case <-r.ctx.Done():
			r.parent.logger.Info("Redis worker stopped")
			return
		}
	}
}

func (r *redisClientWorker) handleEvent(event any) {
	switch v := event.(type) {

	case *chatpb.UpdateRoomMessage:
		fmt.Println("EVENT IS UpdateRoomMessage")
		members, err := r.parent.participantUsecases.InternalGetListParticipant(r.ctx, v.Data.RoomId)
		if err != nil {
			r.parent.logger.Errorf("Get participants error: %v", err)
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
	fmt.Printf("REDIS PUBLISSHHSHSSHSHHSH WORKER HANDLE EVENT: %T\n", event)

	data, err := json.Marshal(event)
	if err != nil {
		r.parent.logger.Errorf("Marshal error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.requestRedisTimeout)
	defer cancel()

	if err := r.redisClient.Publish(ctx, constants.API_WS_CHANNEL, string(data)); err != nil {
		r.parent.logger.Errorf("Redis publish error: %v", err)
	}
}
