package admin_handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/dto"
	_errors "common/errors"
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"github.com/jinzhu/copier"
)

type AdminAssetHandler struct {
	bdspropb.UnimplementedAdminAssetServiceServer
	AdminAssetService *admin_usecases.AdminAssetUsecase
	AssetMapper       *mapper.AssetMapper
	AuthGrpcClient    *clients.AuthGrpcClient
}

func NewAdminAssetHandler(
	uc *admin_usecases.AdminAssetUsecase,
	assetMapper *mapper.AssetMapper,
	authGrpcClient *clients.AuthGrpcClient,
) *AdminAssetHandler {
	return &AdminAssetHandler{
		AdminAssetService: uc,
		AssetMapper:       assetMapper,
		AuthGrpcClient:    authGrpcClient,
	}
}

// @Summary Get assets
// @Description Get assets
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Success 200 {object} bdspropb.AssetSearchResponse
// @Router /admin/assets [get]
func (h *AdminAssetHandler) GetAssets(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.AssetSearchResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_XEM"}); err != nil {
		return nil, err
	}

	result, total, err := h.AdminAssetService.GetList(ctx,
		&dto.AssetSearchRequest{
			Pagable: _dto.Pagable{
				Page: req.Page,
				Size: req.Size,
			},
		})
	if err != nil {
		return nil, err
	}

	assets := h.AssetMapper.MapAssetPbList(result)
	return &bdspropb.AssetSearchResponse{
		Data:  assets,
		Total: uint64(total),
	}, nil
}

// @Summary Create a new asset
// @Description Create a new asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param asset body bdspropb.SaveAssetRequest true "Asset to create"
// @Success 200 {object} bdspropb.Asset
// @Router /admin/assets [post]
func (h *AdminAssetHandler) CreateAsset(ctx context.Context, req *bdspropb.SaveAssetRequest) (*bdspropb.Asset, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_TAO"}); err != nil {
		return nil, err
	}

	asset := h.AssetMapper.AssetPbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "Tên tài sản không được để trống")
	}

	createdAsset, err := h.AdminAssetService.Create(ctx, asset)
	if err != nil {
		return nil, err
	}

	return h.AssetMapper.MapAssetPb(createdAsset), nil
}

// @Summary Update an asset
// @Description Update an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param asset body bdspropb.SaveAssetRequest true "Asset to update"
// @Success 200 {object} bdspropb.Asset
// @Router /admin/assets [put]
func (h *AdminAssetHandler) UpdateAsset(ctx context.Context, req *bdspropb.SaveAssetRequest) (*bdspropb.Asset, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_SUA"}); err != nil {
		return nil, err
	}

	asset := h.AssetMapper.AssetPbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "Tên tài sản không được để trống")
	}

	updatedAsset, err := h.AdminAssetService.Update(ctx, req.Id, asset)
	if err != nil {
		return nil, err
	}

	return h.AssetMapper.MapAssetPb(updatedAsset), nil
}

// @Summary Delete an asset
// @Description Delete an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/assets/{id} [delete]
func (h *AdminAssetHandler) DeleteAsset(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_XOA"}); err != nil {
		return nil, err
	}

	_, err := h.AdminAssetService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Xóa tài sản thành công",
	}, nil
}

// @Summary Approve an asset
// @Description Approve an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/assets/{id}/approve [post]
func (h *AdminAssetHandler) ApproveAsset(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_DUYET"}); err != nil {
		return nil, err
	}

	err := h.AdminAssetService.Approve(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Duyệt tài sản thành công",
	}, nil
}

// @Summary Reject an asset
// @Description Reject an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/assets/{id}/reject [post]
func (h *AdminAssetHandler) RejectAsset(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_DUYET"}); err != nil {
		return nil, err
	}

	err := h.AdminAssetService.Reject(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Từ chối tài sản thành công",
	}, nil
}

// @Summary Archive an asset
// @Description Archive an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param request body bdspropb.ArchivedRequest true "Archive request"
// @Success 200 {object} bdspropb.Response
// @Router /admin/assets/archived [put]
func (h *AdminAssetHandler) ArchiveAsset(ctx context.Context, req *bdspropb.ArchivedRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_LUU_TRU"}); err != nil {
		return nil, err
	}

	body := dto.ArchivedRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	err := h.AdminAssetService.Archive(ctx, &body)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Lưu trữ tài sản thành công",
	}, nil
}

// @Summary Merge assets
// @Description Merge multiple assets into one
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param request body bdspropb.MergeAssetRequest true "Merge request"
// @Success 200 {object} bdspropb.Asset
// @Router /admin/assets/merge [post]
func (h *AdminAssetHandler) MergeAssets(ctx context.Context, req *bdspropb.MergeAssetRequest) (*bdspropb.Asset, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_HOP_NHAT"}); err != nil {
		return nil, err
	}

	body := dto.MergeAssetRequest{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	mergedAsset, err := h.AdminAssetService.Merge(ctx, &body)
	if err != nil {
		return nil, err
	}

	return h.AssetMapper.MapAssetPb(mergedAsset), nil
}

// @Summary Split asset
// @Description Split an asset into multiple assets
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param request body bdspropb.SplitAssetRequest true "Split request"
// @Success 200 {object} bdspropb.Asset
// @Router /admin/assets/split [post]
func (h *AdminAssetHandler) SplitAsset(ctx context.Context, req *bdspropb.SplitAssetRequest) (*bdspropb.Asset, error) {
	body := dto.SplitAssetDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_TACH"}); err != nil {
		return nil, err
	}

	splitAsset, err := h.AdminAssetService.Split(ctx, &body)
	if err != nil {
		return nil, err
	}

	return h.AssetMapper.MapAssetPb(splitAsset), nil
}
