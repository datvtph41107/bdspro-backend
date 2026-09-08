package handler

import (
	"context"
	"time"

	_utils "common/utils"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateDataHandler struct {
	hubpb.UnimplementedUpdateDataServiceServer
	updateDataUsecase _usecase.IUpdateDataUsecase
}

func NewUpdateDataHandler(updateDataUsecase _usecase.IUpdateDataUsecase) *UpdateDataHandler {
	return &UpdateDataHandler{
		updateDataUsecase: updateDataUsecase,
	}
}

// CheckVersionSync kiểm tra có cần sync không
// @Summary Kiểm tra trạng thái sync
// @Description Kiểm tra xem có thay đổi cần sync cho resource từ lastSync
// @Tags UpdateData
// @Accept json
// @Produce json
// @Param resource path string true "Tên resource (product, asset, ...)"
// @Param id query uint64 false "ID scope (profileId)"
// @Param lastSync query int64 false "Unix timestamp lần sync cuối"
// @Param limit query int64 false "Limit"
// @Success 200 {object} hubpb.CheckVersionSyncResponse
// @Router /v2/hub/update/{resource}/sync/status [get]
func (h *UpdateDataHandler) CheckVersionSync(
	ctx context.Context,
	req *hubpb.CheckVersionSyncRequest,
) (*hubpb.CheckVersionSyncResponse, error) {
	if req == nil || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required")
	}

	ownerID := req.Id
	if ownerID == 0 {
		ownerID = _utils.GetProfileIdWithContext(ctx)
	}
	if ownerID == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required (id in request or profile in context)")
	}

	changedIds, err := h.updateDataUsecase.CheckVersionSync(
		ctx, req.Resource, req.LastSync, req.Limit,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check version sync failed: %v", err)
	}

	resp := &hubpb.CheckVersionSyncResponse{
		NeedSync:       len(changedIds) > 0,
		CurrentVersion: time.Now().Unix(),
		ChangeCount:    int32(len(changedIds)),
		ChangedIds:     changedIds,
	}
	return resp, nil
}

func (h *UpdateDataHandler) CheckVersionSyncById(
	ctx context.Context,
	req *hubpb.CheckVersionSyncRequest,
) (*hubpb.CheckVersionSyncResponse, error) {
	if req == nil || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required (id in request or profile in context)")
	}

	timestamp, err := h.updateDataUsecase.CheckVersionSyncById(
		ctx, req.Resource, req.LastSync, req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check version sync failed: %v", err)
	}

	resp := &hubpb.CheckVersionSyncResponse{
		NeedSync:       req.LastSync < timestamp,
		CurrentVersion: timestamp,
		ChangeCount:    0,
		ChangedIds:     nil,
	}
	return resp, nil
}

// RefreshSyncIds lấy danh sách IDs đã thay đổi cần sync
// @Summary Lấy danh sách IDs cần refresh
// @Description Lấy danh sách resource IDs đã thay đổi từ lastSync
// @Tags UpdateData
// @Accept json
// @Produce json
// @Param resource path string true "Tên resource"
// @Param lastSync query int64 false "Unix timestamp lần sync cuối"
// @Param limit query int64 false "Limit"
// @Success 200 {object} hubpb.RefreshSyncIdsResponse
// @Router /v2/hub/update/{resource}/refresh/ids [get]
func (h *UpdateDataHandler) RefreshSyncIds(
	ctx context.Context,
	req *hubpb.RefreshSyncIdsRequest,
) (*hubpb.RefreshSyncIdsResponse, error) {
	if req == nil || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required")
	}

	ownerID := _utils.GetProfileIdWithContext(ctx)
	if ownerID == 0 {
		return nil, status.Error(codes.Unauthenticated, "profile_id is required")
	}

	needSync, currentVersion, changeCount, changedIds, next, err := h.updateDataUsecase.RefreshSyncIds(
		ctx, ownerID, req.Resource, req.LastSync, req.Limit,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh sync ids failed: %v", err)
	}

	resp := &hubpb.RefreshSyncIdsResponse{
		NeedSync:       needSync,
		CurrentVersion: currentVersion,
		ChangeCount:    changeCount,
		ChangedIds:     changedIds,
		Next:           next,
	}
	return resp, nil
}

// FlushSyncIds xóa các sync IDs cũ (trim)
// @Summary Flush sync IDs cũ
// @Description Xóa các bản ghi sync cũ, giữ lại limit mới nhất
// @Tags UpdateData
// @Accept json
// @Produce json
// @Param resource path string true "Tên resource"
// @Param limit query int64 false "Số bản ghi giữ lại"
// @Success 200 {object} hubpb.FlushSyncIdsResponse
// @Router /v2/hub/update/{resource}/flush/ids [get]
func (h *UpdateDataHandler) FlushSyncIds(
	ctx context.Context,
	req *hubpb.FlushSyncIdsRequest,
) (*hubpb.FlushSyncIdsResponse, error) {
	if req == nil || req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required")
	}

	ownerID := _utils.GetProfileIdWithContext(ctx)
	if ownerID == 0 {
		return nil, status.Error(codes.Unauthenticated, "profile_id is required")
	}

	err := h.updateDataUsecase.FlushSyncIds(ctx, ownerID, req.Resource, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "flush sync ids failed: %v", err)
	}

	return &hubpb.FlushSyncIdsResponse{
		Message: "success",
	}, nil
}
