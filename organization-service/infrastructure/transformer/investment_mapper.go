package transformer

import (
	_models "common/models"
	_utils "common/utils"
	"organization/internal/domain/entity"
	"organization/internal/enums"
	organizationpb "pb/types/organization"
)

type InvestmentMapper struct {
	DocumentMapper *DocumentMapper
}

func NewInvestmentMapper(documentMapper *DocumentMapper) *InvestmentMapper {
	return &InvestmentMapper{
		DocumentMapper: documentMapper,
	}
}

func (m *InvestmentMapper) SubmitInvestRequestToEntity(req *organizationpb.SubmitInvestRequest) *entity.Investment {
	result := &entity.Investment{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		DealID:             req.DealId,
		MemberID:           req.MemberId,
		Amount:             req.Amount,
		TransferTime:       _utils.ParseStringToTime(req.TransferTime),
		Note:               req.Note,
		TransferProofImage: req.TransferProofImage,
		DoneInvestment:     req.DoneInvestment,
		Status:             enums.ApprovedStatus(req.Status),
	}

	// Handle nullable proxy user ID
	if req.ProxyUserId != nil {
		result.ProxyUserID = req.ProxyUserId
	}

	if len(req.Documents) > 0 {
		result.AttachDocument = make([]entity.AttachDocument, len(req.Documents))
		for i, document := range req.Documents {
			if document != nil {
				result.AttachDocument[i] = *m.DocumentMapper.PbToEntity(document)
			}
		}
	}

	return result
}

func (m *InvestmentMapper) EntityToPb(investment *entity.Investment) *organizationpb.Investment {
	result := &organizationpb.Investment{
		Id:                 investment.ID,
		DealId:             investment.DealID,
		MemberId:           investment.MemberID,
		Amount:             investment.Amount,
		TransferTime:       _utils.FormatTimeToString(investment.TransferTime),
		Note:               investment.Note,
		TransferProofImage: investment.TransferProofImage,
		InvestType:         uint32(investment.InvestType),
		Status:             uint32(investment.Status),
		StatusName:         enums.ApprovedStatusMap[investment.Status],
		ChangedAt:          _utils.FormatTimeToString(investment.TransferTime),
		CreatedBy:          investment.CreatedBy,
		Confirmed:          investment.Confirmed,
	}

	// Handle nullable proxy user ID
	if investment.ProxyUserID != nil {
		result.ProxyUserId = investment.ProxyUserID
	}

	if len(investment.AttachDocument) > 0 {
		result.Documents = make([]*organizationpb.DocumentItem, len(investment.AttachDocument))
		for i, document := range investment.AttachDocument {
			result.Documents[i] = m.DocumentMapper.EntityToPb(&document)
		}
	}

	return result
}

func (m *InvestmentMapper) EntityToPbWithConfirmedHistory(investment *entity.Investment, confirmedInvestments []entity.Investment, totalInvestment float64) *organizationpb.Investment {
	result := m.EntityToPb(investment)

	// Map confirmed investments
	result.Histories = m.EntityToPbWithHistories(confirmedInvestments)

	// Set total investment
	result.TotalInvestment = totalInvestment

	return result
}

func (m *InvestmentMapper) EntityToPbWithHistories(histories []entity.Investment) []*organizationpb.ApprovedInvestmentItem {
	result := make([]*organizationpb.ApprovedInvestmentItem, len(histories))
	for i, inv := range histories {
		result[i] = &organizationpb.ApprovedInvestmentItem{
			Id:           inv.ID,
			Amount:       inv.Amount,
			TransferTime: _utils.FormatTimeToString(inv.TransferTime),
			Note:         inv.Note,
			StatusName:   enums.ApprovedStatusMap[inv.Status],
			Confirmed:    inv.Confirmed,
		}
	}

	return result
}
