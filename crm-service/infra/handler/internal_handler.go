package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/repo"
	crmpb "pb/types/crm"
)

type CrmInternalService struct {
	crmpb.UnimplementedCrmInternalServiceServer
	FriendRepo   repo.FriendRepo
	FriendMapper *mapper.FriendMapper
}

func NewCrmInternalService(
	friendRepo repo.FriendRepo,
	friendMapper *mapper.FriendMapper,
) *CrmInternalService {
	return &CrmInternalService{FriendRepo: friendRepo, FriendMapper: friendMapper}
}

func (s *CrmInternalService) FriendStatusByProfileIds(ctx context.Context, req *crmpb.FriendStatusByProfileIdsRequest) (*crmpb.FriendStatusByProfileIdsResponse, error) {
	friendStatuses, err := s.FriendRepo.GetFriendStatusByProfileIds(ctx, req.CurrentProfileId, req.ProfileIds)
	if err != nil {
		return nil, err
	}

	friendStatusesPb := s.FriendMapper.FriendItemToPbs(friendStatuses)
	return &crmpb.FriendStatusByProfileIdsResponse{FriendStatus: friendStatusesPb}, nil
}
