package domain

import "bdspro/internal/enums"

// TransactionCategory represents a category for transactions
type TransactionCategory struct {
	Id         uint32
	Name       string
	Type       enums.TransactionType
	OwnerScope enums.OwnerScope
}
