package transformer

import (
	_utils "common/utils"
	"organization/internal/domain/entity"
	"organization/internal/dto"
)

// ToDealActionDomain converts CreateDealActionRequest to domain.DealAction
func ToDealActionDomain(req *dto.CreateDealActionRequest, createdBy, organizationID uint64) *entity.DealAction {
	return &entity.DealAction{
		DealID:          req.DealID,
		TransactionID:   req.TransactionID,
		CustomerID:      req.CustomerID,
		ProductID:       req.ProductIDs,
		Amount:          req.Amount,
		TransactionStep: req.TransactionStep,
		TransactionType: req.TransactionType,
		Note:            req.Note,
		Timestamp:       &req.Timestamp,
		OrganizationID:  organizationID,
	}
}

// ToDealActionResponse converts domain.DealAction to DealActionResponse
func ToDealActionResponse(dealAction *entity.DealAction) *dto.DealActionResponse {
	return &dto.DealActionResponse{
		ID:              dealAction.ID,
		DealID:          dealAction.DealID,
		TransactionID:   dealAction.TransactionID,
		CustomerID:      dealAction.CustomerID,
		ProductIDs:      dealAction.ProductID,
		Amount:          dealAction.Amount,
		TransactionStep: dealAction.TransactionStep,
		TransactionType: dealAction.TransactionType,
		Note:            dealAction.Note,
		Timestamp:       _utils.FormatTimeToString(dealAction.Timestamp),
		OrganizationID:  dealAction.OrganizationID,
	}
}

// ToDealActionResponseList converts slice of domain.DealAction to slice of DealActionResponse
func ToDealActionResponseList(dealActions []*entity.DealAction) []dto.DealActionResponse {
	var responses []dto.DealActionResponse
	for _, dealAction := range dealActions {
		responses = append(responses, *ToDealActionResponse(dealAction))
	}
	return responses
}
