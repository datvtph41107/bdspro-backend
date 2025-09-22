package handler

import (
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	bdspropb "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetCostTypeService struct {
	bdspropb.UnimplementedAssetCostTypeServiceServer
	Usecase   *shared_usecase.AssetCostTypeUsecase
	Validator *validator.CostTypeValidator
	Mapper    *mapper.CostTypeMapper
}

func NewAssetCostTypeServer(
	AssetCostTypeUc *shared_usecase.AssetCostTypeUsecase,
	validator *validator.CostTypeValidator,
	mapper *mapper.CostTypeMapper,
) *AssetCostTypeService {
	return &AssetCostTypeService{
		Usecase:   AssetCostTypeUc,
		Validator: validator,
		Mapper:    mapper,
	}
}

// @Summary Lấy danh sách loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param body query bdspropb.CostTypeRequest true "Body"
// @Security BearerAuth
// @Param assetId path uint64 true "ID"
// @Param ownerType path enums.EOwnerOf true "Owner Type"
// @Router /v2/bdspro/v2/asset/cost-type/list/{assetId}/{ownerType} [get]
func (s AssetCostTypeService) Search(c context.Context, req *bdspropb.CostTypeRequest) (*bdspropb.AssetCostTypeListDTO, error) {
	searchDTO := &dto.AssetCostTypeSearchDTO{
		Type: req.Type,
	}
	copier.Copy(&searchDTO, &req)
	result, total, err := s.Usecase.Search(c, enums.EOwnerOf(req.OwnerType), searchDTO)

	pbdatas := []*bdspropb.AssetCostTypeDTO{}
	copier.Copy(&pbdatas, &result)

	return &bdspropb.AssetCostTypeListDTO{
		Data:          pbdatas,
		TotalElements: uint32(total),
	}, err
}

// @Summary Lấy danh sách loại chi phí theo nhóm
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param ownerType path enums.EOwnerOf true "Owner Type"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type/group/{id}/{ownerType} [get]
func (s AssetCostTypeService) SearchGroup(c context.Context, req *bdspropb.IdRequest) (*bdspropb.GroupType, error) {
	searchDTO := &dto.AssetCostTypeSearchDTO{}
	incomes, spends, err := s.Usecase.SearchGroup(c, enums.EOwnerOf(req.OwnerType), searchDTO)

	pbincomes := s.Mapper.CostTypeToPbList(incomes)
	pbspends := s.Mapper.CostTypeToPbList(spends)

	return &bdspropb.GroupType{
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
func (s AssetCostTypeService) Detail(c context.Context, req *bdspropb.IdRequest) (*bdspropb.AssetCostTypeDTO, error) {
	r, err := s.Usecase.GetByID(c, req.Id)

	pbdata := s.Mapper.CostTypeToPb(r)

	return pbdata, err
}

// @Summary Tạo loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param body body bdspropb.AssetCostTypeDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type [post]
func (s AssetCostTypeService) Create(c context.Context, req *bdspropb.AssetCostTypeDTO) (*bdspropb.AssetCostTypeDTO, error) {
	dto := s.Mapper.CostTypePbToDomain(req)
	err := s.Validator.ValidateSaveCostType(req)
	if err != nil {
		return nil, err
	}
	err = s.Usecase.Create(c, dto)

	pbdata := s.Mapper.CostTypeToPb(dto)
	return pbdata, err
}

// @Summary Cập nhật loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Param body body bdspropb.AssetCostTypeDTO true "Body"
// @Router /v2/bdspro/v2/asset/cost-type/{id} [put]
func (s AssetCostTypeService) Update(c context.Context, req *bdspropb.AssetCostTypeDTO) (*bdspropb.AssetCostTypeDTO, error) {
	dto := s.Mapper.CostTypePbToDomain(req)

	err := s.Validator.ValidateSaveCostType(req)
	if err != nil {
		return nil, err
	}

	err = s.Usecase.Update(c, req.Id, dto)

	pbdata := &bdspropb.AssetCostTypeDTO{}
	copier.Copy(&pbdata, &dto)
	return pbdata, err
}

// @Summary Xóa loại chi phí
// @Tags Hợp đồng & Thu chi & Loại chi phí
// @Produce json
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset/cost-type/{id} [delete]
func (s AssetCostTypeService) Delete(c context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.Usecase.Delete(c, req.Id)
	return &bdspropb.Response{
		Message: "success",
		Id:      req.Id,
	}, err
}
