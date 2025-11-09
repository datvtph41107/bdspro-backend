package usecase

import (
	_enum "common/domain/enum"
	_utils "common/utils"
	"context"
	"social/internal/domain"
	"social/internal/enums"
	iusecase "social/internal/interface"
	"social/internal/repo"
	"time"
)

type NumberUpdateUsecase struct {
	newsFeedRepo repo.NewsFeedRepo
	likeRepo     repo.LikeRepo
	commentRepo  repo.CommentRepo
	userClient   iusecase.UserClient
	notiClient   iusecase.NotiClient
	userUsecase  *UserUsecase
}

func NewNumberUpdateUsecase(newsFeedRepo repo.NewsFeedRepo,
	likeRepo repo.LikeRepo,
	commentRepo repo.CommentRepo,
	userClient iusecase.UserClient,
	notiClient iusecase.NotiClient,
	userUsecase *UserUsecase,
) *NumberUpdateUsecase {
	return &NumberUpdateUsecase{
		newsFeedRepo: newsFeedRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		userClient:   userClient,
		notiClient:   notiClient,
		userUsecase:  userUsecase,
	}
}

func (u *NumberUpdateUsecase) UpdateLikeNumberNewsFeed(ctx context.Context, newsFeed *domain.NewsFeed) (bool, error) {
	ok := u.RequiredTimeUpdate(newsFeed.UpdatedAt)
	if !ok {
		return false, nil
	}

	totalLike, err := u.likeRepo.Count(ctx, newsFeed.ID, enums.TargetTypeNewsFeed, false)
	if err != nil {
		return false, err
	}

	err = u.newsFeedRepo.UpdateLikeNumber(ctx, newsFeed.ID, totalLike)
	if err != nil {
		return false, err
	}

	likedUserIds, _ := u.likeRepo.GetNearestUserIds(ctx, newsFeed.ID, enums.TargetTypeNewsFeed)
	sumUserLiked, _ := u.likeRepo.SumUserComment(ctx, newsFeed.ID, enums.TargetTypeNewsFeed)
	fullNames := u.userUsecase.getFullNames(ctx, likedUserIds, sumUserLiked)

	u.notiClient.CreateToOwner(ctx,
		"",
		"",
		[]string{"", fullNames, " thích bài viết của bạn: \"" + _utils.ShortenContent(newsFeed.Content, 50) + "\""},
		_enum.NotificationLikeNewsFeed,
		&newsFeed.ID,
		*newsFeed.CreatedBy,
		_enum.EOwnerOfMember,
		[]string{fullNames, _utils.ShortenContent(newsFeed.Content, 50)},
	)

	return true, nil
}

func (u *NumberUpdateUsecase) UpdateLikeNumberComment(ctx context.Context, comment *domain.Comment, dislike bool) (bool, error) {
	totalLike, err := u.likeRepo.Count(ctx, comment.ID, enums.TargetTypeComment, dislike)
	if err != nil {
		return false, err
	}

	ok := u.RequiredTimeUpdate(comment.UpdatedAt)
	if !ok && totalLike > 0 {
		return false, nil
	}

	err = u.commentRepo.UpdateLikeNumber(ctx, comment.ID, totalLike, dislike)
	if err != nil {
		return false, err
	}

	likedUserIds, _ := u.likeRepo.GetNearestUserIds(ctx, comment.ID, enums.TargetTypeComment)
	sumUserLiked, _ := u.likeRepo.SumUserComment(ctx, comment.ID, enums.TargetTypeComment)
	fullNames := u.userUsecase.getFullNames(ctx, likedUserIds, sumUserLiked)

	u.notiClient.CreateToOwner(ctx,
		"",
		"",
		[]string{"", fullNames, " thích bình luận của bạn: \"" + _utils.ShortenContent(comment.Content, 50) + "\""},
		_enum.NotificationLikeComment,
		&comment.ID,
		comment.UserID,
		_enum.EOwnerOfMember,
		[]string{fullNames, _utils.ShortenContent(comment.Content, 50)},
	)

	return true, nil
}

func (u *NumberUpdateUsecase) UpdateNumberReplyComment(ctx context.Context,
	comment *domain.Comment,
	parentComment *domain.Comment,
	newsFeed *domain.NewsFeed,
) (bool, error) {
	// update number reply of comment
	totalReply, err := u.commentRepo.CountReply(ctx, parentComment.ID)
	if err != nil {
		return false, err
	}

	ok := u.RequiredTimeUpdate(newsFeed.UpdatedAt)
	if !ok && totalReply > 0 {
		return false, nil
	}

	err = u.commentRepo.UpdateNumberReply(ctx, parentComment.ID, totalReply)
	if err != nil {
		return false, err
	}

	_, err = u.setNumberComment(ctx, newsFeed.ID)
	commentedUserIds, _ := u.commentRepo.GetNearestUserIds(ctx, newsFeed.ID, &parentComment.ID)
	sumUserCommented, _ := u.commentRepo.SumUserComment(ctx, newsFeed.ID, &parentComment.ID)
	fullNames := u.userUsecase.getFullNames(ctx, commentedUserIds, sumUserCommented)

	u.notiClient.CreateToOwner(ctx,
		"",
		"",
		[]string{"", fullNames, " đã trả lời bình luận của bạn: \"" + _utils.ShortenContent(comment.Content, 50) + "\""},
		_enum.NotificationCommentChildren,
		&comment.ID,
		parentComment.UserID,
		_enum.EOwnerOfMember,
		[]string{fullNames, _utils.ShortenContent(comment.Content, 50)},
	)

	return true, err
}

func (u *NumberUpdateUsecase) RequiredTimeUpdate(t *time.Time) bool {
	if time.Now().Before(t.Add(time.Second * 30)) {
		return false
	}
	return true
}

func (u *NumberUpdateUsecase) UpdateNumberComment(ctx context.Context,
	comment *domain.Comment,
	parentComment *domain.Comment,
	newsFeed *domain.NewsFeed,
) (bool, error) {
	ok := u.RequiredTimeUpdate(newsFeed.UpdatedAt)
	if !ok {
		return false, nil
	}

	_, err := u.setNumberComment(ctx, newsFeed.ID)
	commentedUserIds, _ := u.commentRepo.GetNearestUserIds(ctx, newsFeed.ID, nil)
	sumUserCommented, _ := u.commentRepo.SumUserComment(ctx, newsFeed.ID, nil)
	fullNames := u.userUsecase.getFullNames(ctx, commentedUserIds, sumUserCommented)

	u.notiClient.CreateToOwner(ctx,
		"",
		"",
		[]string{"", fullNames, " đã bình luận bài viết của bạn: \"" + _utils.ShortenContent(newsFeed.Content, 50) + "\""},
		_enum.NotificationCommentNewsFeed,
		&newsFeed.ID,
		*newsFeed.CreatedBy,
		_enum.EOwnerOfMember,
		[]string{fullNames, _utils.ShortenContent(newsFeed.Content, 50)},
	)

	return true, err
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

func (u *NumberUpdateUsecase) setNumberComment(ctx context.Context, newsFeedID uint64) (int, error) {
	totalComment, err := u.commentRepo.CountByNewsFeedID(ctx, newsFeedID)
	if err != nil {
		return totalComment, err
	}
	err = u.newsFeedRepo.UpdateCommentNumber(ctx, newsFeedID, totalComment)
	if err != nil {
		return totalComment, err
	}

	return totalComment, nil
}
