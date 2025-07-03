package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	pb_bdspro "pb/types/bdspro"
)

type AssetLegalMapper struct{}

func NewAssetLegalMapper() *AssetLegalMapper {
	return &AssetLegalMapper{}
}

func (m *AssetLegalMapper) AssetLegalToPb(entity *domain.AssetLegal) *pb_bdspro.AssetLegalDTO {
	return &pb_bdspro.AssetLegalDTO{
		Id:                  entity.ID,
		DocumentName:        entity.DocumentName,
		DocumentUrl:         entity.DocumentURL,
		IssuedDate:          entity.IssuedDate,
		ExpiryDate:          entity.ExpiryDate,
		Description:         entity.Description,
		AssetId:             entity.AssetID,
		RelatedSplitMergeId: entity.RelatedSplitMergeID,
	}
}

func (m *AssetLegalMapper) PbToListAssetLegal(entities []*pb_bdspro.LegalItem) *[]domain.AssetLegal {
	domainLegalItems := make([]domain.AssetLegal, len(entities))
	for i, legalItem := range entities {
		domainLegalItems[i] = m.PbToDomain(legalItem)
	}
	return &domainLegalItems
}

func (m *AssetLegalMapper) PbToDomain(entity *pb_bdspro.LegalItem) domain.AssetLegal {
	return domain.AssetLegal{
		BaseEntity: _models.BaseEntity{
			ID: entity.Id,
		},
		DocumentName: entity.DocumentName,
		DocumentURL:  entity.DocumentUrl,
		IssuedDate:   entity.IssuedDate,
		ExpiryDate:   entity.ExpiryDate,
		Description:  entity.Description,
		AssetID:      entity.AssetId,
	}
}
