package clients

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"log"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	userpb "pb/types/user"

	"google.golang.org/grpc"
)

type CrmGrpcClient struct {
	crmpb.CrmInternalServiceClient
	FriendServiceClient crmpb.FriendServiceClient
}

func BindCrmGrpcClient(conn grpc.ClientConnInterface) *CrmGrpcClient {
	return &CrmGrpcClient{
		CrmInternalServiceClient: crmpb.NewCrmInternalServiceClient(conn),
		FriendServiceClient:      crmpb.NewFriendServiceClient(conn),
	}
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

func (c *CrmGrpcClient) GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendItemDTO, error) {
	return nil, nil
}

func (c *CrmGrpcClient) GetRelationShipV3(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendV3DTO, error) {

	response, err := c.FriendStatusByProfileIds(ctx, &crmpb.FriendStatusByProfileIdsRequest{
		ProfileIds:       []uint64{currentId, targetId},
		CurrentProfileId: currentId,
	})
	if err != nil {
		return nil, err
	}

	if len(response.FriendStatus) == 0 {
		return nil, _errors.ReturnError(404, "chưa kết bạn")
	}

	return &_dto.FriendV3DTO{
		ID:         response.FriendStatus[0].Id,
		ReceiverID: response.FriendStatus[0].ReceiverId,
		CreatedBy:  response.FriendStatus[0].CreatedBy,
		Status:     response.FriendStatus[0].Status,
	}, nil
}

func (c *CrmGrpcClient) GetUserRateStats(ctx context.Context, userId uint64) (*_dto.RateStats, error) {
	// rateStats, err := c.FeedbackGrpcClient.GetRateStats(ctx, userId, 10)
	// if err != nil {
	// 	return nil, err
	// }
	return &_dto.RateStats{
		TotalReviews: 0,
		AverageScore: 0,
		Star1Count:   0,
		Star2Count:   0,
		Star3Count:   0,
		Star4Count:   0,
		Star5Count:   0,
	}, nil
}

func (c *CrmGrpcClient) GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error) {

	resp, err := c.CrmInternalServiceClient.GetRateStats(ctx, &sharepb.IdRequest{Id: userId})
	if err != nil {
		return &_dto.RateStateV3DTO{
			TotalReviews: 0,
			AverageScore: 0,
			Star1Count:   0,
			Star2Count:   0,
			Star3Count:   0,
			Star4Count:   0,
			Star5Count:   0,
		}, nil
	}

	return &_dto.RateStateV3DTO{
		TotalReviews: uint64(resp.TotalReviews),
		AverageScore: float32(resp.AverageScore),
		Star1Count:   uint64(resp.Star1Count),
		Star2Count:   uint64(resp.Star2Count),
		Star3Count:   uint64(resp.Star3Count),
		Star4Count:   uint64(resp.Star4Count),
		Star5Count:   uint64(resp.Star5Count),
	}, nil
}

// GetOriginProfile lấy origin profile từ CRM service theo profileID
func (c *CrmGrpcClient) GetOriginID(ctx context.Context, profileID uint64) (uint64, error) {
	if c.CrmInternalServiceClient == nil {
		log.Printf("CRM client is nil")
		return 0, nil
	}

	resp, err := c.CrmInternalServiceClient.GetOriginProfileByProfileID(ctx, &sharepb.IdRequest{
		Id: profileID,
	})
	if err != nil {
		log.Printf("Failed to get origin profile: %v", err)
		return 0, err
	}

	if resp == nil {
		return 0, nil
	}

	// Trả về origin profile đầu tiên
	return resp.OriginId, nil
}

func (c *CrmGrpcClient) GetFriendInfo(ctx context.Context, profileID uint64) (*_dto.ContactInfoV3DTO, error) {
	if c.CrmInternalServiceClient == nil {
		return nil, _errors.ReturnError(500, "CRM internal service client not available")
	}

	resp, err := c.CrmInternalServiceClient.GetFriendInfoByProfileId(ctx, &sharepb.IdRequest{Id: profileID})
	if err != nil {
		return nil, err
	}

	return &_dto.ContactInfoV3DTO{
		NumFriend:    resp.NumFriend,
		NumFollower:  resp.NumFollower,
		NumFollowing: resp.NumFollowing,
		ContactId:    resp.ContactId,
		BlockId:      resp.BlockId,
		ReportId:     resp.ReportId,
		FollowId:     resp.FollowId,
	}, nil
}

func (c *CrmGrpcClient) MapFriendStatusToSearchItemPb(ctx context.Context, searchItems []*userpb.UserSearchItem) {
	profileIds := make(map[uint64]struct{})
	for _, searchItem := range searchItems {
		profileIds[searchItem.ProfileId] = struct{}{}
	}

	friendStatusMap, err := c.GetMapFriendStatusByProfileIds(ctx, profileIds)
	if err != nil {
		return
	}

	for i, searchItem := range searchItems {
		searchItems[i].FriendStatus = friendStatusMap[searchItem.ProfileId]
	}
}

// GetContactByOriginId retrieves contact information by originId using internal service
func (c *CrmGrpcClient) GetContactByOriginId(ctx context.Context, originId uint64) (*sharepb.ContactDTO, error) {
	if c.CrmInternalServiceClient == nil {
		return nil, errors.New("crm service not initialized")
	}

	resp, err := c.CrmInternalServiceClient.GetContactByOriginId(
		ctx,
		&sharepb.IdRequest{
			Id: originId,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
