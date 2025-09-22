package handler

import (
	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"

	"github.com/jinzhu/copier"
)

type IncomeDocumentServer struct {
	bdspropb.UnimplementedIncomeDocumentServiceServer
	UC usecases.IncomeDocumentUsecase
}

func NewIncomeDocumentServer(IncomeDocumentUc usecases.IncomeDocumentUsecase) *IncomeDocumentServer {
	return &IncomeDocumentServer{
		UC: IncomeDocumentUc,
	}
}

func (s IncomeDocumentServer) Get(c context.Context, dto *bdspropb.IncomeDocumentDTO) (*bdspropb.IncomeDocumentListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*bdspropb.IncomeDocumentDTO{}
	for _, r := range result {
		pbdata := &bdspropb.IncomeDocumentDTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &bdspropb.IncomeDocumentListDTO{
		Data: pbdatas,
	}, err
}
func (s IncomeDocumentServer) Detail(c context.Context, dto *bdspropb.IncomeDocumentDTO) (*bdspropb.IncomeDocumentDTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &bdspropb.IncomeDocumentDTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s IncomeDocumentServer) Create(c context.Context, dto *bdspropb.IncomeDocumentDTO) (*bdspropb.IncomeDocumentDTO, error) {
	d := &domain.IncomeDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s IncomeDocumentServer) Update(c context.Context, dto *bdspropb.IncomeDocumentDTO) (*bdspropb.IncomeDocumentDTO, error) {
	d := &domain.IncomeDocument{}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s IncomeDocumentServer) Delete(c context.Context, dto *bdspropb.IncomeDocumentDTO) (*bdspropb.IncomeDocumentDTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
