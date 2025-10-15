package clients

import (
	"common"
	_middleware "common/middleware"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"log"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

type CrmGrpcClient struct {
	client crmpb.CrmInternalServiceClient
}

func MakeCrmGrpcClient(env_runtime string) *CrmGrpcClient {
	crmClient := &CrmGrpcClient{}
	go func() {
		name := "crm"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8206", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := crmpb.NewCrmInternalServiceClient(conn)
		crmClient.client = client
	}()
	return crmClient
}

func (c *CrmGrpcClient) GetMapFriendStatusByProfileIds(ctx context.Context, profileIds map[uint64]struct{}) (map[uint64]*sharepb.FriendItem, error) {
	if c.client == nil {
		return nil, errors.New("crm client not avaiable")
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	if currentUserId == 0 {
		return nil, errors.New("current user not found")
	}

	profileIDs := make([]uint64, 0, len(profileIds))
	for id := range profileIds {
		profileIDs = append(profileIDs, id)
	}

	response, err := c.client.FriendStatusByProfileIds(ctx, &crmpb.FriendStatusByProfileIdsRequest{
		ProfileIds:       profileIDs,
		CurrentProfileId: currentUserId,
	})
	if err != nil {
		return nil, err
	}

	friendStatusMap := make(map[uint64]*sharepb.FriendItem, len(response.FriendStatus))
	for _, friendStatus := range response.FriendStatus {
		if friendStatus.CreatedBy != nil && *friendStatus.CreatedBy == currentUserId {
			friendStatusMap[friendStatus.ReceiverId] = friendStatus
		} else if friendStatus.CreatedBy != nil {
			friendStatusMap[*friendStatus.CreatedBy] = friendStatus
		}

	}

	return friendStatusMap, nil
}
