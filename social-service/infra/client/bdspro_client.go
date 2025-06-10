package client

import (
	_routes "common/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	pb_bdspro "pb/types/bdspro"
	pb_social "pb/types/social"
	"social/env"
	"social/internal/dto"
	"time"

	"google.golang.org/grpc"
)

// @bind: social/internal/interface.BdsproClient
type BdsproClient struct {
	client pb_bdspro.PostServiceClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{client: nil}
	go func() {
		url := "localhost"
		name := "bdspro_grpc"
		if env.ENV_RUNTIME == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8202", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go monitorConnection(ctx, name, conn)

		client := pb_bdspro.NewPostServiceClient(conn)
		bdsproClient.client = client
	}()

	return bdsproClient
}

func (c *BdsproClient) OwnerPost(ctx context.Context, postID uint64, profileID uint64) error {
	resposne, err := c.client.OwnerPost(ctx, &pb_bdspro.OwnerPostRequest{
		PostId:    postID,
		ProfileId: profileID,
	})
	if err != nil {
		return err
	}
	if !resposne.Success {
		return &_routes.Except{
			Code:    http.StatusForbidden,
			Message: "bạn không có quyền tạo/forward tin đăng này",
		}
	}
	return nil
}

func (c *BdsproClient) GetPostByID(ctx context.Context, postID uint64) (*dto.BdsproPost, error) {
	resposne, err := c.client.GetPostByID(ctx, &pb_bdspro.GetPostByIDRequest{
		Id: postID,
	})
	if err != nil {
		return nil, err
	}
	mediaList := make([]*dto.PostMedia, 0)
	for _, media := range resposne.MediaList {
		mediaList = append(mediaList, &dto.PostMedia{
			Id:          media.Id,
			Url:         media.Url,
			Type:        media.Type,
			OrderNumber: media.OrderNumber,
		})
	}
	post := &dto.BdsproPost{
		Id:              resposne.Id,
		ProductId:       resposne.ProductId,
		ExpiredAt:       resposne.ExpiredAt,
		Status:          resposne.Status,
		Visibility:      resposne.Visibility,
		Hidden:          resposne.Hidden,
		Content:         resposne.Content,
		TransactionType: resposne.TransactionType,
		Title:           resposne.Title,
		NumDate:         resposne.NumDate,
		Type:            resposne.Type,
		Price:           resposne.Price,
		PriceType:       resposne.PriceType,
		PackageVisible:  resposne.PackageVisible,
		Like:            resposne.Like,
		Comment:         resposne.Comment,
		NumView:         resposne.NumView,
		OwnerId:         resposne.OwnerId,
		OwnerType:       resposne.OwnerType,
		MediaList:       mediaList,
		MediaUrl:        resposne.MediaUrl,
		MediaType:       resposne.MediaType,
	}
	return post, nil
}

func (c *BdsproClient) GetPostByIDs(ctx context.Context, postIDs []uint64) ([]*dto.BdsproPost, error) {
	resposne, err := c.client.GetPostByIDs(ctx, &pb_bdspro.GetPostByIDsRequest{
		Ids: postIDs,
	})
	if err != nil {
		return nil, err
	}
	posts := make([]*dto.BdsproPost, 0)
	for _, post := range resposne.Posts {
		mediaList := make([]*dto.PostMedia, 0)
		for _, media := range post.MediaList {
			mediaList = append(mediaList, &dto.PostMedia{
				Id:          media.Id,
				Url:         media.Url,
				Type:        media.Type,
				OrderNumber: media.OrderNumber,
			})
		}
		posts = append(posts, &dto.BdsproPost{
			Id:              post.Id,
			ProductId:       post.ProductId,
			ExpiredAt:       post.ExpiredAt,
			Status:          post.Status,
			Visibility:      post.Visibility,
			Hidden:          post.Hidden,
			Content:         post.Content,
			TransactionType: post.TransactionType,
			Title:           post.Title,
			NumDate:         post.NumDate,
			Type:            post.Type,
			Price:           post.Price,
			PriceType:       post.PriceType,
			PackageVisible:  post.PackageVisible,
			Like:            post.Like,
			Comment:         post.Comment,
			NumView:         post.NumView,
			OwnerId:         post.OwnerId,
			OwnerType:       post.OwnerType,
			MediaList:       mediaList,
		})
	}
	return posts, nil
}

func (c *BdsproClient) PbPostToDomain(ctx context.Context, newsFeeds []*pb_social.NewsFeed) {
	postIDs := make([]uint64, 0)
	for _, nf := range newsFeeds {
		if nf.PostId != nil {
			postIDs = append(postIDs, *nf.PostId)
		}
	}

	if len(postIDs) == 0 {
		return
	}

	posts, err := c.client.GetPostByIDs(ctx, &pb_bdspro.GetPostByIDsRequest{
		Ids: postIDs,
	})
	if err != nil {
		return
	}

	postMap := make(map[uint64]*pb_bdspro.Post, len(posts.Posts))
	for _, post := range posts.Posts {
		postMap[post.Id] = post
	}

	for _, nf := range newsFeeds {
		if nf.PostId != nil {
			if post, ok := postMap[*nf.PostId]; ok {
				nf.Post = post
			}
		}
	}
}
