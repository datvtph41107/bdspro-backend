package usecase

import (
	"context"
	"social/internal/enums"
	"social/internal/repo"
	"time"
)

type NumberUpdateUsecase struct {
	newsFeedRepo repo.NewsFeedRepo
	likeRepo     repo.LikeRepo
	commentRepo  repo.CommentRepo
}

func NewNumberUpdateUsecase(newsFeedRepo repo.NewsFeedRepo, likeRepo repo.LikeRepo, commentRepo repo.CommentRepo) *NumberUpdateUsecase {
	return &NumberUpdateUsecase{newsFeedRepo: newsFeedRepo, likeRepo: likeRepo, commentRepo: commentRepo}
}

func (u *NumberUpdateUsecase) UpdateLikeNumberNewsFeed(ctx context.Context, newsFeedId uint64) (bool, error) {
	if !u.checkTimeNewsFeed(ctx, newsFeedId) {
		return false, nil
	}

	totalLike, err := u.likeRepo.Count(ctx, newsFeedId, enums.TargetTypeNewsFeed, false)
	if err != nil {
		return false, err
	}

	err = u.newsFeedRepo.UpdateLikeNumber(ctx, newsFeedId, totalLike)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *NumberUpdateUsecase) UpdateLikeNumberComment(ctx context.Context, commentID uint64, dislike bool) (bool, error) {
	if !u.checkTimeComment(ctx, commentID) {
		return false, nil
	}

	totalLike, err := u.likeRepo.Count(ctx, commentID, enums.TargetTypeComment, dislike)
	if err != nil {
		return false, err
	}
	err = u.commentRepo.UpdateLikeNumber(ctx, commentID, totalLike, dislike)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *NumberUpdateUsecase) UpdateNumberReplyComment(ctx context.Context, commentID uint64, newsFeedID uint64) (bool, error) {
	if !u.checkTimeComment(ctx, commentID) {
		return false, nil
	}

	// update number reply of comment
	totalReply, err := u.commentRepo.CountReply(ctx, commentID)
	if err != nil {
		return false, err
	}
	err = u.commentRepo.UpdateNumberReply(ctx, commentID, totalReply)
	if err != nil {
		return false, err
	}

	// update number comment of news feed
	return u.setNumberComment(ctx, newsFeedID)
}

func (u *NumberUpdateUsecase) UpdateNumberComment(ctx context.Context, newsFeedID uint64) (bool, error) {
	if !u.checkTimeNewsFeed(ctx, newsFeedID) {
		return false, nil
	}

	return u.setNumberComment(ctx, newsFeedID)
}

func (u *NumberUpdateUsecase) checkTimeNewsFeed(ctx context.Context, newsFeedID uint64) bool {
	updatedTime, err := u.newsFeedRepo.GetLastUpdated(ctx, newsFeedID)
	if err != nil {
		return false
	}
	if time.Now().Before(updatedTime.Add(time.Minute)) {
		return false
	}
	return true
}

func (u *NumberUpdateUsecase) checkTimeComment(ctx context.Context, commentID uint64) bool {
	updatedTime, err := u.commentRepo.GetLastUpdated(ctx, commentID)
	if err != nil {
		return false
	}
	if time.Now().Before(updatedTime.Add(time.Minute)) {
		return false
	}
	return true
}

func (u *NumberUpdateUsecase) setNumberComment(ctx context.Context, newsFeedID uint64) (bool, error) {
	totalComment, err := u.commentRepo.CountByNewsFeedID(ctx, newsFeedID)
	if err != nil {
		return false, err
	}
	err = u.newsFeedRepo.UpdateCommentNumber(ctx, newsFeedID, totalComment)
	if err != nil {
		return false, err
	}

	return true, nil
}
