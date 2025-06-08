package service

import (
	_dto "common/dto"
	_models "common/models"
	"context"
	pb_social "pb/types/social"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"
	"social/internal/usecase"
)

type NewsFeedService struct {
	postUsecase *usecase.NewsFeedUsecase
	pb_social.UnimplementedNewsFeedServiceServer
}

func NewNewsFeedService(postUsecase *usecase.NewsFeedUsecase) *NewsFeedService {
	return &NewsFeedService{postUsecase: postUsecase}
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
	newsFeed, err := s.postUsecase.CreateNewsFeed(ctx, newsFeed)
	if err != nil {
		return nil, err
	}
	return s.DomainPostToPb(newsFeed), nil
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
	post, err := s.postUsecase.UpdatePost(ctx, post)
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
	err := s.postUsecase.DeletePost(ctx, post)
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
	posts, total, err := s.postUsecase.GetListPost(ctx, postSearch)
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

// --- mapping ---

func (s *NewsFeedService) PbPostToDomain(pbPost *pb_social.NewsFeed) *domain.NewsFeed {
	return &domain.NewsFeed{
		BaseEntity: _models.BaseEntity{
			ID: pbPost.Id,
		},
		Title:         pbPost.Title,
		Content:       pbPost.Content,
		Image:         pbPost.Image,
		Link:          pbPost.Link,
		Visibility:    enums.Visibility(pbPost.Visibility),
		LikeNumber:    int(pbPost.LikeNumber),
		CommentNumber: int(pbPost.CommentNumber),
		ShareNumber:   int(pbPost.ShareNumber),
		ViewNumber:    int(pbPost.ViewNumber),
	}
}

func (s *NewsFeedService) DomainPostToPb(post *domain.NewsFeed) *pb_social.NewsFeed {
	return &pb_social.NewsFeed{
		Id:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Image:         post.Image,
		Link:          post.Link,
		Visibility:    int32(post.Visibility),
		LikeNumber:    int32(post.LikeNumber),
		CommentNumber: int32(post.CommentNumber),
		ShareNumber:   int32(post.ShareNumber),
		ViewNumber:    int32(post.ViewNumber),
	}
}
