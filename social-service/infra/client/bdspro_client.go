package client

import (
	_routes "common/routes"
	"context"
	"net/http"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	pb_social "pb/types/social"
	"social/internal/dto"
)

// @bind: social/internal/interface.BdsproClient
type BdsproClient struct {
	client bdspropb.PostServiceClient
}

func NewBdsproClient(rpcClient *clients.BdsproGrpcClient) *BdsproClient {
	return &BdsproClient{client: rpcClient.PostClient}
}

func (c *BdsproClient) OwnerPost(ctx context.Context, postID uint64, profileID uint64) error {
	resposne, err := c.client.OwnerPost(ctx, &bdspropb.OwnerPostRequest{
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
	resposne, err := c.client.GetPostByID(ctx, &bdspropb.GetPostByIDRequest{
		Id: postID,
	})
	if err != nil {
		return nil, err
	}
	mediaList := make([]*dto.PostMedia, 0)
	for _, media := range resposne.MediaList {
		mediaList = append(mediaList, &dto.PostMedia{
			Id:          media.Id,
			Url:         media.MediaUrl,
			Type:        media.MediaType,
			OrderNumber: media.SortOrder,
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
		Price:           resposne.PostPrice,
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
	resposne, err := c.client.GetPostByIDs(ctx, &bdspropb.GetPostByIDsRequest{
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
				Url:         media.MediaUrl,
				Type:        media.MediaType,
				OrderNumber: media.SortOrder,
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
			Price:           post.PostPrice,
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

	posts, err := c.client.GetPostByIDs(ctx, &bdspropb.GetPostByIDsRequest{
		Ids: postIDs,
	})
	if err != nil {
		return
	}

	postMap := make(map[uint64]*bdspropb.Post, len(posts.Posts))
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
