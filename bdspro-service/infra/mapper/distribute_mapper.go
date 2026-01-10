package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type DistributionMapper struct{}

func NewDistributionMapper() *DistributionMapper {
	return &DistributionMapper{}
}

func (m *DistributionMapper) distributionEntityToProto(
	entity *domain.DistributionEntity,
) *bdspropb.DistributionEntity {

	return &bdspropb.DistributionEntity{
		Id:                entity.ID,
		ProductId:         entity.ProductID,
		Status:            bdspropb.DistributionStatus(entity.Status),
		Reason:            bdspropb.DistributionReason(entity.Reason),
		ReasonOther:       entity.ReasonOther,
		DistributedAt:     _utils.FormatTimeToString(entity.DistributedAt),
		RevokedAt:         _utils.FormatTimeToString(entity.RevokedAt),
		FromDate:          _utils.FormatTimeToString(entity.FromDate),
		ToDate:            _utils.FormatTimeToString(entity.ToDate),
		Price:             entity.Price,
		CanDeal:           entity.CanDeal,
		ChannelPrice:      entity.ChannelPrice,
		CommissionType:    bdspropb.CommissionType(entity.CommissionType),
		CommissionValue:   entity.CommissionValue,
		PartnerCommission: entity.PartnerCommission,
		Note:              entity.Note,
		CreatedAt:         _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *DistributionMapper) DistributionWithPartnersToProto(
	dto *dto.DistributionWithPartners,
) *bdspropb.DistributionEntity {
	partnerIds := make([]uint64, len(dto.PartnerIds))
	for i, id := range dto.PartnerIds {
		partnerIds[i] = uint64(id)
	}
	pbResult := m.distributionEntityToProto(&dto.DistributionEntity)
	pbResult.PartnerIds = partnerIds
	pbResult.PartnerCount = int32(dto.PartnerCount)
	return pbResult
}

func (m *DistributionMapper) DistributionWithPartnersListToProto(
	list []*dto.DistributionWithPartners,
) []*bdspropb.DistributionEntity {

	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		result[i] = m.DistributionWithPartnersToProto(e)
	}
	return result
}
