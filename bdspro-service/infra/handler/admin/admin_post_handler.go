package admin_handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminPostHandler struct {
	bdspropb.UnimplementedAdminPostServiceServer
	AdminPostService *admin_usecases.AdminPostUsecase
	PostMapper       *mapper.PostMapper
	AuthGrpcClient   *client.AuthClient
	UserClient       *client.UserClient
}

func NewAdminPostHandler(
	uc *admin_usecases.AdminPostUsecase,
	postMapper *mapper.PostMapper,
	authGrpcClient *client.AuthClient,
	userClient *client.UserClient,
) *AdminPostHandler {
	return &AdminPostHandler{
		AdminPostService: uc,
		PostMapper:       postMapper,
		AuthGrpcClient:   authGrpcClient,
		UserClient:       userClient,
	}
}

// @Summary Approve a post
// @Description Approve a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/posts/{id}/approve [post]
func (h *AdminPostHandler) ApprovePost(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_DUYET"}); err != nil {
		return nil, err
	}

	err := h.AdminPostService.Approve(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Duyệt tin đăng thành công",
	}, nil
}

// @Summary Reject a post
// @Description Reject a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/posts/{id}/reject [post]
func (h *AdminPostHandler) RejectPost(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_TUCHOI"}); err != nil {
		return nil, err
	}

	err := h.AdminPostService.Reject(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Từ chối tin đăng thành công",
	}, nil
}

// @Summary Archive a post
// @Description Archive a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param request body bdspropb.ArchivedRequest true "Archive request"
// @Success 200 {object} bdspropb.Response
// @Router /admin/posts/archived [put]
func (h *AdminPostHandler) ArchivePost(ctx context.Context, req *bdspropb.ArchivedRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_LUU_TRU"}); err != nil {
		return nil, err
	}

	err := h.AdminPostService.Archive(ctx, req.Id, req.Archived)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Lưu trữ tin đăng thành công",
	}, nil
}

// @Summary Hide a post
// @Description Hide a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/posts/{id}/hide [post]
func (h *AdminPostHandler) HidePost(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_AN"}); err != nil {
		return nil, err
	}

	err := h.AdminPostService.Hide(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Ẩn tin đăng thành công",
	}, nil
}

// @Summary Unhide a post
// @Description Unhide a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/posts/{id}/unhide [post]
func (h *AdminPostHandler) UnhidePost(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_HIEN"}); err != nil {
		return nil, err
	}

	err := h.AdminPostService.Unhide(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Hiện tin đăng thành công",
	}, nil
}

// @Summary Get posts
// @Description Get posts for admin
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Param name query string false "Post name"
// @Param status query int false "Post status"
// @Param fromDate query string false "From date"
// @Param toDate query string false "To date"
// @Success 200 {object} bdspropb.AdminPostSearchResponse
// @Router /admin/posts [get]
func (h *AdminPostHandler) GetPosts(ctx context.Context, req *bdspropb.AdminPostSearchRequest) (*bdspropb.AdminPostSearchResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_XEM"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	reqDto := &dto.AdminPostSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Title:            req.Title,
		Status:           req.Status,
		FromDate:         _utils.ParseStringToTime(req.FromDate),
		ToDate:           _utils.ParseStringToTime(req.ToDate),
		TransactionTypes: req.TransactionTypes,
		VisibleStatuses:  req.VisibleStatuses,
	}

	// Get posts from usecase
	posts, total, err := h.AdminPostService.GetPosts(ctx, reqDto)
	if err != nil {
		return nil, err
	}

	// Convert DTO to proto response using mapper
	response := &bdspropb.AdminPostSearchResponse{
		Data:  h.PostMapper.MapAdminPostList(ctx, posts),
		Total: total,
	}

	return response, nil
}

// @Summary Get post detail
// @Description Get detailed information of a post
// @Tags AdminPost
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} bdspropb.Post
// @Router /admin/posts/{id} [get]
func (h *AdminPostHandler) AdminDetailPost(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Post, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TD_XEM"}); err != nil {
		return nil, err
	}

	post, err := h.AdminPostService.GetPostDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	result := h.PostMapper.PostToPbDetail(post)
	if post.OwnerID != 0 {
		user := h.UserClient.GetProfileById(ctx, post.OwnerID)
		result.Owner = user
	}
	return result, nil
}
