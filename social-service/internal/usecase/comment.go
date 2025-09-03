package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"social/internal/domain"
	"social/internal/dto"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type CommentUsecase struct {
	commentRepo         repo.CommentRepo
	permissionUsecase   iusecase.PermissionUsecase
	newsFeedRepo        repo.NewsFeedRepo
	likeRepo            repo.LikeRepo
	numberUpdateUsecase *NumberUpdateUsecase
	userClient          iusecase.UserClient
	notiClient          iusecase.NotiClient
}

func NewCommentUsecase(commentRepo repo.CommentRepo,
	permissionUsecase iusecase.PermissionUsecase,
	newsFeedRepo repo.NewsFeedRepo,
	likeRepo repo.LikeRepo,
	numberUpdateUsecase *NumberUpdateUsecase,
	userClient iusecase.UserClient,
	notiClient iusecase.NotiClient,
) *CommentUsecase {
	return &CommentUsecase{
		commentRepo:         commentRepo,
		permissionUsecase:   permissionUsecase,
		newsFeedRepo:        newsFeedRepo,
		likeRepo:            likeRepo,
		numberUpdateUsecase: numberUpdateUsecase,
		userClient:          userClient,
		notiClient:          notiClient,
	}
}

func (u *CommentUsecase) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	comment.UserID = profileID
	var err error

	var parentComment *domain.Comment
	if comment.ParentID != nil {
		parentComment, err = u.commentRepo.GetByID(ctx, *comment.ParentID)
		if err != nil {
			return nil, err
		}
		if parentComment == nil {
			return nil, errors.New("bình luận không tồn tại")
		}
		comment.NewsFeedID = parentComment.NewsFeedID
	}

	exist, err := u.newsFeedRepo.GetPublicIgnorePreloadByID(ctx, comment.NewsFeedID)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, _errors.ReturnError(400, "bài viết không tồn tại hoặc bị giới hạn bình luận")
	}

	result, err := u.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	newCtx := _utils.CloneContext(ctx)
	if result.ParentID != nil {
		go u.numberUpdateUsecase.UpdateNumberReplyComment(newCtx, result, parentComment, exist)
	} else {
		go u.numberUpdateUsecase.UpdateNumberComment(newCtx, result, parentComment, exist)
	}

	return result, nil
}

func (u *CommentUsecase) UpdateComment(ctx context.Context, dto *domain.Comment) (*domain.Comment, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	dto.UserID = profileID

	comment, err := u.commentRepo.GetByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("bình luận không tồn tại")
	}
	if comment.UserID != profileID {
		return nil, errors.New("bình luận không phải của bạn")
	}
	result, err := u.commentRepo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *CommentUsecase) DeleteComment(ctx context.Context, comment *domain.Comment) error {
	err := u.permissionUsecase.OwnerComment(ctx, comment.ID, comment.UserID)
	if err != nil {
		return err
	}
	return u.commentRepo.Delete(ctx, comment)
}

func (u *CommentUsecase) GetComment(ctx context.Context, dto *dto.CommentRequest) ([]domain.CommentQuery, int64, error) {
	comments, total, err := u.commentRepo.SearchPublic(ctx, dto)
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func (u *CommentUsecase) GetCommentByNewsFeedID(ctx context.Context, newsFeedID uint64) ([]*domain.Comment, error) {
	return u.commentRepo.GetByNewsFeedID(ctx, newsFeedID)
}

func (u *CommentUsecase) GetCommentByParentID(ctx context.Context, parentID uint64) ([]*domain.Comment, error) {
	return u.commentRepo.GetByParentID(ctx, parentID)
}

func (u *CommentUsecase) GetCommentByNewsFeedIDAndUserID(ctx context.Context, newsFeedID uint64, userID uint64) ([]*domain.Comment, error) {
	return u.commentRepo.GetByNewsFeedIDAndUserID(ctx, newsFeedID, userID)
}
