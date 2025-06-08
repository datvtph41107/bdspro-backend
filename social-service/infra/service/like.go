package service

import (
	_utils "common/utils"
	"context"
	pb_social "pb/types/social"
	"social/internal/domain"
	"social/internal/usecase"
)

type LikeService struct {
	likeUsecase *usecase.LikeUsecase
	pb_social.UnimplementedLikeServiceServer
}

func NewLikeService(likeUsecase *usecase.LikeUsecase) *LikeService {
	return &LikeService{likeUsecase: likeUsecase}
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param commentId path int true "CommentId"
// @Param body body pb_social.LikeCommentRequest true "LikeCommentRequest"
// @Success 200 {object} pb_social.LikeResponse
// @Router /like/comment/{commentId} [post]
func (s *LikeService) LikeComment(ctx context.Context, req *pb_social.LikeCommentRequest) (*pb_social.LikeResponse, error) {
	like, err := s.likeUsecase.LikeComment(ctx, req.CommentId, req.Dislike)
	if err != nil {
		return nil, err
	}
	return s.DomainLikeToPb(like), nil
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newsFeedId path int true "NewsFeedId"
// @Success 200 {object} pb_social.LikeResponse
// @Router /like/news-feed/{newsFeedId} [post]
func (s *LikeService) LikeNewsFeed(ctx context.Context, req *pb_social.LikeNewsFeedRequest) (*pb_social.LikeResponse, error) {
	like, err := s.likeUsecase.LikeNewsFeed(ctx, req.NewsFeedId)
	if err != nil {
		return nil, err
	}
	return s.DomainLikeToPb(like), nil
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param commentId path int true "CommentId"
// @Success 200 {object} pb_social.LikeResponse
// @Router /like/comment/{commentId} [delete]
func (s *LikeService) UnlikeComment(ctx context.Context, req *pb_social.LikeCommentRequest) (*pb_social.LikeResponse, error) {
	like, err := s.likeUsecase.UnlikeComment(ctx, req.CommentId)
	if err != nil {
		return nil, err
	}
	return s.DomainLikeToPb(like), nil
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newsFeedId path int true "NewsFeedId"
// @Success 200 {object} pb_social.LikeResponse
// @Router /like/news-feed/{newsFeedId} [delete]
func (s *LikeService) UnlikeNewsFeed(ctx context.Context, req *pb_social.LikeNewsFeedRequest) (*pb_social.LikeResponse, error) {
	like, err := s.likeUsecase.UnlikeNewsFeed(ctx, req.NewsFeedId)
	if err != nil {
		return nil, err
	}
	return s.DomainLikeToPb(like), nil
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param commentId path int true "CommentId"
// @Success 200 {object} pb_social.ListLikeResponse
// @Router /like/comment/list/{commentId} [get]
func (s *LikeService) GetListLikeComment(ctx context.Context, req *pb_social.LikeCommentSearch) (*pb_social.ListLikeResponse, error) {
	return nil, nil
}

// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newsFeedId path int true "NewsFeedId"
// @Success 200 {object} pb_social.ListLikeResponse
// @Router /like/news-feed/list/{newsFeedId} [get]
func (s *LikeService) GetListLikeNewsFeed(ctx context.Context, req *pb_social.LikeNewsFeedSearch) (*pb_social.ListLikeResponse, error) {
	return nil, nil
}

// -- mapping --

// func (s *LikeService) PbLikeToDomain(pbLike *pb_social.LikeCommentRequest) *domain.Like {
// 	return &domain.Like{
// 		TargetID:   pbLike.CommentId,
// 		TargetType: enums.TargetTypeComment,
// 		UserID:     _utils.GetProfileIdWithContext(ctx),
// 	}
// }

func (s *LikeService) DomainLikeToPb(like *domain.Like) *pb_social.LikeResponse {
	return &pb_social.LikeResponse{
		TargetId:   like.TargetID,
		TargetType: uint32(like.TargetType),
		UserId:     like.UserID,
		CreatedAt:  _utils.FormatTimeToString(like.CreatedAt),
	}
}
