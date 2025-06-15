package usecase

import "crm/internal/repo"

type FriendUsecase struct {
	friendRepo repo.FriendRepo
}

func NewFriendUsecase(friendRepo repo.FriendRepo) *FriendUsecase {
	return &FriendUsecase{
		friendRepo: friendRepo,
	}
}
