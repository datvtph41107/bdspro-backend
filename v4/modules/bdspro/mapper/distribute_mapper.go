package mapper

import (
	"context"
	"fmt"

	_provider "common/domain/provider"
	_utils "common/utils"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type DistributionMapper struct {
	CrmClient _provider.CrmProvider
}

func NewDistributionMapper(crmClient _provider.CrmProvider) *DistributionMapper {
	return &DistributionMapper{CrmClient: crmClient}
}

func (m *DistributionMapper) distributionEntityToProto(row *entity.Distribution) *bdspropb.DistributionEntity {
	result := &bdspropb.DistributionEntity{
		Id:            row.ID,
		ProductId:     row.ProductID,
		Reason:        bdspropb.DistributionReason(row.Reason),
		ReasonOther:   row.ReasonOther,
		DistributedAt: _utils.FormatTimeToString(row.DistributedAt),
		RevokedAt:     _utils.FormatTimeToString(row.RevokedAt),
		FromDate:      _utils.FormatTimeToString(row.FromDate),
		ToDate:        _utils.FormatTimeToString(row.ToDate),
		CanDeal:       row.CanDeal,
		Note:          row.Note,
		CreatedAt:     _utils.FormatTimeToString(row.CreatedAt),
		UpdatedAt:     _utils.FormatTimeToString(row.UpdatedAt),
	}
	if row.ProductPrice != nil {
		if row.ProductPrice.SalePrice != nil {
			result.Price = int64(*row.ProductPrice.SalePrice)
		}
		if row.ProductPrice.SaleCommissionType != nil {
			result.CommissionType = bdspropb.CommissionType(*row.ProductPrice.SaleCommissionType)
		}
		if row.ProductPrice.SaleCommission != nil {
			result.CommissionValue = *row.ProductPrice.SaleCommission
		}
		result.ChannelPrice = row.ProductPrice.ChannelPrice
	}
	return result
}

func (m *DistributionMapper) DistributionWithPartnersToProto(ctx context.Context, item *dto.DistributionWithPartners) *bdspropb.DistributionEntity {
	partnerIds := make([]uint64, len(item.PartnerIds))
	for i, id := range item.PartnerIds {
		partnerIds[i] = uint64(id)
	}
	pbResult := m.distributionEntityToProto(&item.Distribution)
	pbResult.PartnerIds = partnerIds
	pbResult.PartnerCount = item.PartnerCount
	pbResult.Status = bdspropb.DistributionStatus(item.Status)
	pbResult.Price = item.Price
	pbResult.CommissionType = bdspropb.CommissionType(item.CommissionType)
	pbResult.CommissionValue = item.CommissionValue
	pbResult.Code = m.generateDistributionCode(item.ID)
	m.attachPartners(ctx, pbResult, partnerIds)
	return pbResult
}

func (m *DistributionMapper) DistributionWithPartnersListToProto(ctx context.Context, list []*dto.DistributionWithPartners) []*bdspropb.DistributionEntity {
	originIDSet := make(map[uint64]struct{})
	for _, e := range list {
		for _, id := range e.PartnerIds {
			if uint64(id) != 0 {
				originIDSet[uint64(id)] = struct{}{}
			}
		}
	}
	var originProfileMap map[uint64]*sharepb.OriginProfile
	if m.CrmClient != nil && len(originIDSet) > 0 {
		originIDs := make([]uint64, 0, len(originIDSet))
		for id := range originIDSet {
			originIDs = append(originIDs, id)
		}
		if profiles, err := m.CrmClient.GetOriginProfileByOriginIds(ctx, originIDs); err == nil {
			originProfileMap = profiles
		}
	}
	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		partnerIds := make([]uint64, len(e.PartnerIds))
		for j, id := range e.PartnerIds {
			partnerIds[j] = uint64(id)
		}
		pbResult := m.distributionEntityToProto(&e.Distribution)
		pbResult.PartnerIds = partnerIds
		pbResult.PartnerCount = e.PartnerCount
		pbResult.Status = bdspropb.DistributionStatus(e.Status)
		pbResult.Price = e.Price
		pbResult.CommissionType = bdspropb.CommissionType(e.CommissionType)
		pbResult.CommissionValue = e.CommissionValue
		pbResult.Code = m.generateDistributionCode(e.ID)
		if originProfileMap != nil {
			partners := make([]*sharepb.OriginProfile, 0)
			for _, partnerID := range partnerIds {
				if partnerID != 0 {
					if profile, ok := originProfileMap[partnerID]; ok {
						partners = append(partners, profile)
					}
				}
			}
			pbResult.Partners = partners
		}
		result[i] = pbResult
	}
	return result
}

func (m *DistributionMapper) attachPartners(ctx context.Context, pbResult *bdspropb.DistributionEntity, partnerIds []uint64) {
	if m.CrmClient == nil || len(partnerIds) == 0 {
		return
	}
	originIDs := make([]uint64, 0, len(partnerIds))
	for _, id := range partnerIds {
		if id != 0 {
			originIDs = append(originIDs, id)
		}
	}
	if len(originIDs) == 0 {
		return
	}
	originProfileMap, err := m.CrmClient.GetOriginProfileByOriginIds(ctx, originIDs)
	if err != nil || len(originProfileMap) == 0 {
		return
	}
	partners := make([]*sharepb.OriginProfile, 0, len(partnerIds))
	for _, partnerID := range partnerIds {
		if partnerID != 0 {
			if profile, ok := originProfileMap[partnerID]; ok {
				partners = append(partners, profile)
			}
		}
	}
	pbResult.Partners = partners
}

func (m *DistributionMapper) generateDistributionCode(distributionID uint64) string {
	if distributionID == 0 {
		return ""
	}
	return fmt.Sprintf("PP%08d", distributionID)
}

func (m *DistributionMapper) ConvertDistributionStatus(status []bdspropb.DistributionStatus) []enums.DistributionStatus {
	result := make([]enums.DistributionStatus, len(status))
	for i, s := range status {
		result[i] = enums.DistributionStatus(s)
	}
	return result
}
