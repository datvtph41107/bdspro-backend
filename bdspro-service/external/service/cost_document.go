package service

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type CostDocumentServer struct {
	pb_bdspro.UnimplementedCostDocumentServiceServer
	UC usecases.CostDocumentUsecase
}

func NewCostDocumentServer(CostDocumentUc usecases.CostDocumentUsecase) *CostDocumentServer {
	return &CostDocumentServer{
		UC: CostDocumentUc,
	}
}

func (s CostDocumentServer) Get(c context.Context, dto *pb_bdspro.CostDocumentDTO) (*pb_bdspro.CostDocumentListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_bdspro.CostDocumentDTO{}
	for _, r := range result {
		pbdata := &pb_bdspro.CostDocumentDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.CostDocumentListDTO{
		Data: pbdatas,
	}, err
}
func (s CostDocumentServer) Detail(c context.Context, dto *pb_bdspro.CostDocumentDTO) (*pb_bdspro.CostDocumentDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_bdspro.CostDocumentDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s CostDocumentServer) Create(c context.Context, dto *pb_bdspro.CostDocumentDTO) (*pb_bdspro.CostDocumentDTO, error) {
	d := &domain.CostDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s CostDocumentServer) Update(c context.Context, dto *pb_bdspro.CostDocumentDTO) (*pb_bdspro.CostDocumentDTO, error) {
	d := &domain.CostDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s CostDocumentServer) Delete(c context.Context, dto *pb_bdspro.CostDocumentDTO) (*pb_bdspro.CostDocumentDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
