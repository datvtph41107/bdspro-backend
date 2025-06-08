package service

import (
	_dto "common/dto"
	_models "common/models"
	_routes "common/routes"
	"context"
	"net/http"
	pb_social "pb/types/social"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"social/internal/usecase"
)

type NewsFeedService struct {
	newsFeedUsecase      *usecase.NewsFeedUsecase
	newsFeedShareUsecase *usecase.NewsFeedShareUsecase
	pb_social.UnimplementedNewsFeedServiceServer
}

func NewNewsFeedService(
	newsFeedUsecase *usecase.NewsFeedUsecase,
	newsFeedShareUsecase *usecase.NewsFeedShareUsecase,
) *NewsFeedService {
	return &NewsFeedService{
		newsFeedUsecase:      newsFeedUsecase,
		newsFeedShareUsecase: newsFeedShareUsecase,
	}
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
func (s *NewsFeedService) CreateNewsFeed(ctx context.Context, req *pb_social.NewsFeed) (*pb_social.NewsFeed, error) {
	newsFeed := s.PbPostToDomain(req)
	result, err := s.newsFeedUsecase.CreateNewsFeed(ctx, newsFeed)
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
// @Param post query pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.NewsFeed
// @Router /news-feed [put]
func (s *NewsFeedService) UpdateNewsFeed(ctx context.Context, req *pb_social.NewsFeed) (*pb_social.NewsFeed, error) {
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
// @Param post body pb_social.NewsFeed true "NewsFeed"
// @Success 200 {object} pb_social.Empty
// @Router /news-feed [delete]
func (s *NewsFeedService) DeleteNewsFeed(ctx context.Context, req *pb_social.NewsFeed) (*pb_social.Empty, error) {
	post := s.PbPostToDomain(req)
	err := s.newsFeedUsecase.DeletePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return &pb_social.Empty{}, nil
}

// @Summary Lấy danh sách news feed
// @Description Get global news feed
// @Tags News Feed
// @Accept json
// @Produce json
// @Param post query pb_social.NewsFeedGlobalSearch true "NewsFeedGlobalSearch"
// @Success 200 {object} pb_social.ListNewsFeed
// @Router /news-feed/public [get]
func (s *NewsFeedService) GetGlobalNewsFeed(ctx context.Context, req *pb_social.NewsFeedGlobalSearch) (*pb_social.ListNewsFeed, error) {
	postSearch := &dto.NewsFeedGlobalSearch{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Text: req.Text,
	}
	posts, total, err := s.newsFeedUsecase.GetListPost(ctx, postSearch)
	if err != nil {
		return nil, err
	}

	result := &pb_social.ListNewsFeed{
		Data:  []*pb_social.NewsFeed{},
		Total: &total,
	}
	for _, post := range posts {
		result.Data = append(result.Data, s.DomainPostToPb(post))
	}
	return result, nil
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
func (s *NewsFeedService) ShareNewsFeed(ctx context.Context, req *pb_social.ShareRequest) (*pb_social.ShareResponse, error) {
	if req.NewsFeedId == 0 || req.ShareType == 0 {
		return nil, &_routes.Except{
			Code:    http.StatusBadRequest,
			Message: "newsFeedId or shareType is required",
		}
	}
	newsFeedShare, err := s.newsFeedShareUsecase.ShareNewsFeed(ctx,
		req.NewsFeedId,
		req.TargetId,
		req.Content,
		enums.ShareType(req.ShareType),
	)
	if err != nil {
		return nil, err
	}
	return &pb_social.ShareResponse{
		Id:      newsFeedShare.ID,
		Success: true,
	}, nil
}

// --- mapping ---

func (s *NewsFeedService) PbPostToDomain(pbPost *pb_social.NewsFeed) *domain.NewsFeed {
	medias := make([]*domain.NewsFeedMedia, len(pbPost.MediaFeeds))
	for i, media := range pbPost.MediaFeeds {
		medias[i] = &domain.NewsFeedMedia{
			NewsFeedID:  &pbPost.Id,
			Url:         media.Url,
			Type:        media.Type,
			OrderNumber: int(media.OrderNumber),
		}
	}
	return &domain.NewsFeed{
		BaseEntity: _models.BaseEntity{
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
	}
}

func (s *NewsFeedService) DomainPostToPb(post *domain.NewsFeed) *pb_social.NewsFeed {
	return &pb_social.NewsFeed{
		Id:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Image:      post.Image,
		Link:       post.Link,
		Visibility: int32(post.Visibility),
		NumLike:    int32(post.NumLike),
		NumComment: int32(post.NumComment),
		NumShare:   int32(post.NumShare),
		NumView:    int32(post.NumView),
	}
}
