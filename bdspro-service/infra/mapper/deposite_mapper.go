package mapper

type DepositeMapper struct{}

func NewDepositeMapper() *DepositeMapper {
	return &DepositeMapper{}
}

// func (m *DepositeMapper) MapDepositePb(deposite *domain.Deposite) *bdspropb.Deposite {
// 	result := &bdspropb.Deposite{
// 		Id:              deposite.ID,
// 		Amount:          float64(deposite.Amount),
// 		CustomerName:    deposite.CustomerName,
// 		CustomerPhone:   deposite.CustomerPhone,
// 		Note:            deposite.Note,
// 		TransactionType: uint32(deposite.TransactionType),
// 	}

// 	if deposite.ProductID != nil {
// 		result.ProductId = *deposite.ProductID
// 	}

// 	return result
// }

// func (m *DepositeMapper) MapDepositePbToDomain(deposite *bdspropb.DepositeRequest) *tx_domain.Tx {
// 	return &tx_domain.Tx{
// 		BaseEntity:      _models.BaseEntity{ID: deposite.Id},
// 		ProductID:       &deposite.ProductId,
// 		DepositeAmount:  &deposite.DepositeAmount,
// 		DepositeNote:    &deposite.DepositeNote,
// 		TransactionType: enums.TransactionType(deposite.TransactionType),
// 		ContactID:       &deposite.ContactId,
// 	}
// }

// func (m *DepositeMapper) DepositeToDTO(deposite *bdspropb.DepositRequest) *domain.Deposite {
// 	return &domain.Deposite{
// 		BaseEntity:    _models.BaseEntity{ID: deposite.Id},
// 		ProductID:     deposite.ProductId,
// 		Amount:        uint64(deposite.Amount),
// 		CustomerName:  deposite.CustomerName,
// 		CustomerPhone: deposite.CustomerPhone,
// 		Note:          deposite.Note,
// 	}
// }
