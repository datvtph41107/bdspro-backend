package client

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"pb/clients"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/env"
	models "user/internal/models"
)

// @bind: user/internal/interface/providers.CrmProvider
type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *CrmClient) RegisterService(ctx context.Context) error {
	_, err := c.CrmGrpcClient.RegisterEntity(ctx, &crmpb.FeedbackRegistry{
		RegistryName: "user",
		Port:         8201,
		Service:      "user",
	})

	return err
}

func (c *CrmClient) MapFriendStatusToSearchItemPb(ctx context.Context, searchItems []*userpb.UserSearchItem) {
	profileIds := make(map[uint64]struct{})
	for _, searchItem := range searchItems {
		profileIds[searchItem.ProfileId] = struct{}{}
	}

	friendStatusMap, err := c.CrmGrpcClient.GetMapFriendStatusByProfileIds(ctx, profileIds)
	if err != nil {
		return
	}

	for i, searchItem := range searchItems {
		searchItems[i].FriendStatus = friendStatusMap[searchItem.ProfileId]
	}
}

func (c *CrmClient) GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*models.FriendItem, error) {
	return nil, nil
}

func (c *CrmClient) GetRelationShipV3(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendV3DTO, error) {

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

// GetFriendInfo lấy thông tin friend (số lượng bạn bè, follower, following) theo profileId
func (c *CrmClient) GetFriendInfo(ctx context.Context, profileID uint64) (*sharepb.ContactInfoV3Proto, error) {
	if c.CrmGrpcClient.CrmInternalServiceClient == nil {
		return nil, _errors.ReturnError(500, "CRM internal service client not available")
	}

	resp, err := c.CrmGrpcClient.GetFriendInfoByProfileId(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CrmClient) GetRateStats(ctx context.Context, req *sharepb.IdRequest) (*sharepb.RateStatsV3Proto, error) {
	resp, err := c.CrmGrpcClient.GetRateStats(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CrmClient) GetUserRateStats(ctx context.Context, userId uint64) (*models.RateStats, error) {
	// rateStats, err := c.FeedbackGrpcClient.GetRateStats(ctx, userId, 10)
	// if err != nil {
	// 	return nil, err
	// }
	return &models.RateStats{
		TotalReviews: 0,
		AverageScore: 0,
		Star1Count:   0,
		Star2Count:   0,
		Star3Count:   0,
		Star4Count:   0,
		Star5Count:   0,
	}, nil
}

func (c *CrmClient) GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error) {

	resp, err := c.CrmGrpcClient.GetRateStats(ctx, &sharepb.IdRequest{Id: userId})
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
