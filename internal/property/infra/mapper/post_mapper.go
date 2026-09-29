package mapper

import (
	"bdspro/infra/client"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type PostMapper struct {
	ProductMapper *ProductMapper // forward reference
	AssetMapper   *AssetMapper   // forward reference
	UserClient    *client.UserClient
}

func NewPostMapper(userClient *client.UserClient) *PostMapper {
	return &PostMapper{
		UserClient: userClient,
	}
}

func (*PostMapper) MapPostMedia(media *domain.PostMediaEntity) *bdspropb.PostMedia {
	return &bdspropb.PostMedia{
		Id:        media.ID,
		MediaUrl:  media.MediaURL,
		MediaType: media.MediaType,
		IsMain:    media.IsMain,
		SortOrder: int32(media.SortOrder),
	}
}

func (m *PostMapper) PostToPbList(posts *[]domain.PostItem) []*bdspropb.Post {
	postList := make([]*bdspropb.Post, len(*posts))
	for i, post := range *posts {
		postList[i] = m.PostToPb(&post)
	}
	return postList
}

func (m *PostMapper) PostPointerToPbList(posts []*domain.PostItem) []*bdspropb.Post {
	postList := make([]*bdspropb.Post, len(posts))
	for i, post := range posts {
		postList[i] = m.PostToPb(post)
	}
	return postList
}

func (m *PostMapper) PostToPb(post *domain.PostItem) *bdspropb.Post {
	mediaList := make([]*bdspropb.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = m.MapPostMedia(&media)
	}

	// Đảm bảo luôn có transactionStatus mặc định
	transactionStatus := uint32(post.TransactionStatus)
	if transactionStatus == 0 {
		// Nếu chưa được set, dùng giá trị mặc định dựa vào TransactionType
		if post.TransactionType == enums.TransactionTypeSale {
			transactionStatus = 100 // Đang bán (mặc định)
		} else {
			transactionStatus = 200 // Đang cho thuê (mặc định)
		}
	}

	result := &bdspropb.Post{
		Id:                post.ID,
		Title:             post.Title,
		PostPrice:         post.PostPrice,
		PriceType:         uint32(post.PriceType),
		PackageVisible:    uint32(post.PackageVisible),
		TransactionType:   uint32(post.TransactionType),
		Visibility:        uint32(post.Visibility),
		Like:              post.Like,
		Comment:           post.Comment,
		NumView:           post.NumView,
		NumDate:           int32(post.NumDate),
		Content:           post.Content,
		OwnerId:           post.OwnerID,
		OwnerType:         uint32(post.OwnerOf),
		MediaList:         mediaList,
		Status:            uint32(post.Status),
		TransactionStatus: transactionStatus, // Map trạng thái giao dịch: 100-Đang bán, 110-Đã bán, 200-Đang cho thuê, 210-Đã cho thuê
		ExpiredAt:         _utils.FormatTimeToString(post.ExpiredAt),
		ProvinceId:        post.ProvinceId,
		WardId:            post.WardId,
		ProvinceName:      post.ProvinceName,
		WardName:          post.WardName,
		Area:              post.Area,
	}

	return result
}

func (m *PostMapper) PostMediaToDTO(media []*bdspropb.PostMedia) []dto.PostMediaItem {
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
func (m *PostMapper) PostSaveRequestToDTO(request *bdspropb.PostSaveRequest) *dto.PostSaveRequest {
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
		Medatadatas:     m.PostMediaToDTO(request.MediaList),
		PriceType:       request.PriceType,
	}
	if request.Hidden != nil {
		result.Hidden = *request.Hidden
	}
	return result
}

func (m *PostMapper) PostUpdateRequestToDTO(request *bdspropb.PostSaveRequest) *dto.PostUpdateRequest {
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

func (m *PostMapper) PostDetailToPb(ctx context.Context, post *dto.PostDetail) *bdspropb.PostDetailResponse {
	result := &bdspropb.PostDetailResponse{
		Post: m.PostToPb(post.Post),
	}
	if post.ProductInfo != nil {
		result.ProductInfo = m.ProductMapper.MapProductPb(post.ProductInfo)
	}
	if post.AssetInfo != nil {
		result.AssetInfo = m.AssetMapper.MapAssetPb(post.AssetInfo)
	}
	// Lấy thông tin owner
	if post.Post != nil && post.Post.OwnerID > 0 {
		owner := m.UserClient.GetProfileById(ctx, post.Post.OwnerID)
		if owner != nil {
			result.Owner = owner
		}
	}
	return result
}

func (m *PostMapper) PostSearchRequestToDTO(request *bdspropb.PostSearchRequest) *dto.PostSearchRequest {
	result := &dto.PostSearchRequest{
		Pagable: _dto.Pagable{
			Size: request.Size,
			Page: request.Page,
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

// RequestV3ProtoToPostSearchDTO convert RequestV3Proto sang PostSearchRequest
func (m *PostMapper) RequestV3ProtoToPostSearchDTO(request *sharepb.RequestV3Proto) *dto.PostSearchRequest {
	result := &dto.PostSearchRequest{
		Pagable: _dto.Pagable{
			Page: request.Page,
			Size: request.Size,
		},
		Text: request.Text,
	}
	return result
}

// Admin methods
// func (m *PostMapper) MapPostPb(post *domain.Post) *bdspropb.Post {
// 	if post == nil {
// 		return nil
// 	}

// 	return &bdspropb.Post{
// 		Id:              post.ID,
// 		Title:           post.Title,
// 		PostPrice:       post.Price,
// 		TransactionType: uint32(post.TransactionType),
// 		Visibility:      uint32(post.Visibility),
// 		PackageVisible:  uint32(post.PackageVisible),
// 		NumDate:         int32(post.NumDate),
// 		Content:         post.Content,
// 		ProductId:       post.ProductID,
// 		Status:          uint32(post.Status),
// 		Hidden:          &post.Hidden,
// 		ExpiredAt:       _utils.FormatTimeToString(post.ExpiredAt),
// 		CreatedAt:       _utils.FormatTimeToString(&post.CreatedAt),
// 		UpdatedAt:       _utils.FormatTimeToString(&post.UpdatedAt),
// 	}
// }

// func (m *PostMapper) MapPostPbList(posts []*domain.Post) []*bdspropb.Post {
// 	if posts == nil {
// 		return nil
// 	}

// 	result := make([]*bdspropb.Post, len(posts))
// 	for i, post := range posts {
// 		result[i] = m.MapPostPb(post)
// 	}
// 	return result
// }

func (m *PostMapper) PostSavePbToDTO(req *bdspropb.PostSaveRequest) *dto.PostSaveRequest {
	if req == nil {
		return nil
	}

	result := &dto.PostSaveRequest{
		ProductID:       req.ProductId,
		Title:           req.Title,
		Content:         req.Content,
		Price:           req.Price,
		TransactionType: enums.TransactionType(req.TransactionType),
		Visibility:      enums.EVisibility(req.Visibility),
		PackageVisible:  uint(req.PackageVisible),
		ExpiredAt:       _utils.ParseStringToTime(req.ExpiredAt),
		NumDate:         int(req.NumDate),
		PriceType:       req.PriceType,
	}

	if req.Hidden != nil {
		result.Hidden = *req.Hidden
	}

	if len(req.MediaList) > 0 {
		result.Medatadatas = m.PostMediaToDTO(req.MediaList)
	}

	return result
}

// MapAdminPostItem - Map domain.Post sang bdspropb.AdminPostItem
func (m *PostMapper) MapAdminPostItem(ctx context.Context, post *domain.Post) *bdspropb.AdminPostItem {
	return &bdspropb.AdminPostItem{
		Id:                post.ID,
		Code:              post.Code,
		Title:             post.Title,
		Status:            uint32(post.Status), // Sử dụng Status field (trạng thái cơ bản)
		StatusName:        m.getStatusName(post.Status),
		PublishedAt:       _utils.FormatTimeToString(post.PublishedAt),
		ViewCount:         int64(post.NumView),
		InterestCount:     int64(post.Like),
		ContactCount:      int64(post.Comment),
		TransactionType:   uint32(post.TransactionType),
		TransactionStatus: uint32(post.TransactionStatus),
		VisibleStatus:     uint32(post.VisibleStatus),
	}
}

// MapAdminPostList - Map danh sách domain.Post sang bdspropb.AdminPostItem
func (m *PostMapper) MapAdminPostList(ctx context.Context, posts []domain.Post) []*bdspropb.AdminPostItem {
	result := make([]*bdspropb.AdminPostItem, len(posts))

	ownerMap := make(map[uint64]struct{})
	for i, post := range posts {
		result[i] = m.MapAdminPostItem(ctx, &post)
		ownerMap[post.OwnerID] = struct{}{}
	}

	profiles, err := m.UserClient.GetMapByIDs(ctx, ownerMap)
	if err != nil {
		return result
	}

	for i, post := range posts {
		if profile, ok := profiles[post.OwnerID]; ok {
			result[i].Owner = profile
		}
	}

	return result
}

// Helper methods
func (m *PostMapper) getOwnerName(post *domain.Post) string {
	// TODO: Implement logic to get owner name from profile
	// This might require additional data or separate query
	return "N/A"
}

func (m *PostMapper) getStatusName(status enums.EPostStatus) string {
	if name, exists := enums.EPostStatusNames[status]; exists {
		return name
	}
	return "Không xác định"
}

func (m *PostMapper) PostToPbDetail(post *domain.Post) *bdspropb.Post {
	result := &bdspropb.Post{
		Id:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		PostPrice:       post.Price,
		TransactionType: uint32(post.TransactionType),
		Visibility:      uint32(post.Visibility),
		PackageVisible:  uint32(post.PackageVisible),
		NumDate:         int32(post.NumDate),
		Hidden:          post.Hidden,
		ExpiredAt:       _utils.FormatTimeToString(post.ExpiredAt),
		// CreatedAt:         _utils.FormatTimeToString(post.CreatedAt),
		// UpdatedAt:         _utils.FormatTimeToString(post.UpdatedAt),
		// MediaList:         m.PostMediaToDTO(post.MediaList),
		OwnerId:   post.OwnerID,
		OwnerType: uint32(post.OwnerOf),
		// PublishedAt:       _utils.FormatTimeToString(post.PublishedAt),
		// PublishedBy:       post.PublishedBy,
		// Code:              post.Code,
		PriceType: uint32(post.PriceType),
		// TransactionStatus: uint32(post.TransactionStatus),
		// VisibleStatus:     uint32(post.VisibleStatus),
		Like:              post.Like,
		Comment:           post.Comment,
		NumView:           post.NumView,
		Status:            uint32(post.Status),
		TransactionStatus: uint32(post.TransactionStatus),
		VisibleStatus:     uint32(post.VisibleStatus),
	}
	if post.Product != nil {
		result.Product = m.ProductMapper.MapProductPb(post.Product)
	}
	return result
}
