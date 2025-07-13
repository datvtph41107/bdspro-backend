package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/dto"
	_utils "common/utils"
	pb_bdspro "pb/types/bdspro"
)

type PostMapper struct {
	ProductMapper *ProductMapper // forward reference
	AssetMapper   *AssetMapper   // forward reference
}

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

func (m *PostMapper) PostToPbList(posts *[]domain.PostItem) []*pb_bdspro.Post {
	postList := make([]*pb_bdspro.Post, len(*posts))
	for i, post := range *posts {
		postList[i] = m.PostToPb(&post)
	}
	return postList
}

func (m *PostMapper) PostPointerToPbList(posts []*domain.PostItem) []*pb_bdspro.Post {
	postList := make([]*pb_bdspro.Post, len(posts))
	for i, post := range posts {
		postList[i] = m.PostToPb(post)
	}
	return postList
}

func (m *PostMapper) PostToPb(post *domain.PostItem) *pb_bdspro.Post {
	mediaList := make([]*pb_bdspro.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = m.MapPostMedia(&media)
	}
	return &pb_bdspro.Post{
		Id:              post.ID,
		Title:           post.Title,
		PostPrice:       post.PostPrice,
		PriceType:       uint32(post.PriceType),
		PackageVisible:  uint32(post.PackageVisible),
		TransactionType: uint32(post.TransactionType),
		Visibility:      uint32(post.Visibility),
		Like:            post.Like,
		Comment:         post.Comment,
		NumView:         post.NumView,
		NumDate:         int32(post.NumDate),
		Content:         post.Content,
		OwnerId:         post.OwnerID,
		OwnerType:       uint32(post.OwnerType),
		MediaList:       mediaList,
		Status:          uint32(post.Status),
		ExpiredAt:       _utils.FormatTimeToString(post.ExpiredAt),
	}
}

func (m *PostMapper) PostMediaToDTO(media []*pb_bdspro.PostMedia) []dto.PostMediaItem {
	mediaList := make([]dto.PostMediaItem, len(media))
	for i, media := range media {
		mediaList[i] = dto.PostMediaItem{
			MediaURL:  media.MediaUrl,
			MediaType: media.MediaType,
			IsMain:    media.IsMain,
			Order:     int(media.SortOrder),
		}
	}
	return mediaList
}
func (m *PostMapper) PostSaveRequestToDTO(request *pb_bdspro.PostSaveRequest) *dto.PostSaveRequest {
	result := &dto.PostSaveRequest{
		ProductID:       request.ProductId,
		Title:           request.Title,
		Content:         request.Content,
		Price:           request.Price,
		TransactionType: enums.TransactionType(request.TransactionType),
		Visibility:      enums.EVisibility(request.Visibility),
		PackageVisible:  uint(request.PackageVisible),
		ExpiredAt:       _utils.ParseStringToTime(request.ExpiredAt),
		NumDate:         int(request.NumDate),
		MediaList:       m.PostMediaToDTO(request.MediaList),
		PriceType:       request.PriceType,
	}
	if request.Hidden != nil {
		result.Hidden = *request.Hidden
	}
	return result
}

func (m *PostMapper) PostUpdateRequestToDTO(request *pb_bdspro.PostSaveRequest) *dto.PostUpdateRequest {
	return &dto.PostUpdateRequest{
		Title:           request.Title,
		Content:         request.Content,
		Price:           request.Price,
		TransactionType: request.TransactionType,
		PackageVisible:  uint(request.PackageVisible),
		Visibility:      enums.EVisibility(request.Visibility),
		PriceType:       request.PriceType,
		NumDate:         int(request.NumDate),
		MediaList:       m.PostMediaToDTO(request.MediaList),
	}
}

func (m *PostMapper) PostDetailToPb(post *dto.PostDetail) *pb_bdspro.PostDetailResponse {
	result := &pb_bdspro.PostDetailResponse{
		Post: m.PostToPb(post.Post),
	}
	if post.ProductInfo != nil {
		result.ProductInfo = m.ProductMapper.MapProductPb(post.ProductInfo)
	}
	if post.AssetInfo != nil {
		result.AssetInfo = m.AssetMapper.AssetToPb(post.AssetInfo)
	}
	return result
}

func (m *PostMapper) PostSearchRequestToDTO(request *pb_bdspro.PostSearchRequest) *dto.PostSearchRequest {
	result := &dto.PostSearchRequest{
		Pagable: _dto.Pagable{
			Size: int(request.Size),
			Page: int(request.Page),
		},
		Content:         request.Content,
		TransactionType: request.TransactionType,
		Text:            request.Text,
		Title:           request.Title,
		Hidden:          request.Hidden,
		Status:          request.Status,
		Types:           request.Types,
		PackageVisibles: request.PackageVisibles,
	}
	if request.ExpiredFrom != nil {
		result.ExpiredFrom = _utils.ParseStringToTime(*request.ExpiredFrom)
	}
	if request.ExpiredTo != nil {
		result.ExpiredTo = _utils.ParseStringToTime(*request.ExpiredTo)
	}
	return result
}
