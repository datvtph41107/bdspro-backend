package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

func MapPostMedia(media *domain.PostMediaEntity) *pb_bdspro.PostMedia {
	return &pb_bdspro.PostMedia{
		Id:        media.ID,
		MediaUrl:  media.MediaURL,
		MediaType: media.MediaType,
		IsMain:    media.IsMain,
		SortOrder: int32(media.SortOrder),
	}
}

func PostToPbList(posts []*domain.Post) []*pb_bdspro.Post {
	postList := make([]*pb_bdspro.Post, len(posts))
	for i, post := range posts {
		postList[i] = PostToPb(post)
	}
	return postList
}

func PostToPb(post *domain.Post) *pb_bdspro.Post {
	mediaList := make([]*pb_bdspro.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = MapPostMedia(&media)
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
