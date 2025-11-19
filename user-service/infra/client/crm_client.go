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

// @bind: user/internal/interface/providers.FriendProvider
type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
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
