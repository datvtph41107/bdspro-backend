package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

type ApartmentAttributeMapper struct {
}

func NewApartmentAttributeMapper() *ApartmentAttributeMapper {
	return &ApartmentAttributeMapper{}
}

func (m *ApartmentAttributeMapper) MapApartmentAttributePb(apartmentAttribute *domain.ApartmentAttribute) *pb_bdspro.ApartmentAttribute {
	return &pb_bdspro.ApartmentAttribute{
		Id:           apartmentAttribute.ID,
		Area:         apartmentAttribute.Area,
		NumBedroom:   apartmentAttribute.NumBedroom,
		NumBathroom:  apartmentAttribute.NumBathroom,
		Furniture:    apartmentAttribute.Furniture,
		BlueprintUrl: apartmentAttribute.BlueprintUrl,
		Price:        apartmentAttribute.Price,
	}
}
