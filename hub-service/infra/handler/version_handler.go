package handler

import (
	"context"

	_dto "common/domain/dto"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type VersionHandler struct {
	hubpb.UnimplementedVersionServiceServer
	versionUsecase _usecase.IVersionUsecase
	versionMapper  *mapper.VersionMapper
}

func NewVersionHandler(
	versionUsecase _usecase.IVersionUsecase,
	versionMapper *mapper.VersionMapper,
) *VersionHandler {
	return &VersionHandler{
		versionUsecase: versionUsecase,
		versionMapper:  versionMapper,
	}
}

// @Summary Tạo phiên bản mới
// @Description Tạo một phiên bản ứng dụng mới
// @Tags Version
// @Accept json
// @Produce json
// @Param body body hubpb.CreateVersionRequest true "Thông tin phiên bản"
// @Success 200 {object} hubpb.CreateVersionResponse
// @Router /v2/hub/admin/versions [post]
func (h *VersionHandler) CreateVersion(ctx context.Context, req *hubpb.CreateVersionRequest) (*hubpb.CreateVersionResponse, error) {
	entity := h.versionMapper.CreateRequestToEntity(req)

	result, err := h.versionUsecase.CreateBundleVersion(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &hubpb.CreateVersionResponse{
		Data: h.versionMapper.EntityToProto(result),
	}, nil
}

// @Summary Cập nhật phiên bản
// @Description Cập nhật thông tin phiên bản theo ID
// @Tags Version
// @Accept json
// @Produce json
// @Param id path uint64 true "Version ID"
// @Param body body hubpb.UpdateVersionRequest true "Thông tin phiên bản"
// @Success 200 {object} hubpb.UpdateVersionResponse
// @Router /v2/hub/admin/versions/{id} [put]
func (h *VersionHandler) UpdateVersion(ctx context.Context, req *hubpb.UpdateVersionRequest) (*hubpb.UpdateVersionResponse, error) {
	updatedData := h.versionMapper.UpdateRequestToEntity(ctx, req)

	result, err := h.versionUsecase.UpdateBundleVersion(ctx, updatedData)
	if err != nil {
		return nil, err
	}

	return &hubpb.UpdateVersionResponse{
		Data: h.versionMapper.EntityToProto(result),
	}, nil
}

// @Summary Xóa phiên bản
// @Description Xóa (soft delete) một phiên bản theo ID
// @Tags Version
// @Accept json
// @Produce json
// @Param id path uint64 true "Version ID"
// @Success 200 {object} hubpb.DeleteVersionResponse
// @Router /v2/hub/admin/versions/{id} [delete]
func (h *VersionHandler) DeleteVersion(ctx context.Context, req *hubpb.DeleteVersionRequest) (*hubpb.DeleteVersionResponse, error) {
	_, err := h.versionUsecase.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &hubpb.DeleteVersionResponse{
		Success: true,
	}, nil
}

// @Summary Lấy chi tiết phiên bản
// @Description Lấy thông tin chi tiết của phiên bản theo ID
// @Tags Version
// @Accept json
// @Produce json
// @Param id path uint64 true "Version ID"
// @Success 200 {object} hubpb.GetVersionResponse
// @Router /v2/hub/admin/versions/{id} [get]
func (h *VersionHandler) GetVersion(ctx context.Context, req *hubpb.GetVersionRequest) (*hubpb.GetVersionResponse, error) {
	result, err := h.versionUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &hubpb.GetVersionResponse{
		Data: h.versionMapper.EntityToProto(result),
	}, nil
}

// @Summary Lấy danh sách phiên bản
// @Description Lấy danh sách phiên bản theo platform và phân trang
// @Tags Version
// @Accept json
// @Produce json
// @Param appName query string false "Tên ứng dụng"
// @Param platform query string false "Nền tảng"
// @Param page query int false "Số trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} hubpb.GetVersionListResponse
// @Router /v2/hub/admin/versions [get]
func (h *VersionHandler) GetVersionList(ctx context.Context, req *hubpb.GetVersionListRequest) (*hubpb.GetVersionListResponse, error) {
	pagable := &_dto.Pagable{
		Page: uint32(req.GetPage()),
		Size: uint32(req.GetSize()),
	}

	data, total, err := h.versionUsecase.GetListWithFilter(ctx, req.GetAppName(), req.GetPlatform(), pagable)
	if err != nil {
		return nil, err
	}

	return &hubpb.GetVersionListResponse{
		Data:  h.versionMapper.EntitiesToProto(data),
		Total: total,
	}, nil
}

// @Summary Kiểm tra cập nhật ứng dụng
// @Description Kiểm tra phiên bản ứng dụng hiện tại và xác định có cần cập nhật không
// @Tags Version/Public
// @Accept json
// @Produce json
// @Param appName path string true "Tên ứng dụng"
// @Param platform query string true "Nền tảng"
// @Param targetVersion query string false "Phiên bản hiện tại của client"
// @Success 200 {object} hubpb.CheckAppUpdateResponse
// @Router /v2/hub/app-bundle/update-check/{appName} [get]
func (h *VersionHandler) CheckAppUpdate(ctx context.Context, req *hubpb.CheckAppUpdateRequest) (*hubpb.CheckAppUpdateResponse, error) {
	latest, needUpdate, forceUpdate, err := h.versionUsecase.CheckAppUpdate(ctx, req.GetAppName(), req.GetPlatform(), req.GetTargetVersion())
	if err != nil {
		return nil, err
	}

	response := &hubpb.CheckAppUpdateResponse{
		NeedUpdate:  needUpdate,
		ForceUpdate: forceUpdate,
	}

	if latest != nil {
		response.LatestVersion = h.versionMapper.EntityToProto(latest)
	}

	return response, nil
}
