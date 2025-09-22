package handler

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	pb_bdspro "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type IncomeDocumentServer struct {
	pb_bdspro.UnimplementedIncomeDocumentServiceServer
	UC usecases.IncomeDocumentUsecase
}

func NewIncomeDocumentServer(IncomeDocumentUc usecases.IncomeDocumentUsecase) *IncomeDocumentServer {
	return &IncomeDocumentServer{
		UC: IncomeDocumentUc,
	}
}

func (s IncomeDocumentServer) Get(c context.Context, dto *pb_bdspro.IncomeDocumentDTO) (*pb_bdspro.IncomeDocumentListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_bdspro.IncomeDocumentDTO{}
	for _, r := range result {
		pbdata := &pb_bdspro.IncomeDocumentDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_bdspro.IncomeDocumentListDTO{
		Data: pbdatas,
	}, err
}
func (s IncomeDocumentServer) Detail(c context.Context, dto *pb_bdspro.IncomeDocumentDTO) (*pb_bdspro.IncomeDocumentDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_bdspro.IncomeDocumentDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s IncomeDocumentServer) Create(c context.Context, dto *pb_bdspro.IncomeDocumentDTO) (*pb_bdspro.IncomeDocumentDTO, error) {
	d := &domain.IncomeDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s IncomeDocumentServer) Update(c context.Context, dto *pb_bdspro.IncomeDocumentDTO) (*pb_bdspro.IncomeDocumentDTO, error) {
	d := &domain.IncomeDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s IncomeDocumentServer) Delete(c context.Context, dto *pb_bdspro.IncomeDocumentDTO) (*pb_bdspro.IncomeDocumentDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
