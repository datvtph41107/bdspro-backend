package transaction

// TransactionCategory represents a category for transactions
type TransactionCategory struct {
	Id           uint32          
	Name         string          
	Type         TransactionType 
	OwnerScope   OwnerScope      
} 