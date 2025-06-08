package usecase

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type NewsFeedUsecase struct {
	postRepo          repo.NewsFeedRepo
	permissionUsecase iusecase.PermissionUsecase
}

func NewNewsFeedUsecase(postRepo repo.NewsFeedRepo, permissionUsecase iusecase.PermissionUsecase) *NewsFeedUsecase {
	return &NewsFeedUsecase{postRepo: postRepo, permissionUsecase: permissionUsecase}
}

func (u *NewsFeedUsecase) CreateNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (*domain.NewsFeed, error) {
	return u.postRepo.Create(ctx, newsFeed)
}

func (u *NewsFeedUsecase) UpdatePost(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	profileId := ctx.Value("profile").(*uint64)
	if err := u.permissionUsecase.OwnerPost(ctx, post.ID, *profileId); err != nil {
		return nil, err
	}
	return u.postRepo.Update(ctx, post)
}

func (u *NewsFeedUsecase) DeletePost(ctx context.Context, post *domain.NewsFeed) error {
	profileId := ctx.Value("profile").(*uint64)
	if err := u.permissionUsecase.OwnerPost(ctx, post.ID, *profileId); err != nil {
		return err
	}
	return u.postRepo.Delete(ctx, post)
}

func (u *NewsFeedUsecase) GetListPost(ctx context.Context, dto *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	return u.postRepo.GetGlobal(ctx, dto)
}

func (u *NewsFeedUsecase) LikePost(ctx context.Context, newFeedID uint64, numLike int) error {
	return u.postRepo.UpdateLike(ctx, newFeedID, numLike)
}

func (u *NewsFeedUsecase) CommentPost(ctx context.Context, newFeedID uint64, numComment int) error {
	return u.postRepo.UpdateComment(ctx, newFeedID, numComment)
}

func (u *NewsFeedUsecase) SharePost(ctx context.Context, newFeedID uint64, numShare int) error {
	return u.postRepo.UpdateShare(ctx, newFeedID, numShare)
}

func (u *NewsFeedUsecase) ViewPost(ctx context.Context, newFeedID uint64, numView int) error {
	return u.postRepo.UpdateView(ctx, newFeedID, numView)
}
