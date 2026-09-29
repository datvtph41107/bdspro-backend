package mapper

// import (
// 	"bdspro/internal/domain/transaction"
// 	bdspropb "pb/types/bdspro"
// )

// type ProductTransactionMapper struct {
// }

// func NewProductTransactionMapper() *ProductTransactionMapper {
// 	return &ProductTransactionMapper{}
// }

// func (m *ProductTransactionMapper) MapProductTransactionToDomain(req *bdspropb.Transaction) *transaction.Transaction {
// 	return &transaction.Transaction{
// 		ID:                   uint64(req.Id),
// 		Name:                 req.TransactionName,
// 		Description:          req.Description,
// 		TransactionType:      req.TransactionType,
// 		CategoryID:           uint64(req.CategoryId),
// 		PaymentMethodID:      uint64(req.PaymentMethodId),
// 		RelatedDealID:        uint64(req.RelatedDealId),
// 		RelatedTransactionID: uint64(req.RelatedTransactionId),
// 	}
// }
