package usecase

import (
	_utils "common/utils"
	"context"
	"social/internal/domain"
	"social/internal/enums"
	iusecase "social/internal/interface"
	"social/internal/repo"
)

type NewsFeedShareUsecase struct {
	newsFeedUsecase   *NewsFeedUsecase
	newsFeedShareRepo repo.NewsFeedShareRepo
	newsFeedRepo      repo.NewsFeedRepo
	chatClient        iusecase.ChatClient
	notiClient        iusecase.NotiClient
	userClient        iusecase.UserClient
	userUsecase       *UserUsecase
}

func NewNewsFeedShareUsecase(
	newsFeedUsecase *NewsFeedUsecase,
	newsFeedShareRepo repo.NewsFeedShareRepo,
	newsFeedRepo repo.NewsFeedRepo,
	chatClient iusecase.ChatClient,
	notiClient iusecase.NotiClient,
	userClient iusecase.UserClient,
	userUsecase *UserUsecase,
) *NewsFeedShareUsecase {
	return &NewsFeedShareUsecase{
		newsFeedUsecase:   newsFeedUsecase,
		newsFeedRepo:      newsFeedRepo,
		newsFeedShareRepo: newsFeedShareRepo,
		chatClient:        chatClient,
		notiClient:        notiClient,
		userClient:        userClient,
		userUsecase:       userUsecase,
	}
}

func (u *NewsFeedShareUsecase) ShareNewsFeed(
	ctx context.Context,
	newsFeedID uint64,
	targetID uint64,
	content string,
	shareType enums.ShareType,
) (*domain.NewsFeedShare, error) {
	if shareType == enums.ShareTypeNewsFeed {
		return u.ShareToNewsFeed(ctx, newsFeedID)
	} else if shareType == enums.ShareTypeFriend || shareType == enums.ShareTypeGroup || shareType == enums.ShareTypeOrganization {
		return u.ShareToChannel(ctx, content, newsFeedID, targetID, shareType)
	}
	return nil, nil
}

func (u *NewsFeedShareUsecase) ShareToNewsFeed(ctx context.Context, newsFeedID uint64) (*domain.NewsFeedShare, error) {
	newsFeed, err := u.RequiredBeforeShare(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}

	// Tạo bài viết cho người dùng
	_, err = u.newsFeedUsecase.CreateUserNewsFeed(ctx, &domain.NewsFeed{
		Title:      newsFeed.Title,
		Content:    newsFeed.Content,
		Image:      newsFeed.Image,
		Link:       newsFeed.Link,
		Visibility: enums.VisibilityPublic,
		PostID:     newsFeed.PostID,
		ParentID:   &newsFeedID,
	})
	if err != nil {
		return nil, err
	}

	// Tạo bản ghi chia sẻ cho bài viết
	newsFeedShare, err := u.CreateNewsFeedShare(ctx, newsFeed, nil, enums.ShareTypeNewsFeed)
	if err != nil {
		return nil, err
	}

	// Tăng số lượng chia sẻ
	_ = u.newsFeedRepo.IncrementNumberShare(ctx, newsFeedID)

	cloneCtx := _utils.CloneContext(ctx)
	go u.sendNoti(cloneCtx, newsFeed)
	return newsFeedShare, nil
}

func (u *NewsFeedShareUsecase) sendNoti(ctx context.Context, newsFeed *domain.NewsFeed) error {
	// profileId := _utils.GetProfileIdWithContext(ctx)
	// user, err := u.userClient.GetUser(ctx, profileId)
	// if err != nil {
	// 	return err
	// }
	sharedUserIds, _ := u.newsFeedShareRepo.GetNearestUserIds(ctx, newsFeed.ID, enums.TargetTypeNewsFeed)
	sumUserShared, _ := u.newsFeedShareRepo.SumUserComment(ctx, newsFeed.ID, enums.TargetTypeNewsFeed)
	fullNames := u.userUsecase.getFullNames(ctx, sharedUserIds, sumUserShared)

	u.notiClient.NotifyToUser(ctx,
		70,
		newsFeed.CreatedBy,
		[]string{fullNames, _utils.ShortenContent(newsFeed.Content, 50)},
		newsFeed.ID,
		true,
	)
	return nil
}

func (u *NewsFeedShareUsecase) ShareToChannel(
	ctx context.Context,
	content string,
	newsFeedID uint64,
	targetID uint64,
	shareType enums.ShareType,
) (*domain.NewsFeedShare, error) {
	newsFeed, err := u.RequiredBeforeShare(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}

	err = u.chatClient.SendMessage(ctx, targetID, content, newsFeedID)
	if err != nil {
		return nil, err
	}

	newsFeedShare, err := u.CreateNewsFeedShare(ctx, newsFeed, &targetID, shareType)
	if err != nil {
		return nil, err
	}

	err = u.newsFeedRepo.IncrementNumberShare(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}

	return newsFeedShare, nil
}

func (u *NewsFeedShareUsecase) RequiredBeforeShare(ctx context.Context, newsFeedID uint64) (*domain.NewsFeed, error) {
	newsFeed, err := u.newsFeedRepo.GetByID(ctx, newsFeedID)
	if err != nil {
		return nil, err
	}
	if newsFeed.ID == 0 {
		return nil, domain.ErrShareNewsFeedNotFound
	}
	if newsFeed.Visibility != enums.VisibilityPublic {
		return nil, domain.ErrShareNewsFeedNotPublic
	}
	return newsFeed, nil
}

func (u *NewsFeedShareUsecase) CreateNewsFeedShare(
	ctx context.Context,
	newsFeed *domain.NewsFeed,
	targetID *uint64,
	shareType enums.ShareType,
) (*domain.NewsFeedShare, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	newsFeedShare := &domain.NewsFeedShare{
		NewsFeedID: newsFeed.ID,
		TargetID:   targetID,
		UserID:     &profileId,
		ShareType:  shareType,
	}
	_, err := u.newsFeedShareRepo.Create(ctx, newsFeedShare)
	if err != nil {
		return nil, err
	}
	return newsFeedShare, nil
}
