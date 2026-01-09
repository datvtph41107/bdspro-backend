package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
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
		Policy:            m.DistributionPolicyToProto(entity.Policy),
		Status:            bdspropb.DistributionStatus(entity.Status),
		Reason:            bdspropb.DistributionReason(entity.Reason),
		ReasonOther:       entity.ReasonOther,
		DistributedAt:     _utils.FormatTimeToString(entity.DistributedAt),
		RevokedAt:         _utils.FormatTimeToString(entity.RevokedAt),
		FromDate:          _utils.FormatTimeToString(entity.FromDate),
		ToDate:            _utils.FormatTimeToString(entity.ToDate),
		PartnerIds:        entity.PartnerIds,
		PartnerCommission: entity.PartnerCommission,
	}
}

func (m *DistributionMapper) DistributionListToProto(list []*domain.DistributionEntity) []*bdspropb.DistributionEntity {
	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		result[i] = m.DistributionEntityToProto(e)
	}
	return result
}

func (m *DistributionMapper) DistributionPolicyToProto(policy domain.DistributionPolicy) *bdspropb.DistributionPolicy {
	return &bdspropb.DistributionPolicy{
		Price:           policy.Price,
		CanDeal:         policy.CanDeal,
		ChannelPrice:    policy.ChannelPrice,
		CommissionType:  bdspropb.CommissionType(policy.CommissionType),
		CommissionValue: policy.CommissionValue,
		Note:            policy.Note,
	}
}

func (m *DistributionMapper) ProtoToDistributionPolicy(protoPolicy *bdspropb.DistributionPolicy) domain.DistributionPolicy {
	if protoPolicy == nil {
		return domain.DistributionPolicy{}
	}
	return domain.DistributionPolicy{
		Price:           protoPolicy.Price,
		CanDeal:         protoPolicy.CanDeal,
		ChannelPrice:    protoPolicy.ChannelPrice,
		CommissionType:  enums.CommissionType(protoPolicy.CommissionType),
		CommissionValue: protoPolicy.CommissionValue,
		Note:            protoPolicy.Note,
	}
}
