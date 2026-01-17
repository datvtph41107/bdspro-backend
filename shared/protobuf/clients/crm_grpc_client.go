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
	crmpb.CrmInternalServiceClient
	FriendServiceClient crmpb.FriendServiceClient
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := crmpb.NewCrmInternalServiceClient(conn)
		friendClient := crmpb.NewFriendServiceClient(conn)
		crmClient.CrmInternalServiceClient = client
		crmClient.FriendServiceClient = friendClient
	}()
	return crmClient
}

func (c *CrmGrpcClient) GetMapFriendStatusByProfileIds(ctx context.Context, profileIds map[uint64]struct{}) (map[uint64]*sharepb.FriendItem, error) {
	if c.CrmInternalServiceClient == nil {
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

	response, err := c.FriendStatusByProfileIds(ctx, &crmpb.FriendStatusByProfileIdsRequest{
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

// GetFriendInfoByProfileId lấy friend info theo profileId
func (c *CrmGrpcClient) GetFriendInfoByProfileId(ctx context.Context, profileId uint64) (*sharepb.ContactInfoV3Proto, error) {
	if c.CrmInternalServiceClient == nil {
		return nil, errors.New("crm client not available")
	}

	resp, err := c.CrmInternalServiceClient.GetFriendInfoByProfileId(ctx, &sharepb.IdRequest{Id: profileId})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CrmGrpcClient) GetContactById(ctx context.Context, id uint64) (*sharepb.ContactDTO, error) {
	if c.CrmInternalServiceClient == nil {
		return nil, errors.New("crm internal service client not available")
	}
	contact, err := c.CrmInternalServiceClient.GetContactById(ctx, &sharepb.IdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (c *CrmGrpcClient) HasContactRelation(ctx context.Context, ownerID uint64, contactOriginProfileID uint64) (bool, error) {
	if c.CrmInternalServiceClient == nil {
		return false, errors.New("crm internal service client not available")
	}

	log.Println("ownerID", ownerID)
	resp, err := c.CrmInternalServiceClient.HasContactRelation(
		ctx,
		&crmpb.HasContactRelationRequest{
			OwnerId:                ownerID,
			ContactOriginProfileId: contactOriginProfileID,
		},
	)
	if err != nil {
		return false, err
	}

	return resp.Exists, nil
}

func (c *CrmGrpcClient) GetOriginProfileByOriginIds(ctx context.Context, originIds []uint64) (map[uint64]*sharepb.OriginProfile, error) {
	if c.CrmInternalServiceClient == nil {
		return nil, errors.New("crm internal service client not available")
	}
	_, err := c.CrmInternalServiceClient.GetOriginProfileByOriginIds(ctx, &sharepb.IdRequest{Ids: originIds})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
