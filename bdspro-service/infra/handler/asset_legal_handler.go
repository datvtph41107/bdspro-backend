package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"
	sharepb "pb/types/shared"

	"github.com/jinzhu/copier"
)

type AssetLegalService struct {
	pb_bdspro.UnimplementedAssetLegalServiceServer
	UC     *usecases.AssetLegalUsecase
	Mapper *mapper.AssetLegalMapper
}

func NewAssetLegalServer(AssetLegalUc *usecases.AssetLegalUsecase, mapper *mapper.AssetLegalMapper) *AssetLegalService {
	return &AssetLegalService{
		UC:     AssetLegalUc,
		Mapper: mapper,
	}
}

// @Summary Lấy danh sách giấy tờ
// @Description Lấy danh sách giấy tờ
// @Tags AssetLegal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body pb_bdspro.AssetLegalDTO true "Data"
// @Success 200 {object} pb_bdspro.AssetLegalListDTO
// @Router /v2/bdspro/v2/asset-legal [get]
func (s AssetLegalService) Get(c context.Context, data *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalListDTO, error) {
	query := dto.AssetLegalGetDTO{}
	copier.Copy(&query, &data)
	result, _, err := s.UC.GetData(c, query)

	pbdatas := []*pb_bdspro.AssetLegalDTO{}
	for _, r := range result {
		pbdata := s.Mapper.AssetLegalToPb(&r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.AssetLegalListDTO{
		Data: pbdatas,
	}, err
}

// @Summary Lấy chi tiết giấy tờ
// @Description Lấy chi tiết giấy tờ
// @Tags AssetLegal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} pb_bdspro.AssetLegalDTO
// @Router /v2/bdspro/v2/asset-legal/{id} [get]
func (s AssetLegalService) Detail(c context.Context, dto *sharepb.IdRequest) (*pb_bdspro.AssetLegalDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := s.Mapper.AssetLegalToPb(r)

	return pbdata, err
}

// @Summary Tạo giấy tờ
// @Description Tạo giấy tờ
// @Tags AssetLegal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body pb_bdspro.AssetLegalDTO true "Data"
// @Success 200 {object} pb_bdspro.AssetLegalDTO
// @Router /v2/bdspro/v2/asset-legal [post]
func (s AssetLegalService) Create(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	d := &domain.AssetLegal{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	pbdata := s.Mapper.AssetLegalToPb(r)
	return pbdata, err
}

// @Summary Cập nhật giấy tờ
// @Description Cập nhật giấy tờ
// @Tags AssetLegal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Param data body pb_bdspro.AssetLegalDTO true "Data"
// @Success 200 {object} pb_bdspro.AssetLegalDTO
// @Router /v2/bdspro/v2/asset-legal/{id} [put]
func (s AssetLegalService) Update(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	d := &domain.AssetLegal{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	pbdata := s.Mapper.AssetLegalToPb(r)
	return pbdata, err
}

// @Summary Xóa giấy tờ
// @Description Xóa giấy tờ
// @Tags AssetLegal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID"
// @Success 200 {object} sharepb.Empty
// @Router /v2/bdspro/v2/asset-legal/{id} [delete]
func (s AssetLegalService) Delete(c context.Context, dto *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := s.UC.Delete(c, dto.Id)
	return &sharepb.Empty{}, err
}
