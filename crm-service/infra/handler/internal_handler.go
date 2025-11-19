package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/repo"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type CrmInternalService struct {
	crmpb.UnimplementedCrmInternalServiceServer
	FriendRepo    repo.FriendRepo
	FriendMapper  *mapper.FriendMapper
	FriendUsecase *usecase.FriendUsecase
}

func NewCrmInternalService(
	friendRepo repo.FriendRepo,
	friendMapper *mapper.FriendMapper,
	friendUsecase *usecase.FriendUsecase,
) *CrmInternalService {
	return &CrmInternalService{
		FriendRepo:    friendRepo,
		FriendMapper:  friendMapper,
		FriendUsecase: friendUsecase,
	}
}

func (s *CrmInternalService) FriendStatusByProfileIds(ctx context.Context, req *crmpb.FriendStatusByProfileIdsRequest) (*crmpb.FriendStatusByProfileIdsResponse, error) {
	friendStatuses, err := s.FriendRepo.GetFriendStatusByProfileIds(ctx, req.CurrentProfileId, req.ProfileIds)
	if err != nil {
		return nil, err
	}

	friendStatusesPb := s.FriendMapper.FriendItemToPbs(friendStatuses)
	return &crmpb.FriendStatusByProfileIdsResponse{FriendStatus: friendStatusesPb}, nil
}

func (s *CrmInternalService) GetFriendInfoByProfileId(ctx context.Context, req *sharepb.IdRequest) (*sharepb.ContactInfoV3Proto, error) {
	result, err := s.FriendUsecase.FriendInfoByProfileId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.ContactInfoV3Proto{
		NumFollower:  result.NumFollower,
		NumFollowing: result.NumFollowing,
		NumFriend:    result.NumFriend,
		ContactId:    result.ContactId,
		BlockId:      result.BlockId,
		FollowId:     result.FollowId,
	}, nil
}
