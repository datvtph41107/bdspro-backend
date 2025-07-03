package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

type PostMapper struct{}

func NewPostMapper() *PostMapper {
	return &PostMapper{}
}

func (*PostMapper) MapPostMedia(media *domain.PostMediaEntity) *pb_bdspro.PostMedia {
	return &pb_bdspro.PostMedia{
		Id:        media.ID,
		MediaUrl:  media.MediaURL,
		MediaType: media.MediaType,
		IsMain:    media.IsMain,
		SortOrder: int32(media.SortOrder),
	}
}

func (m *PostMapper) PostToPbList(posts []*domain.Post) []*pb_bdspro.Post {
	postList := make([]*pb_bdspro.Post, len(posts))
	for i, post := range posts {
		postList[i] = m.PostToPb(post)
	}
	return postList
}

func (m *PostMapper) PostToPb(post *domain.Post) *pb_bdspro.Post {
	mediaList := make([]*pb_bdspro.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = m.MapPostMedia(&media)
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
		Status:         uint32(post.Status),
	}
}
