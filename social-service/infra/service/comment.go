package service

import (
	_dto "common/dto"
	_models "common/models"
	_utils "common/utils"
	"context"
	"errors"
	"social/infra/client"
	"social/internal/domain"
	"social/internal/dto"
	"social/internal/usecase"

	pb_social "pb/types/social"
)

type CommentService struct {
	commentUsecase *usecase.CommentUsecase
	userClient     *client.UserClient
	pb_social.UnimplementedCommentServiceServer
}

func NewCommentService(
	commentUsecase *usecase.CommentUsecase,
	userClient *client.UserClient,
) *CommentService {
	return &CommentService{
		commentUsecase: commentUsecase,
		userClient:     userClient,
	}
}

// @Tags Comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body pb_social.CommentRequest true "CommentRequest"
// @Success 200 {object} pb_social.CommentResponse
// @Router /comment [post]
func (s *CommentService) CreateComment(ctx context.Context, req *pb_social.CommentRequest) (*pb_social.CommentResponse, error) {
	comment := s.PbCommentToDomain(req)
	if comment.Content == "" || comment.NewsFeedID == 0 {
		return nil, errors.New("content and news feed id are required")
	}
	comment, err := s.commentUsecase.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return s.DomainCommentToPb(comment), nil
}

// @Tags Comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body pb_social.CommentRequest true "CommentRequest"
// @Param id path int true "ID"
// @Success 200 {object} pb_social.CommentResponse
// @Router /comment/{id} [put]
func (s *CommentService) UpdateComment(ctx context.Context, req *pb_social.CommentRequest) (*pb_social.CommentResponse, error) {
	comment := s.PbCommentToDomain(req)
	if comment.Content == "" || comment.ID == 0 {
		return nil, errors.New("content, news feed id and id are required")
	}
	comment, err := s.commentUsecase.UpdateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return s.DomainCommentToPb(comment), nil
}

// @Tags Comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment query pb_social.CommentRequest true "CommentRequest"
// @Param id path int true "ID"
// @Success 200 {object} pb_social.CommentResponse
// @Router /comment/{id} [delete]
func (s *CommentService) DeleteComment(ctx context.Context, req *pb_social.CommentRequest) (*pb_social.CommentResponse, error) {
	comment := s.PbCommentToDomain(req)
	if comment.ID == 0 {
		return nil, errors.New("id is required")
	}
	err := s.commentUsecase.DeleteComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return s.DomainCommentToPb(comment), nil
}

// @Tags Comment
// @Accept json
// @Produce json
// @Param newsFeedId path int true "NewsFeedId"
// @Param comment query pb_social.CommentRequest true "CommentRequest"
// @Success 200 {object} pb_social.CommentListResponse
// @Router /comment/list/{newsFeedId} [get]
func (s *CommentService) GetListComment(ctx context.Context, req *pb_social.CommentRequest) (*pb_social.CommentListResponse, error) {
	dto := &dto.CommentRequest{
		Page: _dto.Page{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		NewsFeedID: req.NewsFeedId,
		ParentID:   req.ParentId,
	}
	comments, total, err := s.commentUsecase.GetComment(ctx, dto)
	if err != nil {
		return nil, err
	}
	results := make([]*pb_social.CommentResponse, 0)
	for _, comment := range comments {
		results = append(results, s.DomainCommentToPb(&comment))
	}

	s.userClient.PbUserToComment(ctx, results)
	return &pb_social.CommentListResponse{
		Data:  results,
		Total: uint64(total),
	}, nil
}

// -- mapping --

func (s *CommentService) PbCommentToDomain(pbComment *pb_social.CommentRequest) *domain.Comment {
	return &domain.Comment{
		BaseEntity: _models.BaseEntity{
			ID: pbComment.Id,
		},
		Content:    pbComment.Content,
		ParentID:   pbComment.ParentId,
		NewsFeedID: pbComment.NewsFeedId,
		MediaUrl:   pbComment.MediaUrl,
		MediaType:  pbComment.MediaType,
	}
}

func (s *CommentService) DomainCommentToPb(comment *domain.Comment) *pb_social.CommentResponse {
	return &pb_social.CommentResponse{
		Id:         comment.ID,
		Content:    comment.Content,
		ParentId:   comment.ParentID,
		NewsFeedId: comment.NewsFeedID,
		UpdatedAt:  _utils.FormatTimeToString(comment.UpdatedAt),
		NumLike:    comment.NumLike,
		NumDisLike: comment.NumDisLike,
		NumReply:   comment.NumReply,
		CreatedBy:  comment.CreatedBy,
		MediaUrl:   comment.MediaUrl,
		MediaType:  comment.MediaType,
	}
}
