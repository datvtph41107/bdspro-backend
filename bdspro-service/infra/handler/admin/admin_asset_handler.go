package admin_handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"

	"github.com/jinzhu/copier"
)

type AdminAssetHandler struct {
	bdspropb.UnimplementedAdminAssetServiceServer
	AdminAssetService *admin_usecases.AdminAssetUsecase
	AssetMapper       *mapper.AssetMapper
	AuthGrpcClient    *client.AuthClient
	UserClient        *client.UserClient
}

func NewAdminAssetHandler(
	uc *admin_usecases.AdminAssetUsecase,
	assetMapper *mapper.AssetMapper,
	authGrpcClient *client.AuthClient,
	userClient *client.UserClient,
) *AdminAssetHandler {
	return &AdminAssetHandler{
		AdminAssetService: uc,
		AssetMapper:       assetMapper,
		AuthGrpcClient:    authGrpcClient,
		UserClient:        userClient,
	}
}

// @Summary Get assets
// @Description Get assets with advanced filtering
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Param name query string false "Asset name"
// @Param status query []uint32 false "Asset status"
// @Param archived query bool false "Archived status"
// @Param fromDate query string false "From date (ISO format)"
// @Param toDate query string false "To date (ISO format)"
// @Param text query string false "Search text"
// @Param provinceIds query []uint64 false "Province IDs"
// @Param districtIds query []uint64 false "District IDs"
// @Param wardIds query []uint64 false "Ward IDs"
// @Param propertyTypeIds query []uint64 false "Property type IDs"
// @Success 200 {object} bdspropb.AdminAssetSearchResponse
// @Router /admin/assets [get]
func (h *AdminAssetHandler) GetAssets(ctx context.Context, req *bdspropb.AdminAssetSearchRequest) (*bdspropb.AdminAssetSearchResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_XEM"}); err != nil {
		return nil, err
	}

	// Convert protobuf request to DTO
	searchReq := &dto.AdminAssetSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Name:            req.Name,
		Status:          req.Status,
		Archived:        &req.Archived,
		Text:            req.Text,
		ProvinceIds:     req.ProvinceIds,
		PropertyTypeIds: req.PropertyTypeIds,
	}

	// Parse dates if provided
	if req.FromDate != "" {
		if fromDate, err := time.Parse(time.RFC3339, req.FromDate); err == nil {
			searchReq.FromDate = &fromDate
		}
	}
	if req.ToDate != "" {
		if toDate, err := time.Parse(time.RFC3339, req.ToDate); err == nil {
			searchReq.ToDate = &toDate
		}
	}

	result, total, err := h.AdminAssetService.GetAdminAssets(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	assets := h.AssetMapper.MapAdminAssetList(ctx, result)
	return &bdspropb.AdminAssetSearchResponse{
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
		return nil, _errors.ReturnError(service.AssetNameRequired)
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
		return nil, _errors.ReturnError(service.AssetNameRequired)
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

// @Summary Get asset detail
// @Description Get detailed information of an asset
// @Tags AdminAsset
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} bdspropb.Asset
// @Router /admin/assets/{id} [get]
func (h *AdminAssetHandler) AssetDetail(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Asset, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TS_XEM"}); err != nil {
		return nil, err
	}

	asset, err := h.AdminAssetService.GetAssetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Get owner information from user service
	if asset.OwnerID != nil {
		owner := h.UserClient.GetProfileById(ctx, *asset.OwnerID)
		if owner != nil {
			asset.Owner = &domain.ProfileInfo{
				ProfileId: owner.Id,
				FullName:  owner.FullName,
			}
		}
	}

	return h.AssetMapper.MapAssetPb(asset), nil
}
