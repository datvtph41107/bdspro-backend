package admin_handler

import (
	"bdspro/infra/mapper"
	admin_usecases "bdspro/internal/usecases/admin"
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminPostHandler struct {
	bdspropb.UnimplementedAdminPostServiceServer
	AdminPostService *admin_usecases.AdminPostUsecase
	PostMapper       *mapper.PostMapper
	AuthGrpcClient   *clients.AuthGrpcClient
}

func NewAdminPostHandler(
	uc *admin_usecases.AdminPostUsecase,
	postMapper *mapper.PostMapper,
	authGrpcClient *clients.AuthGrpcClient,
) *AdminPostHandler {
	return &AdminPostHandler{
		AdminPostService: uc,
		PostMapper:       postMapper,
		AuthGrpcClient:   authGrpcClient,
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
