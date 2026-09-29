package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"github.com/jinzhu/copier"
)

type PostService struct {
	bdspropb.UnimplementedPostServiceServer
	postRepo   repo.PostRepo
	UC         *shared_usecase.PostUsecase
	PostMapper *mapper.PostMapper
}

func NewGrpcPostService(postRepo repo.PostRepo,
	uc *shared_usecase.PostUsecase,
	postMapper *mapper.PostMapper,
) *PostService {
	return &PostService{
		postRepo:   postRepo,
		UC:         uc,
		PostMapper: postMapper,
	}
}

func (s *PostService) OwnerPost(ctx context.Context, req *bdspropb.OwnerPostRequest) (*bdspropb.OwnerPostResponse, error) {
	success, err := s.postRepo.OwnerPost(ctx, req.ProfileId, req.PostId)
	return &bdspropb.OwnerPostResponse{Success: success}, err
}

func (s *PostService) GetPostByID(ctx context.Context, req *bdspropb.GetPostByIDRequest) (*bdspropb.Post, error) {
	post, err := s.postRepo.GetPostItemByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	mediaList := make([]*bdspropb.PostMedia, len(post.MediaList))
	for i, media := range post.MediaList {
		mediaList[i] = &bdspropb.PostMedia{
			Id:        media.ID,
			MediaUrl:  media.MediaURL,
			MediaType: media.MediaType,
			IsMain:    media.IsMain,
			SortOrder: int32(media.SortOrder),
		}
	}
	return &bdspropb.Post{
		Id:             post.ID,
		Title:          post.Title,
		PostPrice:      post.Price,
		PriceType:      uint32(post.PriceType),
		PackageVisible: uint32(post.PackageVisible),
		Like:           post.Like,
		Comment:        post.Comment,
		NumView:        post.NumView,
		OwnerId:        post.OwnerID,
		OwnerType:      uint32(post.OwnerOf),
		MediaList:      mediaList,
	}, nil
}

func (s *PostService) GetPostByIDs(ctx context.Context, req *bdspropb.GetPostByIDsRequest) (*bdspropb.GetPostByIDsResponse, error) {
	posts, err := s.postRepo.GetByIDsWithStatus(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for i, post := range posts {
		posts[i].VisibleStatus = s.UC.MapPostStatus(post)
	}
	postList := s.PostMapper.PostPointerToPbList(posts)

	return &bdspropb.GetPostByIDsResponse{Posts: postList}, nil
}

func (s *PostService) CreatePost(ctx context.Context, req *bdspropb.PostSaveRequest) (*bdspropb.Post, error) {
	body := s.PostMapper.PostSaveRequestToDTO(req)

	result, err := s.UC.CreatePost(ctx, body)
	if err != nil {
		return nil, err
	}

	post := &bdspropb.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) NewExpiredPost(ctx context.Context, req *bdspropb.PostExpiredRequest) (*bdspropb.Post, error) {
	body := dto.PostExpiredRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.MakeNewExpired(ctx, body)
	if err != nil {
		return nil, err
	}

	post := &bdspropb.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) GetDetail(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.PostDetailResponse, error) {
	result, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return s.PostMapper.PostDetailToPb(ctx, result), nil
}
func (s *PostService) UpdatePost(ctx context.Context, req *bdspropb.PostSaveRequest) (*bdspropb.Post, error) {
	body := s.PostMapper.PostUpdateRequestToDTO(req)

	result, err := s.UC.UpdatePost(ctx, req.Id, body)
	if err != nil {
		return nil, err
	}

	post := &bdspropb.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}
func (s *PostService) Delete(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Message: "success",
	}, nil
}
func (s *PostService) UpdateHidden(ctx context.Context, req *bdspropb.UpdateHiddenRequest) (*bdspropb.Response, error) {
	body := dto.UpdateHiddenRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	_, err := s.UC.UpdateHidden(ctx, body)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Message: "success",
	}, nil
}
func (s *PostService) GetPersonalPost(ctx context.Context, req *bdspropb.PostSearchRequest) (*bdspropb.PostListResponse, error) {
	query := s.PostMapper.PostSearchRequestToDTO(req)

	result, total, err := s.UC.PersonalPost(ctx, query)
	if err != nil {
		return nil, err
	}

	postList := s.PostMapper.PostToPbList(&result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(total),
	}, nil
}
func (s *PostService) GetGlobalPost(ctx context.Context, req *bdspropb.PostSearchRequest) (*bdspropb.PostListResponse, error) {
	query := dto.PostSearchRequest{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GlobalPost(ctx, query)
	if err != nil {
		return nil, err
	}

	postList := s.PostMapper.PostToPbList(&result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(total),
	}, nil
}
func (s *PostService) GetPublish(ctx context.Context, req *bdspropb.PostPublishSearch) (*bdspropb.PostListResponse, error) {
	query := dto.PostPublishSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GetPublishByProfileID(ctx, req.ProfileId, query)
	postList := make([]*bdspropb.Post, len(result))
	copier.Copy(&postList, &result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(total),
	}, err
}

func (s *PostService) CreatePostOrganization(ctx context.Context, req *bdspropb.PostSaveRequest) (*bdspropb.Post, error) {
	body := s.PostMapper.PostSaveRequestToDTO(req)

	result, err := s.UC.CreateOrganizationAsset(ctx, body)
	if err != nil {
		return nil, err
	}

	post := &bdspropb.Post{}
	if err := copier.Copy(post, result); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) CurrentOrganizationPosts(ctx context.Context, req *bdspropb.PostSearchRequest) (*bdspropb.PostListResponse, error) {
	query := dto.PostSearchRequest{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GetPostOfCurrentOrganization(ctx, &query)
	if err != nil {
		return nil, err
	}
	postList := s.PostMapper.PostToPbList(&result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(total),
	}, nil
}

// GetTodayPost trả ra ngẫu nhiên 5 bản ghi tin đăng
// @Summary Lấy ngẫu nhiên 5 tin đăng
// @Description API trả ra ngẫu nhiên 5 bản ghi tin đăng đã được publish
// @Tags Post
// @Accept json
// @Produce json
// @Success 200 {object} bdspropb.PostListResponse
// @Router /v2/bdspro/v2/dashboard/today/post [get]
func (s *PostService) GetTodayPost(ctx context.Context, req *sharepb.Empty) (*bdspropb.PostListResponse, error) {
	result, err := s.UC.GetRandomPosts(ctx, 5)
	if err != nil {
		return nil, err
	}

	postList := s.PostMapper.PostToPbList(&result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(len(result)),
	}, nil
}

// GetPersonalPosts lấy danh sách post của user theo profileId
// @Summary Lấy danh sách post của user
// @Description Lấy danh sách post của user theo profileId với phân trang
// @Tags Post
// @Accept json
// @Produce json
// @Param id path uint64 true "Profile ID"
// @Param page query uint32 false "Trang hiện tại (mặc định: 1)"
// @Param size query uint32 false "Số lượng item trên mỗi trang (mặc định: 20)"
// @Param text query string false "Tìm kiếm theo text"
// @Success 200 {object} bdspropb.PostListResponse
// @Router /v3/bdspro/post/personal/{id} [get]
func (s *PostService) GetPersonalPosts(ctx context.Context, req *sharepb.RequestV3Proto) (*bdspropb.PostListResponse, error) {
	// Convert RequestV3Proto sang PostSearchRequest DTO
	query := s.PostMapper.RequestV3ProtoToPostSearchDTO(req)

	// Lấy posts theo profileId từ request
	result, total, err := s.postRepo.Search(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}

	// Map status cho mỗi post
	for i := range result {
		result[i].VisibleStatus = s.UC.MapPostStatus(&result[i])
		s.UC.MapStatusNames(&result[i])
	}

	postList := s.PostMapper.PostToPbList(&result)

	return &bdspropb.PostListResponse{
		Data:  postList,
		Total: int32(total),
	}, nil
}
