package mapper

import (
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type ApartmentMapper struct {
}

func NewApartmentMapper() *ApartmentMapper {
	return &ApartmentMapper{}
}

func (m *ApartmentMapper) MapApartmentPb(apartment *domain.Apartment) *bdspropb.Apartment {
	return &bdspropb.Apartment{
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

func (m *ApartmentMapper) MapApartmentPbItem(apartment *domain.ApartmentItem) *bdspropb.Apartment {
	return &bdspropb.Apartment{
		Id:   apartment.ID,
		Name: apartment.Name,
		// Note:            apartment.Note,
		Floor:           int32(*apartment.Floor),
		Ordinal:         int32(apartment.Ordinal),
		ApartmentAttrId: apartment.ApartmentAttrID,
		Status:          int32(apartment.Status),
		Archived:        int32(apartment.Archived),
		BuildId:         apartment.BuildID,
	}
}

func (m *ApartmentMapper) MapApartmentPbItems(apartments []domain.ApartmentItem) []*bdspropb.Apartment {
	pbs := make([]*bdspropb.Apartment, len(apartments))
	for i, apartment := range apartments {
		pbs[i] = m.MapApartmentPbItem(&apartment)
	}
	return pbs
}
