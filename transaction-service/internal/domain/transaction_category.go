package domain

import "tx/internal/enums"

// TransactionCategory represents a category for transactions
type TransactionCategory struct {
	Id   uint32
	Name string
	Type enums.EMethod
	// OwnerScope enums.EOwnerOf
}
