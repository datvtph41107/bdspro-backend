package transaction

// TransactionType represents the type of transaction (income/expense)
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

// ApprovalStatus represents the approval status of a transaction
type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusDraft     TransactionStatus = "draft"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
)

// TransactionAction represents the type of action performed on a transaction
type TransactionAction string

const (
	TransactionActionCreated  TransactionAction = "created"
	TransactionActionUpdated  TransactionAction = "updated"
	TransactionActionApproved TransactionAction = "approved"
	TransactionActionCanceled TransactionAction = "canceled"
	TransactionActionDeleted  TransactionAction = "deleted"
)

// OwnerScope represents the scope of ownership
type OwnerScope string

const (
	OwnerScopeSystem      OwnerScope = "system"
	OwnerScopeUserDefined OwnerScope = "user_defined"
)

// PaymentStatus represents the status of a payment method
type PaymentStatus string

const (
	PaymentStatusActive   PaymentStatus = "active"
	PaymentStatusInactive PaymentStatus = "inactive"
) 