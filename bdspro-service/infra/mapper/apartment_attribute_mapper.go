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

func (m *ApartmentAttributeMapper) MapApartmentAttributePbItem(apartmentAttribute *domain.ApartmentAttrItem) *pb_bdspro.ApartmentAttribute {
	result := &pb_bdspro.ApartmentAttribute{
		Id:           apartmentAttribute.ID,
		Area:         apartmentAttribute.Area,
		Furniture:    apartmentAttribute.Furniture,
		BlueprintUrl: apartmentAttribute.BlueprintUrl,
		Price:        apartmentAttribute.Price,
	}

	if apartmentAttribute.NumBedroom != nil {
		// numBedroom := uint32(*apartmentAttribute.NumBedroom)
		// result.NumBedroom = &numBedroom
	}
	if apartmentAttribute.NumBathroom != nil {
		numBathroom := uint32(*apartmentAttribute.NumBathroom)
		result.NumBathroom = &numBathroom
	}

	return result
}

func (m *ApartmentAttributeMapper) MapApartmentAttributePbItems(apartmentAttributes []domain.ApartmentAttrItem) []*pb_bdspro.ApartmentAttribute {
	pbs := make([]*pb_bdspro.ApartmentAttribute, len(apartmentAttributes))
	for i, apartmentAttribute := range apartmentAttributes {
		pbs[i] = m.MapApartmentAttributePbItem(&apartmentAttribute)
	}
	return pbs
}
