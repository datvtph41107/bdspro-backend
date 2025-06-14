package transaction

// TransactionCategory represents a category for transactions
type TransactionCategory struct {
	Id           uint32           `json:"id" gorm:"type:uint32;primary_key"`
	Name         string           `json:"name" gorm:"type:varchar(255)"`
	Type         TransactionType  `json:"type" gorm:"type:varchar(20)"`
	OwnerScope   OwnerScope       `json:"owner_scope" gorm:"type:varchar(20)"`
} 