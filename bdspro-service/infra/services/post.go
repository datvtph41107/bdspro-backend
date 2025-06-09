package service

import (
	"bdspro/internal/repo"
	"context"
	pb_bdspro "pb/types/bdspro"
)

type PostService struct {
	pb_bdspro.UnimplementedPostServiceServer
	postRepo repo.UPostRepo
}

func NewGrpcPostService(postRepo repo.UPostRepo) *PostService {
	return &PostService{
		postRepo: postRepo,
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
			Id:          media.ID,
			Url:         media.MediaURL,
			Type:        media.MediaType,
			OrderNumber: int32(media.Order),
		}
	}
	return &pb_bdspro.Post{
		Id:             post.ID,
		Title:          post.Title,
		Price:          post.Price,
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
	postList := make([]*pb_bdspro.Post, len(posts))
	for i, post := range posts {
		mediaList := make([]*pb_bdspro.PostMedia, len(post.MediaList))
		for i, media := range post.MediaList {
			mediaList[i] = &pb_bdspro.PostMedia{
				Id:          media.ID,
				Url:         media.MediaURL,
				Type:        media.MediaType,
				OrderNumber: int32(media.Order),
			}
		}
		postList[i] = &pb_bdspro.Post{
			Id:             post.ID,
			Title:          post.Title,
			Price:          post.Price,
			PriceType:      uint32(post.PriceType),
			PackageVisible: uint32(post.PackageVisible),
			Like:           post.Like,
			Comment:        post.Comment,
			NumView:        post.NumView,
			OwnerId:        post.OwnerID,
			OwnerType:      uint32(post.OwnerType),
			MediaList:      mediaList,
			Status:         uint32(post.Status),
		}
		// if len(post.MediaList) > 0 {
		// 	postList[i].MediaUrl = post.MediaList[0].MediaURL
		// 	postList[i].MediaType = post.MediaList[0].MediaType
		// }
	}
	return &pb_bdspro.GetPostByIDsResponse{Posts: postList}, nil
}
