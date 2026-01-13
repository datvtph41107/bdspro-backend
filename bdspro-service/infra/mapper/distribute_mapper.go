package mapper

import (
	"bdspro/infra/client"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type DistributionMapper struct {
	CrmClient *client.CrmClient
}

func NewDistributionMapper(crmClient *client.CrmClient) *DistributionMapper {
	return &DistributionMapper{
		CrmClient: crmClient,
	}
}

func (m *DistributionMapper) distributionEntityToProto(
	entity *domain.DistributionEntity,
) *bdspropb.DistributionEntity {

	return &bdspropb.DistributionEntity{
		Id:                entity.ID,
		ProductId:         entity.ProductID,
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
	pbResult.Status = bdspropb.DistributionStatus(dto.Status)
	return pbResult
}

func (m *DistributionMapper) DistributionWithPartnersListToProto(
	ctx context.Context,
	list []*dto.DistributionWithPartners,
) []*bdspropb.DistributionEntity {
	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		result[i] = m.DistributionWithPartnersToProto(e)
	}

	// Collect unique origin IDs
	originIDSet := make(map[uint64]struct{})
	for _, e := range result {
		for _, id := range e.PartnerIds {
			if id != 0 {
				originIDSet[id] = struct{}{}
			}
		}
	}

	// Convert map to slice for gRPC request
	if len(originIDSet) == 0 {
		return result
	}

	originIDs := make([]uint64, 0, len(originIDSet))
	for id := range originIDSet {
		originIDs = append(originIDs, id)
	}

	// Call gRPC to get origin profiles
	originProfiles, err := m.CrmClient.GetOriginProfileByOriginIds(ctx, originIDs)
	if err != nil {
		// Log error but don't fail the entire request
		return result
	}

	// Create map for quick lookup
	originProfileMap := make(map[uint64]*sharepb.OriginProfile, len(originProfiles))
	for _, profile := range originProfiles {
		if profile != nil && profile.OriginId != 0 {
			originProfileMap[profile.OriginId] = profile
		}
	}

	// Map partners to each distribution entity
	for _, e := range result {
		partners := make([]*sharepb.OriginProfile, 0)
		for _, partnerId := range e.PartnerIds {
			if partnerId != 0 {
				if profile, ok := originProfileMap[partnerId]; ok {
					partners = append(partners, profile)
				}
			}
		}
		// TODO: Uncomment sau khi regenerate proto với field partners = 21
		// e.Partners = partners
		// e.PartnerCount = int32(len(partners))
		e.Partners = partners
	}
	return result
}
