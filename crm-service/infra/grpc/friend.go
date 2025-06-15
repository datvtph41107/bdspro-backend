package grpc_service

import (
	"crm/internal/usecase"
	crmpb "pb/types/crm"
)

type FriendService struct {
	UC *usecase.FriendUsecase
	crmpb.UnimplementedFriendServiceServer
}

func NewFriendService(uc *usecase.FriendUsecase) *FriendService {
	return &FriendService{
		UC: uc,
	}
}
