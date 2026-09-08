package domain

import "bdspro/internal/enums"

// TransactionCategory represents a category for transactions
type TxTransactionCategory struct {
	Id   uint32
	Name string
	Type enums.TxMethod
	// OwnerScope enums.EOwnerOf
}
