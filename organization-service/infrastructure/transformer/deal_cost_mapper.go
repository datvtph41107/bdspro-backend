package transformer

import (
	_models "common/models"
	_utils "common/utils"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/enums"
	organizationpb "pb/types/organization"
)

type DealCostMapper struct {
}

func NewDealCostMapper() *DealCostMapper {
	return &DealCostMapper{}
}

func (m *DealCostMapper) DealCostModelToProto(e *entity.DealCost) *organizationpb.DealCostDTO {
	result := &organizationpb.DealCostDTO{
		Id:           e.ID,
		DealId:       e.DealID,
		Amount:       e.Amount,
		ConfirmedAt:  _utils.FormatTimeToString(e.ConfirmedAt),
		Note:         e.Note,
		CostTypeId:   e.CostTypeID,
		CostTypeName: e.CostName,
		Status:       uint32(e.Status),
	}

	if e.CostType != nil {
		result.CostTypeName = e.CostType.TypeName
	}

	return result
}

// DealCostProtoToModel maps protobuf DTO to DealCost model
func (m *DealCostMapper) DealCostProtoToModel(p *organizationpb.DealCostDTO) *entity.DealCost {
	var confirmedAt *time.Time
	if p.ConfirmedAt != "" {
		confirmedAt = _utils.ParseStringToTime(p.ConfirmedAt)
	}

	return &entity.DealCost{
		BaseEntity:  _models.BaseEntity{ID: p.Id},
		DealID:      p.DealId,
		Amount:      p.Amount,
		ConfirmedAt: confirmedAt,
		Note:        p.Note,
		CostTypeID:  p.CostTypeId,
		CostName:    p.CostName,
		Status:      enums.ApprovedStatus(p.Status),
	}
}
