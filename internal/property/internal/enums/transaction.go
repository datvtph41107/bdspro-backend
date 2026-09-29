package enums

// TransactionType represents the type of transaction (income/expense)
type TransactionType uint32

const (
	TransactionTypeSale TransactionType = 10 // Bán
	TransactionTypeRent TransactionType = 20 // Thuê
)

// ApprovalStatus represents the approval status of a transaction
type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus uint32

const (
	TransactionStatusDraft     TransactionStatus = 10 // Nháp
	TransactionStatusDeposite  TransactionStatus = 20 // Đặt cọc
	TransactionStatusSigned    TransactionStatus = 30 // Đã ký
	TransactionStatusCompleted TransactionStatus = 40 // Hoàn thành
	TransactionStatusCanceled  TransactionStatus = 50 // Hủy
)

var TransactionStatusNames = map[TransactionStatus]string{
	TransactionStatusDraft:     "Nháp",
	TransactionStatusDeposite:  "Đặt cọc",
	TransactionStatusSigned:    "Đã ký",
	TransactionStatusCompleted: "Hoàn thành",
	TransactionStatusCanceled:  "Hủy",
}

var TransactionTypeNames = map[TransactionType]string{
	TransactionTypeSale: "Bán",
	TransactionTypeRent: "Thuê",
}

// TransactionAction represents the type of action performed on a transaction
type TransactionAction uint32

const (
	TransactionActionSign     TransactionAction = 10
	TransactionActionUpdated  TransactionAction = 20
	TransactionActionApproved TransactionAction = 30
	TransactionActionCanceled TransactionAction = 40
	TransactionActionDeleted  TransactionAction = 50
)

var TransactionActionMap = map[TransactionAction]string{
	TransactionActionSign:     "Ký",
	TransactionActionUpdated:  "updated",
	TransactionActionApproved: "approved",
	TransactionActionCanceled: "canceled",
}

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
