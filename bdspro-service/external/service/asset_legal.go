package service

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type AssetLegalServer struct {
	pb_bdspro.UnimplementedAssetLegalServiceServer
	UC usecases.AssetLegalUsecase
}

func NewAssetLegalServer(AssetLegalUc usecases.AssetLegalUsecase) *AssetLegalServer {
	return &AssetLegalServer{
		UC: AssetLegalUc,
	}
}

func (s AssetLegalServer) Get(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_bdspro.AssetLegalDTO{}
	for _, r := range result {
		pbdata := &pb_bdspro.AssetLegalDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.AssetLegalListDTO{
		Data: pbdatas,
	}, err
}
func (s AssetLegalServer) Detail(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_bdspro.AssetLegalDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s AssetLegalServer) Create(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	d := &domain.AssetLegal{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetLegalServer) Update(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	d := &domain.AssetLegal{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s AssetLegalServer) Delete(c context.Context, dto *pb_bdspro.AssetLegalDTO) (*pb_bdspro.AssetLegalDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
