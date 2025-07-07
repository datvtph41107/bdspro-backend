package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

type ApartmentMapper struct {
}

func NewApartmentMapper() *ApartmentMapper {
	return &ApartmentMapper{}
}

func (m *ApartmentMapper) MapApartmentPb(apartment *domain.Apartment) *pb_bdspro.Apartment {
	return &pb_bdspro.Apartment{
		Id:              apartment.ID,
		Name:            apartment.Name,
		Note:            apartment.Note,
		Floor:           int32(*apartment.Floor),
		Ordinal:         int32(apartment.Ordinal),
		ApartmentAttrId: apartment.ApartmentAttrID,
		Status:          int32(apartment.Status),
		Archived:        int32(apartment.Archived),
		BuildId:         apartment.BuildID,
	}
}
