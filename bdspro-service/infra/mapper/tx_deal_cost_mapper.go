package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	"time"
)

type DealCostMapper struct {
}

func NewDealCostMapper() *DealCostMapper {
	return &DealCostMapper{}
}

func (m *DealCostMapper) DealCostModelToProto(e *domain.TxDealCost) *bdspropb.DealCostDTO {
	result := &bdspropb.DealCostDTO{
		Id:          e.ID,
		DealId:      e.DealID,
		Amount:      e.Amount,
		ConfirmedAt: _utils.FormatTimeToString(e.ConfirmedAt),
		Note:        e.Note,
		CostTypeId:  e.CostTypeID,
		Status:      uint32(e.Status),
		CostName:    e.CostName,
	}

	if e.CostType != nil {
		result.CostTypeName = e.CostType.TypeName
	}

	return result
}

// DealCostProtoToModel maps protobuf DTO to DealCost model
func (m *DealCostMapper) DealCostProtoToModel(p *bdspropb.DealCostDTO) *domain.TxDealCost {
	var confirmedAt *time.Time
	if p.ConfirmedAt != "" {
		confirmedAt = _utils.ParseStringToTime(p.ConfirmedAt)
	}

	return &domain.TxDealCost{
		BaseEntity:  _models.BaseEntity{ID: p.Id},
		DealID:      p.DealId,
		Amount:      p.Amount,
		ConfirmedAt: confirmedAt,
		Note:        p.Note,
		CostTypeID:  p.CostTypeId,
		CostName:    p.CostName,
		Status:      enums.TxApprovedStatus(p.Status),
	}
}
