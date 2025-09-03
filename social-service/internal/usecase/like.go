package usecase

import (
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"time"

	"social/internal/domain"
	"social/internal/enums"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type LikeUsecase struct {
	likeRepo            repo.LikeRepo
	commentRepo         repo.CommentRepo
	newsFeedRepo        repo.NewsFeedRepo
	numberUpdateUsecase *NumberUpdateUsecase
	userClient          iusecase.UserClient
	notiClient          iusecase.NotiClient
}

func NewLikeUsecase(likeRepo repo.LikeRepo,
	commentRepo repo.CommentRepo,
	newsFeedRepo repo.NewsFeedRepo,
	numberUpdateUsecase *NumberUpdateUsecase,
	userClient iusecase.UserClient,
	notiClient iusecase.NotiClient,
) *LikeUsecase {
	return &LikeUsecase{
		likeRepo:            likeRepo,
		commentRepo:         commentRepo,
		newsFeedRepo:        newsFeedRepo,
		numberUpdateUsecase: numberUpdateUsecase,
		userClient:          userClient,
		notiClient:          notiClient,
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
	existComment, err := u.commentRepo.GetByID(ctx, targetId)
	if err != nil {
		return nil, err
	}
	if existComment == nil {
		return nil, errors.New("comment not found")
	}

	result, err := u.Like(ctx, targetId, enums.TargetTypeComment, dislike)
	if err != nil {
		return nil, err
	}

	newCtx := _utils.CloneContext(ctx)
	go u.numberUpdateUsecase.UpdateLikeNumberComment(newCtx, existComment, dislike)

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
	// bg := context.Background()
	// newCtx := context.WithValue(bg, "profileId", profileId)
	// go u.numberUpdateUsecase.UpdateLikeNumberComment(newCtx, targetId, false)
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
	existNewsFeed, err := u.newsFeedRepo.GetByID(ctx, targetId)
	if err != nil {
		return nil, err
	}
	if existNewsFeed == nil {
		return nil, errors.New("news feed not found")
	}

	result, err := u.Like(ctx, targetId, enums.TargetTypeNewsFeed, false)
	if err != nil {
		return nil, err
	}

	newCtx := _utils.CloneContext(ctx)
	go u.numberUpdateUsecase.UpdateLikeNumberNewsFeed(newCtx, existNewsFeed)
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
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	liked, _ := u.likeRepo.GetByTargetIdAndTypeAndUserId(ctx,
		targetId,
		targetType,
		profileId,
	)
	// if err != nil {
	// 	return nil, err
	// }
	if liked != nil &&
		liked.TargetID == targetId &&
		liked.TargetType == targetType &&
		liked.UserID == profileId &&
		liked.DisLike == dislike {
		return nil, errors.New("already like")
	}

	now := time.Now()
	like := &domain.Like{
		TargetID:   targetId,
		TargetType: targetType,
		UserID:     profileId,
		DisLike:    dislike,
		CreatedAt:  &now,
	}

	if liked != nil {
		liked.DisLike = dislike
		return u.likeRepo.Update(ctx, like)
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
