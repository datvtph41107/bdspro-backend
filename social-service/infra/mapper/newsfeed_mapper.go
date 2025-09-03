package mapper

import (
	_utils "common/utils"
	"context"
	socialpb "pb/types/social"
	"social/internal/domain"
)

type NewsFeedMapper struct {
}

func NewNewsFeedMapper() *NewsFeedMapper {
	return &NewsFeedMapper{}
}

func (m *NewsFeedMapper) DomainToNewsFeedPbList(ctx context.Context, newsFeed []*domain.NewsFeed) []*socialpb.NewsFeed {
	result := make([]*socialpb.NewsFeed, len(newsFeed))
	for i, news := range newsFeed {
		result[i] = m.DomainToNewsFeedPb(ctx, news)
	}

	return result
}

func (m *NewsFeedMapper) DomainToNewsFeedPb(ctx context.Context, newsFeed *domain.NewsFeed) *socialpb.NewsFeed {
	result := &socialpb.NewsFeed{
		Id:         newsFeed.ID,
		Content:    newsFeed.Content,
		Image:      newsFeed.Image,
		Link:       newsFeed.Link,
		Visibility: int32(newsFeed.Visibility),
		OwnerOf:    socialpb.OwnerOf_user,
		PostId:     newsFeed.PostID,
		// Post: ,
		// GroupId:    newsFeed.GroupId,
		Title:      newsFeed.Title,
		UpdatedAt:  _utils.FormatTimeToString(newsFeed.UpdatedAt),
		CreatedBy:  newsFeed.CreatedBy,
		LikedAt:    _utils.FormatTimeToString(newsFeed.LikedAt),
		DisLike:    newsFeed.DisLike,
		NumLike:    int32(newsFeed.NumLike),
		NumComment: int32(newsFeed.NumComment),
		NumShare:   int32(newsFeed.NumShare),
		NumView:    int32(newsFeed.NumView),
		IsReel:     newsFeed.IsReel,
	}

	if newsFeed.NewsFeedMedias != nil {
		result.MediaFeeds = make([]*socialpb.NewsFeedMedia, len(newsFeed.NewsFeedMedias))
		for i, media := range newsFeed.NewsFeedMedias {
			result.MediaFeeds[i] = &socialpb.NewsFeedMedia{
				Id:          media.ID,
				Url:         media.Url,
				Type:        media.Type,
				OrderNumber: int32(media.OrderNumber),
			}
		}
	}

	return result
}

func (m *NewsFeedMapper) DomainToReelResponseList(ctx context.Context, newsFeed []*domain.ReelEntity) []*socialpb.ReelResponse {
	result := make([]*socialpb.ReelResponse, len(newsFeed))
	for i, news := range newsFeed {
		result[i] = m.DomainToReelResponse(news)
	}
	return result
}

func (m *NewsFeedMapper) DomainToReelResponse(newsFeed *domain.ReelEntity) *socialpb.ReelResponse {
	result := &socialpb.ReelResponse{
		Id:      newsFeed.ID,
		Content: newsFeed.Content,
		ReelUrl: newsFeed.ReelUrl,
		// AuthorAvatar: newsFeed.AuthorAvatar,
		// AuthorName:   newsFeed.AuthorName,
		// CreatedBy:  *newsFeed.CreatedBy,
		CreatedAt:  _utils.FormatTimeToString(newsFeed.CreatedAt),
		NumLike:    uint64(newsFeed.NumLike),
		NumComment: uint64(newsFeed.NumComment),
		NumShare:   uint64(newsFeed.NumShare),
		NumView:    uint64(newsFeed.NumView),
	}

	if newsFeed.CreatedBy != nil {
		result.CreatedBy = *newsFeed.CreatedBy
	}

	return result
}
