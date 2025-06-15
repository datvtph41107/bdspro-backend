package service

import (
	"bdspro/infra/mapper"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetCostTypeService struct {
	pb_bdspro.UnimplementedAssetCostTypeServiceServer
	UC *shared_usecase.AssetCostTypeUsecase
}

func NewAssetCostTypeServer(AssetCostTypeUc *shared_usecase.AssetCostTypeUsecase) *AssetCostTypeService {
	return &AssetCostTypeService{
		UC: AssetCostTypeUc,
	}
}

// @Summary Lấy danh sách loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type/{id}/{ownerType} [get]
func (s AssetCostTypeService) Search(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AssetCostTypeListDTO, error) {
	searchDTO := &dto.AssetCostTypeSearchDTO{}
	result, total, err := s.UC.Search(c, enums.EOwnerType(req.OwnerType), searchDTO)

	pbdatas := []*pb_bdspro.AssetCostTypeDTO{}
	copier.Copy(&pbdatas, &result)

	return &pb_bdspro.AssetCostTypeListDTO{
		Data:          pbdatas,
		TotalElements: uint32(total),
	}, err
}

func (s AssetCostTypeService) SearchGroup(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.GroupType, error) {
	searchDTO := &dto.AssetCostTypeSearchDTO{}
	incomes, spends, err := s.UC.SearchGroup(c, enums.EOwnerType(req.OwnerType), searchDTO)

	pbincomes := mapper.CostTypeToPbList(incomes)
	pbspends := mapper.CostTypeToPbList(spends)

	return &pb_bdspro.GroupType{
		Incomes: pbincomes,
		Spends:  pbspends,
	}, err
}

// @Summary Lấy chi tiết loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type/detail/{id} [get]
func (s AssetCostTypeService) Detail(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.AssetCostTypeDTO, error) {
	r, err := s.UC.GetByID(c, req.Id)

	pbdata := &pb_bdspro.AssetCostTypeDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}

// @Summary Tạo loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param body body pb_bdspro.AssetCostTypeDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type [post]
func (s AssetCostTypeService) Create(c context.Context, req *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	dto := &domain.AssetCostType{}
	copier.Copy(&dto, &req)

	err := s.UC.Create(c, dto)

	pbdata := &pb_bdspro.AssetCostTypeDTO{}
	copier.Copy(&pbdata, &dto)
	return pbdata, err
}

// @Summary Cập nhật loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param body body pb_bdspro.AssetCostTypeDTO true "Body"
// @Router /v2/bdspro/v2/asset/cost-type/{id} [put]
func (s AssetCostTypeService) Update(c context.Context, req *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	dto := &domain.AssetCostType{}
	copier.Copy(&dto, &req)

	err := s.UC.Update(c, req.Id, dto)

	pbdata := &pb_bdspro.AssetCostTypeDTO{}
	copier.Copy(&pbdata, &dto)
	return pbdata, err
}

// @Summary Xóa loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type/{id} [delete]
func (s AssetCostTypeService) Delete(c context.Context, req *pb_bdspro.IdRequest) (*pb_bdspro.Response, error) {
	err := s.UC.Delete(c, req.Id)
	return &pb_bdspro.Response{
		Message: "success",
		Id:      req.Id,
	}, err
}
