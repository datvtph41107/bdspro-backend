package mapper

import (
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type ApartmentAttributeMapper struct {
}

func NewApartmentAttributeMapper() *ApartmentAttributeMapper {
	return &ApartmentAttributeMapper{}
}

func (m *ApartmentAttributeMapper) MapApartmentAttributePb(apartmentAttribute *domain.ApartmentAttribute) *bdspropb.ApartmentAttribute {
	return &bdspropb.ApartmentAttribute{
		Id:           apartmentAttribute.ID,
		Area:         apartmentAttribute.Area,
		NumBedroom:   apartmentAttribute.NumBedroom,
		NumBathroom:  apartmentAttribute.NumBathroom,
		Furniture:    apartmentAttribute.Furniture,
		BlueprintUrl: apartmentAttribute.BlueprintUrl,
		Price:        apartmentAttribute.Price,
	}
}

func (m *ApartmentAttributeMapper) MapApartmentAttributePbItem(apartmentAttribute *domain.ApartmentAttrItem) *bdspropb.ApartmentAttribute {
	result := &bdspropb.ApartmentAttribute{
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

func (m *ApartmentAttributeMapper) MapApartmentAttributePbItems(apartmentAttributes []domain.ApartmentAttrItem) []*bdspropb.ApartmentAttribute {
	pbs := make([]*bdspropb.ApartmentAttribute, len(apartmentAttributes))
	for i, apartmentAttribute := range apartmentAttributes {
		pbs[i] = m.MapApartmentAttributePbItem(&apartmentAttribute)
	}
	return pbs
}
