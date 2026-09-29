package usecase

import (
	"social/internal/repo"
)

type FriendTagUsecase struct {
	newsFeedRepo  repo.NewsFeedRepo
	friendTagRepo repo.FriendTagRepo
}

func NewFriendTagUsecase(
	newsFeedRepo repo.NewsFeedRepo,
	friendTagRepo repo.FriendTagRepo,
) *FriendTagUsecase {
	return &FriendTagUsecase{
		newsFeedRepo:  newsFeedRepo,
		friendTagRepo: friendTagRepo,
	}
}

// func (u *FriendTagUsecase) SetFriendTag(ctx context.Context, newsFeedID uint64, friendTagIds []uint64) error {
// 	for _, friendTagId := range friendTagIds {
// 		_, err := u.friendTagRepo.Create(ctx, &domain.TagFriend{
// 			NewsFeedID: newsFeedID,
// 			UserID:     friendTagId,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
