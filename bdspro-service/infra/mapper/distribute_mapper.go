package mapper

import (
	"bdspro/internal/domain"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type DistributionMapper struct{}

func NewDistributionMapper() *DistributionMapper {
	return &DistributionMapper{}
}

func (m *DistributionMapper) DistributionEntityToProto(entity *domain.DistributionEntity) *bdspropb.DistributionEntity {
	return &bdspropb.DistributionEntity{
		Id:                entity.ID,
		ProductId:         entity.ProductID,
		Price:             entity.Price,
		CanDeal:           entity.CanDeal,
		ChannelPrice:      entity.ChannelPrice,
		CommissionType:    bdspropb.CommissionType(entity.CommissionType),
		CommissionValue:   entity.CommissionValue,
		Note:              entity.Note,
		Status:            bdspropb.DistributionStatus(entity.Status),
		Reason:            bdspropb.DistributionReason(entity.Reason),
		ReasonOther:       entity.ReasonOther,
		DistributedAt:     _utils.FormatTimeToString(entity.DistributedAt),
		RevokedAt:         _utils.FormatTimeToString(entity.RevokedAt),
		FromDate:          _utils.FormatTimeToString(entity.FromDate),
		ToDate:            _utils.FormatTimeToString(entity.ToDate),
		PartnerIds:        entity.PartnerIds,
		PartnerCommission: entity.PartnerCommission,
		CreatedAt:         _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *DistributionMapper) DistributionListToProto(list []*domain.DistributionEntity) []*bdspropb.DistributionEntity {
	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		result[i] = m.DistributionEntityToProto(e)
	}
	return result
}
