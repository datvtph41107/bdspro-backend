package service

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetCostTypeServer struct {
	pb_bdspro.UnimplementedAssetCostTypeServiceServer
	UC usecases.AssetCostTypeUsecase
}

func NewAssetCostTypeServer(AssetCostTypeUc usecases.AssetCostTypeUsecase) *AssetCostTypeServer {
	return &AssetCostTypeServer{
		UC: AssetCostTypeUc,
	}
}

func (s AssetCostTypeServer) Get(c context.Context, dto *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_bdspro.AssetCostTypeDTO{}
	for _, r := range result {
		pbdata := &pb_bdspro.AssetCostTypeDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.AssetCostTypeListDTO{
		Data: pbdatas,
	}, err
}
func (s AssetCostTypeServer) Detail(c context.Context, dto *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_bdspro.AssetCostTypeDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s AssetCostTypeServer) Create(c context.Context, dto *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	d := &domain.AssetCostType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetCostTypeServer) Update(c context.Context, dto *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	d := &domain.AssetCostType{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetCostTypeServer) Delete(c context.Context, dto *pb_bdspro.AssetCostTypeDTO) (*pb_bdspro.AssetCostTypeDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
