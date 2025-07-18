package service

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	pb_bdspro "pb/types/bdspro"
	sharepb "pb/types/shared"

	"github.com/jinzhu/copier"
)

type AssetService struct {
	pb_bdspro.UnimplementedAssetServiceServer
	UC            *shared_usecase.AssetUsecase
	Mapper        *mapper.AssetMapper
	ProfileClient *client.UserClient
}

func NewGrpcAssetService(uc *shared_usecase.AssetUsecase,
	mapper *mapper.AssetMapper,
	profileClient *client.UserClient) *AssetService {
	return &AssetService{
		UC:            uc,
		Mapper:        mapper,
		ProfileClient: profileClient,
	}
}

// @Summary Lấy danh sách tài sản của người dùng
// @Tags User: Tài sản
// @Produce json
// @Param query query dto.AssetSearchDTO true "Size"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/search [get]
func (s *AssetService) SearchInternal(ctx context.Context, req *pb_bdspro.AssetSearchRequest) (*pb_bdspro.AssetSearchResponse, error) {
	body := dto.AssetSearchDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.SearchAssetMe(ctx, body)
	if err != nil {
		return nil, err
	}

	assets := make([]*pb_bdspro.Asset, 0)
	copier.Copy(&assets, &result)

	return &pb_bdspro.AssetSearchResponse{
		Data:  assets,
		Total: uint64(total),
	}, nil
}

// @Summary Lấy danh sách tài sản chia sẻ
// @Tags User: Tài sản
// @Produce json
// @Param query query dto.AssetSearchDTO true "Size"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/search-share [get]
func (s *AssetService) SearchShare(ctx context.Context, req *pb_bdspro.AssetSearchRequest) (*pb_bdspro.AssetSearchResponse, error) {
	body := dto.AssetSearchDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.SearchAssetShare(ctx, body)
	if err != nil {
		return nil, err
	}

	assets := make([]*pb_bdspro.Asset, 0)
	copier.Copy(&assets, &result)

	return &pb_bdspro.AssetSearchResponse{
		Data:  assets,
		Total: uint64(total),
	}, nil
}

// @Summary Lấy danh sách tài sản có số lượng con
// @Tags User: Tài sản
// @Produce json
// @Param body query dto.AssetSearchDTO true "body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/split [get]
func (s *AssetService) SearchSplit(ctx context.Context, req *pb_bdspro.AssetSearchRequest) (*pb_bdspro.AssetSearchResponse, error) {
	body := dto.AssetSearchDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, total, err := s.UC.GetAssetsWithChildCount(ctx, body)
	if err != nil {
		return nil, err
	}

	assets := make([]*pb_bdspro.Asset, 0)
	copier.Copy(&assets, &result)

	return &pb_bdspro.AssetSearchResponse{
		Data:  assets,
		Total: uint64(total),
	}, nil
}

// @Summary Lấy theo ID
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body dto.AssetArchivedDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/archived/{id} [put]
func (s *AssetService) Archived(ctx context.Context, req *pb_bdspro.ArchivedRequest) (*pb_bdspro.Response, error) {
	body := dto.AssetArchivedDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	err := s.UC.ArchiveAsset(ctx, body.ID, body.Archived)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
		Id:      body.ID,
	}, nil
}

// @Summary Tạo mới
// @Tags User: Tài sản
// @Produce json
// @Param body body pb_bdspro.SaveAssetRequest true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset [post]
func (s *AssetService) CreateAsset(ctx context.Context, req *pb_bdspro.SaveAssetRequest) (*pb_bdspro.Response, error) {
	body := s.Mapper.AssetPbToDomain(req)

	result, err := s.UC.CreateAsset(ctx, body)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
		Id:      result.ID,
	}, nil
}

// @Summary Cập nhật
// @Tags User: Tài sản
// @Produce json
// @Param body body pb_bdspro.SaveAssetRequest true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset [put]
func (s *AssetService) UpdateAsset(ctx context.Context, req *pb_bdspro.SaveAssetRequest) (*pb_bdspro.Response, error) {
	// body := dto.AssetSaveRequest{}
	// body := s.Mapper.AssetPbToDomain(req)
	// if err := copier.Copy(&body, req); err != nil {
	// 	return nil, err
	// }
	body := s.Mapper.AssetPbToDomain(req)

	result, err := s.UC.UpdateAsset(ctx, req.Id, body)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
		Id:      result.ID,
	}, nil
}

// @Summary Xóa
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/{id} [delete]
func (s *AssetService) DeleteAsset(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	err := s.UC.DeleteAsset(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.Response{
		Message: "success",
	}, nil
}

// @Summary Xem chi tiết
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/detail/{id} [get]
func (s *AssetService) GetDetail(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.Asset, error) {
	result, err := s.UC.GetAssetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	asset := s.Mapper.AssetToPb(result)

	return asset, nil
}

// @Summary Gộp tài sản
// @Tags User: Tài sản
// @Produce json
// @Param body body dto.MergeAssetDTO true "body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/merge [post]
func (s *AssetService) Merge(ctx context.Context, req *pb_bdspro.MergeAssetRequest) (*pb_bdspro.Asset, error) {
	body := dto.MergeAssetDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.MergeAsset(ctx, body)
	if err != nil {
		return nil, err
	}

	asset := &pb_bdspro.Asset{}
	copier.Copy(asset, result)

	return asset, nil
}

// @Summary Tách tài sản
// @Tags User: Tài sản
// @Produce json
// @Param body body dto.SplitAssetDTO true "body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/split [post]
func (s *AssetService) Split(ctx context.Context, req *pb_bdspro.SplitAssetRequest) (*pb_bdspro.Asset, error) {
	body := dto.SplitAssetDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.SplitAsset(ctx, body)
	if err != nil {
		return nil, err
	}

	asset := &pb_bdspro.Asset{}
	copier.Copy(asset, result)

	return asset, nil
}

// // @Summary Lịch sử tách gộp tài sản dành cho user
// // @Tags User: Tài sản
// // @Produce json
// // @Param body query dto.AssetHistoryDTO true "body"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/asset/history-owner [get]
// func (s *AssetService) HistoryOwner(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.HistoryResponse, error) {
// 	body := dto.AssetHistoryDTO{}
// 	if err := copier.Copy(&body, req); err != nil {
// 		return nil, err
// 	}

// 	result, err := s.UC.HistoryOwner(ctx, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	history := make([]*pb_bdspro.History, 0)
// 	copier.Copy(&history, &result)

// 	return &pb_bdspro.HistoryResponse{
// 		Data: history,
// 	}, nil
// }

// @Summary Lấy danh sách tài sản đã đăng
// @Tags User: Tài sản
// @Produce json
// @Param profileId path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/publish/{profileId} [get]
func (s *AssetService) GetPublish(ctx context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AssetSearchResponse, error) {
	body := dto.AssetSearchDTO{}
	if err := copier.Copy(&body, req); err != nil {
		return nil, err
	}

	result, err := s.UC.GetPublish(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	assets := make([]*pb_bdspro.Asset, 0)
	copier.Copy(&assets, &result)

	return &pb_bdspro.AssetSearchResponse{
		Data: assets,
	}, nil
}
