package service

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type PostService struct {
	pb_bdspro.UnimplementedPostServiceServer
	postRepo   repo.UPostRepo
	UC         *shared_usecase.PostUsecase
	PostMapper *mapper.PostMapper
}

func NewGrpcPostService(postRepo repo.UPostRepo,
	uc *shared_usecase.PostUsecase,
	postMapper *mapper.PostMapper,
) *PostService {
	return &PostService{
		postRepo:   postRepo,
		UC:         uc,
		PostMapper: postMapper,
	}
}

func (s *PostService) OwnerPost(ctx context.Context, req *pb_bdspro.OwnerPostRequest) (*pb_bdspro.OwnerPostResponse, error) {
	success, err := s.postRepo.OwnerPost(ctx, req.ProfileId, req.PostId)
	return &pb_bdspro.OwnerPostResponse{Success: success}, err
}

func (s *PostService) GetPostByID(ctx context.Context, req *pb_bdspro.GetPostByIDRequest) (*pb_bdspro.Post, error) {
	post, err := s.postRepo.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	mediaList := make([]*pb_bdspro.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = &pb_bdspro.PostMedia{
			Id:        media.ID,
			MediaUrl:  media.MediaURL,
			MediaType: media.MediaType,
			IsMain:    media.IsMain,
			SortOrder: int32(media.SortOrder),
		}
	}
	return &pb_bdspro.Post{
		Id:             post.ID,
		Title:          post.Title,
		PostPrice:      post.Price,
		PriceType:      uint32(post.PriceType),
		PackageVisible: uint32(post.PackageVisible),
		Like:           post.Like,
		Comment:        post.Comment,
		NumView:        post.NumView,
		OwnerId:        post.OwnerID,
		OwnerType:      uint32(post.OwnerType),
		MediaList:      mediaList,
	}, nil
}

func (s *PostService) GetPostByIDs(ctx context.Context, req *pb_bdspro.GetPostByIDsRequest) (*pb_bdspro.GetPostByIDsResponse, error) {
	posts, err := s.postRepo.GetByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	postList := s.PostMapper.PostToPbList(posts)
	return &pb_bdspro.GetPostByIDsResponse{Posts: postList}, nil
}

func (s *PostService) CreatePost(ctx context.Context, req *pb_bdspro.PostSaveRequest) (*pb_bdspro.Post, error) {
	body := dto.PostSaveRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.CreatePost(ctx, &body)
	if err != nil {
		return nil, err
	}

	post := &pb_bdspro.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) NewExpiredPost(ctx context.Context, req *pb_bdspro.PostExpiredRequest) (*pb_bdspro.Post, error) {
	body := dto.PostExpiredRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.MakeNewExpired(ctx, body)
	if err != nil {
		return nil, err
	}

	post := &pb_bdspro.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) GetDetail(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.PostDetailResponse, error) {
	result, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	post := &pb_bdspro.PostDetailResponse{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) UpdatePost(ctx context.Context, req *pb_bdspro.PostSaveRequest) (*pb_bdspro.Post, error) {
	body := dto.PostUpdateRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.UpdatePost(ctx, req.Id, body)
	if err != nil {
		return nil, err
	}

	post := &pb_bdspro.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) Delete(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}
func (s *PostService) UpdateHidden(ctx context.Context, req *pb_bdspro.UpdateHiddenRequest) (*pb_bdspro.Response, error) {
	body := dto.UpdateHiddenRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	_, err := s.UC.UpdateHidden(ctx, body)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}
func (s *PostService) GetPersonalPost(ctx context.Context, req *pb_bdspro.PostSearchRequest) (*pb_bdspro.PostListResponse, error) {
	query := dto.PostSearchRequest{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.PersonalPost(ctx, query)
	if err != nil {
		return nil, err
	}

	postList := make([]*pb_bdspro.Post, len(result))
	copier.Copy(&postList, &result)

	return &pb_bdspro.PostListResponse{
		Data:  postList,
		Total: uint64(total),
	}, nil
}
func (s *PostService) GetGlobalPost(ctx context.Context, req *pb_bdspro.PostSearchRequest) (*pb_bdspro.PostListResponse, error) {
	query := dto.PostSearchRequest{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GlobalPost(ctx, query)
	if err != nil {
		return nil, err
	}

	postList := make([]*pb_bdspro.Post, len(result))
	copier.Copy(&postList, &result)

	return &pb_bdspro.PostListResponse{
		Data:  postList,
		Total: uint64(total),
	}, nil
}
func (s *PostService) GetPublish(ctx context.Context, req *pb_bdspro.PostPublishSearch) (*pb_bdspro.PostListResponse, error) {
	query := dto.PostPublishSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GetPublishByProfileID(ctx, req.ProfileId, query)
	postList := make([]*pb_bdspro.Post, len(result))
	copier.Copy(&postList, &result)

	return &pb_bdspro.PostListResponse{
		Data:  postList,
		Total: uint64(total),
	}, err
}
