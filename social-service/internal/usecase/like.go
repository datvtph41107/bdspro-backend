package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"social/internal/domain"
	"social/internal/enums"
	"social/internal/repo"
)

type LikeUsecase struct {
	likeRepo            repo.LikeRepo
	commentRepo         repo.CommentRepo
	newsFeedRepo        repo.NewsFeedRepo
	numberUpdateUsecase *NumberUpdateUsecase
}

func NewLikeUsecase(likeRepo repo.LikeRepo,
	commentRepo repo.CommentRepo,
	newsFeedRepo repo.NewsFeedRepo,
	numberUpdateUsecase *NumberUpdateUsecase,
) *LikeUsecase {
	return &LikeUsecase{
		likeRepo:            likeRepo,
		commentRepo:         commentRepo,
		newsFeedRepo:        newsFeedRepo,
		numberUpdateUsecase: numberUpdateUsecase,
	}
}

func (u *LikeUsecase) LikeComment(ctx context.Context, targetId uint64, dislike bool) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "unauthorized",
		}
	}
	existComment, err := u.commentRepo.ExistByID(ctx, targetId)
	if err != nil {
		return nil, err
	}
	if !existComment {
		return nil, errors.New("comment not found")
	}

	result, err := u.Like(ctx, targetId, enums.TargetTypeComment, dislike)
	if err != nil {
		return nil, err
	}

	bg := context.Background()
	newCtx := context.WithValue(bg, "profileId", profileId)
	go u.numberUpdateUsecase.UpdateLikeNumberComment(newCtx, targetId, dislike)
	return result, nil
}

func (u *LikeUsecase) UnlikeComment(ctx context.Context, targetId uint64) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "unauthorized",
		}
	}
	result, err := u.Unlike(ctx, targetId, enums.TargetTypeComment)
	if err != nil {
		return nil, err
	}
	bg := context.Background()
	newCtx := context.WithValue(bg, "profileId", profileId)
	go u.numberUpdateUsecase.UpdateLikeNumberComment(newCtx, targetId, false)
	return result, nil
}

func (u *LikeUsecase) LikeNewsFeed(ctx context.Context, targetId uint64) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "unauthorized",
		}
	}
	existNewsFeed, err := u.newsFeedRepo.ExistByID(ctx, targetId)
	if err != nil {
		return nil, err
	}
	if !existNewsFeed {
		return nil, errors.New("news feed not found")
	}

	result, err := u.Like(ctx, targetId, enums.TargetTypeNewsFeed, false)
	if err != nil {
		return nil, err
	}

	bg := context.Background()
	newCtx := context.WithValue(bg, "profileId", profileId)
	go u.numberUpdateUsecase.UpdateLikeNumberNewsFeed(newCtx, targetId)
	return result, nil
}

func (u *LikeUsecase) UnlikeNewsFeed(ctx context.Context, targetId uint64) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "unauthorized",
		}
	}
	result, err := u.Unlike(ctx, targetId, enums.TargetTypeNewsFeed)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *LikeUsecase) Like(ctx context.Context, targetId uint64, targetType enums.TargetType, dislike bool) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    401,
			Message: "unauthorized",
		}
	}
	liked, err := u.likeRepo.ExistByTargetIdAndTypeAndUserId(ctx,
		targetId,
		targetType,
		profileId,
	)
	if err != nil {
		return nil, err
	}
	if liked {
		return nil, errors.New("already like")
	}

	like := &domain.Like{
		TargetID:   targetId,
		TargetType: targetType,
		UserID:     profileId,
		DisLike:    dislike,
		CreatedAt:  _utils.NowPtr(),
	}
	return u.likeRepo.Create(ctx, like)
}

func (u *LikeUsecase) Unlike(ctx context.Context, targetId uint64, targetType enums.TargetType) (*domain.Like, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	liked, err := u.likeRepo.ExistByTargetIdAndTypeAndUserId(ctx,
		targetId,
		targetType,
		profileId,
	)
	if err != nil {
		return nil, err
	}
	if !liked {
		return nil, errors.New("not like")
	}

	err = u.likeRepo.Delete(ctx, &domain.Like{
		TargetID:   targetId,
		TargetType: targetType,
		UserID:     profileId,
	})

	return &domain.Like{
		TargetID:   targetId,
		TargetType: targetType,
		UserID:     profileId,
	}, err
}
