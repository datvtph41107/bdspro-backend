package service

import (
	"bdspro/infra/mapper"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetCostService struct {
	pb_bdspro.UnimplementedAssetCostServiceServer
	UC *shared_usecase.AssetCostUsecase
}

func NewAssetCostService(AssetCostUc *shared_usecase.AssetCostUsecase) *AssetCostService {
	return &AssetCostService{
		UC: AssetCostUc,
	}
}

// @Summary Lấy danh sách chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost/{id}/{ownerType} [get]
func (s AssetCostService) Search(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AssetCostListDTO, error) {
	searchDTO := &dto.AssetCostSearchDTO{}
	result, total, err := s.UC.Find(c, enums.EOwnerType(req.OwnerType), searchDTO)

	pbdatas := mapper.AssetCostToPbList(result)

	return &pb_bdspro.AssetCostListDTO{
		Data:          pbdatas,
		TotalElements: uint32(total),
	}, err
}

// @Summary Lấy chi tiết chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost/detail/{id} [get]
func (s AssetCostService) Detail(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AssetCostDTO, error) {
	r, err := s.UC.GetByID(c, req.Id)

	pbdata := &pb_bdspro.AssetCostDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}

// @Summary Tạo chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param body body pb_bdspro.AssetCostDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost [post]
func (s AssetCostService) Create(c context.Context, req *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
	dto := &domain.AssetCost{}

	if dto.OwnerType == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ownerType is required",
		}
	}
	copier.Copy(&dto, &req)

	err := s.UC.Create(c, dto)

	pbdata := &pb_bdspro.AssetCostDTO{}
	copier.Copy(&pbdata, &dto)
	return pbdata, err
}

// @Summary Cập nhật chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param body body pb_bdspro.AssetCostDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost/{id} [put]
func (s AssetCostService) Update(c context.Context, req *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
	dto := &domain.AssetCost{}
	copier.Copy(&dto, &req)

	err := s.UC.Update(c, req.Id, dto)

	pbdata := &pb_bdspro.AssetCostDTO{}
	copier.Copy(&pbdata, &dto)
	return pbdata, err
}

// @Summary Xóa chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost/{id} [delete]
func (s AssetCostService) Delete(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	err := s.UC.Delete(c, req.Id)
	return &pb_bdspro.Response{
		Message: "success",
		Id:      req.Id,
	}, err
}
