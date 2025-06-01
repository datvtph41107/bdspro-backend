package service

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetIncomeTypeServer struct {
	pb_bdspro.UnimplementedAssetIncomeTypeServiceServer
	UC usecases.AssetIncomeTypeUsecase
}

func NewAssetIncomeTypeServer(AssetIncomeTypeUc usecases.AssetIncomeTypeUsecase) *AssetIncomeTypeServer {
	return &AssetIncomeTypeServer{
		UC: AssetIncomeTypeUc,
	}
}

func (s AssetIncomeTypeServer) Get(c context.Context, dto *pb_bdspro.AssetIncomeTypeDTO) (*pb_bdspro.AssetIncomeTypeListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_bdspro.AssetIncomeTypeDTO{}
	for _, r := range result {
		pbdata := &pb_bdspro.AssetIncomeTypeDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.AssetIncomeTypeListDTO{
		Data: pbdatas,
	}, err
}
func (s AssetIncomeTypeServer) Detail(c context.Context, dto *pb_bdspro.AssetIncomeTypeDTO) (*pb_bdspro.AssetIncomeTypeDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_bdspro.AssetIncomeTypeDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s AssetIncomeTypeServer) Create(c context.Context, dto *pb_bdspro.AssetIncomeTypeDTO) (*pb_bdspro.AssetIncomeTypeDTO, error) {
	d := &domain.AssetIncomeType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetIncomeTypeServer) Update(c context.Context, dto *pb_bdspro.AssetIncomeTypeDTO) (*pb_bdspro.AssetIncomeTypeDTO, error) {
	d := &domain.AssetIncomeType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetIncomeTypeServer) Delete(c context.Context, dto *pb_bdspro.AssetIncomeTypeDTO) (*pb_bdspro.AssetIncomeTypeDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
