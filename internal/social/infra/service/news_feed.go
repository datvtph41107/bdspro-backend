package service

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_models "common/models"
	_utils "common/utils"
	"context"
	stderrors "errors"
	sharepb "pb/types/shared"
	socialpb "pb/types/social"
	"social/internal"

	"social/infra/client"
	"social/infra/mapper"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"social/internal/usecase"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NewsFeedService struct {
	newsFeedUsecase      *usecase.NewsFeedUsecase
	newsFeedShareUsecase *usecase.NewsFeedShareUsecase
	bdsproClient         *client.BdsproClient
	userClient           *client.UserClient
	newsFeedMapper       *mapper.NewsFeedMapper
	socialpb.UnimplementedNewsFeedServiceServer
}

func NewNewsFeedService(
	newsFeedUsecase *usecase.NewsFeedUsecase,
	newsFeedShareUsecase *usecase.NewsFeedShareUsecase,
	bdsproClient *client.BdsproClient,
	userClient *client.UserClient,
	mapper *mapper.NewsFeedMapper,
) *NewsFeedService {
	return &NewsFeedService{
		newsFeedUsecase:      newsFeedUsecase,
		newsFeedShareUsecase: newsFeedShareUsecase,
		bdsproClient:         bdsproClient,
		userClient:           userClient,
		newsFeedMapper:       mapper,
	}
}

func (s *NewsFeedService) RequiredOwner(ctx context.Context, ownerOf socialpb.OwnerOf, ownerId *uint64) error {
	if ownerOf != socialpb.OwnerOf_group && ownerOf != socialpb.OwnerOf_user && ownerOf != socialpb.OwnerOf_organization {
		return _errors.ReturnError(service.OwnerTypeInvalid)
	}
	if ownerOf == socialpb.OwnerOf_group && ownerId == nil {
		return _errors.ReturnError(service.GroupIDRequired)
	}
	if ownerOf == socialpb.OwnerOf_organization && ownerId == nil {
		return _errors.ReturnError(service.OrganizationIDRequired)
	}
	// if ownerOf == socialpb.OwnerOf_user && ownerId == nil {
	// 	return _errors.ReturnError(service.UserIDRequired)
	// }
	return nil
}

// @Summary Create a new news feed
// @Description Create a new news feed
// @Tags News Feed
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed [post]
func (s *NewsFeedService) CreateNewsFeed(ctx context.Context, req *socialpb.NewsFeed) (*socialpb.NewsFeed, error) {
	if err := s.RequiredOwner(ctx, req.OwnerOf, req.GroupId); err != nil {
		return nil, err
	}

	newsFeed := s.PbPostToDomain(req)
	if req.OwnerOf == socialpb.OwnerOf_user {
		result, err := s.newsFeedUsecase.CreateUserNewsFeed(ctx, newsFeed)
		if err != nil {
			return nil, err
		}
		return s.DomainPostToPb(result), nil
	}
	result, err := s.newsFeedUsecase.CreateGroupNewsFeed(ctx, newsFeed, req.GroupId)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(result), nil
}

// @Summary Get news feed by owner
// @Description Get news feed by owner
// @Tags News Feed
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ownerOf path string true "OwnerOf"
// @Param ownerId path uint64 true "OwnerId"
// @Param query query pb_social.NewsFeedGlobalSearch true "NewsFeedGlobalSearch"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/{ownerOf}/{ownerId} [get]
func (s *NewsFeedService) GetNewsFeedByOwner(ctx context.Context, req *socialpb.NewsFeedSearch) (*socialpb.ListNewsFeed, error) {
	if err := s.RequiredOwner(ctx, req.OwnerOf, &req.OwnerId); err != nil {
		return nil, err
	}

	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:   req.Text,
		IsReel: req.IsReel,
		IsPost: req.IsPost,
	}

	var posts []*dto.NewsFeedPublic
	var total int32
	var err error

	if req.OwnerOf == socialpb.OwnerOf_group {
		posts, total, err = s.newsFeedUsecase.GetListPostGroup(ctx, req.OwnerId, postSearch)
		if err != nil {
			return nil, err
		}
	} else if req.OwnerOf == socialpb.OwnerOf_user {
		posts, total, err = s.newsFeedUsecase.GetListPostUser(ctx, postSearch)
		if err != nil {
			return nil, err
		}
	} else if req.OwnerOf == socialpb.OwnerOf_organization {
		posts, total, err = s.newsFeedUsecase.GetListPostOrganization(ctx, req.OwnerId, postSearch)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, _errors.ReturnError(service.OwnerTypeInvalid)
	}

	return s.responseGlobalNewsFeed(ctx, posts, total)
	// data := s.newsFeedMapper.DomainToNewsFeedPbList(ctx, posts)
	// // s.userClient.PbUserToNewsFeed(ctx, data)

	// g, ctx := errgroup.WithContext(ctx)
	// g.Go(func() error {
	// 	s.bdsproClient.PbPostToDomain(ctx, data)
	// 	return nil
	// })

	// g.Go(func() error {
	// 	s.userClient.PbUserToNewsFeed(ctx, data)
	// 	return nil
	// })

	// if err := g.Wait(); err != nil {
	// 	return nil, err
	// }

	// return &socialpb.ListNewsFeed{
	// 	Data:  data,
	// 	Total: &total,
	// }, nil
}

// @Summary Create a new news feed
// @Description Create a new group news feed
// @Tags News Feed
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed/group [post]
func (s *NewsFeedService) CreateGroupNewsFeed(ctx context.Context, req *socialpb.NewsFeed) (*socialpb.NewsFeed, error) {
	newsFeed := s.PbPostToDomain(req)
	result, err := s.newsFeedUsecase.CreateGroupNewsFeed(ctx, newsFeed, req.GroupId)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(result), nil
}

// @Summary Update news feed
// @Description Update news feed
// @Security BearerAuth
// @Tags News Feed
// @Accept json
// @Produce json
// @Param id path uint64 true "NewsFeedId"
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed/{id} [put]
func (s *NewsFeedService) UpdateNewsFeed(ctx context.Context, req *socialpb.NewsFeed) (*socialpb.NewsFeed, error) {
	post := s.PbPostToDomain(req)
	post, err := s.newsFeedUsecase.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(post), nil
}

// @Summary Xóa news feed
// @Description Xóa news feed
// @Security BearerAuth
// @Tags News Feed
// @Accept json
// @Produce json
// @Param id path uint64 true "NewsFeedId"
// @Router /news-feed/{id} [delete]
func (s *NewsFeedService) DeleteNewsFeed(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := s.newsFeedUsecase.DeletePost(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// mapShareNewsFeedError keeps the established Social share wire contract while
// moving business-state ownership out of the usecase. The Unknown code and
// exact legacy messages are compatibility behavior, not domain semantics.
func mapShareNewsFeedError(err error) error {
	switch {
	case stderrors.Is(err, domain.ErrShareNewsFeedNotFound):
		return status.Error(codes.Unknown, "404: Bài viết không tồn tại")
	case stderrors.Is(err, domain.ErrShareNewsFeedNotPublic):
		return status.Error(codes.Unknown, "403: Bài viết không thể chia sẻ do không phải công khai")
	default:
		return err
	}
}

// @Summary Share news feed
// @Description Share news feed
// @Security BearerAuth
// @Tags News Feed
// @Accept json
// @Produce json
// @Param newsFeedId path uint64 true "NewsFeedId"
// @Param shareRequest body pb_social.ShareRequest true "ShareRequest"
// @Success 200 {object} pb_social.ShareResponse
// @Router /news-feed/{newsFeedId}/share [post]
func (s *NewsFeedService) ShareNewsFeed(ctx context.Context, req *socialpb.ShareRequest) (*socialpb.ShareResponse, error) {
	if req.NewsFeedId == 0 || req.ShareType == 0 {
		return nil, _errors.ReturnError(service.ShareFieldsRequired)
	}
	newsFeedShare, err := s.newsFeedShareUsecase.ShareNewsFeed(ctx,
		req.NewsFeedId,
		req.TargetId,
		req.Content,
		enums.ShareType(req.ShareType),
	)
	if err != nil {
		return nil, mapShareNewsFeedError(err)
	}
	return &socialpb.ShareResponse{
		Id:      newsFeedShare.ID,
		Success: true,
	}, nil
}

// @Summary Get news feed by id
// @Description Get news feed by id
// @Tags News Feed
// @Accept json
// @Produce json
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed/global/reel [get]
func (s *NewsFeedService) GetReelNewsFeed(ctx context.Context, req *sharepb.IdRequest) (*socialpb.ListReelResponse, error) {
	reqPagable := _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}
	newsFeed, count, err := s.newsFeedUsecase.GetReelNewsFeed(ctx, reqPagable)
	if err != nil {
		return nil, err
	}

	result := s.newsFeedMapper.DomainToReelResponseList(ctx, newsFeed)
	s.userClient.PbUserToReels(ctx, result)

	return &socialpb.ListReelResponse{
		Data:  result,
		Total: int32(count),
	}, nil
}

// @Summary Lấy danh sách news feed theo userID
// @Description Get news feed by userID
// @Tags News Feed
// @Accept json
// @Produce json
// @Param userId path uint64 true "UserId"
// @Security BearerAuth
// @Param post query pb_social.NewsFeedGlobalSearch true "NewsFeedGlobalSearch"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/user/{userId} [get]
func (s *NewsFeedService) GetNewsFeedByUserID(ctx context.Context, req *socialpb.NewsFeedSearch) (*socialpb.ListNewsFeed, error) {
	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:   req.Text,
		IsReel: req.IsReel,
	}
	posts, total, err := s.newsFeedUsecase.GetListPostByUserID(ctx, req.UserId, postSearch)
	if err != nil {
		return nil, err
	}

	result := &socialpb.ListNewsFeed{
		Data:  []*socialpb.NewsFeed{},
		Total: &total,
	}
	for _, post := range posts {
		result.Data = append(result.Data, s.newsFeedMapper.DomainPublicToPb(post))
	}

	s.MappingPostAndUser(ctx, result.Data)
	return result, nil
}

// @Summary Lấy chi tiết news feed
// @Description Get news feed by id
// @Tags News Feed
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newsFeedId path uint64 true "NewsFeedId"
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed/detail/{newsFeedId} [get]
func (s *NewsFeedService) DetailNewsFeed(ctx context.Context, req *socialpb.DetailRequest) (*socialpb.NewsFeed, error) {
	newsFeed, err := s.newsFeedUsecase.GetNewsFeedByID(ctx, req.NewsFeedId)
	if err != nil {
		return nil, err
	}

	newsFeeds := []*socialpb.NewsFeed{s.DomainPostToPb(newsFeed)}
	err = s.MappingPostAndUser(ctx, newsFeeds)
	if err != nil {
		return nil, err
	}

	return newsFeeds[0], nil
}

// @Summary Update visibility of news feed
// @Description Update visibility of news feed
// @Tags News Feed
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newsFeedId path uint64 true "NewsFeedId"
// @Param visibility body pb_social.VisibilityRequest true "VisibilityRequest"
// @Router /news-feed/{newsFeedId}/visibility [put]
func (s *NewsFeedService) UpdateVisibility(ctx context.Context, req *socialpb.VisibilityRequest) (*socialpb.NewsFeed, error) {
	newsFeed, err := s.newsFeedUsecase.UpdateVisibility(ctx, req.NewsFeedId, enums.Visibility(req.Visibility))
	if err != nil {
		return nil, err
	}

	newsFeeds := []*socialpb.NewsFeed{s.DomainPostToPb(newsFeed)}
	err = s.MappingPostAndUser(ctx, newsFeeds)
	if err != nil {
		return nil, err
	}
	return newsFeeds[0], nil
}

func (s *NewsFeedService) MappingPostAndUser(ctx context.Context, req []*socialpb.NewsFeed) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.bdsproClient.PbPostToDomain(ctx, req)
		return nil
	})

	g.Go(func() error {
		s.userClient.PbUserToNewsFeed(ctx, req)
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// --- mapping ---

func (s *NewsFeedService) PbPostToDomain(pbPost *socialpb.NewsFeed) *domain.NewsFeed {
	medias := make([]domain.NewsFeedMedia, len(pbPost.MediaFeeds))
	for i, media := range pbPost.MediaFeeds {
		medias[i] = domain.NewsFeedMedia{
			NewsFeedID:  &pbPost.Id,
			Url:         media.Url,
			Type:        media.Type,
			OrderNumber: int(media.OrderNumber),
		}
	}
	return &domain.NewsFeed{
		BaseEntity: _models.BaseEntity{
			AuditBase: _models.AuditBase{
				CreatedBy: pbPost.CreatedBy,
			},
			ID: pbPost.Id,
		},
		Title:          pbPost.Title,
		Content:        pbPost.Content,
		Image:          pbPost.Image,
		Link:           pbPost.Link,
		Visibility:     enums.Visibility(pbPost.Visibility),
		NumLike:        int(pbPost.NumLike),
		NumComment:     int(pbPost.NumComment),
		NumShare:       int(pbPost.NumShare),
		NumView:        int(pbPost.NumView),
		NewsFeedMedias: medias,
		FriendTagIds:   pbPost.FriendTagIds,
		IsReel:         pbPost.IsReel,
		PostID:         pbPost.PostId,
	}
}

func (s *NewsFeedService) DomainPostToPb(post *domain.NewsFeed) *socialpb.NewsFeed {
	medias := make([]*socialpb.NewsFeedMedia, 0)
	for _, media := range post.NewsFeedMedias {
		medias = append(medias, &socialpb.NewsFeedMedia{
			Url:         media.Url,
			Type:        media.Type,
			OrderNumber: int32(media.OrderNumber),
		})
	}
	return &socialpb.NewsFeed{
		Id:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		Image:        post.Image,
		Link:         post.Link,
		Visibility:   int32(post.Visibility),
		NumLike:      int32(post.NumLike),
		NumComment:   int32(post.NumComment),
		NumShare:     int32(post.NumShare),
		NumView:      int32(post.NumView),
		FriendTagIds: post.FriendTagIds,
		IsReel:       post.IsReel,
		PostId:       post.PostID,
		CreatedBy:    post.CreatedBy,
		MediaFeeds:   medias,
		UpdatedAt:    _utils.FormatTimeToString(post.UpdatedAt),
		LikedAt:      _utils.FormatTimeToString(post.LikedAt),
		DisLike:      post.DisLike,
	}
}

// @Summary Lấy danh sách news feed
// @Description Get global news feed
// @Tags News Feed
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post query pb_social.NewsFeedGlobalSearch true "NewsFeedGlobalSearch"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/global [get]
func (s *NewsFeedService) GetGlobalNewsFeed(ctx context.Context, req *socialpb.NewsFeedGlobalSearch) (*socialpb.ListNewsFeed, error) {
	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:     req.Text,
		IsReel:   req.IsReel,
		IsPost:   req.IsPost,
		IsGlobal: true,
	}

	posts, total, err := s.newsFeedUsecase.GetListPost(ctx, postSearch)
	if err != nil {
		return nil, err
	}

	return s.responseGlobalNewsFeed(ctx, posts, total)
}

// @Summary Lấy danh sách news feed
// @Description Get global news feed
// @Tags News Feed
// @Accept json
// @Produce json
// @Param post query pb_social.NewsFeedGlobalSearch true "NewsFeedGlobalSearch"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/global/{ownerOf}/{ownerId} [get]
func (s *NewsFeedService) GetGlobalNewsFeedByOwner(ctx context.Context, req *socialpb.NewsFeedGlobalSearch) (*socialpb.ListNewsFeed, error) {
	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:   req.Text,
		IsReel: req.IsReel,
		IsPost: req.IsPost,
	}
	posts, total, err := s.newsFeedUsecase.GetListPost(ctx, postSearch)
	if err != nil {
		return nil, err
	}

	return s.responseGlobalNewsFeed(ctx, posts, total)
}

func (s *NewsFeedService) getNewsFeedByOwner(ctx context.Context, req *socialpb.NewsFeedGlobalSearch, isGlobal bool) (*socialpb.ListNewsFeed, error) {
	if err := s.RequiredOwner(ctx, req.OwnerOf, &req.OwnerId); err != nil {
		return nil, err
	}

	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:     req.Text,
		IsReel:   req.IsReel,
		IsPost:   req.IsPost,
		IsGlobal: isGlobal,
	}

	var posts []*dto.NewsFeedPublic
	var total int32
	var err error

	if req.OwnerOf == socialpb.OwnerOf_group {
		posts, total, err = s.newsFeedUsecase.GetListPostGroup(ctx, req.OwnerId, postSearch)
		if err != nil {
			return nil, err
		}
	} else if req.OwnerOf == socialpb.OwnerOf_user {
		posts, total, err = s.newsFeedUsecase.GetListPostUser(ctx, postSearch)
		if err != nil {
			return nil, err
		}
	} else if req.OwnerOf == socialpb.OwnerOf_organization {
		posts, total, err = s.newsFeedUsecase.GetListPostOrganization(ctx, req.OwnerId, postSearch)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, _errors.ReturnError(service.OwnerTypeInvalid)
	}

	// data := s.newsFeedMapper.DomainToNewsFeedPbList(ctx, posts)
	// // s.userClient.PbUserToNewsFeed(ctx, data)

	// g, ctx := errgroup.WithContext(ctx)
	// g.Go(func() error {
	// 	s.bdsproClient.PbPostToDomain(ctx, data)
	// 	return nil
	// })

	// g.Go(func() error {
	// 	s.userClient.PbUserToNewsFeed(ctx, data)
	// 	return nil
	// })

	// if err := g.Wait(); err != nil {
	// 	return nil, err
	// }

	// return &socialpb.ListNewsFeed{
	// 	Data:  data,
	// 	Total: &total,
	// }, nil
	return s.responseGlobalNewsFeed(ctx, posts, total)
}

func (s *NewsFeedService) responseGlobalNewsFeed(ctx context.Context, newsFeeds []*dto.NewsFeedPublic, total int32) (*socialpb.ListNewsFeed, error) {
	result := &socialpb.ListNewsFeed{
		Data:  []*socialpb.NewsFeed{},
		Total: &total,
	}
	for _, newsFeed := range newsFeeds {
		result.Data = append(result.Data, s.newsFeedMapper.DomainPublicToPb(newsFeed))
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.bdsproClient.PbPostToDomain(ctx, result.Data)
		return nil
	})

	g.Go(func() error {
		s.userClient.PbUserToNewsFeed(ctx, result.Data)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

// @Summary Lấy danh sách news feed của chính mình
// @Description Get my news feed (tất cả bài viết của user đang đăng nhập)
// @Tags News Feed
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int32 false "Page number" default(1)
// @Param size query int32 false "Page size" default(10)
// @Param text query string false "Search text"
// @Param isReel query bool false "Filter reel"
// @Param isPost query bool false "Filter post"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/me [get]
func (s *NewsFeedService) GetMyNewsFeed(ctx context.Context, req *socialpb.MyNewsFeedRequest) (*socialpb.ListNewsFeed, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(service.Unauthenticated)
	}

	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:   req.Text,
		IsReel: req.IsReel,
		IsPost: req.IsPost,
	}

	// Lấy tất cả bài viết của user (không lọc visibility)
	posts, total, err := s.newsFeedUsecase.GetListPostUser(ctx, postSearch)
	if err != nil {
		return nil, err
	}

	return s.responseGlobalNewsFeed(ctx, posts, total)
}

// @Summary Lấy danh sách bài viết của user (Public API)
// @Description Lấy danh sách bài viết công khai của một user cụ thể - API public không cần authentication
// @Tags News Feed
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param page query int32 false "Page number" default(1)
// @Param size query int32 false "Page size" default(10)
// @Param text query string false "Search text"
// @Param isReel query bool false "Filter reel"
// @Param isPost query bool false "Filter post"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /v2/social/public/news-feed/user/{userId} [get]
func (s *NewsFeedService) GetUserNewsFeedPublic(ctx context.Context, req *socialpb.PublicNewsFeedSearch) (*socialpb.ListNewsFeed, error) {
	postSearch := &dto.NewsFeedGlobalSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text:   req.Text,
		IsReel: req.IsReel,
		IsPost: req.IsPost,
	}

	// Lấy bài viết public của user
	posts, total, err := s.newsFeedUsecase.GetListPostByUserID(ctx, req.UserId, postSearch)
	if err != nil {
		return nil, err
	}

	result := &socialpb.ListNewsFeed{
		Data:  []*socialpb.NewsFeed{},
		Total: &total,
	}
	for _, post := range posts {
		result.Data = append(result.Data, s.newsFeedMapper.DomainPublicToPb(post))
	}

	s.MappingPostAndUser(ctx, result.Data)
	return result, nil
}

// CountByOwner đếm số tin đăng theo ownerOf và ownerId
func (s *NewsFeedService) CountByOwner(ctx context.Context, req *socialpb.CountByOwnerRequest) (*socialpb.CountByOwnerResponse, error) {
	count, err := s.newsFeedUsecase.CountByOwner(ctx, req.OwnerOf, req.OwnerId)
	if err != nil {
		return nil, err
	}

	return &socialpb.CountByOwnerResponse{
		Count: count,
	}, nil
}

// @Summary Admin tạo bài viết
// @Description Admin tạo bài viết trong mạng xã hội
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /admin/news-feed [post]
func (s *NewsFeedService) AdminCreateNewsFeed(ctx context.Context, req *socialpb.NewsFeed) (*socialpb.NewsFeed, error) {
	if err := s.RequiredOwner(ctx, req.OwnerOf, req.GroupId); err != nil {
		return nil, err
	}

	newsFeed := s.PbPostToDomain(req)
	result, err := s.newsFeedUsecase.AdminCreateNewsFeed(ctx, newsFeed)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(result), nil
}

// @Summary Admin sửa bài viết
// @Description Admin sửa bài viết trong mạng xã hội
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint64 true "NewsFeedId"
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /admin/news-feed/{id} [put]
func (s *NewsFeedService) AdminUpdateNewsFeed(ctx context.Context, req *socialpb.NewsFeed) (*socialpb.NewsFeed, error) {
	post := s.PbPostToDomain(req)
	post, err := s.newsFeedUsecase.AdminUpdateNewsFeed(ctx, post)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(post), nil
}

// @Summary Admin xóa bài viết
// @Description Admin xóa bài viết trong mạng xã hội
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint64 true "NewsFeedId"
// @Success 200 {object} sharepb.Empty
// @Router /admin/news-feed/{id} [delete]
func (s *NewsFeedService) AdminDeleteNewsFeed(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := s.newsFeedUsecase.AdminDeleteNewsFeed(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// @Summary Admin ẩn bài viết
// @Description Admin ẩn/hiện bài viết trong mạng xã hội
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param newsFeedId path uint64 true "NewsFeedId"
// @Param hideRequest body pb_social.AdminHideRequest true "AdminHideRequest"
// @Success 200 {object} pb_social.NewsFeed
// @Router /admin/news-feed/{newsFeedId}/hide [put]
func (s *NewsFeedService) AdminHideNewsFeed(ctx context.Context, req *socialpb.AdminHideRequest) (*socialpb.NewsFeed, error) {
	newsFeed, err := s.newsFeedUsecase.AdminHideNewsFeed(ctx, req.NewsFeedId, req.IsHidden)
	if err != nil {
		return nil, err
	}

	newsFeeds := []*socialpb.NewsFeed{s.DomainPostToPb(newsFeed)}
	err = s.MappingPostAndUser(ctx, newsFeeds)
	if err != nil {
		return nil, err
	}
	return newsFeeds[0], nil
}

// @Summary Admin lấy danh sách bài viết
// @Description Admin lấy danh sách bài viết trong mạng xã hội (có thể xem tất cả bài viết kể cả đã ẩn)
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int32 false "Page number" default(1)
// @Param size query int32 false "Page size" default(20)
// @Param text query string false "Search text"
// @Param isReel query bool false "Filter reel"
// @Param isPost query bool false "Filter post"
// @Param ownerId query uint64 false "Filter by owner ID"
// @Param ownerOf query int false "Filter by owner type (10: user, 20: group, 30: organization)"
// @Param visibility query int false "Filter by visibility"
// @Param isHidden query bool false "Filter by hidden status (removed_at)"
// @Param createdBy query uint64 false "Filter by creator ID"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /admin/news-feed [get]
func (s *NewsFeedService) AdminGetListNewsFeed(ctx context.Context, req *socialpb.AdminNewsFeedSearch) (*socialpb.ListNewsFeed, error) {
	// Convert proto request to DTO
	adminSearch := &dto.AdminNewsFeedSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Text: req.Text,
	}

	if req.IsReel != nil {
		adminSearch.IsReel = req.IsReel
	}
	if req.IsPost != nil {
		adminSearch.IsPost = req.IsPost
	}
	if req.OwnerId != nil {
		adminSearch.OwnerID = req.OwnerId
	}
	if req.OwnerOf != nil {
		ownerOfInt := int(*req.OwnerOf)
		adminSearch.OwnerOf = &ownerOfInt
	}
	if req.Visibility != nil {
		visibilityInt := int(*req.Visibility)
		adminSearch.Visibility = &visibilityInt
	}
	if req.IsHidden != nil {
		adminSearch.IsHidden = req.IsHidden
	}
	if req.CreatedBy != nil {
		adminSearch.CreatedBy = req.CreatedBy
	}

	posts, total, err := s.newsFeedUsecase.AdminGetListNewsFeed(ctx, adminSearch)
	if err != nil {
		return nil, err
	}

	return s.responseGlobalNewsFeed(ctx, posts, total)
}

// @Summary Admin lấy chi tiết bài viết
// @Description Admin lấy chi tiết bài viết trong mạng xã hội (có thể xem cả bài đã ẩn hoặc đã xóa)
// @Tags News Feed Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param newsFeedId path uint64 true "NewsFeedId"
// @Success 200 {object} pb_social.NewsFeed
// @Router /admin/news-feed/detail/{newsFeedId} [get]
func (s *NewsFeedService) AdminDetailNewsFeed(ctx context.Context, req *socialpb.DetailRequest) (*socialpb.NewsFeed, error) {
	newsFeed, err := s.newsFeedUsecase.AdminGetNewsFeedByID(ctx, req.NewsFeedId)
	if err != nil {
		return nil, err
	}

	newsFeeds := []*socialpb.NewsFeed{s.DomainPostToPb(newsFeed)}
	err = s.MappingPostAndUser(ctx, newsFeeds)
	if err != nil {
		return nil, err
	}

	return newsFeeds[0], nil
}
