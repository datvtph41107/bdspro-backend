package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_provider "common/domain/provider"
	_usecase "common/domain/usecase"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type DistributionMapper struct {
	CrmClient       _provider.CrmProvider
	CodeDataUsecase *_usecase.CodeDataUsecase
}

func NewDistributionMapper(crmClient _provider.CrmProvider, codeDataUsecase *_usecase.CodeDataUsecase) *DistributionMapper {
	return &DistributionMapper{
		CrmClient:       crmClient,
		CodeDataUsecase: codeDataUsecase,
	}
}

func (m *DistributionMapper) distributionEntityToProto(
	entity *domain.DistributionEntity,
) *bdspropb.DistributionEntity {

	result := &bdspropb.DistributionEntity{
		Id:            entity.ID,
		ProductId:     entity.ProductID,
		Reason:        bdspropb.DistributionReason(entity.Reason),
		ReasonOther:   entity.ReasonOther,
		DistributedAt: _utils.FormatTimeToString(entity.DistributedAt),
		RevokedAt:     _utils.FormatTimeToString(entity.RevokedAt),
		FromDate:      _utils.FormatTimeToString(entity.FromDate),
		ToDate:        _utils.FormatTimeToString(entity.ToDate),
		CanDeal:       entity.CanDeal,
		// Price:             entity.Price,
		// CommissionType:    bdspropb.CommissionType(entity.CommissionType),
		// CommissionValue:   entity.CommissionValue,
		// PartnerCommission: entity.PartnerCommission,
		Note:      entity.Note,
		CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(entity.UpdatedAt),
	}

	if entity.ProductPrice != nil {
		result.Price = int64(*entity.ProductPrice.SalePrice)
		result.CommissionType = bdspropb.CommissionType(*entity.ProductPrice.SaleCommissionType)
		result.CommissionValue = *entity.ProductPrice.SaleCommission
		result.ChannelPrice = entity.ProductPrice.ChannelPrice
	}

	return result
}

func (m *DistributionMapper) DistributionWithPartnersToProto(
	ctx context.Context,
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
	pbResult.Price = dto.Price
	pbResult.CommissionType = bdspropb.CommissionType(dto.CommissionType)
	pbResult.CommissionValue = dto.CommissionValue

	// Generate distribution code
	pbResult.Code = m.generateDistributionCode(ctx, dto.ID)

	// Call CRM service to get origin profiles for partners
	if len(partnerIds) > 0 {
		originIDs := make([]uint64, 0)
		for _, id := range partnerIds {
			if id != 0 {
				originIDs = append(originIDs, id)
			}
		}

		if len(originIDs) > 0 {
			originProfiles, err := m.CrmClient.GetOriginProfileByOriginIds(ctx, originIDs)
			if err == nil && len(originProfiles) > 0 {
				// Create map for quick lookup
				originProfileMap := make(map[uint64]*sharepb.OriginProfile, len(originProfiles))
				for _, profile := range originProfiles {
					if profile != nil && profile.OriginId != 0 {
						originProfileMap[profile.OriginId] = profile
					}
				}

				// Map partners to distribution entity
				partners := make([]*sharepb.OriginProfile, 0)
				for _, partnerId := range partnerIds {
					if partnerId != 0 {
						if profile, ok := originProfileMap[partnerId]; ok {
							partners = append(partners, profile)
						}
					}
				}
				pbResult.Partners = partners
			}
		}
	}

	return pbResult
}

func (m *DistributionMapper) DistributionWithPartnersListToProto(
	ctx context.Context,
	list []*dto.DistributionWithPartners,
) []*bdspropb.DistributionEntity {
	// Collect unique origin IDs first for batch call
	originIDSet := make(map[uint64]struct{})
	for _, e := range list {
		for _, id := range e.PartnerIds {
			if uint64(id) != 0 {
				originIDSet[uint64(id)] = struct{}{}
			}
		}
	}

	// Convert map to slice for gRPC request
	originIDs := make([]uint64, 0, len(originIDSet))
	for id := range originIDSet {
		originIDs = append(originIDs, id)
	}

	// Call gRPC to get origin profiles once for all items
	var originProfileMap map[uint64]*sharepb.OriginProfile
	if len(originIDs) > 0 {
		originProfiles, err := m.CrmClient.GetOriginProfileByOriginIds(ctx, originIDs)
		if err == nil && len(originProfiles) > 0 {
			originProfileMap = make(map[uint64]*sharepb.OriginProfile, len(originProfiles))
			for _, profile := range originProfiles {
				if profile != nil && profile.OriginId != 0 {
					originProfileMap[profile.OriginId] = profile
				}
			}
		}
	}

	// Map each distribution entity with partners
	result := make([]*bdspropb.DistributionEntity, len(list))
	for i, e := range list {
		partnerIds := make([]uint64, len(e.PartnerIds))
		for j, id := range e.PartnerIds {
			partnerIds[j] = uint64(id)
		}
		pbResult := m.distributionEntityToProto(&e.DistributionEntity)
		pbResult.PartnerIds = partnerIds
		pbResult.PartnerCount = int32(e.PartnerCount)
		pbResult.Status = bdspropb.DistributionStatus(e.Status)
		pbResult.Price = e.Price
		pbResult.CommissionType = bdspropb.CommissionType(e.CommissionType)
		pbResult.CommissionValue = e.CommissionValue

		// Generate distribution code
		pbResult.Code = m.generateDistributionCode(ctx, e.ID)

		// Map partners from pre-fetched origin profiles
		if originProfileMap != nil {
			partners := make([]*sharepb.OriginProfile, 0)
			for _, partnerId := range partnerIds {
				if partnerId != 0 {
					if profile, ok := originProfileMap[partnerId]; ok {
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

// generateDistributionCode generates distribution code based on ID
func (m *DistributionMapper) generateDistributionCode(ctx context.Context, distributionId uint64) string {
	if m.CodeDataUsecase == nil {
		return ""
	}
	return m.CodeDataUsecase.GetDistributionCode(ctx, distributionId)
}

func (m *DistributionMapper) ConvertDistributionStatus(status []bdspropb.DistributionStatus) []enums.DistributionStatus {
	result := make([]enums.DistributionStatus, len(status))
	for i, s := range status {
		result[i] = enums.DistributionStatus(s)
	}
	return result
}
