package transformer

import (
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
		DealID:       req.DealId,
		Amount:       req.Amount,
		TransferTime: _utils.ParseStringToTime(req.TransferTime),
		Note:         req.Note,
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
		Id:           investment.ID,
		DealId:       investment.DealID,
		Amount:       investment.Amount,
		TransferTime: _utils.FormatTimeToString(investment.TransferTime),
		Note:         investment.Note,
		InvestType:   uint32(investment.InvestType),
		Status:       uint32(investment.Status),
		StatusName:   enums.ApprovedStatusMap[investment.Status],
		ChangedAt:    _utils.FormatTimeToString(investment.ChangedAt),
		CreatedBy:    investment.CreatedBy,
	}

	if len(investment.AttachDocument) > 0 {
		result.Documents = make([]*organizationpb.DocumentItem, len(investment.AttachDocument))
		for i, document := range investment.AttachDocument {
			result.Documents[i] = m.DocumentMapper.EntityToPb(&document)
		}
	}

	return result
}
