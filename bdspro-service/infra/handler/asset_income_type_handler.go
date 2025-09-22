package handler

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetIncomeTypeServer struct {
	bdspropb.UnimplementedAssetIncomeTypeServiceServer
	UC usecases.AssetIncomeTypeUsecase
}

func NewAssetIncomeTypeServer(AssetIncomeTypeUc usecases.AssetIncomeTypeUsecase) *AssetIncomeTypeServer {
	return &AssetIncomeTypeServer{
		UC: AssetIncomeTypeUc,
	}
}

func (s AssetIncomeTypeServer) Get(c context.Context, dto *bdspropb.AssetIncomeTypeDTO) (*bdspropb.AssetIncomeTypeListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*bdspropb.AssetIncomeTypeDTO{}
	for _, r := range result {
		pbdata := &bdspropb.AssetIncomeTypeDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &bdspropb.AssetIncomeTypeListDTO{
		Data: pbdatas,
	}, err
}
func (s AssetIncomeTypeServer) Detail(c context.Context, dto *bdspropb.AssetIncomeTypeDTO) (*bdspropb.AssetIncomeTypeDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &bdspropb.AssetIncomeTypeDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s AssetIncomeTypeServer) Create(c context.Context, dto *bdspropb.AssetIncomeTypeDTO) (*bdspropb.AssetIncomeTypeDTO, error) {
	d := &domain.AssetIncomeType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetIncomeTypeServer) Update(c context.Context, dto *bdspropb.AssetIncomeTypeDTO) (*bdspropb.AssetIncomeTypeDTO, error) {
	d := &domain.AssetIncomeType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetIncomeTypeServer) Delete(c context.Context, dto *bdspropb.AssetIncomeTypeDTO) (*bdspropb.AssetIncomeTypeDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
