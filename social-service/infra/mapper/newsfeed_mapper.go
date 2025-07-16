package mapper

import (
	_utils "common/utils"
	socialpb "pb/types/social"
	"social/internal/domain"
)

type NewsFeedMapper struct {
}

func NewNewsFeedMapper() *NewsFeedMapper {
	return &NewsFeedMapper{}
}

func (m *NewsFeedMapper) DomainToReelResponseList(newsFeed []*domain.ReelEntity) []*socialpb.ReelResponse {
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
